package merkledag

import (
	"errors"
	"hash"
)

// 两个栈：A用于遍历树，B用于存储悬而未决的DIR节点
// 判断节点类型，若为FILE节点，则将其hash值写入KVStore，放回其hash
// 若为DIR节点，将此节点放入B，遍历其子节点
// 若其子节点为FILE则直接存入KVStore，若为DIR则放入A，B
// 弹出A最上层元素，用同样的方法处理，直到A为空
// 弹出B最上层元素，计算其哈希，直到B为空
// 当B为空时，返回根节点的哈希值
func add(store KVStore, node Node, h hash.Hash) []byte {
	A := New()
	B := New()
	if node.Type() == FILE {
		store.Put(h.Sum(node.Bytes()), node.Bytes())
		return h.Sum(node.Bytes())
	} else {
		traverse(node, store, h, A, B)
		for !A.Empty() {
			temp, _ := A.Pop()
			traverse(temp.(Node), store, h, A, B)
		}
		for !B.Empty() {
			temp, _ := B.Pop()
			it = temp.(Node).It()
			result := make([]byte, 0)
			for it.Next() {
				aahash := h.Sum(it.Node().Bytes())
				result = append(result, aahash...)
			}
			store.Put(h.Sum(result), temp.(Node).Bytes())
			if B.Empty() {
				return h.Sum(result)

			}
		}

	}
	return nil
}
func traverse(node Node, store KVStore, h hash.Hash, A *ArrayStack, B *ArrayStack) {
	it = node.It()
	for it.Next() {
		if it.Node().Type() == FILE {
			store.Put(h.Sum(it.Node().Bytes()), it.Node().Bytes())
		} else {
			A.Push(it.Node())
			B.Push(it.Node())
		}
	}
}

type ArrayStack struct {
	elements []interface{}
}

// New creates a new array stack.
func New() *ArrayStack {
	return &ArrayStack{}
}

// Size returns the number of elements in the stack.
func (s *ArrayStack) Size() int {
	return len(s.elements)
}

// Empty returns true or false whether the stack has zero elements or not.
func (s *ArrayStack) Empty() bool {
	return len(s.elements) == 0
}

// Clear clears the stack.
func (s *ArrayStack) Clear() {
	s.elements = make([]interface{}, 0, 10)
}

// Push adds an element to the stack.
func (s *ArrayStack) Push(e interface{}) {
	s.elements = append(s.elements, e)
}

// Pop fetches the top element of the stack and removes it.
func (s *ArrayStack) Pop() (interface{}, error) {
	if s.Empty() {
		return nil, errors.New("Pop: the stack cannot be empty")
	}
	result := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return result, nil
}

// Top returns the top of element from the stack, but does not remove it.
func (s *ArrayStack) Top() (interface{}, error) {
	if s.Empty() {
		return nil, errors.New("Top: stack cannot be empty")
	}
	return s.elements[len(s.elements)-1], nil
}
