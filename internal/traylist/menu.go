package traylist

import (
	"slices"

	"fyne.io/fyne/v2"
)

// Menu is a simple wrapper for fyne menu, allowing to easily delete and add items.
type Menu[T value] struct {
	menu *fyne.Menu
}

func (m *Menu[T]) Menu() *fyne.Menu {
	return m.menu
}

func (m *Menu[T]) Insert(new *trayItem[T]) {
	m.menu.Items = slices.Insert(m.menu.Items, 2, new.menuItem)
}

func (m *Menu[T]) RemoveItem(itm *fyne.MenuItem) {
	for i, it := range m.menu.Items {
		if it == itm {
			m.menu.Items = slices.Delete(m.menu.Items, i, i+1)
		}
	}
}

func (m *Menu[T]) SetTitle(title string) {
	m.menu.Items[0].Label = title
}

func (m *Menu[T]) OnSettingsClick(f func()) {
	m.menu.Items[len(m.menu.Items)-1].Action = f
}

func (m *Menu[T]) Refresh() {
	m.menu.Refresh()
}
