# Wrapper around execute_process() that will fail if the given command fails.
# This function expects:
#  COMMAND              Same meaning as for execute_process()
#  FAIL_MESSAGE         The message to print if COMMAND returns non-zero
#  WORKING_DIRECTORY    Same meaning as for execute_process() (optional, defaults to CMAKE_SOURCE_DIR)
function(safe_execute_process)
    set(oneValueArgs "WORKING_DIRECTORY" "FAIL_MESSAGE")
    set(multiValueArgs "COMMAND")
    cmake_parse_arguments(S "" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})

    # check arguments
    if (NOT S_FAIL_MESSAGE)
        message(FATAL_ERROR "FAIL_MESSAGE is unset")
    endif()
    if (NOT S_WORKING_DIRECTORY)
        set(S_WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})
    endif()
    if (NOT S_COMMAND)
        message(FATAL_ERROR "COMMAND is unset")
    endif()

    # run command
    execute_process(
        COMMAND ${S_COMMAND}
        WORKING_DIRECTORY ${S_WORKING_DIRECTORY}
        RESULTS_VARIABLE results
    )

    # if the process failed, quit
    foreach(result "${results}")
        if (NOT "${result}" STREQUAL "0")
            message(FATAL_ERROR "${S_FAIL_MESSAGE}")
        endif()
    endforeach()
endfunction()
