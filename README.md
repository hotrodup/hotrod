# Hot Rod (CLI)

[![Gobuild Download](http://gobuild.io/badge/github.com/hotrodup/hotrod/downloads.svg)](http://gobuild.io/github.com/hotrodup/hotrod)

<img align="right" height="250" src="http://i.imgur.com/2gYTTc8.png" alt="Hot Rod">

> :checkered_flag: Turbocharge your Node.js development cycle.

Hot Rod is a CLI that provisions a remote development server on Google Cloud and beams up your source code to the server after every local file change.  The traditional development cycle (edit code, preview locally, commit changes, deploy app, refresh webpage, and verify changes) is turbocharged.  With Hot Rod, just edit your source, hit save, and--that's it! Your code is running live on a remote server.

[http://hotrodup.com](http://hotrodup.com)

## Key Features

- **Blazing Fast**: Local changes to source appear online in under `100ms`.  The delay is impercetible on a fast connection.
- **Auto-Refresh**: No need to hit refresh in your browser, Hot Rod automatically reloads your page after every file change.
- **Short URL**: Every development server gets a short URL for easy sharing.
- **Local Editors**: No need to use a clumsy web-based IDE.  Use the tools you love to edit source locally, and preview the changes remotely.

## Dependencies

- [Google Cloud SDK](https://cloud.google.com/sdk/)
- [Git](http://git-scm.com/)

## Installation

1. Download the correct binary:

  |  | Linux | OSX | Windows |
  |:------:|----------------------------------------------------------------------------------------------------------------------------|-----|---------|
  | 32-bit | [hotrod-linux-386.tar.gz](http://gobuild3.qiniudn.com/github.com/hotrodup/hotrod/branch-v-master/hotrod-linux-386.tar.gz) | [hotrod-darwin-386.zip](http://gobuild3.qiniudn.com/github.com/hotrodup/hotrod/branch-v-master/hotrod-darwin-386.zip) | [hotrod-windows-386.zip](http://gobuild3.qiniudn.com/github.com/hotrodup/hotrod/branch-v-master/hotrod-windows-386.zip) |
  | 64-bit | [hotrod-linux-amd64.tar.gz](http://gobuild3.qiniudn.com/github.com/hotrodup/hotrod/branch-v-master/hotrod-linux-amd64.tar.gz ) | [hotrod-darwin-amd64.zip](http://gobuild3.qiniudn.com/github.com/hotrodup/hotrod/branch-v-master/hotrod-darwin-amd64.zip) | [hotrod-windows-amd64.zip](http://gobuild3.qiniudn.com/github.com/hotrodup/hotrod/branch-v-master/hotrod-windows-amd64.zip) |

2. Unzip the package
  ```sh
  $ unzip hotrod-darwin-amd64.zip
  ```

3. Move the binary to your `bin`
  ```sh
  $ sudo mv hotrod /usr/local/bin
  ```

4. Verify that Hot Rod is correctly installed
  ```sh
  $ hotrod
  ```

## Usage

1. Create a new project
  ```sh
  $ hotrod create my-project
  ```

2. `cd` into the project directory
  ```sh
  $ cd my-project
  ```

3. Beam up your source code!
  ```sh
  $ hotrod up
  ```

4. Edit changes and preview them in the browser window that pops up.

## Help

```
Hot Rod (v 0.0.1)

usage: hotrod <command> [<flags>] [<args> ...]

Turbocharge your Node.js development cycle

Flags:
  --help  Show help.

Commands:
  help [<command>]
    Show help for a command.

  create <name>
    Create a new Hot Rod app.

  up
    Beam up the source to your Hot Rod app.
```
