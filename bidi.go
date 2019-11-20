package main

import (
	"sync"
)

// BidiIntToStrDict is used to store mapping of integer to string and vice versa.
type BidiIntToStrDict struct {
	int2str map[int]string
	str2int map[string]int

	sync.Mutex
}

// NewBidiIntToStrDict creates new BidiIntToStrDict object.
func NewBidiIntToStrDict() *BidiIntToStrDict {
	obj := new(BidiIntToStrDict)
	obj.int2str = make(map[int]string)
	obj.str2int = make(map[string]int)

	return obj
}

// Bind is used to define specified integer to string mapping and vice versa.
func (obj *BidiIntToStrDict) Bind(i int, s string) {
	obj.Lock()
	obj.int2str[i] = s
	obj.str2int[s] = i
	obj.Unlock()
}

// UnBind is used to undefine specified keys.
func (obj *BidiIntToStrDict) UnBind(keys ...interface{}) {
	obj.Lock()

	for i := range keys {
		switch v := keys[i].(type) {
		case int:
			if k, ok := obj.GetString(v); ok {
				delete(obj.str2int, k)
			}

			delete(obj.int2str, v)
		case string:
			if k, ok := obj.GetInteger(v); ok {
				delete(obj.int2str, k)
			}

			delete(obj.str2int, v)
		}
	}

	obj.Unlock()
}

// GetInteger returns integer value based on specified string key with bool that specifies successful lookup.
func (obj *BidiIntToStrDict) GetInteger(key string) (int, bool) {
	val, ok := obj.str2int[key]

	return val, ok
}

// GetString returns string value based on specified integer key with bool that specifies successful lookup.
func (obj *BidiIntToStrDict) GetString(key int) (string, bool) {
	val, ok := obj.int2str[key]

	return val, ok
}
