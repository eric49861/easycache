package lru

import "testing"

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestCache_Add(t *testing.T) {

}

func TestCache_Len(t *testing.T) {

}

func TestCache_RemoveOldest(t *testing.T) {

}

func TestOnEvicted(t *testing.T) {

}
