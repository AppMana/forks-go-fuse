// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import "syscall"

// Windows is a compile-only target for this package: the FUSE kernel
// protocol server never runs on Windows. These definitions exist so
// that filesystem implementations layered on the fuse types (e.g. a
// WinFsp adapter) can build with GOOS=windows.
const (
	ENODATA = Status(syscall.ENODATA)
	// ENOATTR is an alias for ENODATA on Linux; Windows syscall has no
	// ENOATTR, so mirror the Linux aliasing.
	ENOATTR = Status(syscall.ENODATA)

	// EREMOTEIO Remote I/O error
	EREMOTEIO = Status(syscall.EREMOTEIO)
)

// Capability bits match the Linux FUSE protocol definitions.
const (
	CAP_NO_OPENDIR_SUPPORT  = (1 << 24)
	CAP_EXPLICIT_INVAL_DATA = (1 << 25)

	CAP_MAP_ALIGNMENT      = (1 << 26)
	CAP_SUBMOUNTS          = (1 << 27)
	CAP_HANDLE_KILLPRIV_V2 = (1 << 28)
	CAP_SETXATTR_EXT       = (1 << 29)
	CAP_INIT_EXT           = (1 << 30)
	CAP_INIT_RESERVED      = (1 << 31)

	// CAP_RENAME_SWAP only exists on OSX.
	CAP_RENAME_SWAP = 0x0
)

func (o *InitOut) setFlags(flags uint64) {
	o.Flags = uint32(flags) | CAP_INIT_EXT
	o.Flags2 = uint32(flags >> 32)
}
