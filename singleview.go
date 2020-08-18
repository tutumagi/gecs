package gecs

import "fmt"

// SingleView Single component view specialization.
type SingleView struct {
	Pool *Storage
}

// NewSingleView new single view
func newSingleView(pool *Storage) *SingleView {
	return &SingleView{
		Pool: pool,
	}
}

// Size of the entities that has the given component
func (v *SingleView) Size() int {
	return v.Pool.Size()
}

// Empty check if the entities is empty
func (v *SingleView) Empty() bool {
	return v.Pool.Empty()
}

// Raw return the array of components
func (v *SingleView) Raw() []interface{} {
	return v.Pool.Raw()
}

// Data direct access to the list of entities
func (v *SingleView) Data() []EntityID {
	return v.Pool.Data()
}

// Begin an iterator to the first entity that has the given component
func (v *SingleView) Begin() *_EntityIDIterator {
	return v.Pool.SparseSet.Begin()
}

// End an iterator to the entity following the last entity that has the given component
func (v *SingleView) End() *_EntityIDIterator {
	return v.Pool.SparseSet.End()
}

// Contains check if a view contains en entity
func (v *SingleView) Contains(entity EntityID) bool {
	return v.Pool.Has(entity)
}

// Get the component data assigned the entity
func (v *SingleView) Get(entity EntityID) interface{} {
	if !v.Contains(entity) {
		panic(fmt.Sprintf("entity(%v)should have the component", entity))
	}
	return v.Pool.Get(entity)
}

// Each Iterates entities and components and applies the given function object to them
// func (v *SingleView) Each(fn func(entity EntityID, data interface{})) {
func (v *SingleView) Each(fn func(entity EntityID, datas map[ComponentID]interface{})) {
	comIter := v.Pool.Iterator()
	entityIter := v.Pool.SparseSet.Iterator()
	Each(entityIter, func(data interface{}) {
		// fn(data.(EntityID), comIter.Data())
		fn(data.(EntityID), map[ComponentID]interface{}{v.Pool.com: comIter.Data()})
		comIter.Next()
	})
}
