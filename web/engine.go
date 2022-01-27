package web

import (
	"net/http"
	"strings"
)

type Engine struct {
	router     map[string]struct{}
	tire       *Tire
	routeGroup map[string]*RouteGroup
	ctx        map[string]*Context
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := request.Method + "-" + request.URL.Path
	ctx := e.context(key, writer, request)
	defer func() {
		delete(e.ctx, key)
	}()

	if _, ok := e.router[key]; ok {
		// 如果存在根路由，先执行根路由
		if rootGrp, exist := e.routeGroup["/"]; exist {
			for _, srv := range rootGrp.middleware {
				srv(ctx)
			}
		}
		for prefix, group := range e.routeGroup {
			if prefix != "/" && strings.HasPrefix(request.URL.Path, prefix) {
				for _, srv := range group.middleware {
					srv(ctx)
				}
			}
		}
		e.tire.GetHandler(request.URL.Path)(ctx)
	}
}

func NewEngine() *Engine {
	return &Engine{
		router:     make(map[string]struct{}, 0),
		tire:       NewTire(),
		routeGroup: make(map[string]*RouteGroup, 0),
		ctx:        make(map[string]*Context),
	}
}

func (e *Engine) Group(prefix string) *RouteGroup {
	group := NewRouteGroup(prefix, e)
	e.routeGroup[prefix] = group
	return group
}

type Service func(ctx *Context) error

// Get register router
func (e *Engine) Get(pattern string, srv Service) {
	e.router["GET"+"-"+pattern] = struct{}{}
	e.tire.AddHandler(pattern, srv)
}

func (e *Engine) context(pattern string, writer http.ResponseWriter, request *http.Request) *Context {
	context, ok := e.ctx[pattern]
	if ok {
		return context
	}
	e.ctx[pattern] = &Context{
		w: writer,
		r: request,
	}
	return e.ctx[pattern]
}

func (e *Engine) Run() error {
	return http.ListenAndServe(":8080", e)
}
