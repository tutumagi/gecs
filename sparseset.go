package entt

import (
	"fmt"
	"gonut/engine/algo"
)

const _EnttPageSize int = 2 << 14

// EntityID for uint32 entityid

type _VersionType uint16
type _DifferenceType int64

const entity_mask = 0xFFFFF // 20bits for entity number
const version_mask = 0xFFF  // 12bits for the version(reset in (0-4095))
const entity_shift = 20     // 20bits for entity number

type SparseSet struct {
	direct  []EntityID
	reverse []EntityID
}

func NewSparseSet() *SparseSet {
	return &SparseSet{
		direct:  make([]EntityID, 0),
		reverse: make([]EntityID, 0),
	}
}

func (s *SparseSet) String() string {
	return fmt.Sprintf("<SparseSet>(direct:%+v, reserve:%+v)", s.direct, s.reverse)
}

// Reserve 预设容量
func (s *SparseSet) Reserve(cap int) {
	s.direct = extendEntitySlice(s.direct, cap)
}

func (s *SparseSet) Construct(entity EntityID) {
	pos := int(entity & entity_mask)
	if !(pos < len(s.reverse)) {
		s.reverse = extendEntitySliceWithValue(s.reverse, pos+1, EntityID(DefaultPlaceholder))
	}
	s.reverse[pos] = EntityID(len(s.direct))
	s.direct = append(s.direct, entity)
}

// Has 是否有这个entity
func (s *SparseSet) Has(entity EntityID) bool {
	pos := int(entity & entity_mask)
	return pos < len(s.reverse) && s.reverse[pos] != DefaultPlaceholder
}

// Fast 快速检查是否还有此entity，不会做pos的越界检查
func (s *SparseSet) Fast(entity EntityID) bool {
	pos := int(entity & entity_mask)
	return s.reverse[pos] != DefaultPlaceholder
}

// Get 返回entity 在此稀疏数组中的索引
func (s *SparseSet) Get(entity EntityID) int {
	if !s.Has(entity) {
		panic(fmt.Sprintf("should have the entity %v", entity))
	}
	return int(s.reverse[int(entity&entity_mask)])
}

// Destroy remove entity from sparseset
//	确保每次移除后，该spareset 仍然是紧凑的
func (s *SparseSet) Destroy(entity EntityID) {
	if !s.Has(entity) {
		panic(fmt.Sprintf("should have the entity %v", entity))
	}
	// 最后一个元素
	back := s.direct[len(s.direct)-1]

	entityPos := entity & entity_mask
	backEntityPos := back & entity_mask

	// 将删掉的元素跟最后一个元素进行交换，保证 direct 永远是紧凑的
	candicate := s.reverse[entityPos]
	s.reverse[backEntityPos] = candicate
	s.direct[int(candicate)] = back

	s.reverse[entityPos] = DefaultPlaceholder

	s.direct = s.direct[:len(s.direct)-1] // remove direct last element
}

// Reset the sparse set
func (s *SparseSet) Reset() {
	s.direct = s.direct[0:0]
	s.reverse = s.reverse[0:0]
}

// Cap 返回cap
func (s *SparseSet) Cap() int {
	return cap(s.direct)
}

// Size 返回数量
func (s *SparseSet) Size() int {
	return len(s.direct)
}

// Empty 是否是空的稀疏数组
func (s *SparseSet) Empty() bool {
	return len(s.direct) == 0
}

// Extent 返回此 sparse set的范围
func (s *SparseSet) Extent() int {
	return len(s.reverse)
}

// Data 返回 packed array
func (s *SparseSet) Data() []EntityID {
	return s.direct
}

func (s *SparseSet) Begin() *_EntityIDIterator {
	return &_EntityIDIterator{
		datas: s.direct,
		pos:   len(s.direct),
	}
}

func (s *SparseSet) End() *_EntityIDIterator {
	return &_EntityIDIterator{
		datas: s.direct,
		pos:   0,
	}
}

func (s *SparseSet) Iterator() *_EntityIDIterator {
	// return &_EntityIDIterator{
	// 	datas: s.direct,
	// 	pos:   0,
	// }
	return s.Begin()
}

// ------------------ Iterator -----------------

// _EntityIDIterator for EntityID type
type _EntityIDIterator struct {
	datas []EntityID
	pos   int
}

// Equal other
func (i _EntityIDIterator) Equal(other algo.IIterator) bool {
	tOther := other.(*_EntityIDIterator)
	if i.pos != tOther.pos {
		return false
	}
	if len(i.datas) != len(tOther.datas) {
		return false
	}
	if len(i.datas) == len(tOther.datas) && (len(i.datas) == 0 || &i.datas[0] == &tOther.datas[0]) {
		return true
	}
	return false
}

// Next iterator
func (i *_EntityIDIterator) Next() algo.IIterator {
	i.pos--
	return algo.IIterator(i)
}

// Prev iterator
func (i *_EntityIDIterator) Prev() algo.IIterator {
	i.pos++
	return algo.IIterator(i)
}

// Begin iterator
func (i _EntityIDIterator) Begin() algo.IIterator {
	i.pos = len(i.datas)
	return algo.IIterator(&i)
}

// End iterator
func (i _EntityIDIterator) End() algo.IIterator {
	i.pos = 0
	return algo.IIterator(&i)
}

// Data with the iterator value
func (i _EntityIDIterator) Data() interface{} {
	return i.datas[i.pos-1]
}
