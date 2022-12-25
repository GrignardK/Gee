package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	// 存放每种请求方法的路径根节点
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 只实现了词首单*的适配
func (r *router) parsePattern(pattern string) []string {
	splitS := strings.Split(pattern, "/")
	// 分割的字符串需要考虑几种情况：1.* 2.//test 3./*/t
	parts := make([]string, 0)
	for _, s := range splitS {
		if s == "" {
			continue
		}
		parts = append(parts, s)
		if s[0] == '*' { // 如果是词首通配符*也需要添加，后面直接省略
			break
		}
	}
	return parts
}

// 添加格式化路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := r.parsePattern(pattern)
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 对于传入的地址进行模式化的匹配，如/user/114/name -> /user/:userId/name, /user/114/name -> /user/*FILEPATH
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	searchParts := r.parsePattern(pattern)
	n := root.search(searchParts, 0)
	if n == nil {
		return nil, nil
	}
	parts := r.parsePattern(n.pattern)
	params := make(map[string]string)
	// 将传入串和已有串库对比
	for idx, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[idx]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[idx:], "/")
			break
		}
	}
	return n, params
}

func (r *router) Handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 not found: %s\n", c.Path)
		})
	}
	c.Next()
}
