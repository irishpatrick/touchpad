cmake_minimum_required(VERSION 3.20)

project(driver C)

option(MINGW "build for mingw" OFF)

if (${MINGW})
    set(sources
        driver.h
        driver_windows.c
    )
else ()
    set(sources
        driver.h
        driver_linux.c
    )

    find_package(PkgConfig REQUIRED)
    pkg_check_modules(EVDEV libevdev REQUIRED)

    set(incl ${EVDEV_INCLUDE_DIRS})
    set(libs ${EVDEV_LIBRARIES})
endif ()

add_library(${PROJECT_NAME} STATIC ${sources})
target_include_directories(${PROJECT_NAME} PRIVATE ${incl})
target_link_libraries(${PROJECT_NAME} PRIVATE ${libs})

install(TARGETS ${PROJECT_NAME}
        DESTINATION ${CMAKE_SOURCE_DIR}/artifacts/)
