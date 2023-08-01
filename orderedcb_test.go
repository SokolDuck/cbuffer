package cbuffer

import "testing"

type D struct {
	value int
}

func (d D) Less(other D) bool {
	return d.value < other.value
}

func (d D) Equal(other D) bool {
	return d == other
}

func TestOrderedIntCB(t *testing.T) {
	cb := NewOrderedCB[D](3)

	t.Logf("CircuitBuffer cap: %v", cb.Cap())
	t.Logf("CircuitBuffer len: %v", cb.Len())

	cb.Add(D{1})
	cb.Add(D{2})
	cb.Add(D{3})
	cb.Add(D{4})

	zeroItem := D{2}

	if cb.GetItem(0) != zeroItem {
		t.Fatalf("%v != %v", cb.GetItem(0), zeroItem)
	} else {
		t.Logf("item at index 0 == %v", cb.GetItem(0))
	}
}

func TestOCBIteration(t *testing.T) {
	cb := NewOrderedCB[D](3)

	cb.Add(D{1})
	cb.Add(D{2})
	cb.Add(D{3})
	cb.Add(D{4})

	index := 0
	check_list := []D{{2}, {3}, {4}}

	// Check that no iteration object related to this OCB
	if cb.iter != nil {
		t.Fatalf("OCB %v must not contain an iter object", cb)
	}

	// Check iteration function
	for item := range cb.Iter() {
		t.Logf("Element %v: %v", index, *item)

		// Check that no iteration object created and related
		if cb.iter == nil {
			t.Fatalf("OCB %v must contain an iter object, but got nil", cb)
		}

		if !item.Equal(check_list[index]) {
			t.Fatalf("Index %v: %v != %v (%v)", index, *item, check_list[index], cb)
		}

		index++
	}

	// Check that iteration object removed correctly
	if cb.iter != nil {
		t.Fatalf("OCB %v must not contain an iter object after end of iteration", cb)
	}

}

// func TestOCBBreakIteration(t *testing.T) {
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
// }
