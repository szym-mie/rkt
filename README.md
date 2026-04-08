# rkt

![sshot1](./promo/sshot1.png)
![sshot2](./promo/sshot2.png)
![sshot3](./promo/sshot3.png)
![sshot4](./promo/sshot4.png)

## 1. Quickstart

Use `make` from the project root to get a playable release.

Move the mouse to rotate the camera around the craft.

| keys  | action                                                     |
| ----- | ---------------------------------------------------------- |
| W, S  | pitch down/up                                              |
| A, D  | yaw left/right                                             |
| Q, E  | roll left/right                                            |
| Space | activate the next stage (depends on vehicle configuration) |
| -, =  | zoom camera out/in                                         |
| Esc   | quit                                                       |

## 2. Building

Available make actions:

|      Action       | Command      |
| :---------------: | :----------- |
|   Build release   | `make`       |
| Development build | `make devel` |
|  Clean all exes   | `make clean` |

> [!IMPORTANT]
> Please note that this project uses CGO and therefore make sure you have the right version of the GCC compiler installed on your machine (please refer to https://go.dev/doc/install/gccgo). On Windows Winlibs Mingw64 toolchain is confirmed to work properly.

## 3. Resources

The resources are stored in a single .ZIP folder, not unlike the .PK3 format. Resources other than bitmaps (and possibly sound files in the future) follow the convention of `<name>.<type>.<ext>` - allowing the game to discern between different resource types and thus to use an correct loader.

## 4. Progress

A short and hopefully up-to-date list of things to implement:

#### OK

- Staging
- Fix camera projection matrices
- Haze implemented as OpenGL fog
- Parachutes
- Impulse can create angular momentum
- Calculate Centre of Mass
- Angular movement around the CoM
- Inertia calculated based on the simple geom. shape description

#### WIP

- Part drag calculation based on AoA (Angle of Attack)
- Mouse-based vehicle editor

#### TODO

- Altitude-dependent skybox
- Proper terrain collisions
- Try out orbital mechanics
- Fuel flow and distribution
- Lift
- Basic lift control devices (all moving fins)
- Ailerons, rudders, elevators
- Airbrakes and spoilers
- Cockpit view
- Migrate to OpenGL3

#### IDEA

- Key-based (with some mouse usage) vehicle editor
- Electrical circuits: actuators, generators, batteries
- Hydraulic circuits: servos, actuators, accumulators, pumps
