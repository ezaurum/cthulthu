package snowflake

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGetNew(t *testing.T) {
	kg := New(0)

	s := kg.Generate()

	assert.NotNil(t, s)
	assert.NotEmpty(t, s, "ID cannot be empty")
}

func TestGetNewSerial(t *testing.T) {
	kg := New(0)
	keyCount := 5000
	waitGroups := 5000

	keys := make(map[string]bool)
	var wg sync.WaitGroup
	wg.Add(waitGroups + 1)

	n := func(c chan string) {
		c <- kg.Generate()
	}

	c := make(chan string, keyCount)
	go func() {
		defer wg.Done()
		for i:=0; i< waitGroups * keyCount ; i++ {
			key := <-c
			keys[key] = true
		}
	}()

	//when
	for i := 0; i < waitGroups; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < keyCount; i++ {
				n(c)
			}
		}()
	}

	wg.Wait()

	//then
	assert.Equal(t, waitGroups*keyCount, len(keys))
}
