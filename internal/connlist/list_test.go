package connlist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/goxray/desktop/internal/connlist/mocks"
)

const sampleVlessLink = "vless://h1px412i-9138-s9m5-9b86-d47d74dd8541@127.0.0.1:8080?type=tcp&security=reality&pbk=4442383675fc0fb574c3e50abbe7d4c5&fp=chrome&sni=yahoo.com&sid=0c&spx=%2F&flow=xtls-rprx-vision#Myremark"

func TestList_AddDelete(t *testing.T) {
	c := New()

	require.Empty(t, c.All())
	require.Equal(t, len(c.All()), len(*c.AllUntyped()))
	require.ErrorContains(t, c.AddItem("Test 1", "link"), "invalid xray link: invalid protocol type")

	var createdItem *Item
	c.OnAdd(func(item *Item) {
		require.Equal(t, item.Label(), "Test 2")
		require.Equal(t, sampleVlessLink, item.Link())
		createdItem = item
	})

	require.NoError(t, c.AddItem("Test 2", sampleVlessLink))
	require.Len(t, c.All(), 1)
	require.Equal(t, len(c.All()), len(*c.AllUntyped()))

	c.RemoveItem(createdItem)
	require.Empty(t, c.All())
	require.Equal(t, len(c.All()), len(*c.AllUntyped()))
}

func TestList_Recorder(t *testing.T) {
	rec := mocks.NewMockNetworkRecorder(gomock.NewController(t))
	rec.EXPECT().RecordInterval().Return(time.Millisecond * 123)
	rec.EXPECT().BytesRead().Return(1234)
	rec.EXPECT().BytesWritten().Return(4321)
	rec.EXPECT().Read().Return([]float64{1, 2, 3})
	rec.EXPECT().Written().Return([]float64{3, 2, 1})

	c, err := newItem("test", sampleVlessLink, New())
	require.NoError(t, err)
	c.recorder = rec

	require.Equal(t, time.Millisecond*123, c.RecordInterval())
	require.Equal(t, 1234, c.BytesRead())
	require.Equal(t, 4321, c.BytesWritten())
	require.Equal(t, []float64{1, 2, 3}, c.Read())
	require.Equal(t, []float64{3, 2, 1}, c.Written())
}
