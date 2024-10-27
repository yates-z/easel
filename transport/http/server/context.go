package server

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/yates-z/easel/logger"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"
)

type HandlerFunc func(*Context) error

var _ context.Context = (*Context)(nil)

// Context is an HTTP request Context. It defines core functions sets of this http.server.
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter

	server   *Server
	fullPath string
	params   []string
	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite
}

func newContext(req *http.Request, resp http.ResponseWriter, s *Server) *Context {
	return &Context{
		Request:  req,
		Response: resp,
		server:   s,
	}
}

/********************************************/
/********* implement context.Context ********/
/********************************************/

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	ctx := c.Request.Context()
	return ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	ctx := c.Request.Context()
	return ctx.Done()
}

func (c *Context) Err() error {
	ctx := c.Request.Context()
	return ctx.Err()
}

func (c *Context) Value(key any) any {
	ctx := c.Request.Context()
	return ctx.Value(key)
}

/***************************/
/********* REQUEST ********/
/***************************/

func (c *Context) Param(key string) string {
	return c.Request.PathValue(key)
}

func (c *Context) FullPath() string {
	return c.fullPath
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	content := c.Request.Header.Get("Content-Type")
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

// IsWebsocket returns true if the request headers indicate that a websocket
// handshake is being initiated by the client.
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.Request.Header.Get("Connection")), "upgrade") &&
		strings.EqualFold(c.Request.Header.Get("Upgrade"), "websocket") {
		return true
	}
	return false
}

// RemoteIP parses the IP from Request.RemoteAddr, normalizes and returns the IP (without the port).
func (c *Context) RemoteIP() string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

// HasBody checks if the http request has a request body.
func (c *Context) HasBody() bool {
	if slices.Contains([]string{http.MethodGet, http.MethodDelete, http.MethodHead}, c.Request.Method) {
		return false
	}
	if c.Request.ContentLength == 0 {
		return false
	}
	if c.Request.Body == nil {
		return false
	}
	return true
}

// HasParams checks if the path has parameters.
func (c *Context) HasParams() bool {
	if strings.HasSuffix(c.FullPath(), "/") {
		logger.Warnf("Path %s should not end with \"/\"", c.FullPath())
	}
	pattern := regexp.MustCompile(`(?i){([a-z.0-9_\s]*)=?([^{}]*)}`)
	matches := pattern.FindAllStringSubmatch(c.FullPath(), -1)
	res := make(map[string]*string, len(matches))
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if len(name) > 1 && len(m[2]) > 0 {
			res[name] = &m[2]
		} else {
			res[name] = nil
		}
	}
	return len(res) > 0
}

/***************************/
/********* RESPONSE ********/
/***************************/

// SetStatus sets the HTTP response code.
func (c *Context) SetStatus(code int) {
	c.Response.WriteHeader(code)
}

// SetSameSite with cookie
func (c *Context) SetSameSite(samesite http.SameSite) {
	c.sameSite = samesite
}

// SetHeader is an intelligent shortcut for c.Response.Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `c.Response.Header().Del(key)`
func (c *Context) SetHeader(key, value string) {
	if value == "" {
		c.Response.Header().Del(key)
		return
	}
	c.Response.Header().Set(key, value)
}

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// GetRawData returns stream data.
func (c *Context) GetRawData() ([]byte, error) {
	if c.Request.Body == nil {
		return nil, errors.New("cannot read nil body")
	}
	return io.ReadAll(c.Request.Body)
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Response, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: c.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// GetCookie returns the named cookie provided in the request or
// ErrNoCookie if not found.
func (c *Context) GetCookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

// String writes the given string into the response body.
func (c *Context) String(code int, text string) error {
	c.SetContentType([]string{"text/plain; charset=utf-8"})
	c.SetStatus(code)

	_, err := c.Response.Write([]byte(text))
	if err != nil {
		return err
	}
	return nil
}

// JSON serializes the given struct as JSON into the response body.
func (c *Context) JSON(code int, obj any) error {
	c.SetContentType([]string{"application/json; charset=utf-8"})
	c.SetStatus(code)

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = c.Response.Write(jsonBytes)
	return err
}

// HTML renders the HTTP template specified by its file name.
func (c *Context) HTML(code int, name string, data any) error {
	c.SetContentType([]string{"text/html; charset=utf-8"})
	c.SetStatus(code)
	if name == "" {
		return c.server.Template().Execute(c.Response, data)
	}
	return c.server.Template().ExecuteTemplate(c.Response, name, data)
}

// XML renders the HTTP template specified by its file name.
func (c *Context) XML(code int, data any) error {
	c.SetContentType([]string{"application/xml; charset=utf-8"})
	c.SetStatus(code)
	return xml.NewEncoder(c.Response).Encode(data)
}

// Proto serializes the given struct as ProtoBuf into the response body.
func (c *Context) Proto(code int, data any) error {
	c.SetContentType([]string{"application/x-protobuf"})
	c.SetStatus(code)
	if value, ok := data.(proto.Message); ok {
		bytes, err := proto.Marshal(value)
		if err != nil {
			return err
		}
		_, err = c.Response.Write(bytes)
		return err
	}
	return errors.New("not a proto message")
}

// Redirect returns an HTTP redirect to the specific location.
func (c *Context) Redirect(url string) {
	http.Redirect(c.Response, c.Request, url, http.StatusMovedPermanently)
}

// File writes the specified file into the body stream in an efficient way.
func (c *Context) File(filepath string) {
	http.ServeFile(c.Response, c.Request, filepath)
}

// FileFromFS writes the specified file from http.FileSystem into the body stream in an efficient way.
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.Request.URL.Path = old
	}(c.Request.URL.Path)

	c.Request.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.Response, c.Request)
}

// SetContentType writes ContentType.
func (c *Context) SetContentType(contentType []string) {
	header := c.Response.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = contentType
	}
}
