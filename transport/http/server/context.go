package server

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/yates-z/easel/logger"
	"google.golang.org/protobuf/proto"
)

type HandlerFunc func(*Context) error

var _ context.Context = (*Context)(nil)

// Context is an HTTP request Context. It defines core functions sets of this http.server.
type Context struct {
	Request  *http.Request
	Response *response

	ctx context.Context

	server   *Server
	fullPath string
	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite

	// storage is a key/value pair.
	storage map[string]any
	// This mutex protects storage map.
	mu sync.RWMutex
}

func newContext(s *Server) *Context {
	ctx := &Context{server: s, ctx: context.Background()}
	return ctx
}

func (c *Context) WithBaseContext(ctx context.Context) {
	if ctx == nil {
		return
	}
	c.ctx = ctx
}

func (c *Context) init(req *http.Request, resp http.ResponseWriter) {
	c.Request = req
	c.Response = &response{ResponseWriter: resp, statusCode: http.StatusOK}
}

func (c *Context) reset() {
	c.Request = nil
	c.Response.reset(nil)
	c.ctx = context.Background()
	c.fullPath = ""
	c.sameSite = 0
	c.storage = nil
}

func (c *Context) Logger() logger.Logger {
	return c.server.log
}

/********************************************/
/********* implement context.Context ********/
/********************************************/

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.storage == nil {
		c.storage = make(map[string]any)
	}

	c.storage[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
func (c *Context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.storage[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) any {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

func getTyped[T any](c *Context, key string) (res T) {
	if val, ok := c.Get(key); ok && val != nil {
		res, _ = val.(T)
	}
	return
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	return getTyped[string](c, key)
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	return getTyped[bool](c, key)
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	return getTyped[int](c, key)
}

// GetInt8 returns the value associated with the key as an integer 8.
func (c *Context) GetInt8(key string) (i8 int8) {
	return getTyped[int8](c, key)
}

// GetInt16 returns the value associated with the key as an integer 16.
func (c *Context) GetInt16(key string) (i16 int16) {
	return getTyped[int16](c, key)
}

// GetInt32 returns the value associated with the key as an integer 32.
func (c *Context) GetInt32(key string) (i32 int32) {
	return getTyped[int32](c, key)
}

// GetInt64 returns the value associated with the key as an integer 64.
func (c *Context) GetInt64(key string) (i64 int64) {
	return getTyped[int64](c, key)
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint(key string) (ui uint) {
	return getTyped[uint](c, key)
}

// GetUint8 returns the value associated with the key as an unsigned integer 8.
func (c *Context) GetUint8(key string) (ui8 uint8) {
	return getTyped[uint8](c, key)
}

// GetUint16 returns the value associated with the key as an unsigned integer 16.
func (c *Context) GetUint16(key string) (ui16 uint16) {
	return getTyped[uint16](c, key)
}

// GetUint32 returns the value associated with the key as an unsigned integer 32.
func (c *Context) GetUint32(key string) (ui32 uint32) {
	return getTyped[uint32](c, key)
}

// GetUint64 returns the value associated with the key as an unsigned integer 64.
func (c *Context) GetUint64(key string) (ui64 uint64) {
	return getTyped[uint64](c, key)
}

// GetFloat32 returns the value associated with the key as a float32.
func (c *Context) GetFloat32(key string) (f32 float32) {
	return getTyped[float32](c, key)
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	return getTyped[float64](c, key)
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	return getTyped[time.Time](c, key)
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	return getTyped[time.Duration](c, key)
}

// GetIntSlice returns the value associated with the key as a slice of integers.
func (c *Context) GetIntSlice(key string) (is []int) {
	return getTyped[[]int](c, key)
}

// GetInt8Slice returns the value associated with the key as a slice of int8 integers.
func (c *Context) GetInt8Slice(key string) (i8s []int8) {
	return getTyped[[]int8](c, key)
}

// GetInt16Slice returns the value associated with the key as a slice of int16 integers.
func (c *Context) GetInt16Slice(key string) (i16s []int16) {
	return getTyped[[]int16](c, key)
}

// GetInt32Slice returns the value associated with the key as a slice of int32 integers.
func (c *Context) GetInt32Slice(key string) (i32s []int32) {
	return getTyped[[]int32](c, key)
}

// GetInt64Slice returns the value associated with the key as a slice of int64 integers.
func (c *Context) GetInt64Slice(key string) (i64s []int64) {
	return getTyped[[]int64](c, key)
}

// GetUintSlice returns the value associated with the key as a slice of unsigned integers.
func (c *Context) GetUintSlice(key string) (uis []uint) {
	return getTyped[[]uint](c, key)
}

// GetUint8Slice returns the value associated with the key as a slice of uint8 integers.
func (c *Context) GetUint8Slice(key string) (ui8s []uint8) {
	return getTyped[[]uint8](c, key)
}

// GetUint16Slice returns the value associated with the key as a slice of uint16 integers.
func (c *Context) GetUint16Slice(key string) (ui16s []uint16) {
	return getTyped[[]uint16](c, key)
}

// GetUint32Slice returns the value associated with the key as a slice of uint32 integers.
func (c *Context) GetUint32Slice(key string) (ui32s []uint32) {
	return getTyped[[]uint32](c, key)
}

// GetUint64Slice returns the value associated with the key as a slice of uint64 integers.
func (c *Context) GetUint64Slice(key string) (ui64s []uint64) {
	return getTyped[[]uint64](c, key)
}

// GetFloat32Slice returns the value associated with the key as a slice of float32 numbers.
func (c *Context) GetFloat32Slice(key string) (f32s []float32) {
	return getTyped[[]float32](c, key)
}

// GetFloat64Slice returns the value associated with the key as a slice of float64 numbers.
func (c *Context) GetFloat64Slice(key string) (f64s []float64) {
	return getTyped[[]float64](c, key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	return getTyped[[]string](c, key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]any) {
	return getTyped[map[string]any](c, key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	return getTyped[map[string]string](c, key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	return getTyped[map[string][]string](c, key)
}

/***************************/
/********* REQUEST ********/
/***************************/

// Param returns the value for the named path wildcard in the [ServeMux] pattern
// that matched the request.
func (c *Context) Param(key string) string {
	return c.Request.PathValue(key)
}

// Query gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string.
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) QueryArray(key string) ([]string, bool) {
	values, ok := c.Request.URL.Query()[key]
	return values, ok
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
