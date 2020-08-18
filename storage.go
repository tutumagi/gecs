package gecs

// Storage 单个类型组件绑定实体的稀疏数组
type Storage struct {
	com ComponentID
	// comID ComponentID
	*SparseSet

	/*
		instances 中存储的组件的数据的排序 跟 实体的排序是完全一致的
		也就是说是跟 SparseSet.direct 中的排序是完全一致的
	*/
	instances []interface{}
}

// NewStorage 构造组件-实体的稀疏数组
func NewStorage(com ComponentID) *Storage {
	return &Storage{
		com:       com,
		SparseSet: newSparseSet(),
		instances: make([]interface{}, 0),
	}
}

// Com return data type the storage binded
func (s *Storage) Com() ComponentID {
	return s.com
}

// Reserve Increases the capacity of a storage
func (s *Storage) Reserve(cap int) {
	s.SparseSet.Reserve(cap)
	s.instances = extendInterfaceSlice(s.instances, cap)
}

// Raw Direct access to the array of objects
func (s *Storage) Raw() []interface{} {
	return s.instances
}

// Get returns the object associated with en entity
func (s *Storage) Get(entity EntityID) interface{} {
	return s.instances[s.SparseSet.Index(entity)]
}

// TryGet returns the object associated with an entity. maybe nil.
func (s *Storage) TryGet(entity EntityID) interface{} {
	if s.Has(entity) {
		return s.instances[s.SparseSet.Index(entity)]
	}
	return nil
}

// Emplace assigns an entity to a storage and constructs its object.
//	This version accept both types that can be constructed in place directly
//	and types like aggregates that do not work well with a placement new as
//	performed usually under the hood during an _emplace back_.
func (s *Storage) Emplace(entity EntityID, data interface{}) interface{} {
	s.instances = append(s.instances, data)
	s.SparseSet.Emplace(entity)
	return data
}

// Destroy remove该实体绑定的组件数据，并且从 实体稀疏数组中remove掉该实体
func (s *Storage) Destroy(entity EntityID) {

	back := s.instances[len(s.instances)-1]
	s.instances[s.SparseSet.Index(entity)] = back
	s.instances = s.instances[:len(s.instances)-1]
	s.SparseSet.Destroy(entity)
}

// Replace 替换数据
func (s *Storage) Replace(entity EntityID, data interface{}) interface{} {
	s.instances[s.SparseSet.Index(entity)] = data
	return data
}

// Reset 组件稀疏数组，以及实体稀疏数组
func (s *Storage) Reset() {
	s.SparseSet.Reset()
	s.instances = make([]interface{}, 0)
}

// Begin 组件的迭代器
func (s *Storage) Begin() *ComponentIterator {
	return &ComponentIterator{
		datas: s.instances,
		pos:   len(s.instances),
	}
}

// End 组件的迭代器
func (s *Storage) End() *ComponentIterator {
	return &ComponentIterator{
		datas: s.instances,
		pos:   0,
	}
}

// Iterator 组件的迭代器
func (s *Storage) Iterator() *ComponentIterator {
	return s.Begin()
}

// ComponentIterator for interface{} type
type ComponentIterator struct {
	datas []interface{}
	pos   int
}

// Equal other
func (i ComponentIterator) Equal(other IIterator) bool {
	tOther := other.(*ComponentIterator)
	return i.pos == tOther.pos && &i.datas[0] == &tOther.datas[0]
}

// Next iterator
func (i *ComponentIterator) Next() IIterator {
	i.pos--
	return IIterator(i)
}

// Prev iterator
func (i *ComponentIterator) Prev() IIterator {
	i.pos++
	return IIterator(i)
}

// Begin iterator
func (i ComponentIterator) Begin() IIterator {
	i.pos = len(i.datas)
	return IIterator(&i)
}

// End iterator
func (i ComponentIterator) End() IIterator {
	i.pos = 0
	return IIterator(&i)
}

// Data with the iterator value
func (i ComponentIterator) Data() interface{} {
	return i.datas[i.pos-1]
}
