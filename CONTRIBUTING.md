# How to contribute

We'd love to accept your patches and contributions to this project. There are just a few small guidelines you need to follow.

1. Code should be `go fmt` formatted.
2. Exported types, constants, variables and functions should be documented.
3. Changes must be covered with tests.
4. All tests must pass constantly `go test .`.

## Versioning

Direct Decisions Go client follows semantic versioning. New functionality should be accompanied by increment to the minor version number.

## Releasing

Any code which is complete, tested, reviewed, and merged to master can be released.

1. Update the `version` number in `directdecisions.go`.
2. Make a pull request with these changes.
3. Once the pull request has been merged, visit [https://github.com/directdecisions/client-go/releases](https://github.com/directdecisions/client-go/release) and click `Draft a new release`.
4. Update the `Tag version` and `Release title` field with the new Direct Decisions Go client version. Be sure the version has a `v` prefixed in both places, e.g. `v1.25.0`.
5. Publish the release.
