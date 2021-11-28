package app

import (
	"context"
	"testing"

	"github.com/inspectorvitya/fibonacci_service/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestApp_GetFibSlice(t *testing.T) {
	db := memory.New()
	fibApp := New(db)
	ctx := context.Background()
	t.Run("basic", func(t *testing.T) {
		actual, err := fibApp.GetFibSlice(ctx, 1, 10)
		require.NoError(t, err)
		require.Equal(t, []string{"0", "1", "1", "2", "3", "5", "8", "13", "21", "34"}, actual)

		actual, err = fibApp.GetFibSlice(ctx, 100, 101)
		require.NoError(t, err)
		require.Equal(t, []string{"218922995834555169026", "354224848179261915075"}, actual)
	})
	t.Run("with errors", func(t *testing.T) {
		t.Run("x more y", func(t *testing.T) {
			_, err := fibApp.GetFibSlice(ctx, 10, 2)
			require.ErrorIs(t, err, ErrXMoreY)
		})
		t.Run("x\\y less zero", func(t *testing.T) {
			_, err := fibApp.GetFibSlice(ctx, 0, 10)
			require.ErrorIs(t, err, ErrLessOrGreaterZero)
		})
	})
}
