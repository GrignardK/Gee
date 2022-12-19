package gee

import "strings"

// 动态路由依赖的前缀树

type node struct {
	pattern  string  // 全路径
	part     string  // 本节点路径
	children []*node // 子节点
	isWild   bool    // 是否模糊匹配
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	ret := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			ret = append(ret, child)
		}
	}
	return ret
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		if n.part == parts[len(parts)-1] {
			panic(n.pattern + "exist same routing :" + n.part)
		}
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}
