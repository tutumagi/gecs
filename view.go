package gecs

// --------------- Multiple component view ----------------

// View multiple component view
type View struct {
	Pools []*Storage

	// value is the index of componentID in pools
	coms map[ComponentID]uint8

	filter []*Storage
}

func newView(pools ...*Storage) *View {
	v := &View{}
	v.coms = make(map[ComponentID]uint8)
	for i, pool := range pools {
		v.coms[pool.com] = uint8(i)
	}
	v.Pools = pools

	return v
}

// func (v *View) withInclude(include ...*Storage) *View {
// 	v.Pools = include
// 	return v
// }

func (v *View) withExclude(filter ...*Storage) *View {
	v.filter = filter
	return v
}

func (v *View) getSinglePool(com ComponentID) *Storage {
	return v.Pools[int(v.coms[com])]
}

// candicate 找到组件集 中，最少组件的那个实体列表和 组件集合
func (v *View) candicate() (*SparseSet, *Storage) {
	var minPool *Storage = v.Pools[0]
	for _, pool := range v.Pools[1:] {
		if pool.Size() < minPool.Size() {
			minPool = pool
		}
	}

	return minPool.SparseSet, minPool
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

// func (v *View) get(com *Storage, other *Storage, entity EntityID) interface{} {
// 	if com == other {
// 		return nil
// 	} else {
// 		return other.Get(entity)
// 	}
// }

func (v *View) traverse(com ComponentID, fn func(EntityID, map[ComponentID]interface{})) {
	// checkComInTypes := func() bool {
	// 	for _, t := range types {
	// 		if t == com {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }()

	// if checkComInTypes {
	// 	specifyStorage := v.getSinglePool(com)
	// 	it := specifyStorage.Begin()
	// 	for _, entt := range specifyStorage.Data() {
	// 		if v.inAllSpecify(com, entt) && !v.inExclude(entt) {
	// 			// func(entt, )
	// 		}
	// 		it = it.Next().(*ComponentIterator)
	// 	}

	// } else {
	specifyStorage := v.getSinglePool(com)
	for _, entt := range specifyStorage.Data() {
		if v.inAllSpecify(com, entt) && !v.inExclude(entt) {
			comDatas := make(map[ComponentID]interface{}, len(v.Pools))
			for _, p := range v.Pools {
				comDatas[p.com] = p.Get(entt)
			}
			fn(entt, comDatas)
		}
	}
	// }
}

func (v *View) iterator(fn func(), coms ...ComponentID) {
	view, _ := v.candicate()
	last := view.Iterator().End()
	first := view.Iterator().Begin()

	for first != last {
		entt := first.Data().(EntityID)
		if v.inAllInclude(entt) && !v.inExclude(entt) {

		}
	}
}

// inExclude check whether the entity in at least one filter component storage
//	return true if at least one filter contains the entity, false otherwise.
func (v *View) inExclude(entity EntityID) bool {
	for _, ex := range v.filter {
		if ex.Has(entity) {
			return true
		}
	}
	return false
}

// inAllInclude check whether the entity in all include components storage
//	return true if all the include components contains the entity, falst otherwise
func (v *View) inAllInclude(entity EntityID) bool {
	for _, in := range v.Pools {
		if !in.Has(entity) {
			return false
		}
	}
	return true
}

func (v *View) inSpecify(com ComponentID, entity EntityID) bool {
	if _, ok := v.coms[com]; ok {
		return true
	}
	return v.getSinglePool(com).Has(entity)
}

func (v *View) inAllSpecify(com ComponentID, entity EntityID) bool {
	for _, s := range v.Pools {
		in := func() bool {
			if s.com == com {
				return true
			}
			return s.Has(entity)
		}()
		if !in {
			return false
		}
	}
	return true
}

func (v *View) sizeOfCom(com ComponentID) int {
	return v.getSinglePool(com).Size()
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

func (v *View) empty(coms ...ComponentID) bool {
	if len(coms) == 0 {
		var maxPool *SparseSet = v.Pools[0].SparseSet
		for _, pool := range v.Pools[1:] {
			if pool.Size() > maxPool.Size() {
				maxPool = pool.SparseSet
			}
		}
		return maxPool.Size() <= 0
	} else {
		for _, com := range coms {
			if !v.getSinglePool(com).Empty() {
				return false
			}
		}
		return true
	}
}

// contains check if a view contains an entity
func (v *View) contains(entity EntityID) bool {
	return v.inAllInclude(entity) && !v.inExclude(entity)
	// sz := int(entity & entity_mask)
	// extent := v.minExtent()

	// if sz < extent {
	// 	for _, pool := range v.Pools {
	// 		if pool.Has(entity) && pool.Data()[pool.SparseSet.Index(entity)] == entity {
	// 		} else {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// } else {
	// 	return false
	// }
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
func (v *View) GetMulti(entity EntityID, coms ...ComponentID) map[ComponentID]interface{} {
	if !v.contains(entity) {
		panic("view should have entity, but not")
	}
	ret := make(map[ComponentID]interface{}, len(coms))
	for _, com := range coms {
		pool := v.getSinglePool(com)
		ret[pool.com] = pool.Get(entity)
	}
	return ret
}

func (v *View) each(cpool *Storage, fn func(entity EntityID, comDatas map[ComponentID]interface{})) {

	v.traverse(cpool.com, fn)
	// other := v.unchecked(cpool.SparseSet)
	// minExtend := v.minExtent()

	// datas := make([]*_EntityIDIterator, 0, len(v.indexSeq))
	// for _, idx := range v.indexSeq {
	// 	datas = append(datas, other[idx].Begin())
	// }
	// raw := make([]*ComponentIterator, 0, len(v.Pools))
	// for _, pool := range v.Pools {
	// 	raw = append(raw, pool.Begin())
	// }

	// end := cpool.SparseSet.End()
	// begin := cpool.SparseSet.Iterator()

	// for !begin.Equal(IIterator(end)) {
	// 	ordered := true
	// 	for _, data := range datas {
	// 		if data.Data() != begin.Data() {
	// 			ordered = false
	// 			break
	// 		}
	// 	}

	// 	if ordered {
	// 		comDatas := make(map[ComponentID]interface{}, len(v.Pools))
	// 		entity := begin.Data().(EntityID)
	// 		for idx, rawData := range raw {
	// 			comDatas[v.Pools[idx].com] = rawData.Data()
	// 		}

	// 		fn(entity, comDatas)

	// 		begin.Next()
	// 	} else {
	// 		break
	// 	}
	// }

	// for !begin.Equal(IIterator(end)) {

	// 	entity := begin.Data().(EntityID)
	// 	// it := cpool.Begin()
	// 	sz := int(entity & entity_mask) // entity index

	// 	if sz < minExtend && allOf(other, entity) {
	// 		comDatas := make(map[ComponentID]interface{}, len(v.Pools))
	// 		for _, pool := range v.Pools {
	// 			if pool.Has(entity) {
	// 				comDatas[pool.com] = pool.Get(entity)
	// 			}
	// 		}

	// 		fn(entity, comDatas)
	// 	}

	// 	begin.Next()
	// }
}

// Each Iterates entities and components and applies the given function object to them
func (v *View) Each(fn func(entity EntityID, datas map[ComponentID]interface{})) {
	_, pool := v.candicate()
	v.each(pool, fn)
}
