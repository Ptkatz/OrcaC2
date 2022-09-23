// Copyright 2012 Google Inc. All rights reserved.
// Author: Ric Szopa (Ryszard) <ryszard.szopa@gmail.com>

// Package skiplist implements skip list based maps and sets.
//
// Skip lists are a data structure that can be used in place of
// balanced trees. Skip lists use probabilistic balancing rather than
// strictly enforced balancing and as a result the algorithms for
// insertion and deletion in skip lists are much simpler and
// significantly faster than equivalent algorithms for balanced trees.
//
// Skip lists were first described in Pugh, William (June 1990). "Skip
// lists: a probabilistic alternative to balanced
// trees". Communications of the ACM 33 (6): 668–676
package skiplist

import (
	"math/rand"
)

// p is the fraction of nodes with level i pointers that also have
// level i+1 pointers. p equal to 1/4 is a good value from the point
// of view of speed and space requirements. If variability of running
// times is a concern, 1/2 is a better value for p.
const p = 0.25

const DefaultMaxLevel = 32

// A Node is a container for key-value pairs that are stored in a skip
// list.
type Node struct {
	forward    []*Node
	backward   *Node
	key, value interface{}
}

// Next returns the Next Node in the skip list containing n.
func (n *Node) Next() *Node {
	if len(n.forward) == 0 {
		return nil
	}
	return n.forward[0]
}

// Previous returns the Previous Node in the skip list containing n.
func (n *Node) Previous() *Node {
	return n.backward
}

// hasNext returns true if n has a Next Node.
func (n *Node) hasNext() bool {
	return n.Next() != nil
}

// hasPrevious returns true if n has a Previous Node.
func (n *Node) hasPrevious() bool {
	return n.Previous() != nil
}

func (n *Node) Key() interface{} {
	return n.key
}

func (n *Node) Value() interface{} {
	return n.value
}

// A SkipList is a map-like data structure that maintains an ordered
// collection of key-value pairs. Insertion, lookup, and deletion are
// all O(log n) operations. A SkipList can efficiently store up to
// 2^MaxLevel items.
//
// To iterate over a skip list (where s is a
// *SkipList):
//
//	for i := s.Iterator(); i.Next(); {
//		// do something with i.Key() and i.Value()
//	}
type SkipList struct {
	lessThan func(l, r interface{}) bool
	header   *Node
	footer   *Node
	length   int
	update   []*Node
}

// Len returns the length of s.
func (s *SkipList) Len() int {
	return s.length
}

func (s *SkipList) Front() *Node {
	return s.header.Next()
}

func (s *SkipList) Last() *Node {
	current := s.footer
	if current == nil {
		return nil
	}
	return current
}

// Seek returns a bidirectional iterator starting with the first element whose
// key is greater or equal to key; otherwise, a nil iterator is returned.
func (s *SkipList) Seek(key interface{}) *Node {
	current := s.getPath(s.header, nil, key)
	if current == nil {
		return nil
	}

	return current
}

func (s *SkipList) level() int {
	return len(s.header.forward) - 1
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (s *SkipList) effectiveMaxLevel() int {
	return maxInt(s.level(), DefaultMaxLevel)
}

// Returns a new random level.
func (s SkipList) randomLevel() (n int) {
	for n = 0; n < s.effectiveMaxLevel() && rand.Float64() < p; n++ {
	}
	return
}

// Get returns the value associated with key from s (nil if the key is
// not present in s). The second return value is true when the key is
// present.
func (s *SkipList) Get(key interface{}) (value interface{}, ok bool) {
	candidate := s.getPath(s.header, nil, key)

	if candidate == nil || candidate.key != key {
		return nil, false
	}

	return candidate.value, true
}

// GetGreaterOrEqual finds the Node whose key is greater than or equal
// to min. It returns its value, its actual key, and whether such a
// Node is present in the skip list.
func (s *SkipList) GetGreaterOrEqual(min interface{}) (actualKey, value interface{}, ok bool) {
	candidate := s.getPath(s.header, nil, min)

	if candidate != nil {
		return candidate.key, candidate.value, true
	}
	return nil, nil, false
}

// getPath populates update with nodes that constitute the path to the
// Node that may contain key. The candidate Node will be returned. If
// update is nil, it will be left alone (the candidate Node will still
// be returned). If update is not nil, but it doesn't have enough
// slots for all the nodes in the path, getPath will panic.
func (s *SkipList) getPath(current *Node, update []*Node, key interface{}) *Node {
	depth := len(current.forward) - 1

	for i := depth; i >= 0; i-- {
		for current.forward[i] != nil && s.lessThan(current.forward[i].key, key) {
			current = current.forward[i]
		}
		if update != nil {
			update[i] = current
		}
	}
	return current.Next()
}

// Sets set the value associated with key in s.
func (s *SkipList) Set(key, value interface{}) {
	if key == nil {
		return
	}
	// s.level starts from 0, so we need to allocate one.
	update := s.update[:s.level()+1]
	candidate := s.getPath(s.header, update, key)

	if candidate != nil && candidate.key == key {
		candidate.value = value
		return
	}

	newLevel := s.randomLevel()

	if currentLevel := s.level(); newLevel > currentLevel {
		// there are no pointers for the higher levels in
		// update. Header should be there. Also add higher
		// level links to the header.
		for i := currentLevel + 1; i <= newLevel; i++ {
			update = append(update, s.header)
			s.header.forward = append(s.header.forward, nil)
		}
	}

	newNode := &Node{
		forward: make([]*Node, newLevel+1, s.effectiveMaxLevel()+1),
		key:     key,
		value:   value,
	}

	if previous := update[0]; previous.key != nil {
		newNode.backward = previous
	}

	for i := 0; i <= newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	s.length++

	if newNode.forward[0] != nil {
		if newNode.forward[0].backward != newNode {
			newNode.forward[0].backward = newNode
		}
	}

	if s.footer == nil || s.lessThan(s.footer.key, key) {
		s.footer = newNode
	}
}

// Delete removes the Node with the given key.
//
// It returns the old value and whether the Node was present.
func (s *SkipList) Delete(key interface{}) (value interface{}, ok bool) {
	if key == nil {
		return nil, false
	}
	update := s.update[:s.level()+1]
	candidate := s.getPath(s.header, update, key)

	if candidate == nil || candidate.key != key {
		return nil, false
	}

	previous := candidate.backward
	if s.footer == candidate {
		s.footer = previous
	}

	next := candidate.Next()
	if next != nil {
		next.backward = previous
	}

	for i := 0; i <= s.level() && update[i].forward[i] == candidate; i++ {
		update[i].forward[i] = candidate.forward[i]
	}

	for s.level() > 0 && s.header.forward[s.level()] == nil {
		s.header.forward = s.header.forward[:s.level()]
	}
	s.length--

	return candidate.value, true
}

// NewCustomMap returns a new SkipList that will use lessThan as the
// comparison function. lessThan should define a linear order on keys
// you intend to use with the SkipList.
func NewCustomMap(lessThan func(l, r interface{}) bool) *SkipList {
	return &SkipList{
		lessThan: lessThan,
		header: &Node{
			forward: []*Node{nil},
		},
		update: make([]*Node, DefaultMaxLevel+1),
	}
}

// NewIntKey returns a SkipList that accepts int keys.
func NewIntMap() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
}

// NewIntKey returns a SkipList that accepts int keys.
func NewInt32Map() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(int32) < r.(int32)
	})
}

// NewStringMap returns a SkipList that accepts string keys.
func NewStringMap() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(string) < r.(string)
	})
}
