package traylist

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/require"
)

type deskMock struct {
	onIconSet func(fyne.Resource)
}

func (d deskMock) SetSystemTrayMenu(menu *fyne.Menu) {}

func (d deskMock) SetSystemTrayIcon(icon fyne.Resource) {
	if d.onIconSet != nil {
		d.onIconSet(icon)
	}
}

type mockItem struct {
	l string
	a bool
}

func (m mockItem) Label() string {
	return m.l
}

func (m mockItem) SetActive(b bool) {
	m.a = b
}

func TestTrayList_Operations(t *testing.T) {
	list := setupList(deskMock{})

	baseMenuLen := len(list.menu.menu.Items)
	require.Len(t, list.menu.Menu().Items, baseMenuLen)

	// Operations on empty list.
	require.False(t, list.HasActive())
	require.Nil(t, list.GetActive())
	require.False(t, list.IsActive(123))
	require.Nil(t, list.Get(123))
	require.NoError(t, list.Remove(nil))
	require.ErrorIs(t, list.Remove(&mockItem{l: "Test 1", a: false}), ErrItemNotFound)
	require.Len(t, list.menu.Menu().Items, baseMenuLen)

	newItems := []*mockItem{
		1: {l: "Test 2", a: false},
		2: {l: "Test 3", a: false},
		3: {l: "Must be deleted", a: false},
		4: {l: "Test 4", a: false},
	}

	// Add and delete some items.
	require.Equal(t, 1, list.Add(newItems[1]))
	require.Equal(t, 2, list.Add(newItems[2]))
	require.Equal(t, 3, list.Add(newItems[3]))
	require.Len(t, list.menu.Menu().Items, baseMenuLen+3) // Should increase total items num by one item.
	require.NoError(t, list.Remove(newItems[3]))
	require.Len(t, list.menu.Menu().Items, baseMenuLen+2)
	require.Equal(t, 4, list.Add(newItems[4]))
	require.Len(t, list.menu.Menu().Items, baseMenuLen+3)

	// New item is accessible.
	require.Equal(t, newItems[1], list.Get(1))
	require.Equal(t, newItems[2], list.Get(2))
	require.Equal(t, (*mockItem)(nil), list.Get(3))
	require.True(t, list.Get(3) == nil, "deleted item must be nil")
	require.Equal(t, newItems[4], list.Get(4))
	require.Len(t, list.menu.Menu().Items, baseMenuLen+3)

	require.False(t, list.HasActive())

	calledOnClickForNew := -1
	list.OnItemClick(func(i int) error {
		calledOnClickForNew = i
		newItems[i].a = !newItems[i].a

		return nil
	})

	require.Equal(t, "Test 2", list.getItem(1).menuItem.Label)
	require.Equal(t, "Test 3", list.getItem(2).menuItem.Label)
	require.Equal(t, "Test 4", list.getItem(4).menuItem.Label)

	assertActive := func(id int) {
		require.Equal(t, id, calledOnClickForNew)
		require.True(t, list.HasActive())
		require.Equal(t, newItems[id], list.GetActive())
		require.True(t, newItems[id].a)
		calledOnClickForNew = -1
	}

	// Imitate click on the item in tray menu, should invoke List.OnItemClick
	list.getItem(2).menuItem.Action()
	assertActive(2)

	// Try to delete active item.
	require.ErrorContains(t, list.Remove(newItems[2]), "item 2 is active")

	// imitate click on the item in tray menu to disable it.
	list.getItem(2).menuItem.Action()
	require.Equal(t, 2, calledOnClickForNew)
	calledOnClickForNew = -1
	require.False(t, list.HasActive())
	require.Nil(t, list.GetActive())

	// Make 4 item active
	list.getItem(4).menuItem.Action()
	assertActive(4)

	// Switch active by clicking on another tray item.
	list.getItem(2).menuItem.Action()
	assertActive(2)

	// Delete inactive items.
	require.NoError(t, list.Remove(newItems[1]))
	require.NoError(t, list.Remove(newItems[4]))
}

func TestTrayList_Status(t *testing.T) {
	var lastTrayIcon fyne.Resource
	list := setupList(deskMock{
		onIconSet: func(ic fyne.Resource) {
			lastTrayIcon = ic
		},
	})

	baseMenuLen := len(list.menu.menu.Items)
	require.Len(t, list.menu.Menu().Items, baseMenuLen)

	newItems := map[int]*mockItem{
		1: {l: "Test 2", a: false},
		2: {l: "Test 3", a: false},
		3: {l: "Test 4", a: false},
	}
	list.Add(newItems[1])
	list.Add(newItems[2])
	list.Add(newItems[3])

	list.OnItemClick(func(i int) error {
		newItems[i].a = !newItems[i].a

		return nil
	})
	require.Equal(t, nil, list.getItem(2).menuItem.Icon)

	assertIcons := func(mainIcon, listIcon fyne.Resource, id int) {
		require.Equal(t, lastTrayIcon, mainIcon)
		for i := range newItems {
			icn := listIcon
			if i != id {
				icn = nil
			}
			require.Equal(t, icn, list.getItem(i).menuItem.Icon)
		}
	}

	// Imitate click on the item in tray menu, should invoke List.OnItemClick
	list.getItem(2).menuItem.Action()
	assertIcons(theme.MediaPlayIcon(), theme.ConfirmIcon(), 2)

	list.getItem(3).menuItem.Action()
	assertIcons(theme.MediaPlayIcon(), theme.ConfirmIcon(), 3)

	list.getItem(3).menuItem.Action()
	require.Equal(t, nil, list.getItem(1).menuItem.Icon)
	require.Equal(t, nil, list.getItem(2).menuItem.Icon)
	require.Equal(t, nil, list.getItem(3).menuItem.Icon)
	require.Equal(t, lastTrayIcon.Name(), theme.MediaPauseIcon().Name())

	// Check error icon
	list.OnItemClick(func(i int) error {
		return errors.New("test error")
	})

	list.getItem(3).menuItem.Action()
	assertIcons(theme.WarningIcon(), theme.WarningIcon(), 3)

	require.Equal(t, "test error", list.menu.Menu().Items[0].Label)
}

func setupList(desk desktop.App) *List[*mockItem] {
	list := NewDefault[*mockItem]("title", desk, nil)
	list.menu.refresh = func() {} // To not initialize fyne windows.

	return list
}
