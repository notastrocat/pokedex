package pokecache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cases := []struct {
		name string
		updateInterval int
	} {
		{
			name: "valid interval 1",
			updateInterval: 5,
		},
		{
			name: "valid interval 2",
			updateInterval: 10,
		},
		{
			name: "valid interval 3",
			updateInterval: 15,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := NewCache(time.Duration(c.updateInterval) * time.Second)
			if cache == nil {
				t.Errorf("ERR: Test failed!!\ncache is nil")
				return
			}
		})
	}
}

func TestAdd_Get(t *testing.T) {
	cases := []struct {
		name string
		key string
		val []byte
	} {
		{
			name: "add entry 1",
			key: "testkey1",
			val: []byte("testvalue1"),
		},
		{
			name: "add entry 2",
			key: "testkey2",
			val: []byte("testvalue2"),
		},
		{
			name: "add entry 3",
			key: "testkey3",
			val: []byte("testvalue3"),
		},
	}

	cache := NewCache(5 * time.Second)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache.Add(c.key, c.val)
			retrievedVal, exists := cache.Get(c.key)
			if !exists {
				t.Errorf("ERR: Test failed!!\nkey %s does not exist", c.key)
				return
			}
			if string(retrievedVal) != string(c.val) {
				t.Errorf("ERR: Test failed!!\nexpected %s, got %s", string(c.val), string(retrievedVal))
				return
			}
		})
	}

	time.Sleep(6 * time.Second)
	for _, c := range cases {
		t.Run("check expiration "+c.name, func(t *testing.T) {
			_, exists := cache.Get(c.key)
			if exists {
				t.Errorf("ERR: Test failed!!\nkey %s should have expired", c.key)
				return
			}
		})
	}
}
