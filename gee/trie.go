package gee

import (
	"errors"
	"strings"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
	isEnd    bool
}

func (n *node) insert(pattern string, parts []string, level int) {
	if level == len(parts) {
		n.pattern = pattern
		n.isEnd = true
		return
	}

	part := parts[level]
	child, _ := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
			isEnd:  false,
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, level+1)
}

func (n *node) search(parts []string, level int) (res *node) {
	if len(parts) == level || strings.HasPrefix(n.part, "*") {
		if n.pattern != "" {
			res = n
		}
		return
	}
	part := parts[level]
	children, _ := n.matchChildren(part)
	for _, child := range children {
		next := child.search(parts, level+1)
		if next != nil {
			res = next
			return
		}
	}
	return
}

// @Depreciated: do not need
func (n *node) isExist(parts []string, level int) bool {
	if len(parts) == level {
		return n.isEnd
	}
	part := parts[level]
	children, _ := n.matchChildren(part)
	for _, child := range children {
		if child.isExist(parts, level+1) {
			return true
		}
	}
	return false
}

func (n *node) matchChild(part string) (*node, error) {
	if n == nil {
		return nil, errors.New("tried: current node is nil")
	}
	for _, child := range n.children {
		if child.part == part {
			return child, nil
		}
	}
	return nil, nil
}

func (n *node) matchChildren(part string) (children []*node, err error) {
	if n == nil {
		err = errors.New("tried: current node is nil")
	}
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	return
}
