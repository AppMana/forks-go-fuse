// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

// ReadResultData is the read return for returning bytes directly.
type readResultData struct {
	// Raw bytes for the read.
	Data []byte
}

func (r *readResultData) Size() int {
	return len(r.Data)
}

func (r *readResultData) Done() {
}

func (r *readResultData) Bytes(buf []byte) ([]byte, Status) {
	return r.Data, OK
}

func ReadResultData(b []byte) ReadResult {
	return &readResultData{b}
}

