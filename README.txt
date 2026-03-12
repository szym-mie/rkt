rkt
===

1. Controls
-----------

Move the mouse to rotate the camera around the craft

|  W, S   | pitch down/up
|  A, D   | yaw left/right
|  Q, E   | roll left/right
|  Space  | activate the next stage (depends on vehicle configuration)
|  -, =   | zoom camera out/in
|  Esc    | quit

2. Running
----------

Use 'go run .' from the project root. One of the unforeseen problems with the Go + OpenGL is that Windows does not trust the binaries built with 'go build -o rkt.ext main.go', completely blocking the execution of the game. Signing might help, but my first attempt at this didn't work at all.

3. Resources
------------

The resources are stored in a single .ZIP folder, not unlike the .PK3 format. Resources other than bitmaps (and possibly sound files in the future) follow the convention of <name>.<type>.<ext> - allowing the game to discern between different resource types and thus to use an correct loader.
