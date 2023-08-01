package cbuffer

// import (
// 	"testing"
// )

// // TestIter create new CB
// func TestIter(t *testing.T) {

// 	cb := NewIntCircularBuffer(5)

// 	cb.Add(1)
// 	cb.Add(2)
// 	cb.Add(3)

// 	for v := range cb.Iter() {
// 		t.Log(*v)
// 	}
// }

// // TestGetItem create new CB and check TestGetItem function
// func TestGetItem(t *testing.T) {
// 	cb := NewIntCircularBuffer(5)

// 	cb.Add(1)
// 	cb.Add(2)
// 	cb.Add(3)

// 	for index, i := range []int{1, 2, 3} {
// 		if i != cb.GetItem(index) {
// 			t.Fatalf("Expected %v by index %v, got %v", i, index, cb.GetItem(index))
// 		}
// 	}

// }

// // TestLen create new CB and check Len function
// func TestLen(t *testing.T) {
// 	cb := NewIntCircularBuffer(5)

// 	_len := cb.Len()
// 	if _len != 0 {
// 		t.Fatal("Let was", _len, "expected 0")
// 	}

// 	cb.Add(1)
// 	cb.Add(2)
// 	cb.Add(3)

// 	_len = cb.Len()
// 	if _len != 3 {
// 		t.Fatal("Let was", _len, "expected 3")
// 	}
// }

// // TestCap create new CB and check Cap function
// func TestCap(t *testing.T) {
// 	cb := NewIntCircularBuffer(5)

// 	_cap := cb.Cap()
// 	if _cap != 5 {
// 		t.Fatal("Let was", _cap, "expected 5")
// 	}
// }

// func TestAddMoreThanSize(t *testing.T) {
// 	cb := NewIntCircularBuffer(3)

// 	cb.Add(1)
// 	cb.Add(2)
// 	cb.Add(3)
// 	cb.Add(4)

// 	if cb.Len() != cb.Cap() {
// 		t.Fatalf("Len %v != cap %v, but expected", cb.Len(), cb.Cap())
// 	}

// 	for index, v := range []int{2, 3, 4} {
// 		if v != cb.GetItem(index) {
// 			t.Fatalf("Expected %v by index %v, got %v", v, index, cb.GetItem(index))
// 		}
// 	}

// 	cb.Add(5)
// 	cb.Add(6)
// 	cb.Add(67)

// 	for index, v := range []int{5, 6, 67} {
// 		if v != cb.GetItem(index) {
// 			t.Fatalf("Expected %v by index %v, got %v", v, index, cb.GetItem(index))
// 		}
// 	}
// }

// func TestBinarySearchInCB(t *testing.T) {
// 	cb := NewIntCircularBuffer(50)

// 	for _, v := range []int{333, 725, 698, 372, 288, 962, 850, 483, 539, 199, 844, 623, 952, 276, 412, 320, 597, 394, 403, 728, 96, 800, 306, 930, 418, 168, 501, 435, 457, 140, 112, 974, 41, 630, 861, 441, 307, 826, 2, 613, 196, 827, 373, 93, 451, 980, 1, 483, 346, 214} {
// 		cb.Add(v)
// 	}

// }
