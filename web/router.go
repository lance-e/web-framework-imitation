package web

import (
	"fmt"
	"strings"
)

// 用来支持路由树的操作
// 代表路由树（森林）
type router struct {
	//一个方法对应一棵树
	trees map[string]*node
}
type node struct {
	route string
	path  string
	//静态节点
	//子 path 到子节点的映射
	children map[string]*node
	//通配符匹配
	starChild *node
	//路径参数
	paramChild *node

	//每个节点可能会有业务逻辑，需加入handler
	handler HandleFunc
}
type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func NewRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 用来注册路由,提供handleFunc
// addRoute 为私有方法，可以避免用户传入一个错误的http方法
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	//进行一些限制：
	if path == "" {
		panic("web 路径不能为空字符串")
	}
	if path[0] != '/' {
		panic("根目录不为 /")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("路径结尾不能为/")
	}
	//首先找到树
	root, ok := r.trees[method]
	if !ok {
		//说明没有根节点,没有则创建
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	//根节点特殊处理
	if path == "/" {
		if root.handler != nil {
			panic("重复注册根节点/")
		}
		root.handler = handleFunc
		root.route = "/"
		return
	}

	//切割path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {

		if seg == "" {
			panic("不能有连续的 /")
		}
		//递归下去，找准位置
		//如果中途有节点不存在，你就要创建出来
		child := root.childrenOrCreat(seg)
		//root.children[seg] = children
		root = child
	}
	if root.handler != nil {
		panic(fmt.Sprintf("路由冲突，重复注册，path[%s]", path))
	}
	root.handler = handleFunc
	root.route = path
}

// 查看一个子节点,没有则创建一个传入children中
func (n *node) childrenOrCreat(seg string) *node {
	if seg[0] == ':' {
		if n.starChild != nil {
			panic("web 不允许同时注册路径参数和通配符匹配,已经存在通配符匹配")
		}
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}
	//通配符匹配
	if seg == "*" {
		if n.paramChild != nil {
			panic("web 不允许同时注册路径参数和通配符，已经存在路径参数")
		}
		n.starChild = &node{
			path: seg,
		}
		return n.starChild
	}

	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		// 要新建一个
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	//对根节点特殊处理一下
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}

	//这里把前置和后置的/ 去掉
	path = strings.Trim(path, "/")
	//按照斜杠切割
	segs := strings.Split(path, "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		//命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			//因为path是 :id 这个格式
			pathParams[child.path[1:]] = seg
		}
		root = child

	}
	//返回一个true，就代表了我确实有这个节点，但是这个节点有没有业务逻辑（handler）就不一定了
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
	//return root,root.handler!=nil
}

// childOf 优先静态匹配，匹配不上，再考虑通配符匹配
// 第一个返回值是子节点
// 第二个返回值标记是否是路径参数
// 第三个返回值标记路由是否命中
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, ok
}
