package connlist

import (
	"errors"
)

// Collection represents a collection of items.
// Is used to easily pass events and update the UI state in one place (on{*} methods).
type Collection struct {
	items []*Item

	onAdd    func(*Item)
	onDelete func(*Item)
	onSwap   func(*Item, *Item)
	onChange func()
}

func New() *Collection {
	items := &Collection{items: make([]*Item, 0)}
	items.OnAdd(func(item *Item) {})
	items.OnDelete(func(item *Item) {})
	items.OnChange(func() {})

	return items
}

func (l *Collection) AllUntyped() *[]any {
	bindItems := make([]any, len(l.All()))
	for i, item := range l.All() {
		bindItems[i] = item
	}

	return &bindItems
}

func (l *Collection) OnAdd(onAdd func(item *Item)) {
	l.onAdd = func(i *Item) {
		onAdd(i)
		l.onChange()
	}
}

func (l *Collection) OnSwap(onSwap func(*Item, *Item)) {
	l.onSwap = func(i1 *Item, i2 *Item) {
		onSwap(i1, i2)
		l.onChange()
	}
}

// OnDelete note: provided method is called before the actual deletion of the item.
func (l *Collection) OnDelete(onDelete func(item *Item)) {
	// We do not wrap onChange because item is not actually deleted when l.onDelete called
	l.onDelete = onDelete
}

func (l *Collection) OnChange(onChange func()) {
	l.onChange = onChange
}

func (l *Collection) All() []*Item {
	res := make([]*Item, 0, len(l.items))
	for _, item := range l.items {
		if item == nil {
			continue
		}
		res = append(res, item)
	}

	return res
}

func (l *Collection) AddItem(label, link string) error {
	item, err := newItem(label, link, l)
	if err != nil {
		return err
	}

	l.items = append(l.items, item)
	l.onAdd(item)

	return nil
}

func (l *Collection) RemoveItem(del *Item) {
	for i, item := range l.items {
		if item == del {
			l.remove(i)
		}
	}
}

func (l *Collection) SwapItems(itm1 *Item, itm2 *Item) error {
	id1, id2 := -1, -1
	for i, item := range l.items {
		if item == itm1 {
			id1 = i
		}
		if item == itm2 {
			id2 = i
		}
	}
	if id1 == -1 || id2 == -1 {
		return errors.New("cannot swap items")
	}

	l.items[id1], l.items[id2] = l.items[id2], l.items[id1]
	l.onSwap(itm1, itm2)

	return nil
}

func (l *Collection) remove(i int) {
	if len(l.items) < i {
		return
	}

	l.onDelete(l.items[i])
	l.items[i] = nil
	l.onChange()
}
