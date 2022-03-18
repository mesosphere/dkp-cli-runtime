// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreAndRestoreOutoutInsideContext(t *testing.T) {
	output := NewNonInteractiveShell(os.Stdout, os.Stdout, 10)

	output.Info("hello")
	ctx := context.Background()
	ctx = WithContext(ctx, output)
	restoredOutput := Ctx(ctx)
	_, isRightType := restoredOutput.(*nonInteractiveShellOutput)
	assert.True(t, isRightType, "restored Output is not of the expected type")
	restoredOutput.Info("world")
}

func ExampleContext() {
	ctx := context.Background()
	ctx = WithContext(
		ctx,
		NewInteractiveShell(os.Stdout, os.Stdout, 10),
	)

	output := Ctx(ctx)
	output.Info("hello world")

	// Output:
	// hello world
}
