// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

// arbitrary values, matching constants_freebsd.go
const syscall_O_LARGEFILE = 1 << 29
const syscall_O_NOATIME = 1 << 30

// Linux open(2) flag values for flag-name printing; the Windows syscall
// package does not define these.
const (
	syscall_O_ASYNC     = 0x2000
	syscall_O_NOCTTY    = 0x100
	syscall_O_NONBLOCK  = 0x800
	syscall_O_CLOEXEC   = 0x80000
	syscall_O_DIRECTORY = 0x10000
)
