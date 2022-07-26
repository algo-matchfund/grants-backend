# find git
find_program(Git_EXECUTABLE git)
if (NOT Git_EXECUTABLE)
    message(FATAL_ERROR "Could not find git")
endif()

# if a container engine is not given, default to docker
if (NOT DEFINED CONTAINER_ENGINE)
    find_program(CONTAINER_ENGINE docker)
    if (NOT CONTAINER_ENGINE)
        message(FATAL_ERROR "Could not find docker")
    endif()
endif()
