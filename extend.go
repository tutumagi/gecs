package gecs

// extendEntitySlice extend slice cap to newCap
func extendEntitySlice(slice []EntityID, newCap int) []EntityID {
	if cap(slice) >= newCap {
		return slice
	}

	newSlice := make([]EntityID, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	return slice
}

// extendEntitySliceWithValue extend slice cap to newCap
func extendEntitySliceWithValue(slice []EntityID, newCap int, value EntityID) []EntityID {
	if cap(slice) >= newCap {
		return slice
	}
	oldLen := len(slice)
	newSlice := make([]EntityID, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	for i := oldLen; i < cap(slice); i++ {
		slice = append(slice, value)
	}

	return slice
}

// extendInterfaceSlice extend slice cap to newCap
func extendInterfaceSlice(slice []interface{}, newCap int) []interface{} {
	if cap(slice) >= newCap {
		return slice
	}

	newSlice := make([]interface{}, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	return slice
}

// extendInterfaceSliceWithValue extend slice cap to newCap
func extendInterfaceSliceWithValue(slice []interface{}, newCap int, value interface{}) []interface{} {
	if cap(slice) >= newCap {
		return slice
	}
	oldLen := len(slice)
	newSlice := make([]interface{}, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	for i := oldLen; i < cap(slice); i++ {
		slice = append(slice, value)
	}

	return slice
}

// extendSparseSetSlice extend slice cap to newCap
func extendSparseSetSlice(slice []*SparseSet, newCap int) []*SparseSet {
	if cap(slice) >= newCap {
		return slice
	}

	newSlice := make([]*SparseSet, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	return slice
}

// extendSparseSetSliceWithValue extend slice cap to newCap
func extendSparseSetSliceWithValue(slice []*SparseSet, newCap int, value *SparseSet) []*SparseSet {
	if cap(slice) >= newCap {
		return slice
	}
	oldLen := len(slice)
	newSlice := make([]*SparseSet, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	for i := oldLen; i < cap(slice); i++ {
		slice = append(slice, value)
	}

	return slice
}

// extendPoolHandlerWithValue extend slice cap to newCap
func extendPoolHandlerWithValue(slice []*_PoolHandler, newCap int, constructor func() *_PoolHandler) []*_PoolHandler {
	if cap(slice) >= newCap {
		return slice
	}
	oldLen := len(slice)
	newSlice := make([]*_PoolHandler, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	for i := oldLen; i < cap(slice); i++ {
		slice = append(slice, constructor())
	}

	return slice
}

// extend2DEntityArrayWithValue extend slice cap to newCap
func extend2DEntityArrayWithValue(slice [][]EntityID, newCap int, value []EntityID) [][]EntityID {
	if cap(slice) >= newCap {
		return slice
	}

	oldLen := len(slice)

	newSlice := make([][]EntityID, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	for i := oldLen; i < cap(slice); i++ {
		newValue := make([]EntityID, len(value), cap(value))
		copy(newValue, value)
		slice = append(slice, newValue)
	}

	return slice
}
