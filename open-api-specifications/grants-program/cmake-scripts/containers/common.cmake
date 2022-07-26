include(${CMAKE_CURRENT_LIST_DIR}/find_dependencies.cmake)

# if a tag suffix is not given, default to the empty string
if (NOT DEFINED TAG_SUFFIX)
    set(TAG_SUFFIX "")
endif()

if (NOT DEFINED PROD)
    set(PROD OFF)
endif()

# determine the commit we are on
execute_process(
    COMMAND ${Git_EXECUTABLE} rev-parse --short HEAD
    OUTPUT_VARIABLE commit
    OUTPUT_STRIP_TRAILING_WHITESPACE
)

# Determine the tag for the current commit, if there is one. Note that only
# annotated commits are considered, and if there are multiple annotated tags on
# this commit, only one will be chosen.
execute_process(
    COMMAND ${Git_EXECUTABLE} describe --exact-match HEAD
    OUTPUT_VARIABLE git_tag
    RESULT_VARIABLE git_describe_result
    OUTPUT_STRIP_TRAILING_WHITESPACE
    ERROR_QUIET
)
if (git_describe_result EQUAL 0)
    set(on_git_tag ON)
else()
    set(on_git_tag OFF)
endif()


# determine the branch we are on
execute_process(
    COMMAND ${Git_EXECUTABLE} rev-parse --abbrev-ref HEAD
    OUTPUT_VARIABLE branch
    OUTPUT_STRIP_TRAILING_WHITESPACE
)

# If the branch is 'HEAD', then we are in a detached HEAD state and it doesn't
# make sense to push a branch-tagged image. If a branch with the same name as a
# tag has been checked out, then the branch ref name will be
# 'heads/<tag_name>'. This is an invalid tag for a Docker image and it isn't
# needed anyway, so in this case we simply don't push a branch-tagged image.
if (NOT "${branch}" STREQUAL "HEAD" AND NOT "${branch}" MATCHES "^heads/.*$")
    set(on_branch ON)
else()
    set(on_branch OFF)
endif()

# Get a list of git tags sorted in descending order by version
execute_process(
    COMMAND ${Git_EXECUTABLE} -c versionsort.suffix=- tag --sort=-version:refname
    OUTPUT_VARIABLE git_tag_versions
    OUTPUT_STRIP_TRAILING_WHITESPACE
)

# Convert string output to a CMake list
string(REGEX REPLACE "[ \t\r\n]+" ";" sorted_versions "${git_tag_versions}")

# Determine if we are on the latest version
set(on_latest_version OFF)
list(LENGTH sorted_versions version_count)
if (version_count GREATER 0)
    list(GET sorted_versions 0 latest_version)
    if (on_git_tag AND "${git_tag}" STREQUAL "${latest_version}")
        set(on_latest_version ON)
    endif()
endif()

if (PROD)
    set(image "registry.matchfund.app/prod/${IMAGE}")
else()
    set(image "registry.matchfund.app/${IMAGE}")
endif()

set(commit_tagged_image "${image}:${commit}${TAG_SUFFIX}")
set(branch_tagged_image "${image}:${branch}${TAG_SUFFIX}")
set(latest_tagged_image "${image}:latest${TAG_SUFFIX}")
set(git_tagged_image    "${image}:${git_tag}${TAG_SUFFIX}")
