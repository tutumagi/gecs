package gecs

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

type Count struct {
	count int
}

func newEntityID(init int) func() EntityID {
	return func() EntityID {
		init++
		return EntityID(init)
	}
}

func Test_View(t *testing.T) {
	newEntity := newEntityID(0)

	nameID := ComponentID(1)
	nameStorage := NewStorage(nameID)
	ageID := ComponentID(2)
	ageStorage := NewStorage(ageID)
	countID := ComponentID(3)
	countStorage := NewStorage(countID)
	timeID := ComponentID(4)
	timeStorage := NewStorage(timeID)

	singleView := newView(nameStorage)
	multiView := newView(nameStorage, ageStorage, countStorage)
	nameNotAgeView := newView(nameStorage).withExclude(ageStorage)
	ageTimeNotCountNameView := newView(ageStorage, timeStorage).withExclude(countStorage, nameStorage)

	Equal(t, singleView.empty(), true)
	Equal(t, multiView.empty(), true)
	Equal(t, nameNotAgeView.empty(), true)
	Equal(t, ageTimeNotCountNameView.empty(), true)

	nameEntity := newEntity()
	nameStorage.Emplace(nameEntity, "nameEntity")
	emptyEntity := newEntity()

	Equal(t, singleView.contains(nameEntity), true)
	Equal(t, singleView.contains(emptyEntity), false)
	Equal(t, singleView.size(), 1)
	Equal(t, singleView.Get(nameEntity, nameID).(string), "nameEntity")
	PanicMatches(t, func() { singleView.Get(emptyEntity, nameID) }, "view should have entity, but not")

	Equal(t, multiView.contains(nameEntity), false)
	equalCountView(t, multiView, 0)
	PanicMatches(t, func() { multiView.Get(nameEntity, nameID) }, "view should have entity, but not")

	nameAgeEntity := newEntity()
	nameStorage.Emplace(nameAgeEntity, "nameAgeEntity")

	Equal(t, singleView.contains(nameEntity), true)
	Equal(t, singleView.contains(nameAgeEntity), true)
	Equal(t, singleView.size(), 2)
	Equal(t, singleView.Get(nameEntity, nameID).(string), "nameEntity")
	Equal(t, singleView.Get(nameAgeEntity, nameID).(string), "nameAgeEntity")

	equalCountView(t, multiView, 0)
	PanicMatches(t, func() { multiView.Get(nameEntity, nameID) }, "view should have entity, but not")

	nameStorage.Destroy(nameEntity)

	Equal(t, singleView.contains(nameEntity), false)
	Equal(t, singleView.contains(nameAgeEntity), true)
	Equal(t, singleView.size(), 1)
	PanicMatches(t, func() { singleView.Get(nameEntity, nameID) }, "view should have entity, but not")
	Equal(t, singleView.Get(nameAgeEntity, nameID).(string), "nameAgeEntity")

	equalCountView(t, multiView, 0)

	nameAgeCountEntity := newEntity()
	nameStorage.Emplace(nameAgeCountEntity, "nameAgeCountEntity")
	ageStorage.Emplace(nameAgeCountEntity, 18)
	countStorage.Emplace(nameAgeCountEntity, &Count{count: 100})

	equalCountView(t, multiView, 1)
	Equal(t, multiView.contains(nameAgeCountEntity), true)
	Equal(t, multiView.Get(nameAgeCountEntity, nameID).(string), "nameAgeCountEntity")
	Equal(t, multiView.Get(nameAgeCountEntity, ageID).(int), 18)
	Equal(t, multiView.Get(nameAgeCountEntity, countID).(*Count), &Count{count: 100})
	Equal(t, multiView.GetMulti(nameAgeCountEntity, nameID, ageID, countID), map[ComponentID]interface{}{
		nameID:  "nameAgeCountEntity",
		ageID:   18,
		countID: &Count{count: 100},
	})
	Equal(t, multiView.GetMulti(nameAgeCountEntity, nameID, ageID), map[ComponentID]interface{}{
		nameID: "nameAgeCountEntity",
		ageID:  18,
	})
	Equal(t, multiView.GetMulti(nameAgeCountEntity, countID, ageID), map[ComponentID]interface{}{
		ageID:   18,
		countID: &Count{count: 100},
	})

	nameAgeCountEntity2 := newEntity()
	nameStorage.Emplace(nameAgeCountEntity2, "nameAgeCountEntity2")
	ageStorage.Emplace(nameAgeCountEntity2, 22)
	countStorage.Emplace(nameAgeCountEntity2, &Count{count: 85})

	equalCountView(t, multiView, 2)
	Equal(t, multiView.contains(nameAgeCountEntity2), true)
	Equal(t, multiView.Get(nameAgeCountEntity2, nameID).(string), "nameAgeCountEntity2")
	Equal(t, multiView.Get(nameAgeCountEntity2, ageID).(int), 22)
	Equal(t, multiView.Get(nameAgeCountEntity2, countID).(*Count), &Count{count: 85})
	Equal(t, multiView.GetMulti(nameAgeCountEntity2, nameID, ageID, countID), map[ComponentID]interface{}{
		nameID:  "nameAgeCountEntity2",
		ageID:   22,
		countID: &Count{count: 85},
	})
	Equal(t, multiView.GetMulti(nameAgeCountEntity2, nameID, ageID), map[ComponentID]interface{}{
		nameID: "nameAgeCountEntity2",
		ageID:  22,
	})
	Equal(t, multiView.GetMulti(nameAgeCountEntity2, countID, ageID), map[ComponentID]interface{}{
		ageID:   22,
		countID: &Count{count: 85},
	})

	ageStorage.Destroy(nameAgeCountEntity2)
	Equal(t, multiView.size(), 1)
	Equal(t, multiView.contains(nameAgeCountEntity2), false)
	equalCountView(t, multiView, 1)
	PanicMatches(t, func() { multiView.GetMulti(nameAgeCountEntity2, nameID, ageID, countID) }, "view should have entity, but not")

	nameStorage.Destroy(nameAgeCountEntity)
	equalCountView(t, multiView, 0)
}

func equalCountView(t *testing.T, view *View, count int) {
	cal := 0
	view.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		cal++
	})
	EqualSkip(t, 2, cal, count)
}
