# Creates container-build and container-push targets. Expects a Dockerfile in
# the project root directory.

include(ProcessorCount)
include(${CMAKE_CURRENT_LIST_DIR}/../util/strings.cmake)

find_program(Git_EXECUTABLE git)
find_program(Docker_EXECUTABLE docker)

# check if we're in a repo
execute_process(
    COMMAND ${Git_EXECUTABLE} rev-parse --is-inside-work-tree
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    RESULT_VARIABLE is_repo
    OUTPUT_QUIET
    ERROR_QUIET
)

if (NOT Docker_EXECUTABLE)
    message(WARNING "Could not find docker executable. Targets for building container images will not be created")
elseif (NOT Git_EXECUTABLE)
    message(WARNING "Could not find git. Targets for building container images will not be created")
elseif (NOT "${is_repo}" STREQUAL "0")
    message(WARNING "Not in a git repo. Targets for building container images will not be created")
else()
    # if the number of build jobs is not given, set a reasonable default
    if (NOT DEFINED CONTAINER_BUILD_JOBS)
        ProcessorCount(cores)
        if (cores EQUAL 0)
            message(WARNING "Could not determine processor core count, defaulting to 1")
            set(cores 1)
        endif()
    endif()

    set(CONTAINER_BUILD_JOBS "${cores}" CACHE STRING "Number of build jobs, passed as a build arg.")
    set(CONTAINER_IMAGE ${CMAKE_PROJECT_NAME} CACHE STRING "Name of the image to build (automatically prefixed by the registry name).")
    set(CONTAINER_BUILD_OPTIONS "" CACHE STRING "List of extra options to pass to the container engine for the build.")
    set(PROD OFF CACHE BOOL "Whether to put images under the prod namespace.")

    string_to_list(build_options "${CONTAINER_BUILD_OPTIONS}" SEMICOLON_ESCAPE)

    # build target
    add_custom_target(container-build
        COMMAND ${CMAKE_COMMAND}
            -DCONTAINERFILE=Dockerfile
            -DIMAGE=${CONTAINER_IMAGE}
            -DBUILD_JOBS=${CONTAINER_BUILD_JOBS}
            -DBUILD_OPTIONS="${build_options}"
            -DPROD=${PROD}
            -P ${CMAKE_CURRENT_SOURCE_DIR}/build-image.cmake
    )

    # push target
    add_custom_target(container-push
        COMMAND ${CMAKE_COMMAND}
            -DIMAGE=${CONTAINER_IMAGE}
            -DPROD=${PROD}
            -P ${CMAKE_CURRENT_SOURCE_DIR}/push-image.cmake
    )

    # before we can push images we must have built them
    add_dependencies(container-push container-build)
endif()
