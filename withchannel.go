package main

type KeyValue struct {
	Key   string
	Value int
}

type Maybe struct {
	Value int
	Ok    bool
}

type Filter struct {
	Key    string
	Result chan Maybe
}

type ChannelDB struct {
	add        chan KeyValue
	delete     chan string
	find       chan *Filter
	repository map[string]int
}

func NewChannelDB() *ChannelDB {
	newChannelDB := ChannelDB{
		repository: make(map[string]int),
		add:        make(chan KeyValue),
		delete:     make(chan string),
		find:       make(chan *Filter),
	}
	return &newChannelDB
}

func (m *ChannelDB) Add(key string, value int) {
	m.add <- KeyValue{
		Key:   key,
		Value: value,
	}
}

func (m *ChannelDB) Delete(key string) {
	m.delete <- key
}

func (m *ChannelDB) Find(done <-chan interface{}, key string) <-chan Maybe {
	newFilter := Filter{
		Key:    key,
		Result: make(chan Maybe),
	}
	go func() {
		select {
		case <-done:
			return
		case m.find <- &newFilter:
		}
	}()
	return newFilter.Result
}

func (c *ChannelDB) Start(done <-chan interface{}) {
	for {
		select {
		case <-done:
			return
		case newPair := <-c.add:
			c.repository[newPair.Key] = newPair.Value
		case key := <-c.delete:
			delete(c.repository, key)
		case filter := <-c.find:
			value, ok := c.repository[filter.Key]
			result := Maybe{
				Value: value,
				Ok:    ok,
			}
			go func(done <-chan interface{}) {
				select {
				case <-done:
					return
				default:
					defer close(filter.Result)
					filter.Result <- result
				}
			}(done)
		}
	}
}
