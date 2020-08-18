package gecs

import (
	"testing"
)

// IntIterator for int type
type IntIterator struct {
	datas []int
	pos   int
}

// Equal other
func (i IntIterator) Equal(other IIterator) bool {
	tOther := other.(*IntIterator)
	return i.pos == tOther.pos && &i.datas[0] == &tOther.datas[0]
}

// Next iterator
func (i *IntIterator) Next() IIterator {
	i.pos++
	return IIterator(i)
}

// Prev iterator
func (i *IntIterator) Prev() IIterator {
	i.pos--
	return IIterator(i)
}

// Begin iterator
func (i IntIterator) Begin() IIterator {
	i.pos = 0
	return IIterator(&i)
}

// End iterator
func (i IntIterator) End() IIterator {
	i.pos = len(i.datas)
	return IIterator(&i)
}

// Data with the iterator value
func (i IntIterator) Data() interface{} {
	return i.datas[i.pos]
}

func TestIterator(t *testing.T) {
	items1 := []int{
		1, 2, 3, 4, 5,
	}
	items3 := items1[:]
	items2 := []int{
		1, 2, 3, 4, 5,
	}

	if &items1[0] == &items3[0] {
		t.Log("items1 equal items3")
	}
	if &items2 == &items3 {
		t.Log("items2 equal items3")
	}

	Each(&IntIterator{
		datas: items1,
		pos:   0,
	}, func(data interface{}) {
		t.Logf("item is %v", data)
	})
	ReverseEach(&IntIterator{
		datas: items1,
		pos:   0,
	}, func(data interface{}) {
		t.Logf("item is %v", data)
	})
}

// type People struct {
// 	Name string
// }

// func TestModify(t *testing.T) {
// 	datas := []*People{
// 		{
// 			Name: "1",
// 		}, {
// 			Name: "2",
// 		}, {
// 			Name: "3",
// 		}, {
// 			Name: "4",
// 		}}
// 	fn := func(i int) *People {
// 		return datas[i]
// 	}
// 	t.Logf("datas: %+v", datas)
// 	modifyFn := func(i int) {
// 		a := fn(i)
// 		t.Logf("index datas[%v]=%+v : %+v", i, datas[i], a)
// 		b := &a
// 		*b = &People{Name: "100"}
// 		t.Logf("index datas[%v]=%+v : %+v", i, datas[i], a)
// 	}

// 	modifyFn(1)
// }
