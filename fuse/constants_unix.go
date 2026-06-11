//go:build !windows

// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import "syscall"

const (
	syscall_O_ASYNC     = syscall.O_ASYNC
	syscall_O_NOCTTY    = syscall.O_NOCTTY
	syscall_O_NONBLOCK  = syscall.O_NONBLOCK
	syscall_O_CLOEXEC   = syscall.O_CLOEXEC
	syscall_O_DIRECTORY = syscall.O_DIRECTORY
)
