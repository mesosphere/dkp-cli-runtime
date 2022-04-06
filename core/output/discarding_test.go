package output_test

import (
	"testing"

	"github.com/mesosphere/dkp-cli-runtime/core/output"
)

func TestDiscardingOutput(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		o := output.NewDiscardingOutput()
		o.Info("test")
	})
}
