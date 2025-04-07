package skiplist

import (
	"bytes"
	"math"
	"math/rand"
	"time"
)

const (
	maxLevel    int     = 18
	probability float64 = 1 / math.E
)

type (
	Node struct {
		next []*Element
	}

	Element struct {
		Node
		key   []byte
		value interface{}
	}

	Skiplist struct {
		Node
		maxLevel       int
		Len            int
		randSource     rand.Source // 随机数种子
		probability    float64     // 概率系数
		probTable      []float64   // 落在每一层的概率表
		prevNodesCache []*Node
	}
)

func NewSkipList() *Skiplist {
	return &Skiplist{
		Node:           Node{next: make([]*Element, maxLevel)},
		prevNodesCache: make([]*Node, maxLevel),
		maxLevel:       maxLevel,
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
		probability:    probability,
		probTable:      probabilityTable(probability, maxLevel),
	}
}

func (e *Element) Key() []byte {
	return e.key
}

func (e *Element) Value() interface{} {
	return e.value
}

func (e *Element) Next() *Element {
	return e.next[0]
}

func (t *Skiplist) Front() *Element {
	return t.next[0]
}

func (t *Skiplist) Get(key []byte) *Element {
	var prev = &t.Node
	var next *Element

	for i := t.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]

		for next != nil && bytes.Compare(key, next.key) > 0 {
			prev = &next.Node
			next = next.next[i]
		}
	}

	if next != nil && bytes.Compare(next.key, key) <= 0 {
		return next
	}
	return nil
}

func (t *Skiplist) Exist(key []byte) bool {
	return t.Get(key) != nil
}

func (t *Skiplist) backNodes(key []byte) []*Node {
	var prev = &t.Node
	var next *Element

	prevs := t.prevNodesCache
	for i := t.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && bytes.Compare(key, next.key) > 0 {
			prev = &next.Node
			next = next.next[i]
		}
		prevs[i] = prev
	}
	return prevs
}

func (t *Skiplist) Remove(key []byte) *Element {
	prev := t.backNodes(key)

	if element := prev[0].next[0]; element != nil && bytes.Compare(element.key, key) <= 0 {
		for k, v := range element.next {
			prev[k].next[k] = v
		}
		t.Len--
		return element
	}

	return nil
}

func (t *Skiplist) Put(key []byte, value interface{}) *Element {
	var element *Element
	prev := t.backNodes(key)

	if element = prev[0].next[0]; element != nil && bytes.Compare(element.key, key) <= 0 {
		element.value = value
		return element
	}

	element = &Element{
		Node: Node{
			next: make([]*Element, t.randomLevel()),
		},
		key:   key,
		value: value,
	}

	for i := range element.next {
		element.next[i] = prev[i].next[i]
		prev[i].next[i] = element
	}

	t.Len++
	return element
}

func probabilityTable(probability float64, maxLevel int) (table []float64) {
	for i := 1; i <= maxLevel; i++ {
		prob := math.Pow(probability, float64(i-1))
		table = append(table, prob)
	}
	return table
}

func (t *Skiplist) randomLevel() (level int) {
	r := float64(t.randSource.Int63()) / (1 << 63)
	level = 1
	for level < t.maxLevel && r < t.probTable[level] {
		level++
	}
	return
}
