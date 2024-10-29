package server

import (
	"fmt"
	"net/http"
	pathlib "path"
	"regexp"
	"strings"

	"github.com/yates-z/easel/logger"
)

var (
	// regEnLetter matches english letters for http method name
	regEnLetter = regexp.MustCompile("^[A-Z]+$")

	AnyMethods = []string{
		http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete,
		http.MethodHead, http.MethodConnect, http.MethodOptions, http.MethodTrace,
	}
)

type IRouter interface {
	IRoute
	Group(path string, middleware ...Middleware) IRoute
}

type IRoute interface {
	Handle(method, path string, handler HandlerFunc, middlewares ...Middleware)
	ANY(path string, handler HandlerFunc, middlewares ...Middleware)
	GET(path string, handler HandlerFunc, middlewares ...Middleware)
	POST(path string, handler HandlerFunc, middlewares ...Middleware)
	PUT(path string, handler HandlerFunc, middlewares ...Middleware)
	PATCH(path string, handler HandlerFunc, middlewares ...Middleware)
	DELETE(path string, handler HandlerFunc, middlewares ...Middleware)
	HEAD(path string, handler HandlerFunc, middlewares ...Middleware)
	OPTIONS(path string, handler HandlerFunc, middlewares ...Middleware)

	StaticFile(string, string)
	StaticFileFS(string, string, http.FileSystem)
	Static(string, string)
	StaticFS(string, http.FileSystem)
}

var _ IRouter = (*Router)(nil)

type Router struct {
	basePath    string
	server      *Server
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewRouter(s *Server) *Router {
	r := &Router{
		server: s,
		mux:    http.NewServeMux(),
	}
	return r
}

func (r *Router) Handle(method, path string, handler HandlerFunc, middlewares ...Middleware) {
	if matched := regEnLetter.MatchString(method); !matched {
		panic("http method " + method + " is not valid")
	}

	fullPath := r.joinPaths(r.basePath, path)
	pattern := fmt.Sprintf("%s %s", method, fullPath)

	if r.server.showInfo {
		logger.Info(pattern)
	}

	middlewares = append(middlewares, r.middlewares...)
	_handler := chain(middlewares...)(handler)

	entrance := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		ctx := req.Context().(*Context)
		ctx.fullPath = fullPath
		if err := _handler(ctx); err != nil {
			r.server.errorHandler(ctx, err)
		}
	})

	r.mux.Handle(pattern, entrance)
}

// Group implements IRouter.
func (r *Router) Group(path string, middleware ...Middleware) IRoute {
	return &Router{
		basePath:    r.joinPaths(r.basePath, path),
		server:      r.server,
		mux:         r.mux,
		middlewares: append(r.middlewares, middleware...),
	}
}

func (r *Router) HEAD(path string, handler HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodHead, path, handler, middlewares...)
}

func (r *Router) GET(path string, handler HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodGet, path, handler, middlewares...)
}

func (r *Router) POST(path string, handler HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPost, path, handler, middlewares...)
}

func (r *Router) PUT(path string, handler HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPut, path, handler, middlewares...)
}

func (r *Router) PATCH(path string, handler HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPatch, path, handler, middlewares...)
}

func (r *Router) DELETE(path string, handler HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodDelete, path, handler, middlewares...)
}

func (r *Router) OPTIONS(path string, handler HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodOptions, path, handler, middlewares...)
}

func (r *Router) ANY(path string, handler HandlerFunc, middlewares ...Middleware) {
	r.Handle("", path, handler, middlewares...)
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
func (r *Router) StaticFile(path, filePath string) {
	r.staticFileHandler(path, func(c *Context) error {
		c.File(filePath)
		return nil
	})
}

// StaticFileFS works just like `StaticFile` but a custom `http.FileSystem` can be used instead..
func (r *Router) StaticFileFS(path, filePath string, fs http.FileSystem) {
	r.staticFileHandler(path, func(c *Context) error {
		c.FileFromFS(filePath, fs)
		return nil
	})
}

func (r *Router) staticFileHandler(path string, handler HandlerFunc) {
	if strings.Contains(path, ":") || strings.Contains(path, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	r.GET(path, handler)
	r.HEAD(path, handler)
}

// Static serves files from the given file system root.
func (r *Router) Static(path, root string) {
	r.StaticFS(path, Dir(root, false))
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
func (r *Router) StaticFS(path string, fs http.FileSystem) {
	if strings.Contains(path, ":") || strings.Contains(path, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := r.createStaticHandler(path, fs)
	urlPattern := pathlib.Join(path, "/{filepath...}")

	// Register GET and HEAD handlers
	r.GET(urlPattern, handler)
	r.HEAD(urlPattern, handler)
}

func (r *Router) createStaticHandler(path string, fs http.FileSystem) HandlerFunc {
	absolutePath := r.joinPaths(r.basePath, path)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) error {
		if _, noListing := fs.(*OnlyFilesFS); noListing {
			c.Response.WriteHeader(http.StatusNotFound)
		}

		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			c.Response.WriteHeader(http.StatusNotFound)
			return err
		}
		f.Close()

		fileServer.ServeHTTP(c.Response, c.Request)
		return nil
	}
}

// utils from https://github.com/gin-gonic/gin
func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func (r *Router) joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := pathlib.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}
