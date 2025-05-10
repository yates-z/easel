package main

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yates-z/easel/registry"
	"github.com/yates-z/easel/transport"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var APP_KEY = struct{}{}

// AppInfo is application context value.
type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoints() []string
}

// Option is an application option.
type Option func(app *Application)

// ID with service id.
func ID(id string) Option {
	return func(app *Application) { app.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(app *Application) { app.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(app *Application) { app.version = version }
}

// Metadata with service metadata.
func Metadata(md map[string]string) Option {
	return func(app *Application) { app.metadata = md }
}

// Endpoint with service endpoint.
func Endpoint(endpoints ...*url.URL) Option {
	return func(app *Application) { app.endpoints = endpoints }
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(app *Application) { app.ctx = ctx }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(app *Application) { app.sigs = sigs }
}

// Registrar with service registry.
func Registrar(r registry.Registrar) Option {
	return func(app *Application) { app.registrar = r }
}

// RegistrarTimeout with registrar timeout.
func RegistrarTimeout(t time.Duration) Option {
	return func(app *Application) { app.registrarTimeout = t }
}

// StopTimeout with app stop timeout.
func StopTimeout(t time.Duration) Option {
	return func(app *Application) { app.stopTimeout = t }
}

// Server with transport servers.
func Server(srv ...transport.Server) Option {
	return func(app *Application) { app.servers = srv }
}

// BeforeStart run funcs before app starts
func BeforeStart(fn func(context.Context) error) Option {
	return func(app *Application) {
		app.beforeStart = append(app.beforeStart, fn)
	}
}

// BeforeStop run funcs before app stops
func BeforeStop(fn func(context.Context) error) Option {
	return func(app *Application) {
		app.beforeStop = append(app.beforeStop, fn)
	}
}

// AfterStart run funcs after app starts
func AfterStart(fn func(context.Context) error) Option {
	return func(app *Application) {
		app.afterStart = append(app.afterStart, fn)
	}
}

// AfterStop run funcs after app stops
func AfterStop(fn func(context.Context) error) Option {
	return func(app *Application) {
		app.afterStop = append(app.afterStop, fn)
	}
}

// Application is an application components lifecycle manager.
type Application struct {
	baseCtx context.Context
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.Mutex

	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL
	instance  *registry.ServiceInstance

	sigs []os.Signal

	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration

	servers []transport.Server
	// Before and After funcs
	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// New create an application lifecycle manager.
func New(opts ...Option) *Application {

	app := &Application{
		baseCtx:          context.Background(),
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		app.id = id.String()
	}
	for _, opt := range opts {
		opt(app)
	}
	ctx, cancel := context.WithCancel(app.baseCtx)
	app.ctx = ctx
	app.cancel = cancel
	return app
}

// ID returns app instance id.
func (app *Application) ID() string { return app.id }

// Name returns service name.
func (app *Application) Name() string { return app.name }

// Version returns app version.
func (app *Application) Version() string { return app.version }

// Metadata returns service metadata.
func (app *Application) Metadata() map[string]string { return app.metadata }

// Endpoints returns endpoints.
func (app *Application) Endpoints() []string {
	if app.instance != nil {
		return app.instance.Endpoints
	}
	return nil
}

func (app *Application) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0, len(app.endpoints))
	for _, e := range app.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range app.servers {
			e, err := srv.Endpoint()
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, e.String())
		}
	}
	return &registry.ServiceInstance{
		ID:        app.id,
		Name:      app.name,
		Version:   app.version,
		Metadata:  app.metadata,
		Endpoints: endpoints,
	}, nil
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (app *Application) Run() error {
	instance, err := app.buildInstance()
	if err != nil {
		return err
	}
	app.mu.Lock()
	app.instance = instance
	app.mu.Unlock()
	sctx := context.WithValue(app.ctx, APP_KEY, app)
	eg, ctx := errgroup.WithContext(sctx)
	wg := sync.WaitGroup{}

	for _, fn := range app.beforeStart {
		if err = fn(sctx); err != nil {
			return err
		}
	}

	octx := context.WithValue(app.baseCtx, APP_KEY, app)
	for _, srv := range app.servers {
		server := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx := octx
			if app.stopTimeout > 0 {
				var cancel context.CancelFunc
				stopCtx, cancel = context.WithTimeout(stopCtx, app.stopTimeout)
				defer cancel()
			}
			return server.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done() // here is to ensure server start has begun running before register, so defer is not needed
			return server.Start(octx)
		})
	}
	wg.Wait()

	if app.registrar != nil {
		rctx, rcancel := context.WithTimeout(ctx, app.registrarTimeout)
		defer rcancel()
		if err = app.registrar.Register(rctx, instance); err != nil {
			return err
		}
	}
	for _, fn := range app.afterStart {
		if err = fn(sctx); err != nil {
			return err
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, app.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return app.Stop()
		}
	})
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	err = nil
	for _, fn := range app.afterStop {
		err = fn(sctx)
	}
	return err
}

// Stop gracefully stops the application.
func (app *Application) Stop() (err error) {
	sctx := context.WithValue(app.ctx, APP_KEY, app)
	for _, fn := range app.beforeStop {
		err = fn(sctx)
	}

	app.mu.Lock()
	instance := app.instance
	app.mu.Unlock()
	if app.registrar != nil && instance != nil {
		ctx, cancel := context.WithTimeout(context.WithValue(app.ctx, APP_KEY, app), app.registrarTimeout)
		defer cancel()
		if err = app.registrar.Deregister(ctx, instance); err != nil {
			return err
		}
	}
	if app.cancel != nil {
		app.cancel()
	}
	return err
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(APP_KEY).(AppInfo)
	return
}
