// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import "context"

type key struct{}

// WithContext allows to store an instance of Output in context.Context.
func WithContext(ctx context.Context, o Output) context.Context {
	return context.WithValue(ctx, key{}, o)
}

// Ctx restores Output from a given Context.
func Ctx(ctx context.Context) Output {
	if o, ok := ctx.Value(key{}).(Output); ok {
		return o
	}

	return &noopOutput{}
}
