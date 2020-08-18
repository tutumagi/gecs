package gecs

import (
	"math"
	"strconv"
)

// // IEntity entity interface
// type IEntity interface {
// 	SetRID(EntityID)
// 	RID() EntityID
// }

// //ComponentID component interface
// typeComponentID interface {
// 	SetRID(ComponentID)
// 	RID() ComponentID
// }

// EntityID entity runtime id
type EntityID uint32

func (e EntityID) String() string {
	if DefaultPlaceholder == e {
		return "^"
	}
	return strconv.FormatUint(uint64(e), 10)
}

// DefaultPlaceholder default entity id
const DefaultPlaceholder EntityID = math.MaxUint32

// ComponentID component runtime id
type ComponentID uint32

// // Interfacer of ecs
// type Interfacer interface {
// 	RegisterComponent(name string) ComponentID
// 	Create() EntityID
// 	Assign(entity EntityID, com ComponentID, data interface{}) interface{}
// 	Has(entity EntityID, com ComponentID) bool
// 	Remove(entity EntityID, com ComponentID)
// 	Destroy(entity EntityID)
// 	Get(entity EntityID, coms ...ComponentID) []interface{}
// 	GetSingle(entity EntityID, com ComponentID) (interface{}, string, bool)
// 	Replace(entity EntityID, id ComponentID, data interface{})
// }

// // IteratorView iterator view
// type IteratorView interface {
// 	Each(fn func(entity EntityID, components map[ComponentID]interface{}))
// }
