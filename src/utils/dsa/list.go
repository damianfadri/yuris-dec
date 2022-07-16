package dsa

import (
	"sort"
)

type Comparator func(int, int) bool

type List[T any] struct {
	Items		[]T
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

func (l *List[T]) Count() int {
	return len(l.Items)
}

func (l *List[T]) Add(item T) {
	l.Items = append(l.Items, item)
}

func (l *List[T]) Remove(i int) {
	if i != l.Count() {
		l.Items = append(l.Items[:i], l.Items[i+1:]...)
	}
}

func (l *List[T]) Sort(fn Comparator) {
	sort.Slice(l.Items, fn)
}

func (l *List[T]) Reverse() {
	n := l.Count()
	for i := 0; i < n / 2; i++ {
		l.Items[i], l.Items[n - 1 - i] = l.Items[n - 1 - i], l.Items[i]
	}
}