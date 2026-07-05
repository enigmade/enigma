cmake_minimum_required(VERSION 3.21)

project(EnigmaCenter VERSION 0.1.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTOMOC ON)

# Enigma Center: thin Qt/QML client over the enigma daemon API (SPEC §18).
# Zero duplicated logic — every panel renders api.State fetched from the
# daemon's GET /v1/state, and actions POST back to the same socket.

find_package(Qt6 REQUIRED COMPONENTS Quick Network Test)

qt_standard_project_setup(REQUIRES 6.5)

qt_add_executable(enigma-center
    src/main.cpp
    src/daemonclient.cpp
    src/daemonclient.h
)

qt_add_qml_module(enigma-center
    URI EnigmaCenter
    VERSION 1.0
    QML_FILES
        qml/Main.qml
        qml/tabs/RuntimesTab.qml
        qml/tabs/ServicesTab.qml
        qml/tabs/AITab.qml
        qml/tabs/ProjectsTab.qml
        qml/tabs/WindowsTab.qml
        qml/tabs/SystemTab.qml
)

target_link_libraries(enigma-center PRIVATE Qt6::Quick Qt6::Network)

# Headless CI tests against a mocked daemon (SPEC §18: "Center must pass
# headless tests against a mocked enigma daemon").
enable_testing()
add_subdirectory(tests)
