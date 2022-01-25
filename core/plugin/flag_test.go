// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package plugin_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/dkp-cli-runtime/core/plugin"
)

func TestFlagsToSpecAndBack(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ExitOnError)

	flags.String("simple", "", "how to use")
	flags.StringP("with-shorthand", "s", "default value", "how to use")
	flags.StringArray("array", []string{"default", "values"}, "how to use")
	flags.Bool("boolean", true, "how to use")
	flags.String("annotated", "", "how to use")
	flags.SetAnnotation("annotated", "key", []string{"value1", "value2"})
	flags.String("hidden", "", "how to use")
	flags.MarkHidden("hidden")
	flags.String("deprecated", "", "how to use")
	flags.MarkDeprecated("deprecated", "use that other flag instead")
	flags.String("shorthand-deprecated", "", "how to use")
	flags.MarkShorthandDeprecated("shorthand-deprecated", "use that other shorthand instead")

	flags.VisitAll(func(flag *pflag.Flag) {
		t.Run(flag.Name, func(t *testing.T) {
			spec := plugin.SpecFromFlag(flag)
			assertFlagsEqual(t, flag, spec.ToFlag())
		})
	})
}

func assertFlagsEqual(t *testing.T, expected, actual *pflag.Flag) {
	assert := assert.New(t)
	assert.Equal(expected.Value.Type(), actual.Value.Type())
	// compare everything but the Value
	actual.Value = expected.Value
	assert.Equal(expected, actual)
}
