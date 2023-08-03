package cbuffer

import "testing"

type IntComparable int

func (i IntComparable) Less(other IntComparable) bool {
	return i < other
}

func (i IntComparable) Equal(other IntComparable) bool {
	return i == other
}

func TestOrderedIntCB(t *testing.T) {
	cb := NewOrderedCircuitBuffer[IntComparable](3)

	t.Logf("CircuitBuffer cap: %v", cb.Cap())
	t.Logf("CircuitBuffer len: %v", cb.Len())

	cb.Add(1)
	cb.Add(2)
	cb.Add(3)
	cb.Add(4)

	t.Logf("CircuitBuffer cap: %v", cb.Cap())
	t.Logf("CircuitBuffer len: %v", cb.Len())

	zeroItem := IntComparable(2)

	if cb.GetItem(0) != zeroItem {
		t.Fatalf("%v != %v", cb.GetItem(0), zeroItem)
	} else {
		t.Logf("item at index 0 == %v", cb.GetItem(0))
	}
}

func TestOCBIteration(t *testing.T) {
	cb := NewOrderedCircuitBuffer[IntComparable](3)

	cb.Add(1)
	cb.Add(2)
	cb.Add(3)
	cb.Add(4)

	index := 0
	check_list := []IntComparable{2, 3, 4}

	// Check that no iteration object related to this OCB
	if cb.iter != nil {
		t.Fatalf("OCB %v must not contain an iter object", cb)
	}

	// Check iteration function
	for item := range cb.Iter() {
		t.Logf("Element %v: %v", index, *item)

		// Check that no iteration object created and related
		if cb.iter == nil && *item != 4 {
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

func TestCBBreakIteration(t *testing.T) {
	cb := NewCircuitBuffer[int](3)

	for i := range []int{1, 3, 5, 4} {
		cb.Add(i)
	}

	index := 0

	// Check iteration function
	for item := range cb.Iter() {
		t.Logf("Element %v: %v", index, *item)

		// Check that no iteration object created and related
		if *item != 3 {
			cb.Break()
			break
		}

		index++
	}

	// Check that iteration object removed correctly
	if cb.iter != nil {
		t.Fatalf("OCB %v must not contain an iter object after Break of iteration", cb)
	}
}

func TestOCBWrongBreakIteration(t *testing.T) {
	ocb := NewOrderedCircuitBuffer[IntComparable](3)

	for i := range []int{1, 3, 5, 4} {
		err := ocb.Add(IntComparable(i))

		if err != nil && i != 4 {
			t.Fatalf("%v can't be inserted into ocb %v", i, ocb)
		}
	}

}

func TestOCBSearch(t *testing.T) {
	ocb := NewOrderedCircuitBuffer[IntComparable](100)

	for i := 0; i < 100; i++ {
		err := ocb.Add(IntComparable(i))

		if err != nil {
			t.Fatalf("%v can't be inserted into ocb %s", i, ocb)
		}
	}

	expected := 55
	index, found := ocb.Search(IntComparable(expected))
	if !found {
		t.Fatalf("%v can't be found in ocb %s", expected, ocb)
	}

	if index != expected {
		t.Fatalf("%v have wrong index %v (expected 55) founded in ocb %s", expected, index, ocb)
	}

}

func TestOCBSearchNotFound(t *testing.T) {
	ocb := NewOrderedCircuitBuffer[IntComparable](50)

	for i := 0; i < 100; i++ {
		err := ocb.Add(IntComparable(i))

		if err != nil {
			t.Fatalf("%v can't be inserted into ocb %s", i, ocb)
		}
	}

	expected := IntComparable(20)
	index, found := ocb.Search(expected)
	if found {
		t.Fatalf("%v couldn't be found in ocb %s", expected, ocb)
	}

	if index != -1 {
		t.Fatalf("Expected index -1, got %v, while searching %v in %s", index, expected, ocb)
	}

}
