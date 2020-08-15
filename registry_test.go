package entt

import (
	"fmt"
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

func Test_Registry(t *testing.T) {

	registry := newRegistry()

	name := registry.RegisterComponent("name", false)
	age := registry.RegisterComponent("age", false)
	position := registry.RegisterComponent("position", false)
	mapIntInt := registry.RegisterComponent("mapintint", false)

	printNameAgePositionMap := func() {
		fmt.Println()
		t.Log("Name & Age & Position & MapIntInt")
		registry.View(name, age, position, mapIntInt).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
			t.Logf("entity %v %+v", entity, datas[name].(Name))
			t.Logf("entity %v %+v", entity, datas[age].(Age))
			t.Logf("entity %v %+v", entity, datas[position].(*Position))
			t.Logf("entity %v %+v", entity, datas[mapIntInt].(MapIntInt))
		})
	}
	printNameAgePosition := func() {
		fmt.Println()
		t.Log("Name & Age & Position")
		registry.View(name, age, position, mapIntInt).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
			t.Logf("entity %v %+v", entity, datas[name].(Name))
			t.Logf("entity %v %+v", entity, datas[age].(Age))
			t.Logf("entity %v %+v", entity, datas[position].(*Position))
		})
	}
	printNameAge := func() {
		fmt.Println()
		t.Log("Name & Age")
		registry.View(name, age).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
			t.Logf("entity %v %+v", entity, datas[name].(Name))
			t.Logf("entity %v %+v", entity, datas[age].(Age))
		})
	}
	printNamePosition := func() {
		fmt.Println()
		t.Log("Name & Position")
		registry.View(name, position).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
			t.Logf("entity %v %+v", entity, datas[name].(Name))
			t.Logf("entity %v %+v", entity, datas[position].(*Position))
		})
	}
	printAgePosition := func() {
		fmt.Println()
		t.Log("Age & Position")
		registry.View(age, position).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
			t.Logf("entity %v %+v", entity, datas[age].(Age))
			t.Logf("entity %v %+v", entity, datas[position].(*Position))
		})
	}
	printAge := func() {
		fmt.Println()
		t.Log("Age")
		registry.SingleView(age).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
			t.Logf("entity %v %+v", entity, datas[age].(Age))
		})
	}
	printPosition := func() {
		fmt.Println()
		t.Log("Position")
		registry.SingleView(position).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
			t.Logf("entity %v %+v", entity, datas[position].(*Position))
		})
	}
	printName := func() {
		fmt.Println()
		t.Log("Name")
		registry.SingleView(name).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
			t.Logf("entity %v %+v", entity, datas[name].(Name))
		})
	}

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
		printNameAgePosition()

		Equal(t, registry.View(name).size(), 4)
		Equal(t, registry.View(position).size(), 2)
		Equal(t, registry.View(age).size(), 3)
		Equal(t, registry.View(mapIntInt).size(), 1)

		Equal(t, registry.View(name, position).size(), 2)
		Equal(t, registry.View(name, age).size(), 3)
		Equal(t, registry.View(name, mapIntInt).size(), 1)

		Equal(t, registry.View(name, position, age).size(), 1)
		Equal(t, registry.View(name, position, mapIntInt).size(), 1)
		Equal(t, registry.View(position, age, mapIntInt).size(), 1)
		Equal(t, registry.View(name, age, mapIntInt).size(), 1)

		Equal(t, registry.View(name, age, position, mapIntInt).size(), 1)

		printNameAgePosition()
		printNameAgePositionMap()
		printNameAge()
		printNamePosition()
		printAgePosition()
		printName()
		printAge()
		printPosition()

		registry.Destroy(entity0)
		registry.Destroy(entity1)
		registry.Destroy(entity2)
	}

	printNameAgePosition()
	printNameAge()
	printNamePosition()
	printAgePosition()
	printName()
	printAge()
	printPosition()
	{
		entity0 := registry.Create()
		registry.Assign(entity0, name, Name("zhanglei"))
		registry.Assign(entity0, age, Age(1000))
		registry.Assign(entity0, position, &Position{X: 10, Y: 15, Z: 20})

		entity1 := registry.Create()
		registry.Assign(entity1, name, Name("lili"))
		registry.Assign(entity1, age, Age(2000))

		entity2 := registry.Create()
		registry.Assign(entity2, name, Name("yuheng"))
		registry.Assign(entity2, position, &Position{X: 100, Y: 150, Z: 200})

		entity3 := registry.Create()
		registry.Assign(entity3, name, Name("zhangqi"))
		registry.Assign(entity3, age, Age(3000))

		printNameAgePosition()
		printNameAge()
		printNamePosition()
		printAgePosition()
		printName()
		printAge()
		printPosition()

		registry.Destroy(entity0)
		registry.Destroy(entity1)
		registry.Destroy(entity2)
		registry.Destroy(entity3)

		printNameAgePosition()
		printNameAge()
		printNamePosition()
		printAgePosition()
		printName()
		printAge()
		printPosition()

		entity4 := registry.Create()
		registry.Assign(entity4, name, Name("zhanglei"))
		registry.Assign(entity4, age, Age(1000))
		registry.Assign(entity4, position, &Position{X: 10, Y: 15, Z: 20})

		entity5 := registry.Create()
		registry.Assign(entity5, name, Name("lili"))
		registry.Assign(entity5, age, Age(2000))

		entity6 := registry.Create()
		registry.Assign(entity6, name, Name("yuheng"))
		registry.Assign(entity6, position, &Position{X: 100, Y: 150, Z: 200})

		entity7 := registry.Create()
		registry.Assign(entity7, name, Name("zhangqi"))
		registry.Assign(entity7, age, Age(3000))

		printNameAgePosition()
		printNameAge()
		printNamePosition()
		printAgePosition()
		printName()
		printAge()
		printPosition()
	}

}
