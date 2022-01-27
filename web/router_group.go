package web

type RouteGroup struct {
	// 路由分组前缀
	prefix string
	engine *Engine
	// 路由中间件
	middleware []Service
}

func NewRouteGroup(prefix string, e *Engine) *RouteGroup {
	r := &RouteGroup{
		prefix: prefix,
		engine: e,
	}
	// 添加默认路由分组,panic恢复和日志组建
	r.Use(Recover)
	r.Use(Log)
	return r
}

func (r *RouteGroup) Use(srv Service) {
	r.middleware = append(r.middleware, srv)
}

func (r *RouteGroup) Get(pattern string, srv Service) {
	if pattern == "/" {
		pattern = ""
	}
	prefix := r.prefix + pattern
	r.engine.router["GET"+"-"+prefix] = struct{}{}
	r.engine.tire.AddHandler(prefix, srv)
}
