
# Release Process

The release process is automated when merging a Pull Request.

## How to trigger a release

1. Create a Pull Request.
1. Attach a label [`bump:patch`, `bump:minor`, or `bump:major`]. Cloudgrep uses [haya14busa/action-bumpr](https://github.com/haya14busa/action-bumpr).
1. [The release workflow](.github/workflows/release.yml) automatically tags a
   new version depending on the label and create a new release on merging the
   Pull Request.

If you do not want to create a release for a given PR, do not attach a bump label to it.
