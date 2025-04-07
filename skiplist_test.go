package skiplist

import "testing"

func TestTestSkiplist(t *testing.T) {
	skiplist := NewSkipList()
	skiplist.Put([]byte{123}, 123)
	e := skiplist.Get([]byte{123})
	if e == nil {
		t.Error("结果为nil")
		return
	}
	if e.value != 123 {
		t.Error("结果不符合预期")
	}
}
