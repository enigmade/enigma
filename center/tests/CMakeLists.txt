# Headless test of DaemonClient against a mocked enigma daemon (SPEC §18).
#
# This is written to be buildable either as part of the top-level Center
# build OR standalone (so CI can validate the client logic without pulling
# in the full Qt Quick/QML toolchain). When configured directly, it defines
# its own project.
if (NOT DEFINED PROJECT_NAME)
    cmake_minimum_required(VERSION 3.21)
    project(EnigmaCenterTests LANGUAGES CXX)
    set(CMAKE_CXX_STANDARD 17)
    set(CMAKE_CXX_STANDARD_REQUIRED ON)
    set(CMAKE_AUTOMOC ON)
    enable_testing()
endif()

find_package(Qt6 REQUIRED COMPONENTS Test Network)

add_executable(tst_daemonclient
    tst_daemonclient.cpp
    ../src/daemonclient.cpp
    ../src/daemonclient.h
)

target_include_directories(tst_daemonclient PRIVATE ../src)
target_link_libraries(tst_daemonclient PRIVATE Qt6::Test Qt6::Network)

add_test(NAME tst_daemonclient COMMAND tst_daemonclient)
