// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package plugin_test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mesosphere/dkp-cli-runtime/core/plugin"
)

func TestDiscoveryCommand(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pluginCmd := &cobra.Command{
		Use: "example-plugin",
	}
	pluginCmd.Flags().String("flag1", "", "flag1 usage")
	pluginCmd.Flags().StringP("flag2", "f", "flag2Default", "flag2 usage")
	pluginCmd.PersistentFlags().Bool("global", false, "globalFlagUsage")

	subCmd := &cobra.Command{
		Use:   "action",
		Short: "short description",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	subCmd.Flags().Int("flag3", 42, "flag3 usage")
	pluginCmd.AddCommand(subCmd)

	outBuf := bytes.Buffer{}
	discoveryCmd := plugin.NewDiscoveryCommand(&outBuf, pluginCmd)
	err := discoveryCmd.Execute()
	require.NoError(err)

	assert.JSONEq(`{
		"commands": {
			"use": "example-plugin",
			"local_flags": [
				{
					"type": "string",
					"name": "flag1",
					"usage": "flag1 usage"
				},
				{
					"type": "string",
					"name": "flag2",
					"shorthand": "f",
					"usage": "flag2 usage",
					"default_value": "flag2Default"
				}
			],
			"persistent_flags": [
				{
					"type": "bool",
					"name": "global",
					"usage": "globalFlagUsage",
					"default_value": "false",
					"no_opt_default_value": "true"
				}
			],
			"sub_commands": [
				{
					"use": "action",
					"short": "short description",
					"runnable": true,
					"local_flags": [
						{
							"type": "int",
							"name": "flag3",
							"usage": "flag3 usage",
							"default_value": "42"
						}
					]
				}
			]
		}
	}`, outBuf.String())
}
