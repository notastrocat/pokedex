package pokecache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cases := []struct {
		name           string
		updateInterval int
	}{
		{
			name:           "valid interval 1",
			updateInterval: 5,
		},
		{
			name:           "valid interval 2",
			updateInterval: 10,
		},
		{
			name:           "valid interval 3",
			updateInterval: 15,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := NewCache(time.Duration(c.updateInterval) * time.Second)
			defer cache.Close()
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
		key  string
		val  []byte
	}{
		{
			name: "add entry 1",
			key:  "testkey1",
			val:  []byte("testvalue1"),
		},
		{
			name: "add entry 2",
			key:  "testkey2",
			val:  []byte("testvalue2"),
		},
		{
			name: "add entry 3",
			key:  "testkey3",
			val:  []byte("testvalue3"),
		},
	}

	cache := NewCache(5 * time.Second)
	defer cache.Close()

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

func TestReapLoopTimeInterval(t *testing.T) {
	// Test that reapLoop respects the time interval
	// Entries should NOT expire before the interval
	interval := 2 * time.Second
	cache := NewCache(interval)
	defer cache.Close()

	cache.Add("test_key", []byte("test_value"))

	// Check immediately - should exist
	_, exists := cache.Get("test_key")
	if !exists {
		t.Errorf("ERR: Entry should exist immediately after adding")
		return
	}

	// Check halfway through interval - should still exist
	time.Sleep(time.Duration(float64(interval) * 0.6))
	_, exists = cache.Get("test_key")
	if !exists {
		t.Errorf("ERR: Entry should still exist before interval expires")
		return
	}

	// Wait for expiration to definitely occur
	// Add buffer to account for reapLoop ticker granularity
	time.Sleep(time.Duration(float64(interval)*0.6) + 500*time.Millisecond)
	_, exists = cache.Get("test_key")
	if exists {
		t.Errorf("ERR: Entry should have expired after interval")
		return
	}
}

func TestReapLoopMultipleEntries(t *testing.T) {
	// Test that reapLoop only removes expired entries, not all entries
	interval := 3 * time.Second
	cache := NewCache(interval)
	defer cache.Close()

	// Add first entry
	cache.Add("key1", []byte("value1"))
	time.Sleep(1 * time.Second)

	// Add second entry after a delay
	cache.Add("key2", []byte("value2"))

	// Wait for first entry to expire but not second
	time.Sleep(2200 * time.Millisecond)

	_, key1Exists := cache.Get("key1")
	_, key2Exists := cache.Get("key2")

	if key1Exists {
		t.Errorf("ERR: key1 should have expired (waited 3.2s with 3s interval)")
	}
	if !key2Exists {
		t.Errorf("ERR: key2 should still exist (only 2.2s elapsed)")
	}

	// Wait longer to ensure second entry definitely expires
	// The reapLoop ticker runs at intervals, so we need to account for up to 2x the interval
	time.Sleep(3000 * time.Millisecond)

	_, key2Exists = cache.Get("key2")
	if key2Exists {
		t.Errorf("ERR: key2 should have expired (waited 5.2s with 3s interval)")
	}
}
