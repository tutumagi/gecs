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
	intSliceID := ComponentID(4)
	intSliceStorage := NewStorage(intSliceID)

	nameView := newView(nameStorage)
	nameAgeCountView := newView(nameStorage, ageStorage, countStorage)
	nameNotAgeView := newView(nameStorage).withExclude(ageStorage)
	ageIntSliceNotCountNameView := newView(ageStorage, intSliceStorage).withExclude(countStorage, nameStorage)

	t.Run("empty", func(t *testing.T) {
		Equal(t, nameView.empty(), true)
		Equal(t, nameAgeCountView.empty(), true)
		Equal(t, nameNotAgeView.empty(), true)
		Equal(t, ageIntSliceNotCountNameView.empty(), true)

		Equal(t, nameView.empty(nameID), true)
		Equal(t, nameAgeCountView.empty(nameID, ageID), true)
		Equal(t, nameNotAgeView.empty(nameID), true)
		Equal(t, ageIntSliceNotCountNameView.empty(intSliceID), true)
	})

	nameEntity := newEntity()
	nameStorage.Emplace(nameEntity, "nameEntity")

	emptyEntity := newEntity()

	nameAgeEntity := newEntity()
	nameStorage.Emplace(nameAgeEntity, "nameAgeEntity")
	ageStorage.Emplace(nameAgeEntity, 29)

	nameAgeCountEntity2 := newEntity()
	nameStorage.Emplace(nameAgeCountEntity2, "nameAgeCountEntity2")
	ageStorage.Emplace(nameAgeCountEntity2, 22)
	countStorage.Emplace(nameAgeCountEntity2, &Count{count: 85})

	ageEntity := newEntity()
	ageStorage.Emplace(ageEntity, 158)

	ageIntSliceEntity := newEntity()
	ageStorage.Emplace(ageIntSliceEntity, 89)
	intSliceStorage.Emplace(ageIntSliceEntity, []int{1, 2, 3})

	nameAgeCountEntity := newEntity()
	nameStorage.Emplace(nameAgeCountEntity, "nameAgeCountEntity")
	ageStorage.Emplace(nameAgeCountEntity, 18)
	countStorage.Emplace(nameAgeCountEntity, &Count{count: 100})

	t.Run("contains", func(t *testing.T) {
		Equal(t, nameView.contains(nameEntity), true)
		Equal(t, nameView.contains(emptyEntity), false)
		Equal(t, nameView.contains(nameAgeEntity), true)
		Equal(t, nameView.contains(nameAgeCountEntity), true)
		Equal(t, nameView.contains(nameAgeCountEntity2), true)
		Equal(t, nameView.contains(ageIntSliceEntity), false)
		Equal(t, nameView.contains(ageEntity), false)

		Equal(t, nameAgeCountView.contains(nameEntity), false)
		Equal(t, nameAgeCountView.contains(emptyEntity), false)
		Equal(t, nameAgeCountView.contains(nameAgeEntity), false)
		Equal(t, nameAgeCountView.contains(nameAgeCountEntity), true)
		Equal(t, nameAgeCountView.contains(nameAgeCountEntity2), true)
		Equal(t, nameAgeCountView.contains(ageIntSliceEntity), false)
		Equal(t, nameAgeCountView.contains(ageEntity), false)

		Equal(t, nameNotAgeView.contains(nameEntity), true)
		Equal(t, nameNotAgeView.contains(emptyEntity), false)
		Equal(t, nameNotAgeView.contains(nameAgeEntity), false)
		Equal(t, nameNotAgeView.contains(nameAgeCountEntity), false)
		Equal(t, nameNotAgeView.contains(nameAgeCountEntity2), false)
		Equal(t, nameNotAgeView.contains(ageIntSliceEntity), false)
		Equal(t, nameNotAgeView.contains(ageEntity), false)

		Equal(t, ageIntSliceNotCountNameView.contains(nameEntity), false)
		Equal(t, ageIntSliceNotCountNameView.contains(emptyEntity), false)
		Equal(t, ageIntSliceNotCountNameView.contains(nameAgeCountEntity), false)
		Equal(t, ageIntSliceNotCountNameView.contains(nameAgeCountEntity2), false)
		Equal(t, ageIntSliceNotCountNameView.contains(ageIntSliceEntity), true)
		Equal(t, ageIntSliceNotCountNameView.contains(ageEntity), false)
	})

	t.Run("size", func(t *testing.T) {
		Equal(t, nameView.size() > 0, true)

		Equal(t, nameAgeCountView.size() > 0, true)

		Equal(t, nameNotAgeView.size() > 0, true)

		Equal(t, ageIntSliceNotCountNameView.size() > 0, true)
	})

	t.Run("normalGet", func(t *testing.T) {
		Equal(t, nameView.Get(nameEntity, nameID).(string), "nameEntity")
		Equal(t, nameView.Get(nameAgeCountEntity, nameID).(string), "nameAgeCountEntity")
		Equal(t, nameView.Get(nameAgeCountEntity2, nameID).(string), "nameAgeCountEntity2")

		Equal(t, nameAgeCountView.Get(nameAgeCountEntity, nameID).(string), "nameAgeCountEntity")
		Equal(t, nameAgeCountView.Get(nameAgeCountEntity, ageID).(int), 18)
		Equal(t, nameAgeCountView.Get(nameAgeCountEntity, countID).(*Count), &Count{count: 100})
		Equal(t, nameAgeCountView.Get(nameAgeCountEntity2, nameID).(string), "nameAgeCountEntity2")
		Equal(t, nameAgeCountView.Get(nameAgeCountEntity2, ageID).(int), 22)
		Equal(t, nameAgeCountView.Get(nameAgeCountEntity2, countID).(*Count), &Count{count: 85})

		Equal(t, nameNotAgeView.Get(nameEntity, nameID).(string), "nameEntity")

		Equal(t, ageIntSliceNotCountNameView.Get(ageIntSliceEntity, ageID).(int), 89)
		Equal(t, ageIntSliceNotCountNameView.Get(ageIntSliceEntity, intSliceID).([]int), []int{1, 2, 3})
	})

	t.Run("panicGet", func(t *testing.T) {
		PanicMatches(t, func() { nameView.Get(ageIntSliceEntity, nameID) }, "view should have entity, but not")
		PanicMatches(t, func() { nameView.Get(ageEntity, nameID) }, "view should have entity, but not")
		PanicMatches(t, func() { nameView.Get(emptyEntity, nameID) }, "view should have entity, but not")

		PanicMatches(t, func() { nameAgeCountView.Get(emptyEntity, ageID) }, "view should have entity, but not")
		PanicMatches(t, func() { nameAgeCountView.Get(nameEntity, nameID) }, "view should have entity, but not")
		PanicMatches(t, func() { nameAgeCountView.Get(ageIntSliceEntity, ageID) }, "view should have entity, but not")
		PanicMatches(t, func() { nameAgeCountView.Get(ageEntity, countID) }, "view should have entity, but not")

		PanicMatches(t, func() { nameNotAgeView.Get(nameAgeCountEntity2, nameID) }, "view should have entity, but not")

		PanicMatches(t, func() { ageIntSliceNotCountNameView.Get(nameAgeCountEntity2, ageID) }, "view should have entity, but not")
		PanicMatches(t, func() { ageIntSliceNotCountNameView.Get(ageEntity, intSliceID) }, "view should have entity, but not")
	})

	t.Run("eachCount", func(t *testing.T) {
		equalCountView(t, nameView, 4)
		equalCountView(t, nameAgeCountView, 2)
		equalCountView(t, nameNotAgeView, 1)
		equalCountView(t, ageIntSliceNotCountNameView, 1)
	})

	t.Run("multiGet", func(t *testing.T) {
		Equal(t, nameView.GetMulti(nameEntity, nameID), map[ComponentID]interface{}{nameID: "nameEntity"})
		Equal(t, nameView.GetMulti(nameAgeCountEntity, nameID), map[ComponentID]interface{}{nameID: "nameAgeCountEntity"})
		Equal(t, nameView.GetMulti(nameAgeCountEntity2, nameID), map[ComponentID]interface{}{nameID: "nameAgeCountEntity2"})

		Equal(t, nameAgeCountView.GetMulti(nameAgeCountEntity, nameID), map[ComponentID]interface{}{nameID: "nameAgeCountEntity"})
		Equal(t, nameAgeCountView.GetMulti(nameAgeCountEntity, ageID), map[ComponentID]interface{}{ageID: 18})
		Equal(t, nameAgeCountView.GetMulti(nameAgeCountEntity, countID), map[ComponentID]interface{}{countID: &Count{count: 100}})
		Equal(t, nameAgeCountView.GetMulti(nameAgeCountEntity2, nameID, ageID, countID), map[ComponentID]interface{}{
			nameID:  "nameAgeCountEntity2",
			ageID:   22,
			countID: &Count{count: 85},
		})

		Equal(t, nameNotAgeView.GetMulti(nameEntity, nameID), map[ComponentID]interface{}{nameID: "nameEntity"})

		Equal(t, ageIntSliceNotCountNameView.GetMulti(ageIntSliceEntity, ageID),
			map[ComponentID]interface{}{ageID: 89})
		Equal(t, ageIntSliceNotCountNameView.GetMulti(ageIntSliceEntity, intSliceID),
			map[ComponentID]interface{}{intSliceID: []int{1, 2, 3}})
		Equal(t, ageIntSliceNotCountNameView.GetMulti(ageIntSliceEntity, intSliceID, ageID),
			map[ComponentID]interface{}{
				intSliceID: []int{1, 2, 3},
				ageID:      89},
		)

	})

	nameStorage.Destroy(nameEntity)

	ageStorage.Destroy(nameAgeCountEntity2)

	nameStorage.Destroy(nameAgeCountEntity)

}

func equalCountView(t *testing.T, view *View, count int) {
	cal := 0
	view.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		cal++
	})
	EqualSkip(t, 2, cal, count)
}

func equalCountSingleView(t *testing.T, view *SingleView, count int) {
	cal := 0
	view.Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		cal++
	})
	EqualSkip(t, 2, cal, count)
}
