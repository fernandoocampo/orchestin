package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelAdd(t *testing.T) {
	givenKey := "one"
	givenValue := 1
	givenDB := NewChannelDB()
	done := make(chan interface{})
	defer close(done)

	go givenDB.Start(done)

	givenDB.Add(givenKey, givenValue)
	result := givenDB.Find(done, givenKey)

	got := <-result

	assert.Equal(t, true, got.Ok)
	assert.Equal(t, givenValue, got.Value)
}

func TestChannelAddMultipleProducers(t *testing.T) {
	keyOne, valueOne := "one", 1
	keyTwo, valueTwo := "two", 2

	giveChannelDB := NewChannelDB()
	done := make(chan interface{})
	defer close(done)
	go giveChannelDB.Start(done)

	var wg sync.WaitGroup

	producerOne := func() {
		defer wg.Done()
		giveChannelDB.Add(keyOne, valueOne)
	}
	producerTwo := func() {
		defer wg.Done()
		giveChannelDB.Add(keyTwo, valueTwo)
	}

	wg.Add(2)
	go producerOne()
	go producerTwo()
	wg.Wait()
	resultOne := giveChannelDB.Find(done, keyOne)
	resultTwo := giveChannelDB.Find(done, keyTwo)

	gotOne := <-resultOne
	gotTwo := <-resultTwo

	assert.True(t, gotOne.Ok)
	assert.True(t, gotTwo.Ok)
	assert.Equal(t, valueOne, gotOne.Value)
	assert.Equal(t, valueTwo, gotTwo.Value)

}
