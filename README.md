# multi-ping
simple TUI for monitoring multiple hosts with ICMP ping.

## Features

- Add, edit, and delete hosts.
- View the list of configured targets.
- Ping hosts using ICMP in parallel.
- Continuously monitor all targets until stopped by the user.
- Display `UP`, `DOWN`, or `WARN` status with latency.

## Running from source

Make sure [Go](https://go.dev/) 1.25.5 or later is installed.

```bash
go run .
```

To build locally:

```bash
go build -o multi-ping .
```

## Installation

Prebuilt binaries are available on the [Releases](https://github.com/abrarrayva/multi-url-ping/releases) page.
Release `v1.0.0` includes binaries for Linux, macOS, and Windows.

### Linux (command line)

```bash
mkdir -p "$HOME/.local/bin"
curl -fsSL https://github.com/abrarrayva/multi-url-ping/releases/latest/download/multi-ping-linux-amd64 \
  -o "$HOME/.local/bin/multi-ping"
chmod +x "$HOME/.local/bin/multi-ping"
grep -qxF 'export PATH="$HOME/.local/bin:$PATH"' "$HOME/.profile" 2>/dev/null || \
  echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.profile"
export PATH="$HOME/.local/bin:$PATH"
multi-ping
```

For ARM64 Linux, replace `multi-ping-linux-amd64` with
`multi-ping-linux-arm64`.

Uninstall:

```bash
rm -f "$HOME/.local/bin/multi-ping"
```

### macOS (command line)

Intel Macs:

```bash
mkdir -p "$HOME/.local/bin"
curl -fsSL https://github.com/abrarrayva/multi-url-ping/releases/latest/download/multi-ping-macos-amd64 \
  -o "$HOME/.local/bin/multi-ping"
chmod +x "$HOME/.local/bin/multi-ping"
grep -qxF 'export PATH="$HOME/.local/bin:$PATH"' "$HOME/.zprofile" 2>/dev/null || \
  echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zprofile"
export PATH="$HOME/.local/bin:$PATH"
multi-ping
```

Apple Silicon Macs (M1/M2/M3/M4):

```bash
mkdir -p "$HOME/.local/bin"
curl -fsSL https://github.com/abrarrayva/multi-url-ping/releases/latest/download/multi-ping-macos-arm64 \
  -o "$HOME/.local/bin/multi-ping"
chmod +x "$HOME/.local/bin/multi-ping"
grep -qxF 'export PATH="$HOME/.local/bin:$PATH"' "$HOME/.zprofile" 2>/dev/null || \
  echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zprofile"
export PATH="$HOME/.local/bin:$PATH"
multi-ping
```

macOS may ask you to allow the unsigned binary in **System Settings > Privacy & Security**.

Uninstall:

```bash
rm -f "$HOME/.local/bin/multi-ping"
```

### Windows (PowerShell)

Install the binary into your user programs directory and add that directory to
your user `PATH`:

```powershell
$installDir = Join-Path $env:LOCALAPPDATA "Programs\multi-ping"
New-Item -ItemType Directory -Force -Path $installDir | Out-Null
Invoke-WebRequest `
  -Uri https://github.com/abrarrayva/multi-url-ping/releases/latest/download/multi-ping-windows-amd64.exe `
  -OutFile (Join-Path $installDir "multi-ping.exe")
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (($userPath -split ";") -notcontains $installDir) {
  [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
}
$env:Path = "$installDir;$env:Path"
multi-ping
```

Open a new PowerShell window after installation if `multi-ping` is not found.

Uninstall:

```powershell
$installDir = Join-Path $env:LOCALAPPDATA "Programs\multi-ping"
Remove-Item -Recurse -Force $installDir
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = ($userPath -split ";" | Where-Object { $_ -and $_ -ne $installDir }) -join ";"
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
```

## Controls

- `↑` / `↓` or `j` / `k`: navigate the menu
- `Enter`: select or submit
- `1`–`6`: quick menu selection
- `Esc`: go back or cancel
- `q` / `Ctrl+C`: quit

Edit and delete actions use the same keyboard navigation to select a target;
you no longer need to enter a target ID manually.

When health monitoring is running, each target is checked concurrently every
5 seconds. Press `Enter` or `Esc` to stop monitoring.

Enter a hostname or IP address, with or without `http://` or `https://`.
The app extracts the host and sends ICMP echo requests, so inputs such as
`8.8.8.8`, `1.1.1.1`, and `example.com` are accepted.

ICMP may require elevated privileges depending on the operating system:

- Linux: run with `sudo`, or grant the binary the required raw-socket capability.
- macOS: run with `sudo` if regular execution cannot open an ICMP socket.
- Windows: run the terminal as Administrator if the ping is denied.

`8.8.8.8` and `1.1.1.1` are valid DNS resolver IPs and can be tested with
this app because ICMP ping is different from an HTTP request. A firewall or
network policy may still block ICMP and report the host as down.


## Project Structure

- `main.go`: application state, menu actions, and input handling
- `ui.go`: TUI views and rendering
- `health.go`: concurrent ICMP health checks
