package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutexAdd(t *testing.T) {
	givenKey := "one"
	givenValue := 1
	givenDB := NewMutexDB()

	got := givenDB.Add(givenKey, givenValue)

	assert.NoError(t, got)
	storedValue, ok := givenDB.Find(givenKey)
	assert.Equal(t, true, ok)
	assert.Equal(t, givenValue, storedValue)
}

func TestMutexAddMultipleProducers(t *testing.T) {
	keyOne, valueOne := "one", 1
	keyTwo, valueTwo := "two", 2

	var wg sync.WaitGroup

	giveMutexDB := NewMutexDB()
	producerOne := func() {
		defer wg.Done()
		got := giveMutexDB.Add(keyOne, valueOne)
		assert.NoError(t, got)
	}
	producerTwo := func() {
		defer wg.Done()
		got := giveMutexDB.Add(keyTwo, valueTwo)
		assert.NoError(t, got)
	}

	wg.Add(2)
	go producerOne()
	go producerTwo()
	wg.Wait()

	gotOne, okOne := giveMutexDB.Find(keyOne)
	gotTwo, okTwo := giveMutexDB.Find(keyTwo)

	assert.True(t, okOne)
	assert.True(t, okTwo)
	assert.Equal(t, valueOne, gotOne)
	assert.Equal(t, valueTwo, gotTwo)

}
