package gecs

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

	singleView := newView(nameStorage)
	multiView := newView(nameStorage, ageStorage, countStorage)
	Equal(t, singleView.empty(), true)
	Equal(t, multiView.empty(), true)

	entity0 := EntityID(0)
	nameStorage.Emplace(entity0, "entity0")
	entity1 := EntityID(1)

	Equal(t, singleView.contains(entity0), true)
	Equal(t, singleView.contains(entity1), false)
	Equal(t, singleView.size(), 1)
	Equal(t, singleView.Get(entity0, nameID).(string), "entity0")
	PanicMatches(t, func() { singleView.Get(entity1, nameID) }, "view should have entity, but not")

	Equal(t, multiView.contains(entity0), false)
	equalCountView(t, multiView, 0)
	PanicMatches(t, func() { multiView.Get(entity0, nameID) }, "view should have entity, but not")

	nameStorage.Emplace(entity1, "entity1")

	Equal(t, singleView.contains(entity0), true)
	Equal(t, singleView.contains(entity1), true)
	Equal(t, singleView.size(), 2)
	Equal(t, singleView.Get(entity0, nameID).(string), "entity0")
	Equal(t, singleView.Get(entity1, nameID).(string), "entity1")

	equalCountView(t, multiView, 0)
	PanicMatches(t, func() { multiView.Get(entity0, nameID) }, "view should have entity, but not")

	nameStorage.Destroy(entity0)

	Equal(t, singleView.contains(entity0), false)
	Equal(t, singleView.contains(entity1), true)
	Equal(t, singleView.size(), 1)
	PanicMatches(t, func() { singleView.Get(entity0, nameID) }, "view should have entity, but not")
	Equal(t, singleView.Get(entity1, nameID).(string), "entity1")

	equalCountView(t, multiView, 0)

	entity2 := EntityID(2)
	nameStorage.Emplace(entity2, "entity2")
	ageStorage.Emplace(entity2, 18)
	countStorage.Emplace(entity2, &Count{count: 100})

	equalCountView(t, multiView, 1)
	Equal(t, multiView.contains(entity2), true)
	Equal(t, multiView.Get(entity2, nameID).(string), "entity2")
	Equal(t, multiView.Get(entity2, ageID).(int), 18)
	Equal(t, multiView.Get(entity2, countID).(*Count), &Count{count: 100})
	Equal(t, multiView.GetMulti(entity2, nameID, ageID, countID), map[ComponentID]interface{}{
		nameID:  "entity2",
		ageID:   18,
		countID: &Count{count: 100},
	})
	Equal(t, multiView.GetMulti(entity2, nameID, ageID), map[ComponentID]interface{}{
		nameID: "entity2",
		ageID:  18,
	})
	Equal(t, multiView.GetMulti(entity2, countID, ageID), map[ComponentID]interface{}{
		ageID:   18,
		countID: &Count{count: 100},
	})

	entity3 := EntityID(3)
	nameStorage.Emplace(entity3, "entity3")
	ageStorage.Emplace(entity3, 22)
	countStorage.Emplace(entity3, &Count{count: 85})

	equalCountView(t, multiView, 2)
	Equal(t, multiView.contains(entity3), true)
	Equal(t, multiView.Get(entity3, nameID).(string), "entity3")
	Equal(t, multiView.Get(entity3, ageID).(int), 22)
	Equal(t, multiView.Get(entity3, countID).(*Count), &Count{count: 85})
	Equal(t, multiView.GetMulti(entity3, nameID, ageID, countID), map[ComponentID]interface{}{
		nameID:  "entity3",
		ageID:   22,
		countID: &Count{count: 85},
	})
	Equal(t, multiView.GetMulti(entity3, nameID, ageID), map[ComponentID]interface{}{
		nameID: "entity3",
		ageID:  22,
	})
	Equal(t, multiView.GetMulti(entity3, countID, ageID), map[ComponentID]interface{}{
		ageID:   22,
		countID: &Count{count: 85},
	})

	ageStorage.Destroy(entity3)
	Equal(t, multiView.size(), 1)
	Equal(t, multiView.contains(entity3), false)
	equalCountView(t, multiView, 1)
	PanicMatches(t, func() { multiView.GetMulti(entity3, nameID, ageID, countID) }, "view should have entity, but not")

	nameStorage.Destroy(entity2)
	equalCountView(t, multiView, 0)
}

func equalCountView(t *testing.T, view *View, count int) {
	cal := 0
	view.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		cal++
	})
	EqualSkip(t, 2, cal, count)
}
