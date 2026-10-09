package fuse

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// stuckFS serves one regular file, "file", and blocks Flush on it and Lookup
// of "stuck" until released, like a file system whose backend stopped
// answering.
type stuckFS struct {
	RawFileSystem
	release chan struct{}
}

func (fs *stuckFS) Lookup(cancel <-chan struct{}, header *InHeader, name string, out *EntryOut) Status {
	switch name {
	case "stuck":
		<-fs.release
		return ENOENT
	case "file":
		out.NodeId = 2
		out.Attr.Ino = 2
		out.Attr.Mode = syscall.S_IFREG | 0o644
		out.Attr.Nlink = 1
		return OK
	}
	return ENOENT
}

func (fs *stuckFS) GetAttr(cancel <-chan struct{}, input *GetAttrIn, out *AttrOut) Status {
	out.Attr.Ino = input.NodeId
	out.Attr.Nlink = 1
	if input.NodeId == FUSE_ROOT_ID {
		out.Attr.Mode = syscall.S_IFDIR | 0o755
	} else {
		out.Attr.Mode = syscall.S_IFREG | 0o644
	}
	return OK
}

func (fs *stuckFS) Open(cancel <-chan struct{}, input *OpenIn, out *OpenOut) Status {
	out.Fh = 1
	return OK
}

func (fs *stuckFS) Flush(cancel <-chan struct{}, input *FlushIn) Status {
	<-fs.release
	return OK
}

func (fs *stuckFS) Release(cancel <-chan struct{}, input *ReleaseIn) {}

func mountStuckFS(t *testing.T, opts MountOptions) (string, *stuckFS) {
	t.Helper()
	dir := t.TempDir()
	fs := &stuckFS{RawFileSystem: NewDefaultRawFileSystem(), release: make(chan struct{})}
	srv, err := NewServer(fs, dir, &opts)
	if err != nil {
		t.Fatal(err)
	}
	go srv.Serve()
	if err := srv.WaitMount(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		close(fs.release)
		srv.Unmount()
	})
	return dir, fs
}

// A request the file system never answers must still get a reply within
// MountOptions.RequestTimeout, so the calling task cannot sit in
// uninterruptible sleep until the hung-task watchdog fires.
func TestRequestTimeoutAnswersStuckLookup(t *testing.T) {
	const deadline = 2 * time.Second
	dir, _ := mountStuckFS(t, MountOptions{RequestTimeout: deadline})

	done := make(chan error, 1)
	go func() {
		_, err := os.Lstat(filepath.Join(dir, "stuck"))
		done <- err
	}()
	select {
	case err := <-done:
		if !errors.Is(err, syscall.EIO) {
			t.Fatalf("stuck lookup returned %v, want EIO", err)
		}
	case <-time.After(3 * deadline):
		t.Fatalf("stuck lookup did not return within %v", 3*deadline)
	}
}

// The exit of a SIGKILLed process closes its files, and the kernel sends a
// FLUSH that the task waits for uninterruptibly. That wait must end within
// the request timeout.
func TestRequestTimeoutReleasesKilledProcessFlush(t *testing.T) {
	const deadline = 2 * time.Second
	dir, _ := mountStuckFS(t, MountOptions{RequestTimeout: deadline})

	cmd := exec.Command("sh", "-c", `exec 3<"$1"; echo ready; exec sleep 600`, "sh", filepath.Join(dir, "file"))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	ready := make([]byte, 6)
	if _, err := stdout.Read(ready); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	exited := make(chan struct{})
	go func() {
		cmd.Wait()
		close(exited)
	}()
	select {
	case <-exited:
	case <-time.After(3 * deadline):
		t.Fatalf("killed process holding a file did not exit within %v", 3*deadline)
	}
}

// The kernel aborts a connection whose request outlives the timeout that the
// server negotiated at FUSE_INIT, independently of the server answering.
func TestKernelRequestTimeoutAbortsConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("waits for the kernel's request-timeout timer")
	}
	const timeout = 15 * time.Second
	dir, _ := mountStuckFS(t, MountOptions{KernelRequestTimeout: timeout})

	done := make(chan error, 1)
	go func() {
		_, err := os.Lstat(filepath.Join(dir, "stuck"))
		done <- err
	}()
	// The kernel checks for expired requests on a 15 s timer.
	select {
	case err := <-done:
		if !errors.Is(err, syscall.ENOTCONN) && !errors.Is(err, syscall.ECONNABORTED) {
			t.Fatalf("stuck lookup returned %v, want the aborted connection's error", err)
		}
	case <-time.After(timeout + 20*time.Second):
		t.Fatalf("kernel did not abort the connection within %v", timeout+20*time.Second)
	}
}
