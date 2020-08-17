package entt

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

type Count struct {
	count int
}

func Test_View(t *testing.T) {
	nameID := ComponentID(1)
	nameStorage := NewStorage(nameID)
	ageID := ComponentID(2)
	ageStorage := NewStorage(ageID)
	countID := ComponentID(3)
	countStorage := NewStorage(countID)

	singleView := NewView(nameStorage)
	multiView := NewView(nameStorage, ageStorage, countStorage)

	entity0 := EntityID(0)
	nameStorage.Emplace(entity0, "foo")

	entity1 := EntityID(1)
	Equal(t, singleView.contains(entity0), true)
	Equal(t, singleView.contains(entity1), false)
	Equal(t, singleView.size(), 1)
	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, entity, entity0)
		Equal(t, datas[nameID].(string), "foo")
	})
	Equal(t, multiView.contains(entity0), false)
	Equal(t, multiView.size(), 0)

	nameStorage.Emplace(entity1, "bar")
	Equal(t, singleView.contains(entity0), true)
	Equal(t, singleView.contains(entity1), true)
	Equal(t, singleView.size(), 2)

	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		checkComDataFromMapData(datas, nameID, "foo", "bar")
	})

	nameStorage.Destroy(entity0)

	Equal(t, singleView.contains(entity0), false)
	Equal(t, singleView.contains(entity1), true)
	Equal(t, singleView.size(), 1)
	singleView.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		Equal(t, entity, entity1)
		Equal(t, datas[nameID].(string), "bar")
	})

	entity2 := EntityID(2)
	nameStorage.Emplace(entity2, "entity2")
	ageStorage.Emplace(entity2, 18)
	countStorage.Emplace(entity2, &Count{count: 100})

	Equal(t, multiView.contains(entity2), true)
	equalCountView(t, multiView, 1)

	entity3 := EntityID(3)
	nameStorage.Emplace(entity3, "entity3")
	ageStorage.Emplace(entity3, 22)
	countStorage.Emplace(entity3, &Count{count: 85})

	Equal(t, multiView.contains(entity3), true)
	equalCountView(t, multiView, 2)
}

func equalCountView(t *testing.T, view IteratorView, count int) {
	cal := 0
	view.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		cal++
	})
	EqualSkip(t, 2, cal, count)
}
