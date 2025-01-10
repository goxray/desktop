/*
Package traylist represents a generic tray menu list. It automatically manages all
the UI elements for you and you can provide a generic value type for the items to hold.
*/
package traylist

import (
	"errors"
	"fmt"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

// defaultIconSet is iconSet set with default fyne icons.
var defaultIconSet = func() *IconSet {
	return &IconSet{
		LogoPassive: theme.MediaPauseIcon(),
		LogoActive:  theme.MediaPlayIcon(),

		NotSelected: nil,
		InProgress:  theme.MoreHorizontalIcon(),
		Selected:    theme.ConfirmIcon(),
		Settings:    theme.SettingsIcon(),
		Warning:     theme.WarningIcon(),
	}
}

var (
	ErrItemNotFound = errors.New("item not found")
)

type IconSet struct {
	LogoPassive fyne.Resource
	LogoActive  fyne.Resource
	NotSelected fyne.Resource
	InProgress  fyne.Resource
	Selected    fyne.Resource
	Settings    fyne.Resource
	Warning     fyne.Resource
}

type value interface {
	Label() string
	SetActive(bool)
	comparable
}

type List[T value] struct {
	menu    *Menu[T]
	nextID  atomic.Int64 // Generates external session-persistent IDs for items.
	items   map[int]*trayItem[T]
	onClick func(int) error

	itemsStartIDx int
	footerLen     int
	desk          desktop.App
	iconSet       IconSet
}

func NewDefault[T value](title string, desk desktop.App, icons *IconSet) *List[T] {
	if icons == nil {
		icons = defaultIconSet()
	}

	if icons.LogoActive == nil || icons.LogoPassive == nil {
		panic(fmt.Errorf("icons.LogoActive and icons.LogoPassive must not be nil"))
	}

	header := []*fyne.MenuItem{
		{Label: title, Disabled: true},
		fyne.NewMenuItemSeparator(),
	}
	// New items will be placed in between.
	footer := []*fyne.MenuItem{
		fyne.NewMenuItemSeparator(),
		{
			Label:  lang.L("Configuration"),
			Icon:   icons.Settings,
			Action: func() {},
		},
		fyne.NewMenuItemSeparator(),
		{
			Label:  lang.L("Quit"),
			IsQuit: true,
		},
	}

	mb := fyne.NewMenu(title, append(header, footer...)...)

	return New[T](desk, mb, len(header), len(footer), icons)
}

func New[T value](desk desktop.App, menu *fyne.Menu, insertIDx int, footerLen int, icons *IconSet) *List[T] {
	if icons == nil {
		icons = defaultIconSet()
	}

	if icons.LogoActive == nil || icons.LogoPassive == nil {
		panic(fmt.Errorf("icons.LogoActive and icons.LogoPassive must not be nil"))
	}

	menuBar := &List[T]{
		menu:          &Menu[T]{menu: menu, footerLen: footerLen},
		items:         make(map[int]*trayItem[T]),
		onClick:       func(i int) error { return nil },
		desk:          desk,
		itemsStartIDx: insertIDx + 1,
		footerLen:     footerLen,
		iconSet:       *icons,
	}

	return menuBar
}

func (mb *List[T]) OnSettingsClick(f func()) {
	mb.menu.OnSettingsClick(f)
}

func (mb *List[T]) OnItemClick(f func(int) error) {
	mb.onClick = f
}

func (mb *List[T]) Add(data T) int {
	defer mb.updateValues()
	newID := int(mb.nextID.Add(1))
	item := newTrayItem[T](data, mb.desk, mb.iconSet)

	mb.menu.Insert(item)
	mb.items[newID] = item

	item.menuItem.Action = func() {
		curID := -1
		for k, v := range mb.items {
			if v == item {
				curID = k
			}
		}
		if curID == -1 {
			return
		}

		mb.disableAll(true)
		defer mb.disableAll(false)
		item.setInProgress()

		if err := mb.onClick(curID); err != nil {
			defer item.setWarning()
			mb.setLabel(err.Error())

			if active := mb.getActive(); active != nil {
				active.setActive(false)
			}

			return
		}

		// Render all active items as not active except ours
		for id, itm := range mb.items {
			if curID != id && itm.isActive() {
				itm.setActive(false)
			}
		}

		item.toggle()
	}

	return newID
}

func (mb *List[T]) Remove(i T) error {
	var zero T
	if i == zero {
		return nil
	}
	defer mb.updateValues()

	for id, itm := range mb.items {
		if itm.Value() == i {
			if itm.isActive() {
				return fmt.Errorf("item %d is active", id)
			}

			delete(mb.items, id)
			mb.menu.RemoveItem(itm.menuItem)

			return nil
		}
	}

	return ErrItemNotFound
}

func (mb *List[T]) Swap(i1, i2 T) error {
	var zero T
	if i1 == zero || i2 == zero {
		return nil
	}
	defer mb.Refresh()

	id1, id2 := -1, -1
	for id, itm := range mb.items {
		if itm.Value() != i1 && itm.Value() != i2 {
			continue
		}

		if itm.Value() == i1 {
			id1 = id
		}
		if itm.Value() == i2 {
			id2 = id
		}
	}

	mb.menu.Swap(mb.items[id2].menuItem, mb.items[id1].menuItem)
	mb.items[id1], mb.items[id2] = mb.items[id2], mb.items[id1]

	return nil
}

func (mb *List[T]) Refresh() {
	mb.updateValues()
}

func (mb *List[T]) Get(id int) T {
	itm := mb.getItem(id)
	if itm == nil {
		var zero T
		return zero
	}

	return itm.Value()
}

func (mb *List[T]) getItem(id int) *trayItem[T] {
	if mb.items[id] == nil {
		return nil
	}

	return mb.items[id]
}

func (mb *List[T]) IsActive(id int) bool {
	active := mb.getActive()
	item := mb.items[id]
	if active == nil || item == nil {
		return false
	}

	return active == item
}

func (mb *List[T]) HasActive() bool {
	for _, itm := range mb.items {
		if itm.isActive() {
			return true
		}
	}

	return false
}

func (mb *List[T]) GetActive() T {
	if active := mb.getActive(); active != nil {
		return active.Value()
	}

	var zero T
	return zero
}

func (mb *List[T]) Show() {
	mb.desk.SetSystemTrayIcon(mb.iconSet.LogoPassive)
	mb.desk.SetSystemTrayMenu(mb.menu.Menu())
}

func (mb *List[T]) setLabel(label string) {
	mb.menu.SetTitle(label)
}

func (mb *List[T]) disableAll(disable bool) {
	for _, itm := range mb.items {
		itm.menuItem.Disabled = disable
	}
	mb.menu.Refresh()
}

func (mb *List[T]) updateValues() {
	for _, itm := range mb.items {
		itm.menuItem.Label = itm.Value().Label()
	}
	mb.menu.Refresh()
}

func (mb *List[T]) getActive() *trayItem[T] {
	for _, itm := range mb.items {
		if itm.isActive() {
			return itm
		}
	}

	return nil
}
