package entt

// extendEntitySlice 扩大slice 的cap
func extendEntitySlice(slice []EntityID, newCap int) []EntityID {
	if cap(slice) >= newCap {
		return slice
	}

	newSlice := make([]EntityID, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	return slice
}

// extendEntitySliceWithValue 扩大slice 的cap
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

// extendInterfaceSlice 扩大slice 的cap
func extendInterfaceSlice(slice []interface{}, newCap int) []interface{} {
	if cap(slice) >= newCap {
		return slice
	}

	newSlice := make([]interface{}, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	return slice
}

// extendInterfaceSliceWithValue 扩大slice 的cap
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

// extendSparseSetSlice 扩大slice 的cap
func extendSparseSetSlice(slice []*SparseSet, newCap int) []*SparseSet {
	if cap(slice) >= newCap {
		return slice
	}

	newSlice := make([]*SparseSet, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	return slice
}

// extendSparseSetSliceWithValue 扩大slice 的cap
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

// extendPoolWithValue 扩大slice 的cap
func extendPoolWithValue(slice []*_Pool, newCap int, value *_Pool) []*_Pool {
	if cap(slice) >= newCap {
		return slice
	}
	oldLen := len(slice)
	newSlice := make([]*_Pool, len(slice), newCap)
	copy(newSlice, slice)
	slice = newSlice

	for i := oldLen; i < cap(slice); i++ {
		slice = append(slice, value)
	}

	return slice
}
