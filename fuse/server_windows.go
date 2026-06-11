// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import "errors"

// Server is a compile-only stub on Windows. The FUSE kernel protocol
// server cannot run on Windows; filesystems built on the fuse types
// (e.g. a WinFsp adapter) are served by other means. The stub exists so
// that RawFileSystem.Init(*Server) and code holding a *Server compile.
type Server struct{}

var errServerNotSupported = errors.New("fuse: Server is not supported on windows")

// NewServer always fails on Windows.
func NewServer(fs RawFileSystem, mountPoint string, opts *MountOptions) (*Server, error) {
	return nil, errServerNotSupported
}

func (ms *Server) Serve() {}

func (ms *Server) Wait() {}

func (ms *Server) WaitMount() error { return errServerNotSupported }

func (ms *Server) Unmount() error { return errServerNotSupported }

func (ms *Server) InodeNotify(node uint64, off int64, length int64) Status { return OK }

func (ms *Server) EntryNotify(parent uint64, name string) Status { return OK }

func (ms *Server) DeleteNotify(parent uint64, child uint64, name string) Status { return OK }

func (ms *Server) InodeNotifyStoreCache(node uint64, offset int64, data []byte) Status { return OK }

func (ms *Server) KernelSettings() *InitIn { return &InitIn{} }
