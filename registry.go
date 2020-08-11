package entt

import (
	"fmt"
)

type Registry struct {
	pools    []*_Pool
	handlers []*_Pool
	entities []EntityID

	available int
	next      EntityID

	componentFamily _TypeMap
	handlerFamily   _TypeMap

	checkComponentFamily _NameMap
}

func newRegistry() *Registry {
	return &Registry{
		pools:                make([]*_Pool, 0),
		handlers:             make([]*_Pool, 0),
		entities:             make([]EntityID, 0, 1<<19), // 最多同时52w 左右个实体
		componentFamily:      make(_TypeMap, 0),
		handlerFamily:        make(_TypeMap, 0),
		checkComponentFamily: make(_NameMap, 0),

		available: 0,
		next:      EntityID(0),
	}
}

// Create 新建实体，返回实体id
func (r *Registry) Create() EntityID {
	var entity EntityID
	if r.available > 0 {
		// 有实体销毁时，复用该entityID，更新 Version
		entt := r.next
		// 拿到已经被销毁的entt的 version
		version := r.entities[int(entt)] & (version_mask << entity_shift)
		r.next = r.entities[int(entt)] & entity_mask
		entity = entt | version
		r.entities[int(entt)] = entity
		r.available--
	} else {
		// 在没有任何实体销毁(destroy)时
		entity = EntityID(len(r.entities))
		r.entities = append(r.entities, entity)
	}

	return entity
}

// Assign component data to entity
func (r *Registry) Assign(entity EntityID, component ComponentID, data interface{}) interface{} {
	return r.getSinglePool(component).construct(entity, data)
}

// func (r *Registry) replace(entity EntityID, componentID ComponentID, data interface{}) interface{} {

// }

// Remove component from entity
func (r *Registry) Remove(entity EntityID, component ComponentID) {
	r.getSinglePool(component).destroy(entity)
}

func (r *Registry) managed(componentID ComponentID) bool {
	ctype := r.componentFamily[componentID]
	return ctype < len(r.pools) && r.pools[ctype] != nil
}

// Has check entity has coms set
func (r *Registry) Has(entity EntityID, coms ...ComponentID) bool {
	for _, com := range coms {
		if (r.managed(com) && r.getSinglePool(com).Has(entity)) == false {
			return false
		}
	}
	return true
}

// Get coms from entity
func (r *Registry) Get(entity EntityID, coms ...ComponentID) map[ComponentID]interface{} {
	// ret := make([]interface{}, 0, len(coms))
	// for _, com := range coms {
	// 	ret = append(ret, r.getSinglePool(com).Get(entity))
	// }
	// return ret
	if r.Has(entity, coms...) {
		// TODO  需要检查这里的 gc
		ret := make(map[ComponentID]interface{}, len(coms))
		for _, com := range coms {
			// ret = append(ret, r.getSinglePool(com).Get(entity))
			ret[com] = r.getSinglePool(com).Get(entity)
		}
		return ret
	}
	return nil
}

// GetSingle get single component from entity
func (r *Registry) GetSingle(entity EntityID, component ComponentID) interface{} {
	return r.getSinglePool(component).Get(entity)
}

func (r *Registry) getSinglePool(component ComponentID) *_Pool {
	ctype := r.componentFamily[component]
	return r.pools[ctype]
}

// func (r *Registry) assure(componentID ComponentID) {
// 	ctype := r.componentFamily[componentID]
// 	if !(ctype < len(r.pools)) {
// 		// ExtendSparseSetSliceWithValue(r.pools, ctype+1, newPool(r).SparseSet)
// 		r.pools = extendPoolWithValue(r.pools, ctype+1, nil)
// 	}
// 	if r.pools[ctype] == nil {
// 		r.pools[ctype] = newPool(r, componentID)
// 	}
// }

// RegisterComponent register component by name, name should unique
func (r *Registry) RegisterComponent(name string, persistent bool) ComponentID {
	if _, ok := r.checkComponentFamily[name]; ok {
		panic(fmt.Sprintf("register same name component %v", name))
		// return ComponentID(tid)
	}
	componentID := len(r.checkComponentFamily)
	r.checkComponentFamily[name] = componentID

	r.componentFamily[ComponentID(componentID)] = componentID

	if !(componentID < len(r.pools)) {
		// ExtendSparseSetSliceWithValue(r.pools, ctype+1, newPool(r).SparseSet)
		r.pools = extendPoolWithValue(r.pools, componentID+1, nil)
	}
	if r.pools[componentID] == nil {
		r.pools[componentID] = newPool(r, ComponentID(componentID))
	}

	return ComponentID(componentID)
}

// View by coms
func (r *Registry) View(coms ...ComponentID) *View {
	pools := make([]*SparseSet2, 0, len(coms))
	for _, com := range coms {
		pools = append(pools, r.getSinglePool(com).SparseSet2)
	}
	return NewView(pools...)
}

// SingleView by single com
func (r *Registry) SingleView(com ComponentID) *SingleView {
	return NewSingleView(r.getSinglePool(com).SparseSet2)
}

// Destroy entity
func (r *Registry) Destroy(entity EntityID) {
	for pos := len(r.pools); pos > 0; pos-- {
		pool := r.pools[pos-1]
		if pool != nil && pool.Has(entity) {
			pool.destroy(entity)
		}
	}

	// lengthens the implicit list of destroyed entities
	entt := entity & entity_mask
	// 该实体版本号 +1 后的数据
	version := ((entity >> entity_shift) + 1) << entity_shift
	node := r.next
	// 如果没有可用的 entity时，则更新版本号（使用上面算好的值）
	if r.available == 0 {
		node = ((entt + 1) & entity_mask) | version
	}
	r.entities[entt] = node
	r.next = entt
	r.available++
}

// Replace entity com data with newData
func (r *Registry) Replace(entity EntityID, com ComponentID, newData interface{}) {
	cpool := r.getSinglePool(com)
	cpool.SparseSet2.Replace(entity, newData)
}

// --------------------- pool ---------------------

// _Pool
// entt/registry.hpp 这里会有 实体的component的构造和销毁通知
type _Pool struct {
	*SparseSet2

	r *Registry
}

func newPool(r *Registry, com ComponentID) *_Pool {
	return &_Pool{
		SparseSet2: NewSparseSet2(com),
		r:          r,
	}
}

func (p *_Pool) construct(entity EntityID, data interface{}) interface{} {
	p.SparseSet2.Construct(entity, data)
	return data
}

func (p *_Pool) destroy(entity EntityID) {
	p.SparseSet2.Destroy(entity)
}

type _TypeMap map[ComponentID]int
type _NameMap map[string]int

// DefaultRegistry default registry
var DefaultRegistry = newRegistry()
