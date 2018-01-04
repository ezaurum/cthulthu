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

	var wg sync.WaitGroup
	wg.Add(3)

	n := func(c chan string) {
		c <- kg.Generate()
	}

	c := make(chan string)

	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			go n(c)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			go n(c)
		}
	}()

	go func() {
		defer wg.Done()
		for j := 0; j < 100000; j++ {
			s0 := <-c
			s1 := <-c
			assert.NotEqual(t, s1, s0)
		}
	}()

	wg.Wait()
}
