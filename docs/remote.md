---
title: Remote Hosts
description: Run sessions on a Linux box, a cloud VM or another computer, and see them in the same sidebar as your local ones.
---

# Remote Hosts

tuna can show sessions that run on other machines, such as a Linux dev box, a cloud VM, a Mac mini under your desk or a WSL distribution, in the same window as your local ones. They get the same sidebar, statuses, notifications and `tuna` command. Each machine runs its own tunad, so closing your laptop or quitting tuna leaves remote agents running, as it does locally.

## Quick start

A host is any machine with tuna on it. It can be another computer where you already use the tuna app, or a server that has no display.

### A computer with the tuna app

The app needs nothing installed beside it: its own tunad serves as the host. Pick how your other machines reach it:

- **Over Tailscale.** In the app on the host, turn on **Settings › Hosts › Let other machines connect**. It shows the address to add elsewhere, such as `http://studio.tail1234.ts.net:7373`. See [Tailscale](#tailscale).
- **Over SSH.** If the computer takes SSH logins (on a Mac, **System Settings › General › Sharing › Remote Login**), add it by the address you'd give `ssh`. tuna finds the app there on its own.

Sessions on the host keep running while its app is closed, as they do on your own computer.

### A server

On a Linux box, a cloud VM or any machine without the app, install the server build of tuna. It needs no root or administrator, no display and no GTK. On Linux or macOS:

```sh
curl -fsSL https://github.com/smallnest/tuna/releases/latest/download/install-server.sh | sh
```

On Windows, in PowerShell:

```powershell
irm https://github.com/smallnest/tuna/releases/latest/download/install-server.ps1 | iex
```

### Adding the host

On your computer, open **Settings › Hosts**, type the host's address, such as `me@devbox` for SSH or the address the host's app showed, and press **Add**. The host's projects appear in the sidebar with yours, each tagged with the host's name.

The same from the command line:

```sh
tuna hosts add devbox me@devbox
tuna -H devbox new -n review -- claude
```

::: tip
If tuna isn't installed on a host you add, the sidebar says so and offers **Install tuna…**, which runs the script above on the host over SSH.
:::

## How tuna reaches a host

An address says how to reach a host. tuna takes the same address in Settings, in the configuration file and after `tuna -H`.

| Address | How |
| --- | --- |
| `devbox`, `me@devbox`, `me@devbox:2222`, `ssh://me@devbox` | SSH, with your `~/.ssh/config`, keys and agent |
| `wsl://Ubuntu`, `wsl://` | A WSL distribution of this Windows PC, through `wsl.exe` |
| `http://devbox:7373` | tunad's listener on a tailnet address (see [Tailscale](#tailscale)) |
| `https://devbox.tail1234.ts.net` | tunad behind Tailscale Serve |
| `tailcat://tc…` | A pair code, through a relay: no Tailscale needed (see [Tailcat](#tailcat)) |

tuna opens one connection per host and carries everything over it: the sidebar's state, every pane you look at, and the CLI's commands. The connection is compressed, and a heartbeat notices a dead link within seconds. When a link drops, as when a laptop sleeps, the host's sessions stay in the sidebar, faded, and tuna reconnects on its own, waiting longer between tries up to a minute.

### SSH

Over SSH, tuna runs `tuna __bridge` on the host, which starts tunad there if it isn't running and passes the connection to it. There is nothing to set up on the host beyond tuna itself, the app or the server build, and nothing listens on the network.

tuna never prompts for anything over SSH, as there's no one to answer in the app:

- **Keys with a passphrase** need your SSH agent. tuna uses the agent your shell sets in `SSH_AUTH_SOCK`, as 1Password's, as `ssh` in a terminal does.
- **A host you haven't connected to before** needs its key trusted first. Run `ssh me@devbox` once in a terminal.
- **Password logins** aren't supported. Use a key.

tuna looks for `tuna` on the host in `~/.local/bin` (on Windows, `%LOCALAPPDATA%\Programs\tuna-server`), the desktop app and its command line tool, and your `PATH`. If yours is elsewhere, name it in the address: `me@devbox?tuna=/opt/tuna/tuna`, or on Windows `me@pc?tuna=C:/tuna/tuna.exe`.

### Windows hosts

tuna reaches a Windows machine over SSH, so turn on its OpenSSH Server first (**Settings › System › Optional features**, as an administrator). tuna works with the server's default shell, cmd, and with PowerShell. It finds out which one the first time it connects, which costs that connection a moment.

The install script runs tunad as a task of yours in Task Scheduler, which starts it as you log on to the machine, in your session. While no one is logged on, tunad starts as tuna connects and keeps running after the connection ends. A tunad started that way runs without your Windows password, as programs do over SSH with a key, so programs in its sessions can't reach network shares or the credentials Windows keeps for you. Log on to the machine once, or restart tunad after you do, to run them with your full account.

### WSL

On Windows, `wsl://Ubuntu` reaches the tuna installed in that WSL distribution, as SSH would. `wsl://` is the default distribution. If tuna isn't installed there, tuna offers to install it, as it does over SSH.

WSL may stop a distribution a while after its last program started from Windows ends, and tunad with it. Keep agents running by keeping the distribution up, for example with [systemd enabled](https://learn.microsoft.com/windows/wsl/systemd) and `tuna daemon install-service` inside it.

### Tailscale

tunad can also listen on your tailnet, so that it needs no SSH at all. Tailscale tells tunad who connects, and by default only the Tailscale user who owns the machine gets in. On a host with the app, turn on **Settings › Hosts › Let other machines connect**, which shows the address to add elsewhere. On a server, or to set it by hand, add this to the configuration on the host:

```ini
# ~/.config/tuna/config on the host
listen = tailnet://:7373
```

Then add the host as `http://devbox:7373`. The connection is plain HTTP, which WireGuard encrypts between tailnet devices.

To serve it over HTTPS with Tailscale's certificate instead, let tunad listen on a socket and point Tailscale Serve at it:

```ini
listen = tailnet-serve://
```

```sh
tailscale serve --bg --https=443 unix:$HOME/.config/tuna/serve.sock
```

Then add the host as `https://devbox.tail1234.ts.net`. The socket is `serve.sock` in tuna's data folder: `~/.config/tuna` on Linux, `~/Library/Application Support/tuna` on macOS.

To let others in besides the owner, list them with `listen-allow`: a tailnet login, `tag:ci` for tagged devices, or `*` for anyone on the tailnet. A tagged machine has no owner, so only `listen-allow` lets anyone in to it. Tailscale Serve doesn't tell tunad who tagged devices are: they connect to a `tailnet://` listener.

```ini
listen-allow = alice@example.com
listen-allow = tag:ci
```

tunad follows `listen` and `listen-allow` as they change, and a listener turned off drops the connections it took. A `tailnet-serve://` listener added while tunad runs waits for `tuna daemon restart`. tunad listens only when the configuration asks, and only on Tailscale or through a Tailcat relay. It refuses requests from web browsers, so a web page you visit can't reach it. Funnel's requests from the internet never get in. To keep it listening with the machine, run it as a service (below).

### Tailcat

Tailcat reaches a host with a pair code instead of Tailscale. It needs no Tailscale account, and no Tailscale installation on either machine. It uses [Tailcat](https://github.com/tailscale/tailcat), a library made from Tailscale's parts: it makes a WireGuard tunnel between the two machines, through a public DERP relay that lets them find each other, and directly once they do.

The pair code is the host's address. It contains the host's key, a secret pre-shared key and the relay's region, as a long `tc…` string. **Treat it as a password: anyone who has it can connect to the host.**

**Settings › Hosts › Tailcat** holds all Tailcat settings, in one group. The switch **Use Tailcat to Reach Other Hosts** turns them all on. With it on, the group has two parts:

- **Add a Host** adds a host that this machine reaches. Copy the pair code from the host's Tailcat group, or from `tuna hosts tailcat code` on the host. **Remote Tailcat Port** is the port of the tailcat listener of the hosts that you add: 7373 unless one of them names another.
- **Provide Tailcat Service** makes this machine a host. tuna then shows its pair code, with a button to copy it, and **This Computer's Tailcat Listener Port** is the port that tunad listens on here. **Change the pair code** makes a new code.

The pair code does not have the port of the listener in it. The **Remote Tailcat Port** of the other machines must be the **This Computer's Tailcat Listener Port** of the host.

**Relay** is the DERP map that all machines use to find their relay. It applies to both parts.

The client needs this version of tuna, or a later one.

The same settings from the command line:

```sh
# On the host: turn on Tailcat, then listen.
# tailcat = true goes in the configuration first.
tuna hosts tailcat enable  # or tuna hosts tailcat enable 7400
tuna hosts tailcat code      # the pair code, its secret
tuna hosts tailcat rotate    # a new pair code: the old one stops working

# On the other machine:
tuna hosts add box tailcat://tc…
```

The pair code stays the same across restarts. tunad also logs it at start:

```sh
grep 'listening on tailcat' ~/.config/tuna/tunad.log
```

The log is in `~/Library/Application Support/tuna` on macOS, and in `%LOCALAPPDATA%\tuna` on Windows.

**The switch.** Tailcat is off until you turn it on. With it off, the app lists no Tailcat host and connects to none of them, and tunad starts no `tailcat://` listener. Leave it off if you do not trust a relay: tuna then starts nothing for Tailcat.

**The relay.** By default, Tailcat uses Tailscale's public relays at `https://tailcat.dev/derpmap.json`. Tailscale limits their rate, and does not promise they stay available. The sessions themselves do not go through the relay if the machines can reach each other directly over the network: the relay only carries what the machines use to find each other. To use a relay of your own, as Tailscale's [derper](https://github.com/tailscale/tailscale/tree/main/cmd/derper) provides, set `tailcat-relay` to the URL of your DERP map, or set **Relay** in the app. Set the same map on the host and on the other machines: the pair code keeps the relay's number, and each side looks it up in its own map.

**Changing the pair code.** **Change the pair code** in the app, or `tuna hosts tailcat rotate`, makes a new key. tunad starts its listener again with the new key, and the old code no longer connects. Add this machine again on your other machines. To remove the code from a machine, delete its key file, `tailcat.json` in tuna's data folder.


**The keys.** tunad keeps its Tailcat key in `tailcat.json` in tuna's data folder, which only your user can read. The pair code stays the same across restarts, while the file stays. Delete the file and restart tunad to get a new key, as **Change the pair code** does: every machine that has the old code loses the host.

## Running tunad as a service

A host with the app starts tunad as the app opens, or as tuna connects over SSH. On a server, run tunad as a service of your own, so that it starts with the machine and keeps your sessions running while no one is logged in. The install scripts do this for you. On a Mac, the service is a launchd agent, which runs while you're logged in to the Mac, so keep a session logged in there; the script installs it when you are. On Windows, it's a Task Scheduler task, which starts tunad as you log on.

```sh
tuna daemon install-service      # a systemd user unit, a launchd agent on macOS, a task on Windows
tuna daemon status               # how it runs, and where it listens
tuna daemon uninstall-service
```

On Linux this first enables lingering (`loginctl enable-linger`), which keeps your services running after you log out. Without it, a service would stop as your last SSH session ends, which a tunad started on demand outlives, so `install-service` refuses to go on without it. If your system asks an administrator for lingering, ask them to run `sudo loginctl enable-linger <you>`, or install the service anyway with `--no-linger`.

`tuna daemon serve` runs tunad in the foreground and logs to its standard error, for a service manager of your own. It never exits for being idle. `--listen` takes the place of the configuration's `listen` keys.

## Working on a host

- **The sidebar.** It lists projects, yours first, then each host's in the order of **Settings › Hosts**. A host's projects carry its name after theirs, as `api` `devbox`, so the same repository on two machines sits side by side. A host gets a row of its own only when it needs you: when tuna can't reach it, wants updating, or isn't installed there. A host with no sessions doesn't show; start one there with Open Project. Right-click one of a host's projects for the host's menu: reconnect, restart tunad, update tuna, its color, its settings.
- **Colors.** Each host has a color, picked from its name: its tag in the sidebar and its name in a pane's header wear it, so a remote pane is plain to see. Pick another from the host's menu, **Color**, or with `color=` after its address (`host = devbox=me@devbox color=teal`).
- **New sessions.** <kbd>⌘</kbd> <kbd>T</kbd> and **New Session** start a session on the host of the session you're looking at. A project's **+** starts one in it, on its machine.
- **Open Project.** With hosts added, <kbd>⌘</kbd> <kbd>O</kbd> lists folders in tuna instead of your computer's dialog, starting on the machine of the session you're looking at; the switch at its top lists another's. Type to filter, or type a path such as `~/code`. <kbd>↩</kbd> goes into a folder, and <kbd>⌘</kbd> <kbd>↩</kbd> or **Open** opens it as a project.
- **Paths.** <kbd>⌘</kbd>-clicking a file path in a remote pane opens the file in the host's `$EDITOR`, in a pane below. A folder opens in a new session there. tuna never opens a file of your computer that happens to have the same path. Hovering an image's path with <kbd>⌘</kbd> previews the host's image.
- **Pasting images and files.** Pasting an image or files into a remote pane copies them to a private temporary folder on the host and pastes their paths there, so agents that read image paths, such as Claude Code, keep working. Dragging files onto a remote pane still pastes your computer's paths.
- **Agents.** Inside a remote pane, `tuna` talks to that host's tunad, so the agent skill works there unchanged. Run `tuna skill setup` on the host.
- **Notifications.** A host's notifications and the window's title carry the host's name, so a remote program can't pass for a local one.

### Remote Control

`tuna focus` and `tuna do` let a program drive the app. For a host, that's code on another machine driving your computer's window, so it's off by default. Turn it on per host in **Settings › Hosts**, or with `remote-control` after the address:

```ini
host = devbox=me@devbox remote-control
```

### Settings

Appearance, keyboard and notifications are the app's. A host's shell, `env`, scrollback and session restore come from the configuration on that host, since tunad reads them there. Run `tuna config edit` on the host to change them.

## Versions and updates

The app and a host can run different versions of tuna. When the host runs an older release than the app, its first project in the sidebar shows **Update**. It downloads the app's version of tuna onto the host, checks it against the key tuna's own updates are signed with, for that version and platform, and then restarts tunad there. Restarting starts the host's sessions again where they were, with their history, and the programs running in them end, so tuna asks first.

- tuna updates a host only when you ask, never on a reconnect, and never back to an older version.
- When the host runs a newer tuna than the app, update the app instead (**Settings › General**).
- tuna updates an install of its install scripts. If you installed tuna on the host another way, update it the way you installed it.
- A host running the tuna desktop app takes updates from your other machines once you let it: on that computer, turn on **Settings › Hosts › Let them update tuna here** (`remote-updates = true`). The app there then updates to the newest release, relaunches, and restarts tunad. Until then, update it on that computer.

From the command line:

```sh
tuna -H devbox version           # this tuna's version, and the host's
tuna -H devbox daemon status
tuna -H devbox daemon update     # to this tuna's version
tuna -H devbox daemon restart
tuna daemon update               # on the host, to the latest release
tuna daemon update 2.1.0         # on the host, to a version, older ones too
```

## The command line

`-H` (or `--host`), before the command, points a command that acts on sessions at a host's tunad. It takes a host's name, an address, or `local` for this machine. `$TUNA_HOST` does the same, except inside a pane, where `tuna` acts on the pane's own machine. Commands of this machine alone, such as `tuna config` and `tuna hook`, ignore `$TUNA_HOST`.

```sh
tuna hosts                       # the hosts you added (--json for JSON)
tuna hosts add build ci@build.lab --remote-control
tuna hosts check build           # connect, and say how it went; a name or an address
tuna hosts install build         # install tuna there, over SSH or WSL
tuna hosts remove build
tuna ls --all                    # the sessions here and on every host
tuna -H build ls
tuna -H build prompt review "run the tests and fix them"
tuna clients                     # who is connected to this machine's tunad
```

The app and the CLI share one list of hosts, the configuration's `host` keys:

```ini
host = devbox=me@devbox
host = build=ssh://ci@build.lab:2222 remote-control
host = ubuntu=wsl://Ubuntu
```

| Variable | What it does |
| --- | --- |
| `TUNA_HOST` | The host commands act on, as `-H` |
| `TUNA_SSH` | The `ssh` program tuna runs, in place of the one on your `PATH` |
| `TUNA_RELEASES` | A mirror of tuna's releases, for the install scripts and `tuna daemon update` |

## Security

- **A connection can do anything its account can,** including starting programs, as a local one can. Over SSH it's your own account, as SSH logs you in. Over Tailscale, only who the machine's owner lets in.
- **The socket is yours alone.** tunad's socket is created readable by its user only, and tunad refuses connections from other users' processes.
- **Updates are signed.** The server build updates itself only from archives signed with the key of tuna's releases.
- `tuna clients` lists who is connected to a tunad, with how they came: `uid 501 over unix`, `me from 10.0.0.5 over ssh`, `alice@example.com over tailnet`.

## Limits

- The server build is for Linux, with glibc 2.28 or later, macOS and Windows 10 (1809) or later, on x86-64 and ARM64. A machine with the tuna app needs no server build.
- Dragging files onto a remote pane pastes your computer's paths. Paste them instead.
- One host has one account. Login servers, where several people reach their own accounts on one machine, are planned.
