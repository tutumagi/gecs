package gecs

// IIterator iterator interface
type IIterator interface {
	Prev() IIterator
	Next() IIterator
	Equal(other IIterator) bool
	Begin() IIterator
	End() IIterator
	Data() interface{}
}

// Each iterator the items
func Each(items IIterator, f func(data interface{})) {
	for first := items.Begin(); !first.Equal(items.End()); {
		f(first.Data())
		first = first.Next()
	}
}

// ReverseEach iterator the items by reverse order
func ReverseEach(items IIterator, f func(data interface{})) {
	for first := items.End(); !first.Equal(items.Begin()); {
		first = first.Prev()
		f(first.Data())
	}
}

// func AllOf(items IIterator, f func(data interface{}) bool) bool {
// 	for first := items.Begin(); !first.Equal(items.End()); {
// 		if f(first.Data()) == false {
// 			return false
// 		}
// 		first = first.Next()
// 	}
// 	return true
// }
