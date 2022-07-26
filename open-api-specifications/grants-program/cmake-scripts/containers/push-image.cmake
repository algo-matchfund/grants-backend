# this script expects to receive the following input:
# IMAGE:                The image to push
# TAG_SUFFIX:           A suffix to append to the tag (optional, defaults to the empty string)
# CONTAINER_ENGINE:     The container engine to use (optional, defaults to docker)
# PROD:                 Whether to push images in the prod/ namespace (optional, defaults to false)

include(${CMAKE_CURRENT_LIST_DIR}/../util/processes.cmake)
include(${CMAKE_CURRENT_LIST_DIR}/common.cmake)

# check if there are uncommitted local changes
execute_process(
    COMMAND ${Git_EXECUTABLE} diff --exit-code --quiet
    RESULT_VARIABLE uncommitted
)

# never push if there are uncommitted changes
if (uncommitted)
    message(FATAL_ERROR "Pushing images failed, there are uncommitted changes")
endif()

# push the image
safe_execute_process(
    COMMAND ${CONTAINER_ENGINE} push ${commit_tagged_image}
    FAIL_MESSAGE "Could not push ${commit_tagged_image}")
message(STATUS "Successfully pushed ${commit_tagged_image}")

# do we have a branch name as well?
if (on_branch)
    # push the tagged branch image
    safe_execute_process(
        COMMAND ${CONTAINER_ENGINE} push ${branch_tagged_image}
        FAIL_MESSAGE "Could not push ${branch_tagged_image}")
    message(STATUS "Successfully pushed ${branch_tagged_image}")
endif()

# does this commit have an annotated tag?
if (on_git_tag)
    safe_execute_process(
        COMMAND ${CONTAINER_ENGINE} push ${git_tagged_image}
        FAIL_MESSAGE "Could not push ${git_tagged_image}"
    )
    message(STATUS "Successfully pushed ${git_tagged_image}")

    # push the 'latest' tag this is the latest version
    if (on_latest_version)
        safe_execute_process(
            COMMAND ${CONTAINER_ENGINE} push ${latest_tagged_image}
            FAIL_MESSAGE "Could not push ${latest_tagged_image}")
        message(STATUS "Successfully pushed ${latest_tagged_image}")
    endif()
endif()
