package entt

import (
	"fmt"
)

// SingleView Single component view specialization.
type SingleView struct {
	Pool *SparseSet2
}

func NewSingleView(pool *SparseSet2) *SingleView {
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
	return v.Pool.Has(entity) && v.Pool.Data()[v.Pool.SparseSet.Get(entity)] == entity
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

// --------------- Multiple component view ----------------

// View multiple component view
type View struct {
	Pools []*SparseSet2

	coms     map[ComponentID]uint8
	indexSeq []int
}

func NewView(pools ...*SparseSet2) *View {
	v := &View{}
	v.coms = make(map[ComponentID]uint8)
	v.indexSeq = make([]int, 0, len(pools)-1)
	for i, pool := range pools {
		v.coms[pool.com] = uint8(i)

		if i < cap(v.indexSeq) {
			v.indexSeq = append(v.indexSeq, i)
		}
	}
	v.Pools = pools

	return v
}

func (v *View) indexOfCom(com ComponentID) int {
	return int(v.coms[com])
}

func (v *View) getSinglePool(com ComponentID) *SparseSet2 {
	return v.Pools[v.indexOfCom(com)]
}

func (v *View) unchecked(view *SparseSet) []*SparseSet {
	other := make([]*SparseSet, 0, len(v.Pools))

	for _, pool := range v.Pools {
		if pool.SparseSet != view {
			other = append(other, pool.SparseSet)
		}
	}

	return other
}

// candicate 找到组件集 中，最少组件的那个实体列表和 组件集合
func (v *View) candicate() (*SparseSet, *SparseSet2) {
	var minPool *SparseSet2 = v.Pools[0]
	for _, pool := range v.Pools[1:] {
		if pool.Size() < minPool.Size() {
			minPool = pool
		}
	}

	return minPool.SparseSet, minPool
}

// Estimates the number of entities iterated by the view.
func (v *View) size() int {
	var minPool *SparseSet = v.Pools[0].SparseSet
	for _, pool := range v.Pools[1:] {
		if pool.Size() < minPool.Size() {
			minPool = pool.SparseSet
		}
	}
	return minPool.Size()
}

func (v *View) empty() bool {
	var maxPool *SparseSet = v.Pools[0].SparseSet
	for _, pool := range v.Pools {
		if pool.Size() > maxPool.Size() {
			maxPool = pool.SparseSet
		}
	}
	return maxPool.Size() <= 0
}

// contains check if a view contains an entity
func (v *View) contains(entity EntityID) bool {
	sz := int(entity & entity_mask)
	extent := v.minExtent()

	if sz < extent {
		for _, pool := range v.Pools {
			if pool.Has(entity) && pool.Data()[pool.SparseSet.Get(entity)] == entity {
			} else {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (v *View) minExtent() int {
	var minPool *SparseSet = v.Pools[0].SparseSet
	for _, pool := range v.Pools {
		if pool.Extent() < minPool.Extent() {
			minPool = pool.SparseSet
		}
	}
	return minPool.Extent()
}

// Get 获取绑定在实体身上的指定组件的数据
func (v *View) Get(entity EntityID, com ComponentID) interface{} {
	if !v.contains(entity) {
		panic("view should have entity, but not")
	}
	return v.getSinglePool(com).Get(entity)
}

// GetMulti 获取绑定多个组件的实体身上的组件数据
func (v *View) GetMulti(entity EntityID, coms ...ComponentID) []interface{} {
	if !v.contains(entity) {
		panic("view should have entity, but not")
	}
	ret := make([]interface{}, 0, len(coms))
	for _, com := range coms {
		ret = append(ret, v.getSinglePool(com).Get(entity))
	}
	return ret
}

func (v *View) each(cpool *SparseSet2, fn func(entity EntityID, comDatas map[ComponentID]interface{})) {

	other := v.unchecked(cpool.SparseSet)
	minExtend := v.minExtent()

	datas := make([]*_EntityIDIterator, 0, len(v.indexSeq))
	for _, idx := range v.indexSeq {
		datas = append(datas, other[idx].Begin())
	}
	raw := make([]*ComponentIterator, 0, len(v.Pools))
	for _, pool := range v.Pools {
		raw = append(raw, pool.Begin())
	}

	end := cpool.SparseSet.End()
	begin := cpool.SparseSet.Iterator()

	for !begin.Equal(IIterator(end)) {
		ordered := true
		for _, data := range datas {
			if data.Data() != begin.Data() {
				ordered = false
				break
			}
		}

		if ordered {
			comDatas := make(map[ComponentID]interface{}, len(v.Pools))
			entity := begin.Data().(EntityID)
			for idx, rawData := range raw {
				comDatas[v.Pools[idx].com] = rawData.Data()
			}

			fn(entity, comDatas)

			begin.Next()
		} else {
			break
		}
	}

	for !begin.Equal(IIterator(end)) {

		entity := begin.Data().(EntityID)
		// it := cpool.Begin()
		sz := int(entity & entity_mask) // entity index

		if sz < minExtend && allOf(other, entity) {
			comDatas := make(map[ComponentID]interface{}, len(v.Pools))
			for _, pool := range v.Pools {
				if pool.Fast(entity) {
					comDatas[pool.com] = pool.Get(entity)
				}
			}

			fn(entity, comDatas)
		}

		begin.Next()
	}
}

func allOf(pools []*SparseSet, entity EntityID) bool {
	for _, pool := range pools {
		if pool.Fast(entity) == false {
			return false
		}
	}
	return true
}

func (v *View) Each(fn func(entity EntityID, datas map[ComponentID]interface{})) {
	_, pool := v.candicate()
	v.each(pool, fn)
}
