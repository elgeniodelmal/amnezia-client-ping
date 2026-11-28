# Amnezia Ping V4 native Android integration

find_program(PING_GO_EXECUTABLE NAMES go REQUIRED)

set(PING_GO_SOURCE
    "${CMAKE_CURRENT_LIST_DIR}/../../experiments/ping-v4"
)

if(CMAKE_ANDROID_ARCH_ABI STREQUAL "arm64-v8a")
    set(PING_GO_ARCH arm64)
    set(PING_CLANG_TRIPLE aarch64-linux-android)
    set(PING_GO_ARM "")
elseif(CMAKE_ANDROID_ARCH_ABI STREQUAL "armeabi-v7a")
    set(PING_GO_ARCH arm)
    set(PING_CLANG_TRIPLE armv7a-linux-androideabi)
    set(PING_GO_ARM 7)
elseif(CMAKE_ANDROID_ARCH_ABI STREQUAL "x86_64")
    set(PING_GO_ARCH amd64)
    set(PING_CLANG_TRIPLE x86_64-linux-android)
    set(PING_GO_ARM "")
elseif(CMAKE_ANDROID_ARCH_ABI STREQUAL "x86")
    set(PING_GO_ARCH 386)
    set(PING_CLANG_TRIPLE i686-linux-android)
    set(PING_GO_ARM "")
else()
    message(FATAL_ERROR
        "Ping V4: unsupported Android ABI"
    )
endif()

get_filename_component(
    PING_CLANG_DIR
    "${CMAKE_C_COMPILER}"
    DIRECTORY
)

set(PING_CC
    "${PING_CLANG_DIR}/${PING_CLANG_TRIPLE}${APP_ANDROID_MIN_SDK}-clang"
)

set(PING_CXX
    "${PING_CLANG_DIR}/${PING_CLANG_TRIPLE}${APP_ANDROID_MIN_SDK}-clang++"
)

if(NOT EXISTS "${PING_CC}")
    message(FATAL_ERROR
        "Ping V4 compiler not found: ${PING_CC}"
    )
endif()

set(PING_OUTPUT_DIR
    "${CMAKE_CURRENT_BINARY_DIR}/ping-v4/${CMAKE_ANDROID_ARCH_ABI}"
)

file(MAKE_DIRECTORY "${PING_OUTPUT_DIR}")

set(PING_LIBRARY
    "${PING_OUTPUT_DIR}/libamneziaping.so"
)

execute_process(
    COMMAND "${CMAKE_COMMAND}" -E env
        "GOOS=android"
        "GOARCH=${PING_GO_ARCH}"
        "GOARM=${PING_GO_ARM}"
        "CGO_ENABLED=1"
        "GOTOOLCHAIN=auto"
        "CC=${PING_CC}"
        "CXX=${PING_CXX}"
        "${PING_GO_EXECUTABLE}" build
        -trimpath
        -buildmode=c-shared
        "-ldflags=-s -w"
        -o "${PING_LIBRARY}"
        ./cmd/androidbridge
    WORKING_DIRECTORY "${PING_GO_SOURCE}"
    RESULT_VARIABLE PING_BUILD_RESULT
)

if(NOT PING_BUILD_RESULT EQUAL 0)
    message(FATAL_ERROR
        "Ping V4 native library build failed"
    )
endif()

if(NOT EXISTS "${PING_LIBRARY}")
    message(FATAL_ERROR
        "Ping V4 native library missing"
    )
endif()

set_property(
    TARGET ${PROJECT}
    APPEND PROPERTY QT_ANDROID_EXTRA_LIBS
    "${PING_LIBRARY}"
)

message(STATUS
    "Ping V4 library ready: ${PING_LIBRARY}"
)
