---
title: Command Line
description: Drive tuna sessions and panes from a shell or a script with the tuna CLI.
---

# Command Line

The `tuna` command drives sessions from any shell: yours, a script's or an agent's. Install it with **tuna › Install Command Line Tool…**.

## Targets

Most commands take a **target**: a session or a pane. A target can be:

- an id, or any unique prefix of it
- a session's name
- a session, which stands for its focused pane

Inside tuna, `$TUNA_SESSION` and `$TUNA_PANE` name the session and pane a command runs in, and they're the default target.

## Sessions and panes

```sh
tuna ls                          # every session and its panes
tuna new -n api -d ~/code/api    # start a session, print its id
tuna new -- npm run dev          # start one running a command
tuna split --down -- npm test    # a pane beside this one
tuna split --url localhost:3000  # a browser pane showing a page
tuna rename api backend          # rename a session
tuna focus backend               # show it in the app
tuna close backend               # end a session (or, given a pane, the pane)
tuna close fix-login --remove-worktree   # and remove its own worktree
tuna restart <pane>              # start an exited pane's shell again
```

`tuna new` options:

| Flag                        | Meaning                                         |
| --------------------------- | ----------------------------------------------- |
| `-p`, `--project <name\|id>` | The project to put it in (default: by directory) |
| `-n`, `--name <name>`       | Name the session                                |
| `-d`, `--dir <dir>`         | Start it there (default: the current directory) |
| `-f`, `--focus`             | Show it in the app                              |
| `-w`, `--wait`              | Wait for its command to end, and exit with its status |
| `-u`, `--url <url>`         | Open a browser showing the page instead of a shell |
| `--worktree <branch>`       | Start it in a git worktree of its own, on a new branch ([more](./changes.md#sessions-in-their-own-worktree)) |
| `--base <ref>`              | The branch the worktree's starts from (default: the default branch) |

## Typing and reading

```sh
tuna send api "npm run migrate" --enter   # type text, then press Enter
tuna send api -e 'y\n'                    # interpret \n \r \t \e \xHH escapes
tuna key api ctrl+c                       # press keys
tuna read api -n 40                       # the last 40 lines on screen
tuna read api -a                          # the screen with its scrollback
```

`tuna key` understands `enter`, `tab`, `escape`, `up`, `down`, `left`, `right`, `backspace`, and modifiers such as `ctrl+c` or `ctrl+d`.

## Waiting

```sh
tuna wait api                             # until it stops working
tuna wait api --for waiting,error --timeout 10m
tuna prompt review "summarize the diff"   # send + Enter, wait, print the screen
```

`tuna wait` prints the pane's status when it returns. `tuna prompt` sends text and Enter, waits until the pane stops working, then prints its screen. It's the quickest way to ask another agent something.

### Waiting for a command

A command such as `make test` doesn't show as *working*, so to wait for one to finish, start it with `--wait`. `tuna` then exits with the command's own status:

```sh
if tuna split --down --wait -- make test; then
  echo "tests passed"
fi
```

`tuna wait --exit` does the same for a pane that is already running a command it started with:

```sh
pane=$(tuna split --down -- go test ./...)
# … do other things …
tuna wait "$pane" --exit --timeout 10m
echo "go test exited with $?"
```

The pane stays open after the command, with its output, so you can `tuna read` it. Waiting for a command doesn't work with `cmd` as the shell, which can't tell when one ends.

## Watching

`tuna events` prints what happens to sessions and panes as it happens, one line per event, until you stop it. A script can follow every agent at once this way, without polling `tuna status`.

```sh
tuna events                               # every session
tuna events api                           # one session, or one pane
tuna events -p tuna                       # one project's sessions
tuna events --for waiting,error           # only panes going to these statuses
tuna events --json | jq -r .name          # one JSON object per line
tuna events --for waiting --once          # exit after the first one
```

```
12:04:31  api x78n9ss  working → waiting  Allow Bash: rm -rf dist?
12:04:40  api  renamed api → backend
```

| Event | When | Fields |
| --- | --- | --- |
| `pane.status` | A pane's status changes | `pane`, `from`, `to`, `message`, `kind`, `agent`; `exitCode` when it goes to exited |
| `command.exited` | The command a pane started with ends, and its shell goes on | `pane`, `exitCode` |
| `pane.created`, `pane.closed` | A pane starts or closes | `pane` |
| `session.created`, `session.closed` | A session starts or closes | |
| `session.renamed` | A session's name changes | `from`, `to` |
| `session.moved` | A session moves to another project | `from`, `to`, the projects' ids |
| `notice` | A program in a session sends a desktop notification | `pane`, `title`, `body` |

With `--json`, every event also has `event`, `time` (unix ms), `session`, `project` and `name`, the session's name as the sidebar shows it.

`tuna events` reports what changes from now on: for what is there already, use `tuna ls --json`. Given a target, it stops when the target closes. `--once` exits after the first event that matches, with `1` if the target closes first or `--timeout` passes. tunad sends its state when it settles, so a status that lasts only a moment, such as a blip of working, may not show.

## Browsing

`tuna browser` drives a browser pane's page the way `tuna send` and `tuna read` drive a terminal: it clicks, types, reads the page and takes pictures of it. Its target is a browser pane, or a session, for its browser. With no target, it drives the browser of the session it runs in.

```sh
tuna split --url localhost:3000            # open the page, print the pane's id
tuna browser open localhost:3000/todos     # go to a page, and wait for it to load
tuna browser snapshot                      # what's on the page, each element with a ref
tuna browser type --ref e3 "Buy milk" --enter
tuna browser click --ref e5
tuna browser wait '.todo-list li'          # until an element shows
tuna browser wait --text "1 item left"     # until the text shows
tuna browser text                          # the page's text
tuna browser text h1                       # an element's
tuna browser click 'button[type=submit]'
tuna browser screenshot -o page.png
tuna browser eval 'document.title'
```

| Command                         | What it does |
| ------------------------------- | ------------ |
| `open <url>`                    | Go to a page, and wait for it to load |
| `back`, `forward`, `reload`     | Go back or forward, or load the page again, and wait for it |
| `text [selector]`               | Print the page's visible text, or the first matching element's |
| `html [selector]`               | Print the page's HTML, or the element's |
| `snapshot [selector]`           | Print what the page shows as a tree of roles and names (`button "Add"`, `textbox "Email"`), each element with a ref, or the element's part of it |
| `click <selector>`              | Click the element, as a mouse does |
| `type <selector> <text…>`       | Type text into a field, in place of what it holds; `--enter` presses Enter after it, which submits its form |
| `wait [selector] [--text text]` | Wait until the element shows, or the text does; with neither, until the page loads |
| `eval <js…>`                    | Run JavaScript in the page and print its value: a string as it is, anything else as JSON |
| `screenshot [-o file.png]`      | Save a picture of the page, and print where; `-o -` writes it to stdout |
| `cookies save [file]`           | Save the browser's cookies to a file, or print them; `--domain` keeps a site's |
| `cookies load <file>`           | Load cookies from a file; `--domain` loads a site's |

Selectors are CSS selectors; quote them in the shell. Anywhere a command takes a selector, `--ref` takes a ref from `snapshot` in its place: `tuna browser click --ref e5`. An element keeps its ref until the page goes to another, and a ref whose element is gone fails rather than guess another. Commands that need an element wait for it to appear, and `--timeout` (default `10s`) sets how long they, and the page loads, may take. `type` and `eval` take `-` for text from stdin.

Clicking and typing send the events a real mouse and keyboard would, so pages built with React and other frameworks see them. Every command works on any browser pane, shown or not: a screenshot of a pane not on screen takes no focus from what you are doing.

Browser panes share their cookies, and keep them when tuna restarts. `cookies` saves them in Playwright's storage state, as JSON, or for a file named `.txt`, as a Netscape `cookies.txt`, which curl and browser extensions write; `load` reads either. Cookies sign you in to sites: share their file only with what you'd trust with your logins. They stay on this Mac, so programs on a host can't save or load them.

Like `tuna do`, `tuna browser` needs Remote Control, in **Settings › Agents**, and the app to be open.

## Past conversations

```sh
tuna history                     # agents' conversations panes ran, the last to run first
tuna history -p tuna -n 50       # one project's, up to 50
tuna history login --json        # those matching "login", as JSON
tuna resume 13b9                 # resume one: its id, a prefix of it, or its name
tuna resume fix-login -f         # and show it in the app
```

`tuna resume` starts a session that resumes the conversation, in its directory and project, and prints the session's id. A conversation a pane still runs gives that session instead. See [Past conversations](./sessions.md#past-conversations).

## Status

```sh
tuna status api --json                    # status, program, directory
tuna report working "Running migrations"  # report your own status
tuna notify "Need your input"             # flag this pane as needing you
```

See [Status & Notifications](./status.md) for the states.

## Changes

```sh
tuna diff                                 # the pane's uncommitted changes, as git diff prints them
tuna diff --turn                          # what its agent changed in its last turn
tuna diff review --branch                 # since the default branch
```

See [Changes](./changes.md).

## Configuration and setup

```sh
tuna config path                 # where the configuration file is
tuna config show                 # every key with its value
tuna config show --defaults      # every key as tuna comes
tuna config edit                 # open it in $EDITOR
tuna config check                # report unknown keys and bad values
tuna projects                    # list projects
tuna skill setup                 # teach coding agents to use tuna
tuna daemon status|stop|restart  # manage tunad
tuna daemon install-service      # run tunad with the machine (systemd, launchd)
tuna daemon serve                # run tunad in the foreground, for a service manager
tuna clients                     # who is connected to tunad
tuna pets                        # the pets on the window's edge, the one shown first
tuna pets use <name|tuna|codex>  # show a pet: one in ~/.codex/pets, or tuna's
tuna pets install <pack.zip>     # install a Codex pet pack, a .codex-pet.zip
tuna pets rm <name>              # remove a pet tuna installed
```

## Feishu

```sh
tuna feishu setup                # ask for the app's ID and secret, keep them, connect, and say what's next
tuna feishu status               # whether the bot is connected
tuna feishu test                 # check that Feishu knows the app
tuna feishu send <chat> <text…>  # send a message to a chat, by its ID
```

A Feishu chat drives the coding agent in the session it's bound to: `/new` starts one, `/use` binds an open one, and other text is a prompt. See [Feishu](./feishu.md).

## Other machines

`-H <host>` before the command, or `$TUNA_HOST`, points a command that acts on sessions at tunad on another machine: a host you added, or an address such as `me@devbox` or `wsl://Ubuntu`. See [Remote Hosts](./remote.md).

```sh
tuna hosts                       # the hosts you added
tuna hosts add devbox me@devbox  # add one, which the app shows too
tuna hosts check devbox          # connect, and say how it went
tuna hosts install devbox        # install tuna there, over SSH
tuna ls --all                    # the sessions here and on every host
tuna -H devbox ls
tuna -H devbox new -n review -- claude
tuna -H devbox daemon update     # update tuna there to this one's version
```

Inside a pane, `tuna` talks to the tunad of the machine the pane runs on, wherever the app is.

`tuna ls --all` asks every host at once and lists each one's sessions under the name `-H` takes, `local` for this machine. A host out of reach says why in place of its sessions. With `--json`, it prints a list of `{host, addr, error, sessions}`.

## Running app actions

`tuna do` runs any of the app's actions, the same ones as its menus. It needs Remote Control, which is on by default in **Settings › Agents**.

```sh
tuna actions                     # list them
tuna do pane.split.right
tuna do pane.browser url=localhost:3000
tuna do sidebar.toggle
```

## Exit status

| Code | Meaning                              |
| ---- | ------------------------------------ |
| `0`  | OK                                   |
| `1`  | Failed                               |
| `2`  | Bad usage                            |
| `3`  | No such target, or an ambiguous one  |
| `4`  | tunad is unreachable                 |

Commands that wait for a program, `tuna wait --exit` and `--wait`, exit with the program's status instead, as does `tuna wait` when the pane's shell exits with an error. A wait that times out exits with `1`.

Add `--json` to `ls`, `status` and `projects` for machine-readable output.
