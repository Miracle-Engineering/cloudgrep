<p align="center"><img src="https://user-images.githubusercontent.com/855699/172711194-43330c43-c13e-4b04-9e4a-11eabe8cf850.png" width="250"><br/>
Cloud Asset Explorer</p>

<p align="center">
  <a href="https://github.com/run-x/cloudgrep/releases/latest">
    <img src="https://img.shields.io/github/release/run-x/cloudgrep.svg" alt="Current Release" />
  </a>
  <a href="https://github.com/run-x/cloudgrep/actions/workflows/checks.yml">
    <img src="https://github.com/run-x/cloudgrep/actions/workflows/checks.yml/badge.svg" alt="Tests" />
  </a>

  <a href="http://www.apache.org/licenses/LICENSE-2.0.html">
    <img src="https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg" alt="License" />
  </a>

  <img src="https://img.shields.io/github/commit-activity/w/run-x/cloudgrep.svg?style=plastic" alt="Commit Activity" />

</p>
<p align="center">
<a href="https://slack.cloudgrep.dev">
    Slack Community
  </a>
  </p>

# What is Cloudgrep?
Cloudgrep is an asset explorer for cloud resources. It shows you everything you're running on your cloud account(s) and enables you to slice and dice these based on tags and properties. It is a UI tool built on open source technologies and runs completely client side (so no data leaves your machine).
<p align="center">
<img width="820" alt="Screenshot" src="https://user-images.githubusercontent.com/855699/175440360-d6e759d0-ecd6-4a36-889c-b329563979db.png">
</p>



### Why use Cloudgrep?
Cloudgrep's goal is to help your team ensure every resource follows consistent tagging schema. It helps identify missing tags, misspellings and unowned resources. Consistent tagging leads to better cost attribution and faster incident resolution.

Additionally, Cloudgrep is a great tool to visualize all your resources in one place - across regions, accounts and clouds. 

If you'd like to try it out or have any questions - feel free to join our [Slack](https://slack.cloudgrep.dev/) or download it from the [releases](https://github.com/run-x/cloudgrep/releases) section.


<p align="center">
  <a href="https://www.youtube.com/watch?v=Ip-lY9x7bh4"><img width="478" alt="Group 2" src="https://user-images.githubusercontent.com/855699/175441143-5834e9f2-8e23-471a-95be-7b3388d1f455.png"></a>
  </br>
  <span><i>
Demo video</i></span>

</p>

# Features
* Cross-platform support OSX/Linux/Windows 32/64-bit
* Simple installation (distributed as a single binary)
* Zero dependencies
* Supports AWS. (If you'd like GCP/Azure support, do let us know by filing an issue!)
* Supports for major AWS resources (like EC2, RDS, S3, and many others - please file an issue if something is missing!)

# Installation

- [Precompiled binaries](https://github.com/run-x/cloudgrep/releases) for supported
operating systems are available.

# Basic Usage

Cloudgrep works by using the cloud provider credentials (e.g. AWS) found in your terminal to scan your cloud for data about
your existing resources. For it to work, make sure you set such credentials properly (see
[here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) for AWS). These, and a stable
internet connection, are its only prerequisites.

**NOTE: Cloudgrep only needs ReadOnly credentials-- it creates nothing, it modifies nothing. Moreover, it will
do a best effort scan based on its permission, so the user does not need to have read access to all resources.**

Once downloaded, (and assuming you have AWS credentials properly configured in your CLI) you can start using cloudgrep
immediately by executing the binary (no args needed) on your cli:

```bash
./cloudgrep
```

Cloudgrep will then run with the default behavior which does the following:

1. Scan your current cloud account for global resources and resources on your currently configured AWS region
2. Launch the webapp

## Arguments
You can easily pass cli arguments to cloudgrep for customized behavior, such as multiple/different regions to scan,
what port to serve the webapp on, etc... The cli arguments are all fully documented under the cli's `help` option.
To view documentation for them, simply add the `--help` flag like so:

```bash
./cloudgrep --help
```

## Advanced Usage
Cloudgrep's behavior can further be configured via a user-inputted config yaml. Configs are then resolved at runtime by
considering the cli arguments, the user-passed config  yaml, and the defaults in that order of precedence.

The config yaml can be passed in by using the `-c` or `--config` flag as follows:

```bash
cloudgrep -c my_config.yaml
```
The path is relative to the current working directory. Cloudgrep expects the follow possible values in the yaml
(you do not need to markdown all if passing the file as it will always try to default to the original behavior):

```yaml
# Config represents all the user-configurable settings for cloudgrep. One such structure is loaded at runtime and
# is populated through the cli arguments, user-provider config file, or a preset default, with values resolved
# in that order of precedence. To see the default, please refer to
# https://github.com/run-x/cloudgrep/blob/main/pkg/config/config.yaml

# web represents the specs cloudgrep uses for creating the webapp server
web:
  # host is the host the server is running as
  host: localhost
  # port is the port the server is running in
  port: 8080
  # prefix is the url prefix the server uses
  prefix: "/"
  # skipOpen determines whether to automatically open the webui on startup
  skipOpen: false

# datastore represents the specs cloudgrep uses for creating and/or connecting to the datastore/database used.
datastore:
  # type is the kind of datastore to be used by cloudgrep (currently only supports SQLite)
  type: sqlite
  #  skipRefresh determines whether to refresh the data (i.e. scan the cloud) on startup.
  skipRefresh: false
  # dataSourceName is the Type-specific data source name or uri for connecting to the desired data source
  dataSourceName: "~/cloudgrep_data.db"

# providers represents the cloud providers cloudgrep will scan w/ the current credentials
providers:
  - cloud: aws # cloud is the type of the cloud provider (currently only AWS is supported)
    # regions is the list of different regions within the cloud provider to scan
    # The special "all" region can be specified by itself to scan all available regions
    regions: [us-east-1, global]
```

# Development

We love user contributions! Check out our [Dev guide](https://github.com/run-x/cloudgrep/blob/main/DEVELOP.md) to get started.

# Important Resources
* [The Cloudgrep Team](https://www.runx.dev/about)
* [Check Out our Blog](https://blog.runx.dev/)
