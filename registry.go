package gecs

import (
	"fmt"
)

// Registry of one independent ecs system
type Registry struct {
	pools    []*_PoolHandler
	entities []EntityID

	destroyed EntityID

	// componentName -> componentID
	checkComponentFamily _NameMap
}

func newRegistry() *Registry {
	return &Registry{
		pools:    make([]*_PoolHandler, 0),
		entities: make([]EntityID, 0, 1<<19), // 最多同时52w 左右个实体

		checkComponentFamily: make(_NameMap, 0),

		destroyed: DefaultPlaceholder,
	}
}

// SizeOfCom  Returns the number of existing components of the given type.
func (r *Registry) SizeOfCom(componentID ComponentID) int {
	return r.pools[componentID].Size()
}

// Size Returns the number of entities created so far.
func (r *Registry) Size() int {
	return len(r.entities)
}

func (r *Registry) Alive() int {
	// sz := len(r.entities)
	// curr := destroyed

	return 0
}

// ReserveComs Increases the capacity of the pools for the given components.
func (r *Registry) ReserveComs(cap int, coms ...ComponentID) {
	for _, com := range coms {
		r.pools[com].Reserve(cap)
	}
}

// Reserve Increases the capacity of the registry. that is the number of entities it contains.
func (r *Registry) Reserve(cap int) {
	extendEntitySliceWithValue(r.entities, cap, DefaultPlaceholder)
}

// CapacityOfCom Returns the capacity of the pool for the given component.
func (r *Registry) CapacityOfCom(com ComponentID) int {
	return r.pools[com].Capacity()
}

// Capacity Returns the number of entities that a registry has currently
// allocated space for.
func (r *Registry) Capacity() int {
	return cap(r.entities)
}

// EmptyOfComs Checks whether the pools of the given components are empty.
func (r *Registry) EmptyOfComs(coms ...ComponentID) bool {
	for _, com := range coms {
		if r.pools[com].Empty() == false {
			return false
		}
	}
	return true
}

// Empty checks whether the registry is empty
//	A registry is considered empty when it doesn't contain entities that are
//	still in use.
func (r *Registry) Empty() bool {
	return r.Alive() == 0
}

// RawOfCom Direct access to the list of components of a given pool
func (r *Registry) RawOfCom(com ComponentID) []interface{} {
	return r.pools[com].Raw()
}

// DataOfCom Direct access to the list of entities of a given pool
func (r *Registry) DataOfCom(com ComponentID) []EntityID {
	return r.pools[com].Data()
}

// Data Direct access to the list of entities of a registry
func (r *Registry) Data() []EntityID {
	return r.entities
}

// Valid check if an entity identifier refers to a valid entity
func (r *Registry) Valid(entity EntityID) bool {
	pos := int(entity) & entity_mask
	return pos < len(r.entities) && r.entities[pos] == entity
}

// Version Returns the version stored along with an entity identifier
func (r *Registry) Version(entity EntityID) int {
	return int(entity) >> entity_shift
}

// Current Returns the actual version for an entity identifier
//	make sure the entity is belong to the registry
func (r *Registry) Current(entity EntityID) int {
	pos := int(entity) & entity_mask
	return int(r.entities[pos]) >> entity_shift
}

// Create a new entity and return it
// There are two kinds of possible entity identifiers:
//
//	* Newly created ones in case no entities have been previously destroyed.
//	* Recycled ones with updated versions.
func (r *Registry) Create() EntityID {

	var entt EntityID
	if r.destroyed == DefaultPlaceholder {
		entt = EntityID(len(r.entities))
		r.entities = append(r.entities, entt)
	} else {
		curr := int(r.destroyed)
		version := int(r.entities[curr]) & (version_mask << entity_shift)
		r.destroyed = EntityID(int(r.entities[curr]) & entity_mask)
		entt = EntityID(curr | version)
		r.entities[curr] = entt
	}

	return entt
	// var entity EntityID
	// if r.available > 0 {
	// 	// 有实体销毁时，复用该entityID，更新 Version
	// 	entt := r.next
	// 	// 拿到已经被销毁的entt的 version
	// 	version := r.entities[int(entt)] & (version_mask << entity_shift)
	// 	r.next = r.entities[int(entt)] & entity_mask
	// 	entity = entt | version
	// 	r.entities[int(entt)] = entity
	// 	r.available--
	// } else {
	// 	// 在没有任何实体销毁(destroy)时
	// 	entity = EntityID(len(r.entities))
	// 	r.entities = append(r.entities, entity)
	// }

	// return entity
}

// Destroy an entity
func (r *Registry) Destroy(entity EntityID) {
	version := int(entity)>>entity_shift + 1

	r.RemoveAll(entity)

	// lengthens the implicit list of destroyed entities
	entt := int(entity) & entity_mask
	r.entities[entt] = EntityID(int(r.destroyed) | (version << entity_shift))
	r.destroyed = EntityID(entt)

	// for pos := len(r.pools); pos > 0; pos-- {
	// 	pool := r.pools[pos-1]
	// 	if pool != nil && pool.Has(entity) {
	// 		pool.remove(r, entity)
	// 	}
	// }

	// // lengthens the implicit list of destroyed entities
	// entt := entity & entity_mask
	// // 该实体版本号 +1 后的数据
	// version := ((entity >> entity_shift) + 1) << entity_shift
	// node := r.next
	// // 如果没有可用的 entity时，则更新版本号（使用上面算好的值）
	// if r.available == 0 {
	// 	node = ((entt + 1) & entity_mask) | version
	// }
	// r.entities[entt] = node
	// r.next = entt
	// r.available++
}

// Assign the given component data to an entity
func (r *Registry) Assign(entity EntityID, com ComponentID, data interface{}) interface{} {
	return r.pools[com].emplace(r, entity, data)
}

// Patch the given component for an entity.
// The `modify` function arg, return must be the same type with the `Patch` function return
func (r *Registry) Patch(entity EntityID, com ComponentID, modify func(src interface{}) interface{}) interface{} {
	return r.pools[com].patch(r, entity, modify)
}

// Remove the given components from an entity
func (r *Registry) Remove(entity EntityID, coms ...ComponentID) {
	for _, com := range coms {
		r.pools[com].remove(r, entity)
	}
}

// RemoveIfExist remove the given components from an entity
//	Equivalent to the following snippet (pseudocode):
//	if registry.Has(entity, com[i]) {
//		registry.Remove(enttiy, com[i])
//	}
func (r *Registry) RemoveIfExist(entity EntityID, coms ...ComponentID) int {
	removedComsCount := 0
	for _, com := range coms {
		p := r.pools[com]
		if p.Has(entity) {
			p.remove(r, entity)
			removedComsCount++
		}
	}
	return removedComsCount
}

// RemoveAll removes all the components from an enttiy and makes it orphaned.
func (r *Registry) RemoveAll(entity EntityID) {
	for i := len(r.pools); i > 0; i-- {
		p := r.pools[i-1]
		if p.Has(entity) {
			p.remove(r, entity)
		}
	}
}

// Has checks if an entity has all the given components
func (r *Registry) Has(entity EntityID, coms ...ComponentID) bool {
	for _, com := range coms {
		if r.pools[com].Has(entity) == false {
			return false
		}
	}
	return true
}

// Any checks if an entity has at least one of the given components
func (r *Registry) Any(entity EntityID, coms ...ComponentID) bool {
	for _, com := range coms {
		if r.pools[com].Has(entity) {
			return true
		}
	}
	return false
}

// func (r *Registry) ReserveComponent(cap int, coms ...ComponentID) {

// }

// func (r *Registry) managed(componentID ComponentID) bool {
// 	ctype := r.componentFamily[componentID]
// 	return ctype < len(r.pools) && r.pools[ctype] != nil
// }

// // Has check entity has coms set
// func (r *Registry) Has(entity EntityID, coms ...ComponentID) bool {
// 	for _, com := range coms {
// 		if (r.managed(com) && r.getSinglePool(com).Has(entity)) == false {
// 			return false
// 		}
// 	}
// 	return true
// }

// Get Returns references to the given components for an entity
func (r *Registry) Get(entity EntityID, coms ...ComponentID) map[ComponentID]interface{} {
	if r.Has(entity, coms...) {
		// TODO  should check gc
		ret := make(map[ComponentID]interface{}, len(coms))
		for _, com := range coms {
			ret[com] = r.pools[com].Get(entity)
		}
		return ret
	}
	return nil
}

// TryGet Returns pointers to the given components for an entity.
//	Some component data maybe is nil.
func (r *Registry) TryGet(entity EntityID, coms ...ComponentID) map[ComponentID]interface{} {
	// TODO  should check gc
	ret := make(map[ComponentID]interface{}, len(coms))
	for _, com := range coms {
		ret[com] = r.pools[com].TryGet(entity)
	}
	return ret
}

// GetSingle Returns references to the given components for an entity
func (r *Registry) GetSingle(entity EntityID, com ComponentID) interface{} {
	p := r.pools[com]
	if !p.Has(entity) {
		panic(fmt.Sprintf("should have the entity %v", entity))
	}
	return r.pools[com].Get(entity)
}

// TryGetSingle Returns pointers to the given component for an entity.
//	component data maybe is nil.
func (r *Registry) TryGetSingle(entity EntityID, com ComponentID) interface{} {
	return r.pools[com].TryGet(entity)
}

// Clear a whole registry or the pools for the given components
func (r *Registry) Clear(coms ...ComponentID) {
	// if len(coms) == 0 {

	// } else {
	// 	for _, com := range coms {
	// 		p := r.pools[com]

	// 		p.pool.SparseSet.be
	// 	}
	// }
}

// Each Iterates all the entities that are still in use
func (r *Registry) Each(fn func(e EntityID)) {
	if r.destroyed == DefaultPlaceholder {
		for pos := len(r.entities); pos > 0; pos-- {
			fn(r.entities[pos-1])
		}
	} else {
		for pos := len(r.entities); pos > 0; pos-- {
			entt := r.entities[pos-1]
			if (int(entt) & entity_mask) == (pos - 1) {
				fn(entt)
			}
		}
	}
}

// Orphan Checks if an entity has components assigned
//	@returns True if the entity has no components assigned, false otherwise.
func (r *Registry) Orphan(entity EntityID) bool {
	for _, p := range r.pools {
		if p.Has(entity) {
			return false
		}
	}
	return true
}

// // OnConstruct Returns a sink object for the given component
// //
// // A sink is an opaque object used to connect listeners to components.<br/>
// //  The sink returned by this function can be used to receive notifications
// //  whenever a new instance of the given component is created and assigned to
// //  an entity.
// //
// //
// //  Listeners are invoked **after** the component has been assigned to the entity.
// //
// func (r *Registry) OnConstruct(com ComponentID) {

// }

// RegisterComponent register component by name, name should unique
//	like entt's register.assure
func (r *Registry) RegisterComponent(name string, persistent bool) ComponentID {
	if _, ok := r.checkComponentFamily[name]; ok {
		panic(fmt.Sprintf("register same name component %v", name))
		// return ComponentID(tid)
	}

	return r.assure(name).com
}

func (r *Registry) assure(name string) *_PoolHandler {
	var cid int
	var ok bool
	if cid, ok = r.checkComponentFamily[name]; !ok {
		cid = len(r.checkComponentFamily)
		r.checkComponentFamily[name] = cid
	}
	if !(cid < len(r.pools)) {
		r.pools = extendPoolHandlerWithValue(r.pools, cid+1, func() *_PoolHandler { return newPool(ComponentID(cid)) })
	}
	return r.pools[cid]
}

// View by coms
func (r *Registry) View(coms ...ComponentID) *View {
	pools := make([]*Storage, 0, len(coms))
	for _, com := range coms {
		pools = append(pools, r.pools[com].Storage)
	}
	return NewView(pools...)
}

// SingleView by single com
func (r *Registry) SingleView(com ComponentID) *SingleView {
	return NewSingleView(r.pools[com].Storage)
}

// Replace entity com data with newData
func (r *Registry) Replace(entity EntityID, com ComponentID, newData interface{}) interface{} {
	return r.pools[com].Replace(entity, newData)
}

// --------------------- pool ---------------------

// _PoolHandler
// entt/registry.hpp 这里会有 实体的component的构造和销毁通知
type _PoolHandler struct {
	*Storage
}

func newPool(com ComponentID) *_PoolHandler {
	return &_PoolHandler{
		Storage: NewStorage(com),
	}
}

// func (p *_PoolHandler) onconstruct(entity EntityID, data interface{}) interface{} {
// 	p.Storage.Emplace(entity, data)
// 	return data
// }

// func (p *_PoolHandler) onupdate() interface{} {
// 	return nil
// }

// func (p *_PoolHandler) onremove() interface{} {
// 	return nil
// }

func (p *_PoolHandler) emplace(owner *Registry, entity EntityID, data interface{}) interface{} {
	p.Storage.Emplace(entity, data)
	// publish emplace
	// construction.publish(owner, entity)
	return p.Get(entity)
}

func (p *_PoolHandler) remove(owner *Registry, entity EntityID) {
	// publish remove
	// destruction.publish(owner, entityID)
	p.Destroy(entity)
}

func (p *_PoolHandler) patch(owner *Registry, entity EntityID, modify func(srcData interface{}) interface{}) interface{} {
	// publish update

	return p.Replace(entity, modify(p.Get(entity)))
}

func (p *_PoolHandler) replace(owner *Registry, entity EntityID, data interface{}) interface{} {
	return p.patch(owner, entity, func(srcData interface{}) interface{} {
		return data
	})
}

type _TypeMap map[ComponentID]int
type _NameMap map[string]int

// DefaultRegistry default registry
var DefaultRegistry = newRegistry()
