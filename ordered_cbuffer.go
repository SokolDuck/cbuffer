package cbuffer

import (
	"fmt"
)

type Comparable[T any] interface {
	Less(T) bool
	Equal(T) bool
}

type CircuitBuffer[T Comparable[T]] struct {
	buf        []T
	len        int
	startIndex int
	iter       *CBIterator[T]
}

type CBIterator[T Comparable[T]] struct {
	ocb      *CircuitBuffer[T]
	iterChan chan *T
	index    int
}

func (ocb *CircuitBuffer[T]) iterateCB() *CBIterator[T] {
	iterator := &CBIterator[T]{
		ocb, make(chan *T), 0,
	}

	if ocb.iter != nil {
		err := fmt.Errorf("Circuit buffer already iterated or break iteration incorrect %v", ocb)
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

func NewOrderedCB[T Comparable[T]](size int) *CircuitBuffer[T] {
	return &CircuitBuffer[T]{
		make([]T, size), size, 0, nil,
	}
}

func (ocb *CircuitBuffer[T]) GetItem(index int) T {
	realIndex := (ocb.startIndex + index) % ocb.Cap()
	return ocb.buf[realIndex]
}

func (ocb *CircuitBuffer[T]) Add(item T) bool {
	removeFlag := false
	addIndex := ocb.Len()

	if ocb.Len() == ocb.Cap() {
		addIndex = ocb.startIndex
		ocb.startIndex++
		ocb.startIndex %= ocb.Len()
		removeFlag = true
	} else {
		ocb.len++
	}

	ocb.buf[addIndex] = item

	return removeFlag
}

func (ocb *CircuitBuffer[T]) Len() int {
	return ocb.len
}

func (ocb *CircuitBuffer[T]) Cap() int {
	return cap(ocb.buf)
}

func (ocb *CircuitBuffer[T]) Iter() chan *T {

	iterator := ocb.iterateCB()

	go func(iterator *CBIterator[T]) {
		for iterator.index < iterator.Len() {
			iterator.iterChan <- iterator.Next()
			iterator.index++
		}
		close(iterator.iterChan)
		iterator.ocb.iter = nil
	}(iterator)

	return iterator.iterChan
}

func (ocb *CircuitBuffer[T]) Break() {
	// This method need to kill dead goroutine if we break Iter loop earlier.

	iterator := ocb.iter

	if iterator != nil {
		iterator.index = iterator.Len()
		<-iterator.iterChan

		ocb.iter = nil
	}
}

func (ocb *CircuitBuffer[T]) Search(value T) (index int, found bool) {
	start := 0
	end := ocb.Len() - 1

	for start < end {
		index = start + (end-start)/2
		item := ocb.GetItem(index)

		if item.Equal(value) {
			found = true
			return
		} else if item.Less(value) {
			end = index - 1
		} else {
			start = index + 1
		}

	}

	return -1, found
}
