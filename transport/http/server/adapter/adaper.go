package adapter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/transport/grpc/encoding/form"
	"github.com/yates-z/easel/transport/grpc/encoding/json"
	"github.com/yates-z/easel/transport/grpc/encoding/proto"
	"github.com/yates-z/easel/transport/grpc/encoding/xml"
	"github.com/yates-z/easel/transport/http/server"
	"google.golang.org/grpc/encoding"
	"io"
	"net/url"
	"regexp"
	"strings"
)

var (
	_ = form.Name
	_ = json.Name
	_ = xml.Name
	_ = proto.Name
)

func GRPC[T1 any, T2 any](f func(context.Context, *T1) (*T2, error)) server.HandlerFunc {
	return func(ctx *server.Context) error {
		var in T1
		if ctx.HasBody() {
			if err := bindBody(ctx, &in); err != nil {
				return err
			}
		}
		if err := bindQuery(ctx, &in); err != nil {
			return err
		}

		if err := bindParams(ctx, &in); err != nil {
			return err
		}

		reply, err := f(ctx, &in)
		if err != nil {
			return err
		}
		codec, _ := codecForRequest(ctx, "Accept")
		data, err := codec.Marshal(reply)
		if err != nil {
			return err
		}
		ctx.Response.Header().Set("Content-Type", "application/"+codec.Name())
		_, err = ctx.Response.Write(data)
		if err != nil {
			return err
		}
		return nil
	}
}

// bindBody decodes the request body to object.
func bindBody(ctx *server.Context, v interface{}) error {
	codec, ok := codecForRequest(ctx, "Content-Type")
	if !ok {
		return errors.New(fmt.Sprintf("unregister Content-Type: %s", ctx.GetHeader("Content-Type")))
	}

	data, err := io.ReadAll(ctx.Request.Body)

	// reset body.
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	if err = codec.Unmarshal(data, v); err != nil {
		return errors.New(fmt.Sprintf("body unmarshal %s", err.Error()))
	}
	return nil
}

// bindQuery bind url parameters to object.
func bindQuery(ctx *server.Context, v interface{}) error {
	queries := ctx.Request.URL.Query().Encode()
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(queries), v); err != nil {
		return err
	}
	return nil
}

// bindParams bind path parameters to object.
func bindParams(ctx *server.Context, v interface{}) error {
	if strings.HasSuffix(ctx.FullPath(), "/") {
		logger.Warnf("Path %s should not end with \"/\"", ctx.FullPath())
	}
	pattern := regexp.MustCompile(`(?i){([a-z.0-9_\s]*)=?([^{}]*)}`)
	matches := pattern.FindAllStringSubmatch(ctx.FullPath(), -1)
	res := make(map[string]*string, len(matches))
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if len(name) > 1 && len(m[2]) > 0 {
			res[name] = &m[2]
		} else {
			res[name] = nil
		}
	}
	if len(res) == 0 {
		return nil
	}
	params := make(url.Values, len(res))
	for k, _ := range res {
		params[k] = []string{ctx.Param(k)}
	}

	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(params.Encode()), v); err != nil {
		return err
	}
	return nil
}

// codecForRequest get encoding.Codec via http.Request
func codecForRequest(ctx *server.Context, name string) (encoding.Codec, bool) {
	for _, accept := range ctx.Request.Header[name] {
		codec := encoding.GetCodec(contentSubtype(accept))
		if codec != nil {
			return codec, true
		}
	}
	return encoding.GetCodec("json"), false
}

func contentSubtype(contentType string) string {
	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}
	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}
