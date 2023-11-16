package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	//第一个步骤：构造路由树
	//第二个步骤；验证路由树
	testRouter := []struct { //匿名结构体来存放数据
		Method string
		Path   string
	}{
		{
			Method: http.MethodGet,
			Path:   "/",
		},
		{
			Method: http.MethodGet,
			Path:   "/user/home",
		},
		{
			Method: http.MethodGet,
			Path:   "/user",
		},

		{
			//通配符匹配测试用例
			Method: http.MethodGet,
			Path:   "/user/*",
		},
		{
			Method: http.MethodGet,
			Path:   "/article/list",
		},
		{
			Method: http.MethodPost,
			Path:   "/article/create",
		},
		{
			Method: http.MethodPost,
			Path:   "/login",
		},
		{
			Method: http.MethodGet,
			Path:   "/*/:id",
		},
	} //测试用例

	r := NewRouter()
	//构造路由树
	var mockHandleFunc HandleFunc = func(ctx *Context) {}
	for _, route := range testRouter {
		r.addRoute(route.Method, route.Path, mockHandleFunc)
	}
	//断言路由树与你预期的路由树一样
	wantRoute := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandleFunc,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandleFunc,
						children: map[string]*node{
							"home": &node{
								path:     "home",
								children: map[string]*node{},
								handler:  mockHandleFunc,
							},
						},

						starChild: &node{
							path:    "*",
							handler: mockHandleFunc,
						},
					},
					"article": &node{
						path: "article",
						children: map[string]*node{
							"list": &node{
								handler:  mockHandleFunc,
								path:     "list",
								children: map[string]*node{},
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"article": &node{
						path: "article",
						children: map[string]*node{
							"create": &node{
								path:     "create",
								handler:  mockHandleFunc,
								children: map[string]*node{},
							},
						},
					},
					"login": &node{
						path:     "login",
						handler:  mockHandleFunc,
						children: map[string]*node{},
					},
				},
			},
		},
	}
	message, ok := wantRoute.equal(&r)
	assert.True(t, ok, message)
	r = NewRouter()
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "", mockHandleFunc)
	})
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "login", mockHandleFunc)
	})
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/a///b", mockHandleFunc)
	})
	r = NewRouter()
	r.addRoute(http.MethodGet, "/", mockHandleFunc)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandleFunc)
	}, "重复注册/")

	r = NewRouter()
	r.addRoute(http.MethodGet, "/a/b", mockHandleFunc)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b", mockHandleFunc)
	}, "重复注册/a/b")

	a := NewRouter()
	a.addRoute(http.MethodGet, "/a/*", mockHandleFunc)
	assert.Panicsf(t, func() {
		a.addRoute(http.MethodGet, "/a/:id", mockHandleFunc)
	}, "web 不允许同时注册路径参数和通配符匹配,已经存在通配符匹配")

	a = NewRouter()
	a.addRoute(http.MethodGet, "/a/:id", mockHandleFunc)
	assert.Panicsf(t, func() {
		a.addRoute(http.MethodGet, "/a/*", mockHandleFunc)
	}, "web 不允许同时注册路径参数和通配符匹配,已经存在路径参数")

}

func (r *router) equal(y *router) (string, bool) {
	for i, v := range r.trees {
		dst, ok := y.trees[i]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method"), false
		}
		str, ok := v.equal(dst)
		if !ok {
			return str, false
		}

	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.path) != len(y.path) {
		return fmt.Sprintf("子节点数量不正确"), false
	}

	if n.starChild != nil {
		msg, ok := n.starChild.equal(y.starChild)
		if !ok {
			return msg, ok
		}

	}
	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
	}
	//比较handler
	nhandler := reflect.ValueOf(n.handler)
	yhandler := reflect.ValueOf(y.handler)
	if nhandler != yhandler {
		return fmt.Sprintf("handler 不相等"), false
	}

	for path, node := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点不存在"), false
		}
		str, ok := node.equal(dst)
		if !ok {
			return str, false
		}

	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	//先构造路由树
	testRouter := []struct { //匿名结构体来存放数据
		Method string
		Path   string
	}{
		{
			Method: http.MethodGet,
			Path:   "/",
		},
		{
			Method: http.MethodGet,
			Path:   "/user/home",
		},
		{
			Method: http.MethodGet,
			Path:   "/user",
		},
		{
			//通配符匹配测试用例
			Method: http.MethodGet,
			Path:   "/user/*",
		},
		{
			Method: http.MethodGet,
			Path:   "/article/list",
		},
		{
			Method: http.MethodPost,
			Path:   "/article/create",
		},

		{
			Method: http.MethodPost,
			Path:   "/login",
		},
		{
			Method: http.MethodPost,
			Path:   "/login/:username",
		},
	}
	r := NewRouter()
	var mockHandleFunc HandleFunc = func(ctx *Context) {}
	for _, route := range testRouter {
		r.addRoute(route.Method, route.Path, mockHandleFunc)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		info      *matchInfo
	}{
		{
			//完全命中
			name:      "user home",
			method:    http.MethodGet,
			path:      "/user/home",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandleFunc,
					path:    "home",
				},
			},
		},
		{
			//通配符测试
			name:      "user star",
			method:    http.MethodGet,
			path:      "/user/abc",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandleFunc,
					path:    "*",
				},
			},
		},
		{
			//对根节点进行查找
			name:      "root",
			method:    http.MethodGet,
			path:      "/",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandleFunc,
					path:    "/",
					children: map[string]*node{
						"user": &node{
							path:    "user",
							handler: mockHandleFunc,
							children: map[string]*node{
								"home": &node{
									path:     "home",
									handler:  mockHandleFunc,
									children: map[string]*node{},
								},
							},
						},
					},
				},
			},
		},
		{
			//不完全命中，找到了该节点，但是没有业务逻辑（handler）
			name:      "article",
			method:    http.MethodGet,
			path:      "/article",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path: "article",
				},
			},
		},
		{
			name:   "method not found",
			method: http.MethodPut,
			path:   "/user/home",
		},
		{
			//username
			name:      "login username",
			method:    http.MethodPost,
			path:      "/login/longxu",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    ":username",
					handler: mockHandleFunc,
				},
				pathParams: map[string]string{
					"username": "longxu",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matchinfo, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.info.pathParams, matchinfo.pathParams)
			message, ok := tc.info.n.equal(matchinfo.n)
			assert.True(t, ok, message)

		})
	}

}
