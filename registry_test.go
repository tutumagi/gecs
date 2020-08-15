package entt

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

type Name string
type Age int
type Position struct {
	X int
	Y int
	Z int
}
type MapIntInt map[int32]int32

var name ComponentID
var age ComponentID
var position ComponentID
var mapIntInt ComponentID

func initRegistry() *Registry {
	registry := newRegistry()

	name = registry.RegisterComponent("name", false)
	age = registry.RegisterComponent("age", false)
	position = registry.RegisterComponent("position", false)
	mapIntInt = registry.RegisterComponent("mapintint", false)

	return registry
}

func Test_RegistryCreateAndDestroy(t *testing.T) {

	registry := initRegistry()

	{
		entity0 := registry.Create()
		registry.Assign(entity0, name, Name("tufei"))
		registry.Assign(entity0, age, Age(12))
		registry.Assign(entity0, position, &Position{X: 10, Y: 15, Z: 20})
		registry.Assign(entity0, mapIntInt, MapIntInt{3: 4, 5: 6, 7: 8})

		entity1 := registry.Create()
		registry.Assign(entity1, name, Name("weijie"))
		registry.Assign(entity1, age, Age(20))

		entity2 := registry.Create()
		registry.Assign(entity2, name, Name("yihao"))
		registry.Assign(entity2, position, &Position{X: 100, Y: 150, Z: 200})

		entity3 := registry.Create()
		registry.Assign(entity3, name, Name("zhijun"))
		registry.Assign(entity3, age, Age(200))

		Equal(t, registry.View(name).size(), 4)
		Equal(t, registry.View(position).size(), 2)
		Equal(t, registry.View(age).size(), 3)
		Equal(t, registry.View(mapIntInt).size(), 1)

		equalCount(t, registry, 2, name, position)
		equalCount(t, registry, 3, name, age)
		equalCount(t, registry, 1, name, mapIntInt)

		equalCount(t, registry, 1, name, position, age)
		equalCount(t, registry, 1, name, position, mapIntInt)
		equalCount(t, registry, 1, position, age, mapIntInt)
		equalCount(t, registry, 1, name, age, mapIntInt)

		equalCount(t, registry, 1, name, age, position, mapIntInt)

		registry.Destroy(entity0)

		Equal(t, registry.View(name).size(), 3)
		Equal(t, registry.View(position).size(), 1)
		Equal(t, registry.View(age).size(), 2)
		Equal(t, registry.View(mapIntInt).size(), 0)
		equalCount(t, registry, 1, name, position)
		equalCount(t, registry, 2, name, age)
		equalCount(t, registry, 0, name, mapIntInt)

		equalCount(t, registry, 0, name, position, age)
		equalCount(t, registry, 0, name, position, mapIntInt)
		equalCount(t, registry, 0, position, age, mapIntInt)
		equalCount(t, registry, 0, name, age, mapIntInt)

		equalCount(t, registry, 0, name, age, position, mapIntInt)

		registry.Destroy(entity1)
		Equal(t, registry.View(name).size(), 2)
		Equal(t, registry.View(position).size(), 1)
		Equal(t, registry.View(age).size(), 1)
		Equal(t, registry.View(mapIntInt).size(), 0)
		equalCount(t, registry, 1, name, position)
		equalCount(t, registry, 1, name, age)
		equalCount(t, registry, 0, name, mapIntInt)

		equalCount(t, registry, 0, name, position, age)
		equalCount(t, registry, 0, name, position, mapIntInt)
		equalCount(t, registry, 0, position, age, mapIntInt)
		equalCount(t, registry, 0, name, age, mapIntInt)

		equalCount(t, registry, 0, name, age, position, mapIntInt)

		registry.Destroy(entity2)
		Equal(t, registry.View(name).size(), 1)
		Equal(t, registry.View(position).size(), 0)
		Equal(t, registry.View(age).size(), 1)
		Equal(t, registry.View(mapIntInt).size(), 0)
		equalCount(t, registry, 0, name, position)
		equalCount(t, registry, 1, name, age)
		equalCount(t, registry, 0, name, mapIntInt)

		equalCount(t, registry, 0, name, position, age)
		equalCount(t, registry, 0, name, position, mapIntInt)
		equalCount(t, registry, 0, position, age, mapIntInt)
		equalCount(t, registry, 0, name, age, mapIntInt)

		equalCount(t, registry, 0, name, age, position, mapIntInt)

		entity4 := registry.Create()
		registry.Assign(entity4, name, Name("tt2"))
		registry.Assign(entity4, position, &Position{X: 10, Y: 20, Z: 100})
		Equal(t, registry.View(name).size(), 2)
		Equal(t, registry.View(position).size(), 1)
		Equal(t, registry.View(age).size(), 1)
		Equal(t, registry.View(mapIntInt).size(), 0)
		equalCount(t, registry, 1, name, position)
		equalCount(t, registry, 1, name, age)
		equalCount(t, registry, 0, name, mapIntInt)

		equalCount(t, registry, 0, name, position, age)
		equalCount(t, registry, 0, name, position, mapIntInt)
		equalCount(t, registry, 0, position, age, mapIntInt)
		equalCount(t, registry, 0, name, age, mapIntInt)

		equalCount(t, registry, 0, name, age, position, mapIntInt)

		printNameAgePosition(t, registry)
		printNameAgePositionMap(t, registry)
		printNameAge(t, registry)
		printNamePosition(t, registry)
		printAgePosition(t, registry)
		printName(t, registry)
		printAge(t, registry)
		printPosition(t, registry)
	}
}

func Test_Assign(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, name, Name("zhanglei"))
	registry.Assign(entity0, age, Age(1000))
	registry.Assign(entity0, position, &Position{X: 10, Y: 15, Z: 20})

	Equal(t, registry.GetSingle(entity0, name).(Name), Name("zhanglei"))
	Equal(t, registry.GetSingle(entity0, age).(Age), Age(1000))
	Equal(t, registry.GetSingle(entity0, position).(*Position), &Position{X: 10, Y: 15, Z: 20})

	{
		components := registry.Get(entity0, name, age, position)
		Equal(t, components[name].(Name), Name("zhanglei"))
		Equal(t, components[age].(Age), Age(1000))
		Equal(t, components[position].(*Position), &Position{X: 10, Y: 15, Z: 20})
	}

	{
		components := registry.Get(entity0, name, age)
		Equal(t, components[name].(Name), Name("zhanglei"))
		Equal(t, components[age].(Age), Age(1000))
	}

	{
		components := registry.Get(entity0, name, position)
		Equal(t, components[name].(Name), Name("zhanglei"))
		Equal(t, components[position].(*Position), &Position{X: 10, Y: 15, Z: 20})
	}

	{
		components := registry.Get(entity0, age, position)
		Equal(t, components[age].(Age), Age(1000))
		Equal(t, components[position].(*Position), &Position{X: 10, Y: 15, Z: 20})
	}
}

func Test_Replace(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, name, Name("zhanglei"))
	registry.Assign(entity0, age, Age(1000))
	registry.Assign(entity0, position, &Position{X: 10, Y: 15, Z: 20})

	registry.Replace(entity0, name, Name("tutu"))
	Equal(t, registry.GetSingle(entity0, name).(Name), Name("tutu"))

	registry.Replace(entity0, position, &Position{X: 20, Y: 100, Z: 200})
	Equal(t, registry.GetSingle(entity0, position).(*Position), &Position{X: 20, Y: 100, Z: 200})

	registry.Replace(entity0, age, Age(200000))
	Equal(t, registry.GetSingle(entity0, age).(Age), Age(200000))
}

func Test_Remove(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, name, Name("zhanglei"))
	registry.Assign(entity0, age, Age(1000))
	registry.Assign(entity0, position, &Position{X: 10, Y: 15, Z: 20})

	Equal(t, registry.Has(entity0, name), true)
	registry.Remove(entity0, name)
	Equal(t, registry.Has(entity0, name), false)

	Equal(t, registry.Has(entity0, position), true)
	registry.Remove(entity0, position)
	Equal(t, registry.Has(entity0, position), false)

	Equal(t, registry.Has(entity0, age), true)
	registry.Remove(entity0, age)
	Equal(t, registry.Has(entity0, age), false)
}

func Test_ViewEach(t *testing.T) {

}

func Test_SingleViewEach(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, name, Name("zhanglei"))
	registry.Assign(entity0, age, Age(1000))
	registry.Assign(entity0, position, &Position{X: 10, Y: 15, Z: 20})

	registry.SingleView(name).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[name].(Name), Name("zhanglei"))
	})
	registry.SingleView(age).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[age].(Age), Age(1000))
	})
	registry.SingleView(position).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[position].(*Position), &Position{X: 10, Y: 15, Z: 20})
	})
}

func equalCount(t *testing.T, registry *Registry, count int, coms ...ComponentID) {
	cal := 0
	registry.View(coms...).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		cal++
	})
	EqualSkip(t, 2, cal, count)
}

func printNameAgePositionMap(t *testing.T, registry *Registry) {
	t.Log("Name & Age & Position & MapIntInt")
	registry.View(name, age, position, mapIntInt).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		t.Logf("entity %v %+v", entity, datas[name].(Name))
		t.Logf("entity %v %+v", entity, datas[age].(Age))
		t.Logf("entity %v %+v", entity, datas[position].(*Position))
		t.Logf("entity %v %+v", entity, datas[mapIntInt].(MapIntInt))
	})
}
func printNameAgePosition(t *testing.T, registry *Registry) {
	t.Log("Name & Age & Position")
	registry.View(name, age, position, mapIntInt).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		t.Logf("entity %v %+v", entity, datas[name].(Name))
		t.Logf("entity %v %+v", entity, datas[age].(Age))
		t.Logf("entity %v %+v", entity, datas[position].(*Position))
	})
}
func printNameAge(t *testing.T, registry *Registry) {
	t.Log("Name & Age")
	registry.View(name, age).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		t.Logf("entity %v %+v", entity, datas[name].(Name))
		t.Logf("entity %v %+v", entity, datas[age].(Age))
	})
}
func printNamePosition(t *testing.T, registry *Registry) {
	t.Log("Name & Position")
	registry.View(name, position).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		t.Logf("entity %v %+v", entity, datas[name].(Name))
		t.Logf("entity %v %+v", entity, datas[position].(*Position))
	})
}
func printAgePosition(t *testing.T, registry *Registry) {
	t.Log("Age & Position")
	registry.View(age, position).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		t.Logf("entity %v %+v", entity, datas[age].(Age))
		t.Logf("entity %v %+v", entity, datas[position].(*Position))
	})
}
func printAge(t *testing.T, registry *Registry) {
	t.Log("Age")
	registry.SingleView(age).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		t.Logf("entity %v %+v", entity, datas[age].(Age))
	})
}
func printPosition(t *testing.T, registry *Registry) {
	t.Log("Position")
	registry.SingleView(position).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		t.Logf("entity %v %+v", entity, datas[position].(*Position))
	})
}
func printName(t *testing.T, registry *Registry) {
	t.Log("Name")
	registry.SingleView(name).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		t.Logf("entity %v %+v", entity, datas[name].(Name))
	})
}
