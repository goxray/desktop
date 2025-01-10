package traylist

import (
	"slices"

	"fyne.io/fyne/v2"
)

// Menu is a simple wrapper for fyne menu, allowing to easily delete and add items.
type Menu[T value] struct {
	menu      *fyne.Menu
	footerLen int
	refresh   func() // alternative refresh method
}

func (m *Menu[T]) Menu() *fyne.Menu {
	return m.menu
}

func (m *Menu[T]) Insert(new *trayItem[T]) {
	m.menu.Items = slices.Insert(m.menu.Items, len(m.menu.Items)-m.footerLen, new.menuItem)
}

func (m *Menu[T]) RemoveItem(itm *fyne.MenuItem) {
	for i, it := range m.menu.Items {
		if it == itm {
			m.menu.Items = slices.Delete(m.menu.Items, i, i+1)
		}
	}
}

func (m *Menu[T]) Swap(i1 *fyne.MenuItem, i2 *fyne.MenuItem) {
	id1, id2 := -1, -1
	for i, it := range m.menu.Items {
		if it == i1 {
			id1 = i
		}
		if it == i2 {
			id2 = i
		}
	}
	if id1 == -1 || id2 == -1 {
		return
	}

	m.menu.Items[id1], m.menu.Items[id2] = m.menu.Items[id2], m.menu.Items[id1]
	m.Refresh()
}

func (m *Menu[T]) SetTitle(title string) {
	m.menu.Items[0].Label = title
}

func (m *Menu[T]) OnSettingsClick(f func()) {
	m.menu.Items[len(m.menu.Items)-(m.footerLen-1)].Action = f
}

func (m *Menu[T]) Refresh() {
	if m.refresh != nil {
		m.refresh()

		return
	}

	m.menu.Refresh()
}
