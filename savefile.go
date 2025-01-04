package main

import (
	"encoding/json"
	"log/slog"

	"github.com/goxray/desktop/internal/connlist"
)

const (
	itemsConfigKey = "connections_config"
)

// SaveFile is used to store and load connection items from memory.
type SaveFile struct {
	source Source
}

type Source interface {
	// SetString saves key value pair.
	SetString(key string, value string)
	// StringWithFallback gets string by key, if not found return fallback.
	StringWithFallback(key, fallback string) string
}

type SavedState struct {
	Link  string `json:"link"`
	Label string `json:"label"`
}

func serialize(item *connlist.Item) SavedState {
	return SavedState{
		Link:  item.Link(),
		Label: item.Label(),
	}
}

func NewSaveFile(source Source) *SaveFile {
	return &SaveFile{
		source: source,
	}
}

// Update saves list into config.
func (s *SaveFile) Update(list *connlist.Collection) {
	toSave := make([]SavedState, 0, len(list.All()))
	for _, item := range list.All() {
		if item == nil {
			continue
		}
		toSave = append(toSave, serialize(item))
	}

	b, err := json.MarshalIndent(toSave, "", "  ")
	if err != nil {
		slog.Warn(err.Error())
	}

	s.source.SetString(itemsConfigKey, string(b))
}

// Load loads saved items into list.
func (s *SaveFile) Load(list *connlist.Collection) {
	loadedItems := make([]SavedState, 0)
	if err := json.Unmarshal([]byte(s.source.StringWithFallback(itemsConfigKey, "[]")), &loadedItems); err != nil {
		slog.Error("failed to unmarshal tray items", "error", err)
	}

	for _, item := range loadedItems {
		if err := list.AddItem(item.Label, item.Link); err != nil {
			slog.Error("failed to load new item", "error", err)
		}
	}
}
