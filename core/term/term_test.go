// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTerminal(t *testing.T) {
	// test trivial nil case
	if IsTerminal(nil) {
		t.Fatalf("IsTerminal should be false for nil Writer")
	}
	// test something that isn't even a file
	var buff bytes.Buffer
	if IsTerminal(&buff) {
		t.Fatalf("IsTerminal should be false for bytes.Buffer")
	}
	// test a file
	f, err := os.CreateTemp("", "kind-isterminal")
	if err != nil {
		t.Fatalf("Failed to create tempfile %v", err)
	}
	if IsTerminal(f) {
		t.Fatalf("IsTerminal should be false for nil Writer")
	}
	// TODO: testing an actual PTY would be somewhat tricky to do cleanly
	// but we should maybe do this in the future.
	// At least we know this doesn't trigger on things that are obviously not
	// terminals
	if !IsTerminal(&testFakeTTY{}) {
		t.Fatalf("IsTerminal should be true for testFakeTTY")
	}
}

func TestIsSmartTerminal(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Name    string
		FakeEnv map[string]string
		GOOS    string
		Writer  io.Writer
		IsSmart bool
	}{
		{
			Name:    "tty, no env",
			FakeEnv: map[string]string{},
			GOOS:    "linux",
			IsSmart: true,
			Writer:  &testFakeTTY{},
		},
		{
			Name:    "nil writer, no env",
			FakeEnv: map[string]string{},
			GOOS:    "linux",
			IsSmart: false,
		},
		{
			Name:    "tty, windows, no env",
			FakeEnv: map[string]string{},
			GOOS:    "windows",
			IsSmart: false,
			Writer:  &testFakeTTY{},
		},
		{
			Name: "tty, windows, modern terminal env",
			FakeEnv: map[string]string{
				"WT_SESSION": "baz",
			},
			GOOS:    "windows",
			IsSmart: true,
			Writer:  &testFakeTTY{},
		},
		{
			Name: "tty, TERM=dumb",
			FakeEnv: map[string]string{
				"TERM": "dumb",
			},
			GOOS:    "linux",
			IsSmart: false,
			Writer:  &testFakeTTY{},
		},
		{
			Name: "tty, NO_COLOR=",
			FakeEnv: map[string]string{
				"NO_COLOR": "",
			},
			GOOS:    "linux",
			IsSmart: false,
			Writer:  &testFakeTTY{},
		},
		{
			Name: "tty, Travis CI",
			FakeEnv: map[string]string{
				"TRAVIS":                      "true",
				"HAS_JOSH_K_SEAL_OF_APPROVAL": "true",
			},
			GOOS:    "linux",
			IsSmart: false,
			Writer:  &testFakeTTY{},
		},
		{
			Name: "tty, TERM=st-256color",
			FakeEnv: map[string]string{
				"TERM": "st-256color",
			},
			GOOS:    "linux",
			IsSmart: false,
			Writer:  &testFakeTTY{},
		},
	}
	for _, tc := range cases {
		tc := tc // capture tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			res := isSmartTerminal(tc.Writer, tc.GOOS, func(s string) (string, bool) {
				k, set := tc.FakeEnv[s]
				return k, set
			})
			assert.Equal(t, tc.IsSmart, res)
		})
	}
}
