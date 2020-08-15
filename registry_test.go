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

var NameID ComponentID
var AgeID ComponentID
var PositionID ComponentID
var MapIntIntID ComponentID

func initRegistry() *Registry {
	registry := newRegistry()

	NameID = registry.RegisterComponent("name", false)
	AgeID = registry.RegisterComponent("age", false)
	PositionID = registry.RegisterComponent("position", false)
	MapIntIntID = registry.RegisterComponent("mapintint", false)

	return registry
}

func Test_RegistryCreateAndDestroy(t *testing.T) {

	registry := initRegistry()

	{
		entity0 := registry.Create()
		registry.Assign(entity0, NameID, Name("tufei"))
		registry.Assign(entity0, AgeID, Age(12))
		registry.Assign(entity0, PositionID, &Position{X: 10, Y: 15, Z: 20})
		registry.Assign(entity0, MapIntIntID, MapIntInt{3: 4, 5: 6, 7: 8})

		entity1 := registry.Create()
		registry.Assign(entity1, NameID, Name("weijie"))
		registry.Assign(entity1, AgeID, Age(20))

		entity2 := registry.Create()
		registry.Assign(entity2, NameID, Name("yihao"))
		registry.Assign(entity2, PositionID, &Position{X: 100, Y: 150, Z: 200})

		entity3 := registry.Create()
		registry.Assign(entity3, NameID, Name("zhijun"))
		registry.Assign(entity3, AgeID, Age(200))

		Equal(t, registry.View(NameID).size(), 4)
		Equal(t, registry.View(PositionID).size(), 2)
		Equal(t, registry.View(AgeID).size(), 3)
		Equal(t, registry.View(MapIntIntID).size(), 1)

		equalCount(t, registry, 2, NameID, PositionID)
		equalCount(t, registry, 3, NameID, AgeID)
		equalCount(t, registry, 1, NameID, MapIntIntID)

		equalCount(t, registry, 1, NameID, PositionID, AgeID)
		equalCount(t, registry, 1, NameID, PositionID, MapIntIntID)
		equalCount(t, registry, 1, PositionID, AgeID, MapIntIntID)
		equalCount(t, registry, 1, NameID, AgeID, MapIntIntID)

		equalCount(t, registry, 1, NameID, AgeID, PositionID, MapIntIntID)

		registry.Destroy(entity0)

		Equal(t, registry.View(NameID).size(), 3)
		Equal(t, registry.View(PositionID).size(), 1)
		Equal(t, registry.View(AgeID).size(), 2)
		Equal(t, registry.View(MapIntIntID).size(), 0)
		equalCount(t, registry, 1, NameID, PositionID)
		equalCount(t, registry, 2, NameID, AgeID)
		equalCount(t, registry, 0, NameID, MapIntIntID)

		equalCount(t, registry, 0, NameID, PositionID, AgeID)
		equalCount(t, registry, 0, NameID, PositionID, MapIntIntID)
		equalCount(t, registry, 0, PositionID, AgeID, MapIntIntID)
		equalCount(t, registry, 0, NameID, AgeID, MapIntIntID)

		equalCount(t, registry, 0, NameID, AgeID, PositionID, MapIntIntID)

		registry.Destroy(entity1)
		Equal(t, registry.View(NameID).size(), 2)
		Equal(t, registry.View(PositionID).size(), 1)
		Equal(t, registry.View(AgeID).size(), 1)
		Equal(t, registry.View(MapIntIntID).size(), 0)
		equalCount(t, registry, 1, NameID, PositionID)
		equalCount(t, registry, 1, NameID, AgeID)
		equalCount(t, registry, 0, NameID, MapIntIntID)

		equalCount(t, registry, 0, NameID, PositionID, AgeID)
		equalCount(t, registry, 0, NameID, PositionID, MapIntIntID)
		equalCount(t, registry, 0, PositionID, AgeID, MapIntIntID)
		equalCount(t, registry, 0, NameID, AgeID, MapIntIntID)

		equalCount(t, registry, 0, NameID, AgeID, PositionID, MapIntIntID)

		registry.Destroy(entity2)
		Equal(t, registry.View(NameID).size(), 1)
		Equal(t, registry.View(PositionID).size(), 0)
		Equal(t, registry.View(AgeID).size(), 1)
		Equal(t, registry.View(MapIntIntID).size(), 0)
		equalCount(t, registry, 0, NameID, PositionID)
		equalCount(t, registry, 1, NameID, AgeID)
		equalCount(t, registry, 0, NameID, MapIntIntID)

		equalCount(t, registry, 0, NameID, PositionID, AgeID)
		equalCount(t, registry, 0, NameID, PositionID, MapIntIntID)
		equalCount(t, registry, 0, PositionID, AgeID, MapIntIntID)
		equalCount(t, registry, 0, NameID, AgeID, MapIntIntID)

		equalCount(t, registry, 0, NameID, AgeID, PositionID, MapIntIntID)

		entity4 := registry.Create()
		registry.Assign(entity4, NameID, Name("tt2"))
		registry.Assign(entity4, PositionID, &Position{X: 10, Y: 20, Z: 100})
		Equal(t, registry.View(NameID).size(), 2)
		Equal(t, registry.View(PositionID).size(), 1)
		Equal(t, registry.View(AgeID).size(), 1)
		Equal(t, registry.View(MapIntIntID).size(), 0)
		equalCount(t, registry, 1, NameID, PositionID)
		equalCount(t, registry, 1, NameID, AgeID)
		equalCount(t, registry, 0, NameID, MapIntIntID)

		equalCount(t, registry, 0, NameID, PositionID, AgeID)
		equalCount(t, registry, 0, NameID, PositionID, MapIntIntID)
		equalCount(t, registry, 0, PositionID, AgeID, MapIntIntID)
		equalCount(t, registry, 0, NameID, AgeID, MapIntIntID)

		equalCount(t, registry, 0, NameID, AgeID, PositionID, MapIntIntID)

		// printNameAgePosition(t, registry)
		// printNameAgePositionMap(t, registry)
		// printNameAge(t, registry)
		// printNamePosition(t, registry)
		// printAgePosition(t, registry)
		// printName(t, registry)
		// printAge(t, registry)
		// printPosition(t, registry)
	}
}

func Test_Assign(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, NameID, Name("zhanglei"))
	registry.Assign(entity0, AgeID, Age(1000))
	registry.Assign(entity0, PositionID, &Position{X: 10, Y: 15, Z: 20})

	Equal(t, registry.GetSingle(entity0, NameID).(Name), Name("zhanglei"))
	Equal(t, registry.GetSingle(entity0, AgeID).(Age), Age(1000))
	Equal(t, registry.GetSingle(entity0, PositionID).(*Position), &Position{X: 10, Y: 15, Z: 20})

	{
		components := registry.Get(entity0, NameID, AgeID, PositionID)
		Equal(t, components[NameID].(Name), Name("zhanglei"))
		Equal(t, components[AgeID].(Age), Age(1000))
		Equal(t, components[PositionID].(*Position), &Position{X: 10, Y: 15, Z: 20})
	}

	{
		components := registry.Get(entity0, NameID, AgeID)
		Equal(t, components[NameID].(Name), Name("zhanglei"))
		Equal(t, components[AgeID].(Age), Age(1000))
	}

	{
		components := registry.Get(entity0, NameID, PositionID)
		Equal(t, components[NameID].(Name), Name("zhanglei"))
		Equal(t, components[PositionID].(*Position), &Position{X: 10, Y: 15, Z: 20})
	}

	{
		components := registry.Get(entity0, AgeID, PositionID)
		Equal(t, components[AgeID].(Age), Age(1000))
		Equal(t, components[PositionID].(*Position), &Position{X: 10, Y: 15, Z: 20})
	}
}

func Test_Replace(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, NameID, Name("zhanglei"))
	registry.Assign(entity0, AgeID, Age(1000))
	registry.Assign(entity0, PositionID, &Position{X: 10, Y: 15, Z: 20})

	registry.Replace(entity0, NameID, Name("tutu"))
	Equal(t, registry.GetSingle(entity0, NameID).(Name), Name("tutu"))

	registry.Replace(entity0, PositionID, &Position{X: 20, Y: 100, Z: 200})
	Equal(t, registry.GetSingle(entity0, PositionID).(*Position), &Position{X: 20, Y: 100, Z: 200})

	registry.Replace(entity0, AgeID, Age(200000))
	Equal(t, registry.GetSingle(entity0, AgeID).(Age), Age(200000))
}

func Test_Remove(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, NameID, Name("zhanglei"))
	registry.Assign(entity0, AgeID, Age(1000))
	registry.Assign(entity0, PositionID, &Position{X: 10, Y: 15, Z: 20})

	Equal(t, registry.Has(entity0, NameID), true)
	registry.Remove(entity0, NameID)
	Equal(t, registry.Has(entity0, NameID), false)

	Equal(t, registry.Has(entity0, PositionID), true)
	registry.Remove(entity0, PositionID)
	Equal(t, registry.Has(entity0, PositionID), false)

	Equal(t, registry.Has(entity0, AgeID), true)
	registry.Remove(entity0, AgeID)
	Equal(t, registry.Has(entity0, AgeID), false)
}

func Test_ViewEach(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, NameID, Name("zhanglei"))
	registry.Assign(entity0, AgeID, Age(1000))
	registry.Assign(entity0, PositionID, &Position{X: 10, Y: 15, Z: 20})

	registry.View(NameID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[NameID].(Name), Name("zhanglei"))
	})
	registry.View(AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[AgeID].(Age), Age(1000))
	})
	registry.View(PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[PositionID].(*Position), &Position{X: 10, Y: 15, Z: 20})
	})

	registry.View(NameID, AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[NameID].(Name), Name("zhanglei"))
		Equal(t, datas[AgeID].(Age), Age(1000))
	})
	registry.View(AgeID, PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[AgeID].(Age), Age(1000))
		Equal(t, datas[PositionID].(*Position), &Position{X: 10, Y: 15, Z: 20})
	})
	registry.View(NameID, PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[NameID].(Name), Name("zhanglei"))
		Equal(t, datas[PositionID].(*Position), &Position{X: 10, Y: 15, Z: 20})
	})
	registry.View(NameID, AgeID, PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[AgeID].(Age), Age(1000))
		Equal(t, datas[NameID].(Name), Name("zhanglei"))
		Equal(t, datas[PositionID].(*Position), &Position{X: 10, Y: 15, Z: 20})
	})

	entity1 := registry.Create()
	registry.Assign(entity1, NameID, Name("tutu"))
	registry.Assign(entity1, AgeID, Age(100))
	registry.Assign(entity1, PositionID, &Position{X: 100, Y: 150, Z: 200})

	registry.View(NameID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, NameID, Name("zhanglei"), Name("tutu")); !ok {
			t.Errorf("err name %+v", comData)
		}
	})
	registry.View(AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, AgeID, Age(1000), Age(100)); !ok {
			t.Errorf("err age %+v", comData)
		}
	})
	registry.View(PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, PositionID, &Position{X: 10, Y: 15, Z: 20}, &Position{X: 100, Y: 150, Z: 200}); !ok {
			t.Errorf("err position %+v", comData)
		}
	})

	registry.View(NameID, AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, NameID, Name("zhanglei"), Name("tutu")); !ok {
			t.Errorf("err name %+v", comData)
		}
		if comData, ok := checkComDataFromMapData(datas, AgeID, Age(1000), Age(100)); !ok {
			t.Errorf("err age %+v", comData)
		}
	})
	registry.View(AgeID, PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, AgeID, Age(1000), Age(100)); !ok {
			t.Errorf("err age %+v", comData)
		}
	})
	registry.View(PositionID, NameID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, NameID, Name("zhanglei"), Name("tutu")); !ok {
			t.Errorf("err name %+v", comData)
		}
		if comData, ok := checkComDataFromMapData(datas, PositionID, &Position{X: 10, Y: 15, Z: 20}, &Position{X: 100, Y: 150, Z: 200}); !ok {
			t.Errorf("err position %+v", comData)
		}
	})

	registry.View(PositionID, NameID, AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, NameID, Name("zhanglei"), Name("tutu")); !ok {
			t.Errorf("err name %+v", comData)
		}
		if comData, ok := checkComDataFromMapData(datas, PositionID, &Position{X: 10, Y: 15, Z: 20}, &Position{X: 100, Y: 150, Z: 200}); !ok {
			t.Errorf("err position %+v", comData)
		}
		if comData, ok := checkComDataFromMapData(datas, AgeID, Age(1000), Age(100)); !ok {
			t.Errorf("err age %+v", comData)
		}
	})
}

func Test_SingleViewEach(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, NameID, Name("zhanglei"))
	registry.Assign(entity0, AgeID, Age(1000))
	registry.Assign(entity0, PositionID, &Position{X: 10, Y: 15, Z: 20})

	registry.SingleView(NameID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[NameID].(Name), Name("zhanglei"))
	})
	registry.SingleView(AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[AgeID].(Age), Age(1000))
	})
	registry.SingleView(PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, datas[PositionID].(*Position), &Position{X: 10, Y: 15, Z: 20})
	})

	entity1 := registry.Create()
	registry.Assign(entity1, NameID, Name("tutu"))
	registry.Assign(entity1, AgeID, Age(100))
	registry.Assign(entity1, PositionID, &Position{X: 100, Y: 150, Z: 200})

	registry.SingleView(NameID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, NameID, Name("zhanglei"), Name("tutu")); !ok {
			t.Errorf("err name %+v", comData)
		}
	})
	registry.SingleView(AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, AgeID, Age(1000), Age(100)); !ok {
			t.Errorf("err age %+v", comData)
		}
	})
	registry.SingleView(PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		if comData, ok := checkComDataFromMapData(datas, PositionID, &Position{X: 10, Y: 15, Z: 20}, &Position{X: 100, Y: 150, Z: 200}); !ok {
			t.Errorf("err position %+v", comData)
		}
	})
}

func Test_GetSingle(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, NameID, Name("zhanglei"))
	registry.Assign(entity0, AgeID, Age(1000))
	registry.Assign(entity0, PositionID, &Position{X: 10, Y: 15, Z: 20})

	Equal(t, registry.GetSingle(entity0, NameID).(Name), Name("zhanglei"))
	Equal(t, registry.GetSingle(entity0, AgeID).(Age), Age(1000))
	Equal(t, registry.GetSingle(entity0, PositionID).(*Position), &Position{X: 10, Y: 15, Z: 20})

	entity1 := registry.Create()
	registry.Assign(entity1, NameID, Name("zhanglei"))
	registry.Assign(entity1, AgeID, Age(1000))

	PanicMatches(t, func() { registry.GetSingle(entity1, PositionID) }, fmt.Sprintf("should have the entity %v", entity1))
}

func Test_Contains(t *testing.T) {
	registry := initRegistry()
	entity0 := registry.Create()
	registry.Assign(entity0, NameID, Name("zhanglei"))
	registry.Assign(entity0, AgeID, Age(1000))
	registry.Assign(entity0, PositionID, &Position{X: 10, Y: 15, Z: 20})

	Equal(t, registry.SingleView(NameID).Contains(entity0), true)
	Equal(t, registry.SingleView(AgeID).Contains(entity0), true)
	Equal(t, registry.SingleView(PositionID).Contains(entity0), true)
	Equal(t, registry.SingleView(MapIntIntID).Contains(entity0), false)
	Equal(t, registry.View(NameID).contains(entity0), true)
	Equal(t, registry.View(AgeID).contains(entity0), true)
	Equal(t, registry.View(PositionID).contains(entity0), true)
	Equal(t, registry.View(MapIntIntID).contains(entity0), false)

	Equal(t, registry.View(NameID, AgeID).contains(entity0), true)
	Equal(t, registry.View(AgeID, PositionID).contains(entity0), true)
	Equal(t, registry.View(PositionID, NameID).contains(entity0), true)
	Equal(t, registry.View(NameID, MapIntIntID).contains(entity0), false)
	Equal(t, registry.View(AgeID, MapIntIntID).contains(entity0), false)
	Equal(t, registry.View(PositionID, MapIntIntID).contains(entity0), false)
	Equal(t, registry.View(NameID, AgeID, PositionID).contains(entity0), true)

	entity1 := registry.Create()
	registry.Assign(entity1, NameID, Name("tutu"))
	registry.Assign(entity1, AgeID, Age(100))
	registry.Assign(entity1, MapIntIntID, MapIntInt{3: 100, 4: 150, 5: 200})

	// check entity0
	Equal(t, registry.SingleView(NameID).Contains(entity0), true)
	Equal(t, registry.SingleView(AgeID).Contains(entity0), true)
	Equal(t, registry.SingleView(PositionID).Contains(entity0), true)
	Equal(t, registry.SingleView(MapIntIntID).Contains(entity0), false)
	Equal(t, registry.View(NameID).contains(entity0), true)
	Equal(t, registry.View(AgeID).contains(entity0), true)
	Equal(t, registry.View(PositionID).contains(entity0), true)
	Equal(t, registry.View(MapIntIntID).contains(entity0), false)

	Equal(t, registry.View(NameID, AgeID).contains(entity0), true)
	Equal(t, registry.View(AgeID, PositionID).contains(entity0), true)
	Equal(t, registry.View(PositionID, NameID).contains(entity0), true)
	Equal(t, registry.View(NameID, MapIntIntID).contains(entity0), false)
	Equal(t, registry.View(AgeID, MapIntIntID).contains(entity0), false)
	Equal(t, registry.View(PositionID, MapIntIntID).contains(entity0), false)
	Equal(t, registry.View(NameID, AgeID, PositionID).contains(entity0), true)

	// check entity1
	Equal(t, registry.SingleView(NameID).Contains(entity1), true)
	Equal(t, registry.SingleView(AgeID).Contains(entity1), true)
	Equal(t, registry.SingleView(PositionID).Contains(entity1), false)
	Equal(t, registry.SingleView(MapIntIntID).Contains(entity1), true)

	Equal(t, registry.View(NameID).contains(entity1), true)
	Equal(t, registry.View(AgeID).contains(entity1), true)
	Equal(t, registry.View(PositionID).contains(entity1), false)
	Equal(t, registry.View(MapIntIntID).contains(entity1), true)

	Equal(t, registry.View(NameID, AgeID).contains(entity1), true)
	Equal(t, registry.View(AgeID, PositionID).contains(entity1), false)
	Equal(t, registry.View(PositionID, NameID).contains(entity1), false)
	Equal(t, registry.View(MapIntIntID, AgeID).contains(entity1), true)

}

func equalCount(t *testing.T, registry *Registry, count int, coms ...ComponentID) {
	cal := 0
	registry.View(coms...).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		cal++
	})
	EqualSkip(t, 2, cal, count)
}

func checkComDataFromMapData(datas map[ComponentID]interface{}, comID ComponentID, values ...interface{}) (interface{}, bool) {
	comData := datas[comID]
	for _, expectValue := range values {
		if IsEqual(comData, expectValue) {
			return comData, true
		}
	}

	return comData, false
}

// func printNameAgePositionMap(t *testing.T, registry *Registry) {
// 	t.Log("Name & Age & Position & MapIntInt")
// 	registry.View(NameID, AgeID, PositionID, MapIntIntID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
// 		t.Logf("entity %v %+v", entity, datas[NameID].(Name))
// 		t.Logf("entity %v %+v", entity, datas[AgeID].(Age))
// 		t.Logf("entity %v %+v", entity, datas[PositionID].(*Position))
// 		t.Logf("entity %v %+v", entity, datas[MapIntIntID].(MapIntInt))
// 	})
// }
// func printNameAgePosition(t *testing.T, registry *Registry) {
// 	t.Log("Name & Age & Position")
// 	registry.View(NameID, AgeID, PositionID, MapIntIntID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
// 		t.Logf("entity %v %+v", entity, datas[NameID].(Name))
// 		t.Logf("entity %v %+v", entity, datas[AgeID].(Age))
// 		t.Logf("entity %v %+v", entity, datas[PositionID].(*Position))
// 	})
// }
// func printNameAge(t *testing.T, registry *Registry) {
// 	t.Log("Name & Age")
// 	registry.View(NameID, AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
// 		t.Logf("entity %v %+v", entity, datas[NameID].(Name))
// 		t.Logf("entity %v %+v", entity, datas[AgeID].(Age))
// 	})
// }
// func printNamePosition(t *testing.T, registry *Registry) {
// 	t.Log("Name & Position")
// 	registry.View(NameID, PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
// 		t.Logf("entity %v %+v", entity, datas[NameID].(Name))
// 		t.Logf("entity %v %+v", entity, datas[PositionID].(*Position))
// 	})
// }
// func printAgePosition(t *testing.T, registry *Registry) {
// 	t.Log("Age & Position")
// 	registry.View(AgeID, PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
// 		t.Logf("entity %v %+v", entity, datas[AgeID].(Age))
// 		t.Logf("entity %v %+v", entity, datas[PositionID].(*Position))
// 	})
// }
// func printAge(t *testing.T, registry *Registry) {
// 	t.Log("Age")
// 	registry.SingleView(AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
// 		t.Logf("entity %v %+v", entity, datas[AgeID].(Age))
// 	})
// }
// func printPosition(t *testing.T, registry *Registry) {
// 	t.Log("Position")
// 	registry.SingleView(PositionID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
// 		t.Logf("entity %v %+v", entity, datas[PositionID].(*Position))
// 	})
// }
// func printName(t *testing.T, registry *Registry) {
// 	t.Log("Name")
// 	registry.SingleView(NameID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
// 		t.Logf("entity %v %+v", entity, datas[NameID].(Name))
// 	})
// }
