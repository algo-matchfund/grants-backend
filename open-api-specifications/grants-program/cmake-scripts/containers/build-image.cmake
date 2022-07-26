# this script expects to receive the following input:
# CONTAINERFILE:        The path to the container file to use
# IMAGE:                The name of the image to build (automatically prefixed by the registry name)
# TAG_SUFFIX:           A suffix to append to the tag (optional, defaults to the empty string)
# CONTAINER_ENGINE:     The container engine to use (optional, defaults to docker)
# BUILD_CONTEXT:        The build context directory (optional, defaults to the repo root directory)
# BUILD_JOBS:           The number of build jobs for builders in the container files to use (optional, defaults to the number of cores)
# BUILD_OPTIONS:        List of extra options to pass to the container engine for the build (optional)
# PROD:                 Whether to tag the image into the prod/ namespace (optional, defaults to false)

include(ProcessorCount)

include(${CMAKE_CURRENT_LIST_DIR}/../util/processes.cmake)
include(${CMAKE_CURRENT_LIST_DIR}/common.cmake)

# if a build context is not given, default to the repo root directory
if (NOT DEFINED BUILD_CONTEXT)
    execute_process(
        COMMAND ${Git_EXECUTABLE} rev-parse --show-toplevel
        OUTPUT_VARIABLE BUILD_CONTEXT
        OUTPUT_STRIP_TRAILING_WHITESPACE
    )
endif()

# if the number of build jobs is not given, default to the number of cores if
# possible, else 1
if (NOT DEFINED BUILD_JOBS)
    ProcessorCount(BUILD_JOBS)
    if (BUILD_JOBS EQUAL 0)
        message(WARNING "Could not determine processor core count, defaulting to 1")
        set(BUILD_JOBS 1)
    endif()
endif()

# build the image and tag it with the branch name and commit id
safe_execute_process(
    COMMAND ${CONTAINER_ENGINE} build
        -t ${commit_tagged_image}
        --build-arg build_jobs=${BUILD_JOBS}
        --build-arg git_commit=${commit}
        --build-arg git_branch=${branch}
        ${BUILD_OPTIONS}
        -f ${CONTAINERFILE}
        .
    WORKING_DIRECTORY ${BUILD_CONTEXT}
    FAIL_MESSAGE "Container build for ${commit_tagged_image} failed"
)

# do we have a proper branch name?
if (on_branch)
    # tag the image to the branch name
    safe_execute_process(
        COMMAND ${CONTAINER_ENGINE} tag
            ${commit_tagged_image}
            ${branch_tagged_image}
        FAIL_MESSAGE "Could not create tag ${branch_tagged_image} for image ${commit_tagged_image}"
    )
    message(STATUS "Successfully tagged ${branch_tagged_image}")
else()
    message(STATUS "Repo is not on a branch, will not create a branch tag")
endif()

# does this commit have an annotated tag?
if (on_git_tag)
    safe_execute_process(
        COMMAND ${CONTAINER_ENGINE} tag
            ${commit_tagged_image}
            ${git_tagged_image}
       FAIL_MESSAGE "Could not create tag ${git_tagged_image} for image ${commit_tagged_image}"
    )
    message(STATUS "Successfully tagged ${git_tagged_image}")

    if (on_latest_version)
        # tag the image with 'latest'
        safe_execute_process(
            COMMAND ${CONTAINER_ENGINE} tag
                ${commit_tagged_image}
                ${latest_tagged_image}
            FAIL_MESSAGE "Could not create tag ${latest_tagged_image} for image ${commit_tagged_image}"
        )
        message(STATUS "Successfully tagged ${latest_tagged_image}")
    endif()
else()
    message(STATUS "No annotated git tag found for commit ${commit}, will not create a git-tagged image")
endif()
