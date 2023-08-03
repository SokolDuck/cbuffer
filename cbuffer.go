package cbuffer

import (
	"fmt"
)

type CircuitBuffer[T any] struct {
	buf        []T
	len        int
	startIndex int
	iter       *CBIterator[T]
}

type CBIterator[T any] struct {
	ocb      *CircuitBuffer[T]
	iterChan chan *T
	index    int
}

func NewCircuitBuffer[T any](size int) *CircuitBuffer[T] {
	return &CircuitBuffer[T]{
		make([]T, size), size, 0, nil,
	}
}

func (ocb *CircuitBuffer[T]) String() string {
	return fmt.Sprintf("CircuitBuffer(%v)", ocb.buf)
}

func (ocb *CircuitBuffer[T]) iterateCB() *CBIterator[T] {
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

// Ordered
type Comparable[T any] interface {
	Less(T) bool
	Equal(T) bool
}

type OrderedCircuitBuffer[T Comparable[T]] struct {
	CircuitBuffer[T]
}

func NewOrderedCircuitBuffer[T Comparable[T]](size int) *OrderedCircuitBuffer[T] {
	return &OrderedCircuitBuffer[T]{
		CircuitBuffer[T]{make([]T, size), size, 0, nil},
	}
}

func (ocb *OrderedCircuitBuffer[T]) String() string {
	return fmt.Sprintf("OrderedCircuitBuffer(%v)", ocb.buf)
}

func (ocb *OrderedCircuitBuffer[T]) Add(item T) error {
	addIndex := ocb.Len()

	if ocb.len != 0 {
		lastItem := ocb.GetItem(ocb.len - 1)

		if item.Less(lastItem) {
			err := fmt.Errorf("%v can't be added in Ordered Circuit Buffer. Last element in cb %v", item, lastItem)
			return err
		}
	}

	if ocb.Len() == ocb.Cap() {
		addIndex = ocb.startIndex
		ocb.startIndex++
		ocb.startIndex %= ocb.Len()
	} else {
		ocb.len++
	}

	ocb.buf[addIndex] = item

	return nil
}

func (ocb *OrderedCircuitBuffer[T]) Search(value T) (index int, found bool) {
	start := 0
	end := ocb.Len() - 1

	for start < end {
		index = start + (end-start)/2
		item := ocb.GetItem(index)

		if item.Equal(value) {
			found = true
			return
		} else if item.Less(value) {
			start = index + 1
		} else {
			end = index - 1
		}

	}

	return -1, found
}
