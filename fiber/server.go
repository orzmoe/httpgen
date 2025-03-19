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

func (g *RouterGroup) Use(middleware ...httpgen.HandlerFunc) {
	for _, m := range middleware {
		g.group.Use(wrapHandler(m))
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

var Module = fx.Module("fiber",
	fx.Provide(
		NewServer,
	),
)
