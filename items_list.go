package main

type ItemsList struct {
	items []*Item

	onAdd    func(item *Item)
	onDelete func(item *Item)
	onChange func()
}

func NewItemsList() *ItemsList {
	items := &ItemsList{items: make([]*Item, 0)}
	items.OnAdd(func(item *Item) {})
	items.OnDelete(func(item *Item) {})
	items.OnChange(func() {})

	return items
}

func (l *ItemsList) AllUntyped() *[]any {
	bindItems := make([]any, len(l.All()))
	for i, item := range l.All() {
		bindItems[i] = item
	}

	return &bindItems
}

func (l *ItemsList) OnAdd(onAdd func(item *Item)) {
	l.onAdd = func(i *Item) {
		onAdd(i)
		l.onChange()
	}
}

// OnDelete note: provided method is called before the actual deletion of the item.
func (l *ItemsList) OnDelete(onDelete func(item *Item)) {
	// We do not wrap onChange because item is not actually deleted when l.onDelete called
	l.onDelete = onDelete
}

func (l *ItemsList) OnChange(onChange func()) {
	l.onChange = onChange
}

func (l *ItemsList) All() []*Item {
	return l.items
}

func (l *ItemsList) Add(new *Item) int {
	l.items = append(l.items, new)
	l.onAdd(new)

	return len(l.items) - 1
}

func (l *ItemsList) Get(itm *Item) (*Item, int) {
	for i, item := range l.items {
		if item == itm {
			return l.items[i], i
		}
	}

	return nil, 0
}

func (l *ItemsList) UpdateItem() {
	l.onChange()
}

func (l *ItemsList) RemoveItem(del *Item) {
	for i, item := range l.items {
		if item == del {
			l.remove(i)
		}
	}
}

func (l *ItemsList) remove(i int) {
	if len(l.items) < i {
		return
	}

	l.onDelete(l.items[i])
	l.items[i] = nil
	l.onChange()
}
