rkt
===

1. Quickstart
-------------

Use 'go run .' from the project root.

Move the mouse to rotate the camera around the craft.

|  W, S   | pitch down/up
|  A, D   | yaw left/right
|  Q, E   | roll left/right
|  Space  | activate the next stage (depends on vehicle configuration)
|  -, =   | zoom camera out/in
|  Esc    | quit

2. Building
-----------

Standard Go rules persist - please note that this project uses CGO and therefore make sure you have the right version of the GCC compiler installed on your machine (please refer to https://go.dev/doc/install/gccgo). On Windows Winlibs Mingw64 toolchain is confirmed to work properly.

3. Resources
------------

The resources are stored in a single .ZIP folder, not unlike the .PK3 format. Resources other than bitmaps (and possibly sound files in the future) follow the convention of <name>.<type>.<ext> - allowing the game to discern between different resource types and thus to use an correct loader.
