package httpgen

import (
	"context"
	"time"
)

// RouteRegister 定义统一的路由注册方法
type RouteRegister interface {
	Get(path string, handler HandlerFunc)
	Post(path string, handler HandlerFunc)
	Put(path string, handler HandlerFunc)
	Delete(path string, handler HandlerFunc)
	Group(path string) RouteGroup
	Use(middleware ...any)
	Add(method []string, path string, handler HandlerFunc)
}
type RouteGroup interface {
	Get(path string, handler HandlerFunc)
	Post(path string, handler HandlerFunc)
	Put(path string, handler HandlerFunc)
	Delete(path string, handler HandlerFunc)
	Group(path string) RouteGroup
	Use(middleware ...any)
	Add(method []string, path string, handler HandlerFunc)
}

// HTTPServer 定义 HTTP 服务器接口
type HTTPServer interface {
	Start() error
	Shutdown() error
	RouteRegister() RouteRegister
	NativeEngine() any
}

// 中间件接口
type Middleware interface {
	Wrap(HandlerFunc) HandlerFunc
}

// 原生中间件
type NativeMiddleware interface {
	Native() any
}

// HandlerFunc 通用处理函数类型
type HandlerFunc = func(Context) error

// Context 请求上下文接口
type Context interface {
	// 通用上下文方法
	Param(key string) string
	Query(key string) string
	Body() []byte
	BindJSON(v any) error
	BindQuery(v any) error
	BindURI(v any) error
	BindBody(v any) error
	JSON(code int, v any) error
	String(code int, s string) error
	Status(code int) Context
	GetHeader(key string) string
	SetHeader(key, value string)
	Next() error
	Path() string
	Method() string
	Get(key string) string
	GetReqHeaders() map[string][]string
	Context() context.Context
	Set(key, value string)
}

type HttpConfig interface {
	GetAddr() string
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
}
