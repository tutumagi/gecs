package gecs

import (
	"fmt"
)

// _EnttPageSize must be a power of two
const _EnttPageSize int = 2 << 1

// 4 表示目前 entityID使用的 字节树，目前使用的 uint32,4个字节
const _EnttPerPage int = _EnttPageSize / 4

// EntityID for uint32 entityid

type _VersionType uint16
type _DifferenceType int64

const entity_mask = 0xFFFFF // 20bits for entity number
const version_mask = 0xFFF  // 12bits for the version(reset in (0-4095))
const entity_shift = 20     // 20bits for entity number

type SparseSet struct {
	sparse [][]EntityID
	packed []EntityID
}

func NewSparseSet() *SparseSet {
	return &SparseSet{
		sparse: make([][]EntityID, 0),
		packed: make([]EntityID, 0),
	}
}

func (s *SparseSet) String() string {
	return fmt.Sprintf("<SparseSet>(packed:%+v, sparse:%+v)", s.packed, s.sparse)
}

func (s *SparseSet) page(entity EntityID) int {
	return (int(entity) & entity_mask) / _EnttPerPage
}

func (s *SparseSet) offset(entity EntityID) int {
	return int(entity) & (_EnttPerPage - 1)
}

func (s *SparseSet) assure(pos int) []EntityID {
	if !(pos < len(s.sparse)) {
		newPageSpareset := extendEntitySliceWithValue([]EntityID{}, _EnttPerPage, DefaultPlaceholder)
		s.sparse = extend2DEntityArrayWithValue(s.sparse, pos+1, newPageSpareset)
	}

	return s.sparse[pos]
}

// Reserve increases the capacity of a sparse set
func (s *SparseSet) Reserve(cap int) {
	s.packed = extendEntitySlice(s.packed, cap)
}

// Capacity returns the number of elements that a sparse set has currently allcated spae for
func (s *SparseSet) Capacity() int {
	return cap(s.packed)
}

// Extent returns the extent of a sparse set
func (s *SparseSet) Extent() int {
	return len(s.sparse) * _EnttPerPage
}

// Size returns the number of elements in a sparse set.
func (s *SparseSet) Size() int {
	return len(s.packed)
}

// Empty checks whether a sparse set is empty.
func (s *SparseSet) Empty() bool {
	return len(s.packed) == 0
}

// Data direct access to the internal packed array
func (s *SparseSet) Data() []EntityID {
	return s.packed
}

// Has 是否有这个entity
func (s *SparseSet) Has(entity EntityID) bool {
	curr := s.page(entity)
	return curr < len(s.sparse) && len(s.sparse[curr]) > 0 && s.sparse[curr][s.offset(entity)] != DefaultPlaceholder
}

// Index returns the position of an entity in a sparse set.
func (s *SparseSet) Index(entity EntityID) int {
	// if !s.Has(entity) {
	// 	panic(fmt.Sprintf("should have the entity %v", entity))
	// }
	return int(s.sparse[s.page(entity)][s.offset(entity)])
}

// Emplace assign an entity to a sparse set
func (s *SparseSet) Emplace(entity EntityID) {
	s.assure(s.page(entity))[s.offset(entity)] = EntityID(len(s.packed))
	s.packed = append(s.packed, entity)
}

// Destroy remove entity from sparseset
//	确保每次移除后，该spareset 仍然是紧凑的
func (s *SparseSet) Destroy(entity EntityID) {
	if !s.Has(entity) {
		panic(fmt.Sprintf("should have the entity %v", entity))
	}

	curr := s.page(entity)
	pos := s.offset(entity)

	// 最后push进来的entity 在 sparese 里面的索引
	lastPacked := s.packed[len(s.packed)-1]
	// 将
	s.packed[int(s.sparse[curr][pos])] = lastPacked
	s.sparse[s.page(lastPacked)][s.offset(lastPacked)] = s.sparse[curr][pos]
	s.sparse[curr][pos] = DefaultPlaceholder

	s.packed = s.packed[:len(s.packed)-1]

	// // 最后一个元素
	// back := s.sparse[len(s.sparse)-1]

	// entityPos := entity & entity_mask
	// backEntityPos := back & entity_mask

	// // 将删掉的元素跟最后一个元素进行交换，保证 direct 永远是紧凑的
	// candicate := s.packed[entityPos]
	// s.packed[backEntityPos] = candicate
	// s.sparse[int(candicate)] = back

	// s.packed[entityPos] = DefaultPlaceholder

	// s.sparse = s.sparse[:len(s.sparse)-1] // remove direct last element
}

// swap tow entities in the internal packed array
func (s *SparseSet) swap(lhs EntityID, rhs EntityID) {
	from := s.sparse[s.page(lhs)][s.offset(lhs)]
	to := s.sparse[s.page(rhs)][s.offset(rhs)]

	s.packed[int(from)], s.packed[int(to)] = s.packed[int(to)], s.packed[int(from)]

	s.sparse[s.page(lhs)][s.offset(lhs)], s.sparse[s.page(rhs)][s.offset(rhs)] = s.sparse[s.page(rhs)][s.offset(rhs)], s.sparse[s.page(lhs)][s.offset(lhs)]
}

// Reset the sparse set
func (s *SparseSet) Reset() {
	s.sparse = s.sparse[0:0]
	s.packed = s.packed[0:0]
}

func (s *SparseSet) Begin() *_EntityIDIterator {
	return &_EntityIDIterator{
		datas: s.packed,
		pos:   len(s.packed),
	}
}

func (s *SparseSet) End() *_EntityIDIterator {
	return &_EntityIDIterator{
		datas: s.packed,
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
func (i _EntityIDIterator) Equal(other IIterator) bool {
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
func (i *_EntityIDIterator) Next() IIterator {
	i.pos--
	return IIterator(i)
}

// Prev iterator
func (i *_EntityIDIterator) Prev() IIterator {
	i.pos++
	return IIterator(i)
}

// Begin iterator
func (i _EntityIDIterator) Begin() IIterator {
	i.pos = len(i.datas)
	return IIterator(&i)
}

// End iterator
func (i _EntityIDIterator) End() IIterator {
	i.pos = 0
	return IIterator(&i)
}

// Data with the iterator value
func (i _EntityIDIterator) Data() interface{} {
	return i.datas[i.pos-1]
}
