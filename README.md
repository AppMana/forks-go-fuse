# forks-go-fuse — go-fuse with a compile-only Windows port

AppMana fork of [seaweedfs/go-fuse](https://github.com/seaweedfs/go-fuse)
(v2.9.3). The `fuse` package builds under `GOOS=windows`: the platform stubs
mirror the in-repo FreeBSD port, and the kernel protocol server is stubbed
(it never runs on Windows — filesystems built on the fuse types are served by
other means, e.g. WinFsp via cgofuse).

This exists so [AppMana/forks-seaweedfs](https://github.com/AppMana/forks-seaweedfs)
can compile its mount filesystem for Windows without forking its internals.
The module path is unchanged; consumers use a `replace` directive pointing at
a sibling checkout.

Upstream README: https://github.com/seaweedfs/go-fuse
