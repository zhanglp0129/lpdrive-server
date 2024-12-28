package router

import (
	"github.com/gin-gonic/gin"
)

// Router 自定义路由器
type Router struct {
	*gin.RouterGroup
}

func (r *Router) handle(method, path string, handles ...HandleFunc) {
	hs := make([]gin.HandlerFunc, 0, len(handles))
	for _, h := range handles {
		hs = append(hs, h.toGinHandleFunc())
	}
	r.Handle(method, path, hs...)
}

func (r *Router) Get(path string, handles ...HandleFunc) {
	r.handle("GET", path, handles...)
}

func (r *Router) Post(path string, handles ...HandleFunc) {
	r.handle("POST", path, handles...)
}

func (r *Router) Put(path string, handles ...HandleFunc) {
	r.handle("PUT", path, handles...)
}

func (r *Router) Patch(path string, handles ...HandleFunc) {
	r.handle("PATCH", path, handles...)
}

func (r *Router) Delete(path string, handles ...HandleFunc) {
	r.handle("DELETE", path, handles...)
}
