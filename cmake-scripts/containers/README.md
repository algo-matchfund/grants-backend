Containers
==========
We try to keep our container usage consistent with regard to naming and
pushing. This directory contains some helper CMake scripts to build and push
images according to our best practices. For building images with
`build-image.cmake`:
- An image tagged with the current commit is always built.
- If the source is on a proper branch (instead of detached), the above image
  will also be tagged with the branch name.
- If an annotated git tag is present on the current commit, the image will also
  be tagged with the git tag name. Note that if there are multiple tags on a
  single commit, only one will be used to make the tag.

For pushing images with `push-image.cmake`:
- Pushing images to the registry is allowed unless there are uncommitted
  changes.
- If there are no uncommitted changes, pushing the image tagged with the current
  commit is allowed.
- Pushing the image tagged with a branch name is also allowed if a branch is
  checked out (instead of being detached).

A sample `CMakeLists.txt` is included which defines `container-build` and
`container-push` targets that expect a `Dockerfile` at the source root. For most
repos it should be enough to add this directory as a CMake subdirectory. For
example, if your project does not use CMake for anything else, a plausible
top-level `CMakeLists.txt` would be:
```cmake
cmake_minimum_required(VERSION 123)

project(foo NONE)

add_subdirectory(cmake-scripts/containers)
```

Repos with more complicated images (multiple `Dockerfile`'s, etc) will likely
need to use `build-image.cmake` and `push-image.cmake` directly, see their code
for an explanation of the variables they expect and `CMakeLists.txt` for an
example usage.

Note that these scripts are intentionally generic with regard to the container
engine used. e.g., they use `CONTAINERFILE` instead of `DOCKERFILE` and
`CONTAINER_ENGINE` instead of hardcoding `docker`. This is to make it easier to
migrate to a different container engine in the future if we decide to.

## Writing `Dockerfile`'s/`Containerfile`'s
The `build-image.cmake` script passes certain build args that we recommend using
where possible:
- `build_jobs`: This is the number of build jobs that a builder/compiler in the
  container file should use. For example if a C++ program is being compiled,
  prefer `make -j${build_jobs}` instead of just `make` or a hardcoded `make -j4`.
- `git_commit` and `git_branch`: The current git commit and branch. We recommend
  that all container files use these for labels:
  ```Dockerfile
  ARG git_commit="<unknown>"
  ARG git_branch="HEAD"

  LABEL git.commit=${git_commit}
  LABEL git.branch=${git_branch}
  ```

  This makes it easy to discover the commit and branch an image is built from.
