package entt

import (
	"gonut/engine/algo"
)

// SparseSet2 单个类型组件绑定实体的稀疏数组
type SparseSet2 struct {
	com ComponentID
	// comID ComponentID
	*SparseSet

	/*
		instances 中存储的组件的数据的排序 跟 实体的排序是完全一致的
		也就是说是跟 SparseSet.direct 中的排序是完全一致的
	*/
	instances []interface{}
}

// NewSparseSet2 构造组件-实体的稀疏数组
func NewSparseSet2(com ComponentID) *SparseSet2 {
	return &SparseSet2{
		com:       com,
		SparseSet: NewSparseSet(),
		instances: make([]interface{}, 0),
	}
}

func (s *SparseSet2) Com() ComponentID {
	return s.com
}

// Construct 创建实体，并绑定组件数据data
func (s *SparseSet2) Construct(entity EntityID, data interface{}) interface{} {
	s.SparseSet.Construct(entity)
	s.instances = append(s.instances, data)
	return data
}

// Destroy remove该实体绑定的组件数据，并且从 实体稀疏数组中remove掉该实体
func (s *SparseSet2) Destroy(entity EntityID) {

	back := s.instances[len(s.instances)-1]
	s.instances[s.SparseSet.Get(entity)] = back
	s.instances = s.instances[:len(s.instances)-1]
	s.SparseSet.Destroy(entity)
}

// Get 获取该实体绑定的组件数据
func (s *SparseSet2) Get(entity EntityID) interface{} {
	return s.instances[s.SparseSet.Get(entity)]
}

// Replace 替换数据
func (s *SparseSet2) Replace(entity EntityID, data interface{}) interface{} {
	s.instances[s.SparseSet.Get(entity)] = data
	return data
}

// Reset 组件稀疏数组，以及实体稀疏数组
func (s *SparseSet2) Reset() {
	s.SparseSet.Reset()
	s.instances = make([]interface{}, 0)
}

// Raw 组件原始数据
func (s *SparseSet2) Raw() []interface{} {
	return s.instances
}

// Reserve 给组件和实体的稀疏数组扩容
func (s *SparseSet2) Reserve(cap int) {
	s.SparseSet.Reserve(cap)
	s.instances = extendInterfaceSlice(s.instances, cap)
}

// Begin 组件的迭代器
func (s *SparseSet2) Begin() *ComponentIterator {
	return &ComponentIterator{
		datas: s.instances,
		pos:   len(s.instances),
	}
}

// End 组件的迭代器
func (s *SparseSet2) End() *ComponentIterator {
	return &ComponentIterator{
		datas: s.instances,
		pos:   0,
	}
}

// Iterator 组件的迭代器
func (s *SparseSet2) Iterator() *ComponentIterator {
	return s.Begin()
}

// ComponentIterator for interface{} type
type ComponentIterator struct {
	datas []interface{}
	pos   int
}

// Equal other
func (i ComponentIterator) Equal(other algo.IIterator) bool {
	tOther := other.(*ComponentIterator)
	return i.pos == tOther.pos && &i.datas[0] == &tOther.datas[0]
}

// Next iterator
func (i *ComponentIterator) Next() algo.IIterator {
	i.pos--
	return algo.IIterator(i)
}

// Prev iterator
func (i *ComponentIterator) Prev() algo.IIterator {
	i.pos++
	return algo.IIterator(i)
}

// Begin iterator
func (i ComponentIterator) Begin() algo.IIterator {
	i.pos = len(i.datas)
	return algo.IIterator(&i)
}

// End iterator
func (i ComponentIterator) End() algo.IIterator {
	i.pos = 0
	return algo.IIterator(&i)
}

// Data with the iterator value
func (i ComponentIterator) Data() interface{} {
	return i.datas[i.pos-1]
}
