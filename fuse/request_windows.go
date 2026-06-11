// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

// Values match the Linux FUSE protocol; the protocol server never runs
// on Windows (compile-only target).
const outputHeaderSize = 304

const (
	_FUSE_KERNEL_VERSION   = 7
	_MINIMUM_MINOR_VERSION = 12
	_OUR_MINOR_VERSION     = 28
)
