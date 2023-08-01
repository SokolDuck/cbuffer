package cbuffer

// type CircularBuffer struct {
// 	len int
// 	cap int
// 	buf []int

// 	startIndex int
// 	iterDoneCN chan bool
// }

// type ICircularBuffer interface {
// 	Len() int
// 	Cap() int
// 	GetItem(int) int
// 	Iter() chan *int
// }

// func (cb *CircularBuffer) Len() int {
// 	return cb.len
// }

// func (cb *CircularBuffer) Cap() int {
// 	return cb.cap
// }

// func (cb *CircularBuffer) GetItem(index int) int {
// 	return cb.buf[(cb.startIndex+index)%cb.Len()]
// }

// func (cb *CircularBuffer) Iter() chan *int {

// 	c := make(chan *int)

// 	go func(cb *CircularBuffer) {
// 		for i := 0; i < cb.Len(); i++ {
// 			val := cb.GetItem(i)
// 			c <- &val
// 		}
// 		close(c)
// 	}(cb)

// 	return c
// }

// /*
// size = 5

// [0, 0, 0, 0, 0]

// startIndex = 0
// len = 0
// cap = 5

// Add(1)

// [1, 0, 0, 0, 0]
// len++ // 1
// cap = 5
// startIndex = 0

// ...

// Add(5)
// len ++  == cap// 5
// cap = 5
// startIndex = 0
// [1, 2, 3, 4, 5]

// Add(6)
// len == cap // 5
// startIndex++ // 1

// Array in mem    ->  Iter via CircularBuffer
// [6, 2, 3, 4, 5] -> [2, 3, 4, 5, 6]
// Add(7)
// [6, 7, 3, 4, 5] -> [3, 4, 5, 6, 7]
// Add(8)
// [6, 7, 8, 4, 5] -> [4, 5, 6, 7, 8]
// Add(9)
// [6, 7, 8, 9, 5] -> [5, 6, 7, 8, 9]
// startIndex++ // 4

// Add(10)
// [6, 7, 8, 9, 10] -> [6, 7, 8, 9, 10]
// startIndex++ % len // 0
// */
// func (cb *CircularBuffer) Add(value int) bool {
// 	removeFlag := false
// 	addIndex := cb.Len()

// 	if cb.Len() == cb.Cap() {
// 		addIndex = cb.startIndex
// 		cb.startIndex++
// 		cb.startIndex %= cb.Len()
// 		removeFlag = true
// 	} else {
// 		cb.len++
// 	}
// 	cb.buf[addIndex] = value

// 	return removeFlag
// }

// func NewIntCircularBuffer(size int) *CircularBuffer {
// 	return &CircularBuffer{
// 		0, size, make([]int, size), 0, nil,
// 	}
// }
