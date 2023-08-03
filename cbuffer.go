package cbuffer

import (
	"fmt"
	"sync"
)

type CircularBuffer[T any] struct {
	buf        []T
	len        int
	startIndex int
	iter       *CBIterator[T]
	mu         *sync.RWMutex
}

type CBIterator[T any] struct {
	ocb      *CircularBuffer[T]
	iterChan chan *T
	index    int
}

func NewCircularBuffer[T any](size int) *CircularBuffer[T] {
	var mu sync.RWMutex
	return &CircularBuffer[T]{
		make([]T, size), size, 0, nil, &mu,
	}
}

func (ocb *CircularBuffer[T]) String() string {
	return fmt.Sprintf("CircularBuffer(%v)", ocb.buf)
}

func (ocb *CircularBuffer[T]) iterateCB() *CBIterator[T] {
	iterator := &CBIterator[T]{
		ocb, make(chan *T), 0,
	}

	if ocb.iter != nil {
		err := fmt.Errorf("circuit buffer already iterated or break iteration incorrect %v", ocb)
		panic(err)
	}

	ocb.iter = iterator

	return iterator
}

func (iter *CBIterator[T]) Len() int {
	return iter.ocb.Len()
}

func (iter *CBIterator[T]) Next() *T {
	nextItem := iter.ocb.GetItem(iter.index)
	iter.index++

	return &nextItem
}

func (ocb *CircularBuffer[T]) GetItem(index int) T {
	ocb.mu.RLock()
	defer ocb.mu.RUnlock()

	realIndex := (ocb.startIndex + index) % ocb.Cap()
	return ocb.buf[realIndex]
}

func (ocb *CircularBuffer[T]) Add(item T) bool {
	ocb.mu.Lock()
	defer ocb.mu.Unlock()

	removeFlag := false
	addIndex := ocb.len

	if ocb.len == ocb.Cap() {
		addIndex = ocb.startIndex
		ocb.startIndex++
		ocb.startIndex %= ocb.len
		removeFlag = true
	} else {
		ocb.len++
	}

	ocb.buf[addIndex] = item

	return removeFlag
}

func (ocb *CircularBuffer[T]) Len() int {
	ocb.mu.RLock()
	defer ocb.mu.RUnlock()

	return ocb.len
}

func (ocb *CircularBuffer[T]) Cap() int {

	return cap(ocb.buf)
}

func (ocb *CircularBuffer[T]) Iter() chan *T {

	iterator := ocb.iterateCB()

	go func(iterator *CBIterator[T]) {
		ocb.mu.RLock()
		defer ocb.mu.RUnlock()

		for iterator.index < iterator.Len() {
			iterator.iterChan <- iterator.Next()
		}
		close(iterator.iterChan)
		iterator.ocb.iter = nil
	}(iterator)

	return iterator.iterChan
}

func (ocb *CircularBuffer[T]) Break() {
	// This method need to kill dead goroutine if we break Iter loop earlier.

	iterator := ocb.iter

	if iterator != nil {
		iterator.index = iterator.Len()
		<-iterator.iterChan

		ocb.iter = nil
	}
}

// Ordered
type Comparable[T any] interface {
	Less(T) bool
	Equal(T) bool
}

type OrderedCircularBuffer[T Comparable[T]] struct {
	CircularBuffer[T]
}

func NewOrderedCircularBuffer[T Comparable[T]](size int) *OrderedCircularBuffer[T] {
	return &OrderedCircularBuffer[T]{
		CircularBuffer[T]{make([]T, size), size, 0, nil, new(sync.RWMutex)},
	}
}

func (ocb *OrderedCircularBuffer[T]) String() string {
	return fmt.Sprintf("OrderedCircularBuffer(%v)", ocb.buf)
}

func (ocb *OrderedCircularBuffer[T]) Add(item T) error {

	addIndex := ocb.len

	if ocb.len != 0 {
		lastItem := ocb.GetItem(ocb.len - 1)

		if item.Less(lastItem) {
			err := fmt.Errorf("%v can't be added in Ordered Circuit Buffer. Last element in cb %v", item, lastItem)
			return err
		}
	}

	ocb.mu.Lock()
	if ocb.len == ocb.Cap() {
		addIndex = ocb.startIndex
		ocb.startIndex++
		ocb.startIndex %= ocb.len
	} else {
		ocb.len++
	}

	ocb.buf[addIndex] = item
	ocb.mu.Unlock()

	return nil
}

func (ocb *OrderedCircularBuffer[T]) Search(value T) (index int, found bool) {
	ocb.mu.RLock()
	defer ocb.mu.RUnlock()

	start := 0
	end := ocb.len - 1

	for start < end {
		index = start + (end-start)/2
		item := ocb.GetItem(index)

		if item.Equal(value) {
			found = true
			return
		} else if item.Less(value) {
			start = index + 1
		} else {
			end = index
		}

	}

	return -1, found
}
