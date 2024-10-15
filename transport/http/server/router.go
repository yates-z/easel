package server

import (
	"fmt"
	"net/http"
)

var AnyMethods = []string{
	http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete,
	http.MethodHead, http.MethodConnect, http.MethodOptions, http.MethodTrace,
}

type IRouter interface {
	IRoute
	Group(string, handler HandlerFunc, middleware ...Middleware) IRoute
}

type IRoute interface {
	ANY(path string, handler HandlerFunc, middleware ...Middleware)
	GET(path string, handler HandlerFunc, middleware ...Middleware)
	POST(path string, handler HandlerFunc, middleware ...Middleware)
	PUT(path string, handler HandlerFunc, middleware ...Middleware)
	PATCH(path string, handler HandlerFunc, middleware ...Middleware)
	DELETE(path string, handler HandlerFunc, middleware ...Middleware)
	HEAD(path string, handler HandlerFunc, middleware ...Middleware)
	OPTIONS(path string, handler HandlerFunc, middleware ...Middleware)

	StaticFile(string, string)
	StaticFileFS(string, string, http.FileSystem)
	Static(string, string)
	StaticFS(string, http.FileSystem)
}

var _ IRouter = (*Router)(nil)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	r := &Router{
		mux: http.NewServeMux(),
	}
	return r
}

func (r *Router) Handle(method, path string, handler HandlerFunc, middleware ...Middleware) {
	pattern := fmt.Sprintf("%s %s", method, path)

	final := func(resp http.ResponseWriter, req *http.Request) {
		handler(newContext(req, resp))
	}

	// chain := func(middleware ...Middleware) Middleware {
	// 	return func(next http.Handler) http.Handler {
	// 		for i := len(middleware) - 1; i >= 0; i-- {
	// 			next = middleware[i](next)
	// 		}
	// 		return next
	// 	}
	// }
	// c := chain(middleware...)(http.Handler(http.HandlerFunc(final)))

	// r.mux.Handle(pattern, c)
	r.mux.HandleFunc(pattern, final)
}

// Group implements IRouter.
func (r *Router) Group(string HandlerFunc, handler HandlerFunc, middleware ...Middleware) IRoute {
	panic("unimplemented")
}

func (r *Router) HEAD(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Handle(http.MethodHead, path, handler, middleware...)
}

func (r *Router) GET(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Handle(http.MethodGet, path, handler, middleware...)
}

func (r *Router) POST(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Handle(http.MethodPost, path, handler, middleware...)
}

func (r *Router) PUT(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Handle(http.MethodPut, path, handler, middleware...)
}

func (r *Router) PATCH(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Handle(http.MethodPatch, path, handler, middleware...)
}

func (r *Router) DELETE(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Handle(http.MethodDelete, path, handler, middleware...)
}

func (r *Router) OPTIONS(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Handle(http.MethodOptions, path, handler, middleware...)
}

func (r *Router) ANY(path string, handler HandlerFunc, middleware ...Middleware) {
	for _, method := range AnyMethods {
		r.Handle(method, path, handler, middleware...)
	}

}

func (r *Router) StaticFile(s string, s2 string) {
	//TODO implement me
	panic("implement me")
}

func (r *Router) StaticFileFS(s string, s2 string, system http.FileSystem) {
	//TODO implement me
	panic("implement me")
}

func (r *Router) Static(s string, s2 string) {
	//TODO implement me
	panic("implement me")
}

func (r *Router) StaticFS(s string, system http.FileSystem) {
	//TODO implement me
	panic("implement me")
}
