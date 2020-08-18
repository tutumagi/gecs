package gecs

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
