package fiber

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/orzmoe/httpgen"
	"go.uber.org/fx"
)

// RouteRegister Fiber路由注册实现
type Server struct {
	app    *fiber.App
	config httpgen.HttpConfig
}

func NewServer(lc fx.Lifecycle, config httpgen.HttpConfig) httpgen.HTTPServer {
	app := fiber.New(
		fiber.Config{
			ReadTimeout:  config.GetReadTimeout(),
			WriteTimeout: config.GetWriteTimeout(),
		},
	)
	server := &Server{app: app, config: config}
	lc.Append(fx.StartHook(func() error {
		go server.Start()
		return nil
	}))
	lc.Append(fx.StopHook(func() error {
		return server.Shutdown()
	}))
	return server
}

func (s *Server) RouteRegister() httpgen.RouteRegister {
	return &Server{app: s.app}
}

func (r *Server) Get(path string, handler httpgen.HandlerFunc) {
	r.app.Get(path, wrapHandler(handler))
}

func (r *Server) Post(path string, handler httpgen.HandlerFunc) {
	r.app.Post(path, wrapHandler(handler))
}

func (r *Server) Put(path string, handler httpgen.HandlerFunc) {
	r.app.Put(path, wrapHandler(handler))
}

func (r *Server) Delete(path string, handler httpgen.HandlerFunc) {
	r.app.Delete(path, wrapHandler(handler))
}

func (r *Server) Group(path string) httpgen.RouteGroup {
	return &RouterGroup{group: r.app.Group(path)}
}

func (r *Server) Add(method []string, path string, handler httpgen.HandlerFunc) {
	r.app.Add(method, path, wrapHandler(handler))
}

func (r *Server) Use(middleware ...any) {
	for _, m := range middleware {
		switch middleware := m.(type) {
		case httpgen.NativeMiddleware:
			r.app.Use(middleware.Native())
		case httpgen.HandlerFunc:
			r.app.Use(wrapHandler(middleware))
		}
	}
}

func (r *Server) NativeEngine() any {
	return r.app
}

// 原生中间件适配器
type FiberMiddleware struct {
	handler fiber.Handler
}

func WrapNativeMiddleware(handler fiber.Handler) httpgen.NativeMiddleware {
	return &FiberMiddleware{handler: handler}
}
func (m *FiberMiddleware) Native() interface{} {
	return m.handler
}

type RouterGroup struct {
	group fiber.Router
}

func (g *RouterGroup) Get(path string, handler httpgen.HandlerFunc) {
	g.group.Get(path, wrapHandler(handler))
}

func (g *RouterGroup) Post(path string, handler httpgen.HandlerFunc) {
	g.group.Post(path, wrapHandler(handler))
}

func (g *RouterGroup) Put(path string, handler httpgen.HandlerFunc) {
	g.group.Put(path, wrapHandler(handler))
}

func (g *RouterGroup) Delete(path string, handler httpgen.HandlerFunc) {
	g.group.Delete(path, wrapHandler(handler))
}

func (g *RouterGroup) Group(path string) httpgen.RouteGroup {
	return &RouterGroup{group: g.group.Group(path)}
}

func (g *RouterGroup) Use(middlewares ...any) {
	for _, m := range middlewares {
		switch middleware := m.(type) {
		case httpgen.Middleware:
			// 使用抽象中间件
			g.group.Use(wrapMiddleware(middleware))
		case httpgen.NativeMiddleware:
			// 使用原生中间件
			if handler, ok := middleware.Native().(fiber.Handler); ok {
				g.group.Use(handler)
			}
		case fiber.Handler:
			// 直接使用 Fiber 中间件
			g.group.Use(middleware)
		}
	}

}

func (g *RouterGroup) Add(method []string, path string, handler httpgen.HandlerFunc) {
	g.group.Add(method, path, wrapHandler(handler))
}

func (s *Server) Start() error {
	return s.app.Listen(s.config.GetAddr())
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func wrapHandler(h httpgen.HandlerFunc) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := &Context{ctx: c}
		h(ctx)
		return nil
	}
}

func wrapMiddleware(m httpgen.Middleware) fiber.Handler {
	return func(c fiber.Ctx) error {
		// 创建上下文适配器
		ctx := &Context{ctx: c}

		// 包装 next 处理函数
		var err error
		next := func(ctx httpgen.Context) error {
			err = c.Next()
			return err
		}

		// 执行中间件
		wrapped := m.Wrap(next)
		wrapped(ctx)

		return err
	}
}

// Context 适配器
type Context struct {
	ctx fiber.Ctx
}

func (c *Context) Param(key string) string {
	return c.ctx.Params(key)
}

func (c *Context) Query(key string) string {
	return c.ctx.Query(key)
}

func (c *Context) BindJSON(v any) error {
	return c.ctx.Bind().JSON(v)
}
func (c *Context) BindQuery(v any) error {
	return c.ctx.Bind().Query(v)
}

func (c *Context) BindURI(v any) error {
	return c.ctx.Bind().URI(v)
}

func (c *Context) BindBody(v any) error {
	return c.ctx.Bind().Body(v)
}

func (c *Context) JSON(code int, v any) error {
	return c.ctx.Status(code).JSON(v)
}

func (c *Context) Body() []byte {
	return c.ctx.Body()
}

func (c *Context) GetHeader(key string) string {
	return c.ctx.Get(key)
}

func (c *Context) SetHeader(key, value string) {
	c.ctx.Set(key, value)
}

func (c *Context) Status(code int) httpgen.Context {
	c.ctx.Status(code)
	return c
}

func (c *Context) String(code int, s string) error {
	return c.ctx.Status(code).SendString(s)
}

func (c *Context) Next() error {
	return c.ctx.Next()
}

func (c *Context) Path() string {
	return c.ctx.Path()
}

func (c *Context) Method() string {
	return c.ctx.Method()
}

func (c *Context) Get(key string) string {
	return c.ctx.Get(key)
}

func (c *Context) GetReqHeaders() map[string][]string {
	return c.ctx.GetReqHeaders()
}

func (c *Context) Context() context.Context {
	return c.ctx.Context()
}

func (c *Context) Set(key, value string) {
	c.ctx.Set(key, value)
}

var Module = fx.Module("fiber",
	fx.Provide(
		NewServer,
	),
)
