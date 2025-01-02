package netchart

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/goxray/ui/internal/netchart/mocks"
)

func TestRecorder_NilSource(t *testing.T) {
	rec := NewRecorder(nil)

	require.Equal(t, 0, rec.WrittenSinceLast())
	require.Equal(t, 0, rec.ReadSinceLast())
}

func TestRecorder(t *testing.T) {
	incR, incW := 0, 0
	i := 0
	sourceMock := mocks.NewMockSource(gomock.NewController(t))
	sourceMock.EXPECT().BytesRead().DoAndReturn(func() int {
		i++
		incR += 1 * i * bytesToMb
		return incR
	}).AnyTimes()
	sourceMock.EXPECT().BytesWritten().DoAndReturn(func() int {
		i++
		incW += 1 * i * bytesToMb
		return incW
	}).AnyTimes()

	rec := NewRecorder(sourceMock)
	rec.recordLimit = 10
	rec.interval = time.Millisecond

	require.Equal(t, rec.interval, rec.RecordInterval())

	rec.Start()
	<-time.After(time.Millisecond * 20)
	rec.Stop()

	require.Equal(t, []float64{25, 29, 33, 37, 41, 45, 49, 53, 57, 61, 65}, rec.Read())
	require.Equal(t, []float64{27, 31, 35, 39, 43, 47, 51, 55, 59, 63, 67}, rec.Written())

	require.Equal(t, 150875000, rec.BytesWritten())
	require.Equal(t, 142375000, rec.BytesRead())
}
