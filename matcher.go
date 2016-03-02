package pmatcher

import (
	"strings"
)

type Matcher struct {
	roots map[int]*Node
}

func New() *Matcher {
	return &Matcher{
		roots: make(map[int]*Node),
	}
}

func (matcher *Matcher) Add(pattern string) {
	parts := strings.Split(pattern, "/")

	root, ok := matcher.roots[len(parts)]
	if !ok {
		root = NewNode()
		matcher.roots[len(parts)] = root
	}
	cur := root

	variables := make(map[int]string)
	for i, part := range parts {
		if len(part) > 0 && part[0] == ':' {
			variables[i] = part[1:]
		}
		cur = cur.AddNext(part)
	}

	cur.SetEnd(
		&Pattern{
			full:      pattern,
			variables: variables,
		},
	)
}

func (matcher *Matcher) Match(path string) (matched bool, pattern string, params map[string]string) {
	parts := strings.Split(path, "/")

	root, ok := matcher.roots[len(parts)]
	if !ok {
		return false, "", nil
	}

	match, pat := matcher.match(root, parts)
	if match {
		params := make(map[string]string)
		for pos, name := range pat.variables {
			params[name] = parts[pos]
		}
		return match, pat.full, params
	}
	return match, "", nil
}

func (matcher *Matcher) match(node *Node, parts []string) (bool, *Pattern) {
	if len(parts) == 0 {
		end := node.End()
		if end == nil {
			return false, nil
		}
		return true, end
	}

	next, wildcard := node.Next(parts[0])

	var (
		path *Pattern
		ok   bool
	)

	if next != nil {
		ok, path = matcher.match(next, parts[1:])
	}
	if ok {
		return ok, path
	}

	// less priority
	if wildcard != nil {
		ok, path = matcher.match(wildcard, parts[1:])
	}
	if ok {
		return ok, path
	}

	return false, nil
}

type Node struct {
	nexts    map[string]*Node
	wildcard *Node
	end      *Pattern
}

func NewNode() *Node {
	return &Node{
		nexts: make(map[string]*Node),
	}
}

func (node *Node) AddNext(step string) *Node {
	if len(step) > 0 && step[0] == ':' {
		if node.wildcard != nil {
			return node.wildcard
		}
		node.wildcard = NewNode()
		return node.wildcard
	}

	next, ok := node.nexts[step]
	if ok {
		return next
	}

	next = NewNode()
	node.nexts[step] = next

	return node.nexts[step]
}

func (node *Node) SetEnd(path *Pattern) {
	node.end = path
}

func (node *Node) Next(step string) (next *Node, wildcard *Node) {
	return node.nexts[step], node.wildcard
}

func (node *Node) End() *Pattern {
	return node.end
}

type Pattern struct {
	full      string
	variables map[int]string
}
