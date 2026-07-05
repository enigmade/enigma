# Enigma Center

The visual face of the `enigma` CLI (SPEC §18). A single Qt 6 / QML app,
launchable from the tray and the launcher.

## Strict rule: thin client, zero duplicated logic

Everything the GUI does, the CLI can do, and vice versa. The GUI holds **no**
business logic — it renders whatever the `enigma` daemon reports over its
local HTTP-over-unix-socket API (`GET /v1/state`) and posts actions back to
the same socket. The daemon API is implemented in Go at
[`cli/internal/api`](../cli/internal/api).

## Layout

```
center/
  src/
    main.cpp            QGuiApplication + QQmlApplicationEngine, exposes `daemon`
    daemonclient.h/.cpp DaemonClient: GET /v1/state -> QML properties
  qml/
    Main.qml            ApplicationWindow: TabBar + StackLayout, 3s poll
    tabs/
      RuntimesTab.qml   installed language runtimes (mise)
      ServicesTab.qml   enigma-managed services + status/port
      AITab.qml         installed models + combined size
      ProjectsTab.qml   dev projects + .test URLs + ports
      WindowsTab.qml    Bottles / Wine Tier 1 (DirectX/VC++ redists)
      SystemTab.qml     GPU, snapshots/rollback, staged updates, doctor
  tests/
    tst_daemonclient.cpp  headless test vs a mocked daemon (SPEC §18)
```

## Building

Requires Qt 6.5+ (`Quick`, `Network`, `Test`).

```sh
cmake -S center -B center/build -G Ninja
cmake --build center/build
./center/build/enigma-center
```

Point it at a specific daemon (default is
`$XDG_RUNTIME_DIR/enigma-center.sock`):

```sh
ENIGMA_CENTER_URL=http://127.0.0.1:9999 ./center/build/enigma-center
```

## Testing

The `DaemonClient` networking/parsing logic has a headless test that runs
against an in-process mock HTTP daemon — no display or real daemon needed
(SPEC §18: "Center must pass headless tests against a mocked enigma daemon").
This test is standalone-buildable so CI validates it without the full QML
toolchain:

```sh
cmake -S center/tests -B center/tests/build -G Ninja
cmake --build center/tests/build
ctest --test-dir center/tests/build --output-on-failure
```
