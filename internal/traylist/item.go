package traylist

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// trayItem represents a UI item in the managed list.
type trayItem[T value] struct {
	value    T
	active   bool
	iconSet  IconSet
	desk     desktop.App
	menuItem *fyne.MenuItem
}

func newTrayItem[T value](label string, value T, desk desktop.App, iconSet IconSet) *trayItem[T] {
	return &trayItem[T]{
		value:    value,
		menuItem: &fyne.MenuItem{Label: label, Icon: iconSet.NotSelected},
		desk:     desk,
		iconSet:  iconSet,
	}
}

func (ci *trayItem[T]) Value() T {
	return ci.value
}

func (ci *trayItem[T]) SetValue(new T) {
	ci.value = new
}

func (ci *trayItem[T]) setInProgress() {
	ci.menuItem.Icon = ci.iconSet.InProgress
	ci.desk.SetSystemTrayIcon(ci.iconSet.InProgress)
}

func (ci *trayItem[T]) setWarning() {
	ci.menuItem.Icon = ci.iconSet.Warning
	ci.desk.SetSystemTrayIcon(ci.iconSet.Warning)
}

func (ci *trayItem[T]) setActive(active bool) {
	ci.active = active
	ic := ci.iconSet.Selected
	if !active {
		ic = ci.iconSet.NotSelected
	}
	ci.menuItem.Icon = ic

	if active {
		ci.desk.SetSystemTrayIcon(ci.iconSet.LogoActive)
	} else {
		ci.desk.SetSystemTrayIcon(ci.iconSet.LogoPassive)
	}

	ci.Value().SetActive(ci.active)
}

func (ci *trayItem[T]) toggle() {
	ci.setActive(!ci.isActive())
}

func (ci *trayItem[T]) isActive() bool {
	return ci.active
}
