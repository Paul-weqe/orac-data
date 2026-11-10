# orac-data

Data layer for Orac, a pre-deployment network verification tool. Fetches
routing state from network devices, normalizes it into a vendor-agnostic
model, and encodes it as SMT constraints for verification.

Currently supports FRRouting over gRPC. Interface state is the first
supported data type. BGP routes, route-maps, and VRFs are planned.

## Requirements

- Go 1.21 or newer
- A C toolchain (`gcc` or `clang`)
- CMake 3.15 or newer
- `make`
- `git`
- Python 3 (used by Z3's build scripts)

On Debian or Ubuntu:

    sudo apt install build-essential cmake make git python3

On macOS with Homebrew:

    brew install cmake

## Building Z3 from source

The Go bindings ship inside the Z3 source tree at `src/api/go`. They are
not a published module, so Z3 must be cloned and built locally before
this project compiles.

### 1. Clone Z3

    git clone https://github.com/Z3Prover/z3.git

Pick a location you intend to keep. The build writes into that directory
and this project refers to it by path, so moving it later means updating
the Makefile and `go.mod`.

The examples below assume:

    /home/<user>/z3

### 2. Configure

    cd z3
    mkdir build && cd build
    cmake -DZ3_BUILD_GO_BINDINGS=ON -DZ3_BUILD_LIBZ3_SHARED=ON ..

Both flags matter:

- `Z3_BUILD_GO_BINDINGS=ON` generates the Go binding sources.
- `Z3_BUILD_LIBZ3_SHARED=ON` produces `libz3.so`. The bindings need a
  shared library; a static-only build will not link.

### 3. Build

    make -j"$(nproc)"

On macOS:

    make -j"$(sysctl -n hw.ncpu)"

This takes a few minutes.

### 4. Verify

Confirm the shared library was produced:

    ls -la build/libz3.so*

You should see a `libz3.so` symlink pointing at a versioned file, for
example `libz3.so.5.4`.

Confirm the library contains a symbol the bindings rely on:

    nm -D build/libz3.so | grep Z3_mk_abs

Expected output resembles:

    000000000036ff40 T Z3_mk_abs

If nothing is printed, the build did not enable the Go bindings or the
shared library. Remove the `build/` directory and repeat step 2.

### 5. Install the library (recommended)

Installing `libz3` system-wide removes ambiguity about which copy is
loaded at build time and at run time:

    sudo make install
    sudo ldconfig

The library ends up in `/usr/local/lib` and the loader finds it without
further configuration.

If you prefer not to install system-wide, skip this step. The build
section below embeds an rpath in the binary, which covers most cases.

## Configuring the project

### 1. Set Z3_DIR

`Z3_DIR` points at your Z3 checkout, the directory containing `src/` and
`build/`. The Makefile reads it in three ways, in order of precedence:

Command line:

    make build-z3 Z3_DIR=/home/<user>/z3

Environment:

    export Z3_DIR=/home/<user>/z3

The default in the Makefile itself:

    Z3_DIR ?= /home/<user>/z3

Edit the default so it matches your checkout, or override it per
invocation. On a shared machine it is usually cleaner to override it in
your shell rather than edit the file.

### 2. Point go.mod at the local sources

The Go bindings are not fetchable through the module proxy, so `go.mod`
needs a `replace` directive:

    require github.com/Z3Prover/z3/src/api/go v0.0.0
    replace github.com/Z3Prover/z3/src/api/go => /home/<user>/z3/src/api/go

Adjust the path to match `Z3_DIR`. Without the `replace`, `go build`
fails trying to resolve a module that does not exist.

## Building

    make build-z3

The recipe expands to:

    CGO_CFLAGS="-I$(Z3_DIR)/src/api" \
    CGO_LDFLAGS="$(Z3_DIR)/build/libz3.so -lpthread" \
    go build -ldflags="-r $(Z3_DIR)/build" -o bin/z3 cmd/z3/main.go

A note on each part:

- `CGO_CFLAGS` points the compiler at the Z3 headers in `src/api`.
- `CGO_LDFLAGS` passes the absolute path to your `libz3.so`. Using an
  absolute path rather than `-L... -lz3` keeps the linker from picking up
  a system copy, which is usually older and missing symbols the bindings
  need.
- `-ldflags="-r $(Z3_DIR)/build"` embeds an rpath in the binary so the
  loader resolves `libz3.so` without `LD_LIBRARY_PATH`.

Other targets:

    make build-grpc     # builds cmd/explore
    make frr-grpc       # regenerates gRPC stubs from proto/frr/frr-northbound.proto
    make clean          # removes bin/

## Running

    bin/z3

With the rpath in place this runs directly. If the binary is moved to a
different machine, or if it was built without the rpath, set the loader
path:

    LD_LIBRARY_PATH=$(Z3_DIR)/build:$LD_LIBRARY_PATH bin/z3

On macOS:

    DYLD_LIBRARY_PATH=$(Z3_DIR)/build:$DYLD_LIBRARY_PATH bin/z3

## Troubleshooting

### undefined reference to Z3_mk_abs

The linker resolved a `libz3` older than the binding sources. This
normally means a system-installed copy was picked up.

Check what the linker would find:

    ldconfig -p | grep libz3

Check whether `go env` is overriding the flags:

    go env CGO_LDFLAGS CGO_CFLAGS

If either is set, clear it so the Makefile values apply:

    go env -u CGO_LDFLAGS
    go env -u CGO_CFLAGS

Then run `make build-z3` again. The Makefile passes the absolute path to
your library, so no search path is consulted.

### libz3.so.5.4: cannot open shared object file

The binary was built but the loader cannot find the library at run time.
Any of these resolve it:

- Install Z3 system-wide: `sudo make install && sudo ldconfig`
- Set `LD_LIBRARY_PATH=$(Z3_DIR)/build` when running
- Rebuild; `build-z3` embeds the rpath by default

To see which library the binary expects and whether it resolves:

    ldd bin/z3 | grep z3

### z3.h: No such file or directory

`CGO_CFLAGS` is not reaching the compiler. Confirm `Z3_DIR` is set and
that `$(Z3_DIR)/src/api/z3.h` exists. If the build is being invoked
through a wrapper, check that `Z3_DIR` is not being overwritten.

## Project layout

    .
    ├── Makefile
    ├── go.mod
    ├── cmd/
    │   ├── explore/       # gRPC exploration binary
    │   └── z3/            # Z3 front-end binary
    ├── internal/
    │   ├── adapter/       # vendor transport (gRPC, future: NETCONF)
    │   ├── models/        # unified and vendor-specific data models
    │   └── normalizer/    # vendor data to unified conversion
    ├── proto/
    │   └── frr/           # FRR northbound protobuf definitions
    └── bin/               # build output (gitignored)

## Development state

Working:

- FRR gRPC adapter for interface data (`FetchInterfaces`)
- FRR interface normalizer (`InterfaceData`)

In progress:

- Z3 integration (build harness in place, encoder not yet written)

Not started:

- BGP, route-map, and VRF data types
- Additional vendors (Holo, Cisco, Juniper, Nokia)
- Additional transports (NETCONF, RESTCONF)
