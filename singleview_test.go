package gecs

import (
	"fmt"
	"testing"

	. "github.com/go-playground/assert/v2"
)

func Test_SingleView(t *testing.T) {
	nameID := ComponentID(1)
	nameStorage := NewStorage(nameID)

	singleView := newSingleView(nameStorage)

	Equal(t, singleView.Empty(), true)

	entity0 := EntityID(0)
	nameStorage.Emplace(entity0, "entity0")
	Equal(t, singleView.Empty(), false)
	Equal(t, singleView.Get(entity0).(string), "entity0")
	raw := singleView.Raw()
	Equal(t, len(raw), 1)
	Equal(t, raw[0].(string), "entity0")
	data := singleView.Data()
	Equal(t, len(data), 1)
	Equal(t, data[0], entity0)

	entity1 := EntityID(1)
	Equal(t, singleView.Pool.Has(entity0), true)
	Equal(t, singleView.Get(entity0).(string), "entity0")
	Equal(t, singleView.Pool.Has(entity1), false)
	PanicMatches(t, func() { singleView.Get(entity1) }, fmt.Sprintf("entity(%v)should have the component", entity1))

	Equal(t, singleView.Size(), 1)
	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, entity, entity0)
		Equal(t, datas[nameID].(string), "entity0")
	})

	nameStorage.Emplace(entity1, "entity1")
	Equal(t, singleView.Pool.Has(entity0), true)
	Equal(t, singleView.Pool.Has(entity1), true)
	Equal(t, singleView.Size(), 2)
	raw = singleView.Raw()
	Equal(t, len(raw), 2)
	Equal(t, raw[0].(string), "entity0")
	Equal(t, raw[1].(string), "entity1")
	data = singleView.Data()
	Equal(t, len(data), 2)
	Equal(t, data[0], entity0)
	Equal(t, data[1], entity1)

	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		checkComDataFromMapData(datas, nameID, "entity0", "entity1")
	})

	nameStorage.Destroy(entity0)

	Equal(t, singleView.Pool.Has(entity0), false)
	Equal(t, singleView.Pool.Has(entity1), true)
	Equal(t, singleView.Size(), 1)
	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, entity, entity1)
		Equal(t, datas[nameID].(string), "entity1")
	})
	raw = singleView.Raw()
	Equal(t, len(raw), 1)
	Equal(t, raw[0].(string), "entity1")
	data = singleView.Data()
	Equal(t, len(data), 1)
	Equal(t, data[0], entity1)

	nameStorage.Destroy(entity1)

	Equal(t, singleView.Empty(), true)
	raw = singleView.Raw()
	Equal(t, len(raw), 0)
	data = singleView.Data()
	Equal(t, len(data), 0)
}
