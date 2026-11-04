---
title: Quick Start
description: Install tuna, open a project and start your first agent session.
---

# Quick Start

Install tuna, open a project, and run your first agent in a session that keeps going after you quit.

## Install

Download the latest tuna from the [releases page](https://github.com/smallnest/tuna/releases). On macOS, drag **tuna** into your Applications folder; on Windows, run the **tuna** installer, which needs no administrator rights.

On Windows, sessions run in PowerShell 7 when it is installed, else Windows PowerShell; choose another shell, such as `cmd` or Git Bash's `bash`, in **Settings › Shell**. tuna needs Windows 10 1809 or later, for its pseudo consoles.

### Install the command line tool

Choose **tuna › Install Command Line Tool…** to put `tuna` on your `PATH`: a link in `/usr/local/bin` on macOS, in `~/.local/bin` on Linux, and on Windows tuna's `bin` folder added to your user `PATH`, for terminals you open from then on. Run it in any terminal to check:

```sh
tuna help
```

Typing `tuna` alone in a terminal opens the app.

## Open a project

Press <kbd>⌘</kbd> <kbd>O</kbd> and choose a folder. tuna adds it to the sidebar as a project and opens a session there.

You can also start from the command line. A new session goes into the project that matches its directory:

```sh
cd ~/code/my-app
tuna new --focus
```

## Start an agent

Run your agent in the session the way you would in any terminal:

```sh
claude
```

Or give it a session of its own with a name you'll recognize in the sidebar:

```sh
tuna new -n fix-login -- claude
```

While the agent works, its session shows as **working**. When it stops to ask you something, it turns to **waiting** and tuna notifies you. Press <kbd>⌘</kbd> <kbd>J</kbd> to jump to it.

## Quit without losing anything

Press <kbd>⌘</kbd> <kbd>Q</kbd>. Your sessions keep running in tunad. Open tuna again and they're all there, with their scrollback.

## Next steps

- [Panes](./panes.md): split a session to run a dev server beside your agent
- [Coding Agents](./agents.md): run `tuna skill setup` so agents learn to use tuna
