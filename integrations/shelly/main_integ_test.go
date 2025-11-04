package shelly

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestShellySwitchIntegration(t *testing.T) {
	c := GetClient()

	t.Run("Status when OFF", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		st, err := c.GetStatus(ctx)
		require.NoError(t, err)
		first, err := c.Set(ctx, !st.Output, true, 8*time.Second)
		require.NoError(t, err)
		require.True(t, first.Output == !st.Output, "First failed")
		time.Sleep(2 * time.Second)
		second, err := c.Set(ctx, st.Output, true, 8*time.Second)
		require.NoError(t, err)
		require.True(t, second.Output == st.Output, "Second failed")
	})
}
