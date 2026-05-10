# Building

## 1. Preparation

### 1.1. Software

- GNU Make
- Go 1.25+
- GCC or Clang

### 1.2. Libraries

The `go-gl` and `go-glfw` require specific libraries to be present on a system. As per the [go-glfw README](https://github.com/go-gl/glfw):

- on macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
- on Ubuntu/Debian-like Linux distributions, the default Linux build enables both X11 and Wayland, so you need `libgl1-mesa-dev`, `xorg-dev`, `libwayland-dev`, `libxkbcommon-dev` and `wayland-protocols`.
- on CentOS/Fedora-like Linux distributions, you need `libX11-devel`, `libXcursor-devel`, `libXrandr-devel`, `libXinerama-devel`, `mesa-libGL-devel`, `libXi-devel`, `libXxf86vm-devel` packages.
- on FreeBSD, you need the package `pkgconf`. To build for X, you also need the package `xorg`; and to build for Wayland, you need the package `wayland`.
- on NetBSD, to build for X, you need the X11 sets installed. These are included in all graphical installs, and can be added to the system with `sysinst(8)` on non-graphical systems. Wayland support is incomplete, due to missing `wscons` support in upstream GLFW. To attempt to build for Wayland, you need to install the `wayland`, `libepoll-shim` packages and set the environment variable `PKG_CONFIG_PATH=/usr/pkg/libdata/pkgconfig`.
- on OpenBSD, you need the X11 sets. These are installed by default, and can be added from the ramdisk kernel at any time.

Additional information can be found in the [GLFW docs](https://www.glfw.org/docs/latest/compile.html#compile_deps).

### 1.3. Toolchain

The GCC toolchain is recommended and worked well on all of the tested platforms: Windows, Ubuntu, Fedora, Debian, FreeBSD.

> [!IMPORTANT]
> Please note that this project uses CGO and therefore make sure you have the right version of the GCC compiler installed on your machine (please refer to https://go.dev/doc/install/gccgo). On Windows Winlibs Mingw64 toolchain is confirmed to work properly.

## 2. Build

Available actions:

|            Action             | Command        |
| :---------------------------: | :------------- |
|  Create a release executable  | `make build`   |
|   Run a development session   | `make devel`   |
| create a full release package | `make release` |
|        Clean all exes         | `make clean`   |

The scripts will automatically create a resource pack and compile the game.

### 2.1. Release

To build a release package you have to provide additional parameters like so:

- v9 for amd64 Windows: `make release PLAT=win64 VER=v9`
- v11 for i386 NetBSD: `make release PLAT=netbsd32 VER=v11`
