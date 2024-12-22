package main

import (
	"encoding/json"
	"log/slog"
)

const (
	itemsConfigKey = "connections_config"
)

type SaveFile struct {
	source Source
}

type Source interface {
	// SetString saves key value pair.
	SetString(key string, value string)
	// StringWithFallback gets string by key, if not found return fallback.
	StringWithFallback(key, fallback string) string
}

func NewSaveFile(source Source) *SaveFile {
	return &SaveFile{
		source: source,
	}
}

// Update saves list into config.
func (s *SaveFile) Update(list *ItemsList) {
	updatedItems := make([]*Item, 0, len(list.All()))
	for _, item := range list.All() {
		if item == nil {
			continue
		}
		updatedItems = append(updatedItems, item)
	}

	b, err := json.MarshalIndent(updatedItems, "", "  ")
	if err != nil {
		slog.Warn(err.Error())
	}

	s.source.SetString(itemsConfigKey, string(b))
}

// Load loads saved items into list.
func (s *SaveFile) Load(list *ItemsList) {
	items := make([]*Item, 0)
	if err := json.Unmarshal([]byte(s.source.StringWithFallback(itemsConfigKey, "[]")), &items); err != nil {
		slog.Error("failed to unmarshal tray items", "error", err)
	}

	for _, item := range items {
		list.Add(item)
	}
}

//
// func SaveFileSave(pref fyne.Preferences, list *ItemsList[*Item]) {
// 	tosave := make([]*Item, 0, len(list.All()))
// 	for _, item := range list.All() {
// 		if item == nil {
// 			continue
// 		}
// 		tosave = append(tosave, item)
// 	}
//
// 	b, err := json.MarshalIndent(tosave, "", "  ")
// 	if err != nil {
// 		slog.Warn(err.Error())
// 	}
//
// 	pref.SetString(itemsConfigKey, string(b))
// }
//
// func SaveFileLoad(pref fyne.Preferences, list *ItemsList[*Item]) {
// 	items := make([]*Item, 0)
// 	if err := json.Unmarshal([]byte(pref.StringWithFallback(itemsConfigKey, EmptyItemsConfig)), &items); err != nil {
// 		slog.Error("failed to unmarshal tray items: %v", err)
// 	}
//
// 	for _, item := range items {
// 		list.Add(item)
// 	}
// }
