package entt

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

func Test_SingleView(t *testing.T) {
	nameID := ComponentID(1)
	nameStorage := NewStorage(nameID)
	singleView := NewSingleView(nameStorage)

	entity0 := EntityID(0)
	nameStorage.Emplace(entity0, "foo")

	entity1 := EntityID(1)
	Equal(t, singleView.Pool.Has(entity0), true)
	Equal(t, singleView.Pool.Has(entity1), false)
	Equal(t, singleView.Size(), 1)
	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, entity, entity0)
		Equal(t, datas[nameID].(string), "foo")
	})

	nameStorage.Emplace(entity1, "bar")
	Equal(t, singleView.Pool.Has(entity0), true)
	Equal(t, singleView.Pool.Has(entity1), true)
	Equal(t, singleView.Size(), 2)

	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		checkComDataFromMapData(datas, nameID, "foo", "bar")
	})

	nameStorage.Destroy(entity0)

	Equal(t, singleView.Pool.Has(entity0), false)
	Equal(t, singleView.Pool.Has(entity1), true)
	Equal(t, singleView.Size(), 1)
	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, entity, entity1)
		Equal(t, datas[nameID].(string), "bar")
	})
}
