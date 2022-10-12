package gee

import "strings"

type node struct {
	pattern  string  //待匹配路由
	part     string  //匹配部分
	children []*node //存前缀树
	isWild   bool    //是否精准匹配
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 子匹配，如果是匹配到了part，就把当前这个child（其实是node的一部分）追加到nodes里
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]       //这里需要错误处理吗？不需要，因为当len(parts) == height时就返回了
	child := n.matchChild(part) //查找父节点
	if child == nil {           //如果父节点为空，则把当前的part写入
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //如果匹配到了最后一层，且pattern 不为空时返回
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
