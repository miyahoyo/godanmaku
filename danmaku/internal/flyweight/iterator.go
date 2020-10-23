package flyweight

import (
	"github.com/yohamta/godanmaku/danmaku/internal/linkedlist"
)

type Iterator struct {
	current *linkedlist.Element
}

func (ite *Iterator) HasNext() bool {
	return ite.current != nil
}

func (ite *Iterator) Next() *Object {
	e := ite.current
	ite.current = e.GetNext()

	return (*Object)(e.GetValue())
}
