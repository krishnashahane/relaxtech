# RelaxTech

RelaxTech is a Go command-line client for controlling and inspecting supported Eight Sleep functionality from a terminal.

It provides commands for temperature control, device status, sleep data, alarms, schedules, audio, adjustable-base controls, travel/jet-lag data, presence, metrics, and a local scheduler/daemon.

> Compatibility note: RelaxTech talks directly to Eight Sleep API endpoints. Those endpoints are not guaranteed to be stable or officially supported for third-party clients, so API changes may require future compatibility updates.

## Requirements

- Go 1.24+
- An Eight Sleep account
- A supported Eight Sleep device
- Network access to the Eight Sleep API

## Install

### Run from source

```bash
git clone https://github.com/krishnashahane/relaxtech.git
cd relaxtech
go run ./cmd/relaxtech --help
```

### Build a binary

```bash
go build -o relaxtech ./cmd/relaxtech
./relaxtech --help
```

## Authentication and configuration

RelaxTech accepts credentials and connection settings through persistent flags or Viper configuration/environment variables.

Example configuration:

```yaml
email: "you@example.com"
password: "your-password"
user_id: ""
client_id: "sleep-client"
client_secret: ""
timezone: "local"
output: "table"
verbose: false
```

Default configuration path:

```text
~/.config/relaxtech/config.yaml
```

Environment variables use the RELAXTECH_ prefix, for example:

```bash
export RELAXTECH_EMAIL="you@example.com"
export RELAXTECH_PASSWORD="your-password"
```

The client stores access tokens in the operating-system keyring when available. When the encrypted file backend is used, set RELAXTECH_KEYRING_PASSWORD to a strong secret; RelaxTech does not use a predictable built-in fallback password.

### Security guidance

Do not commit config.yaml, passwords, access tokens, or client secrets to source control.

```bash
chmod 600 ~/.config/relaxtech/config.yaml
```

When possible, prefer environment variables or the operating-system keychain over credentials passed directly on a command line.

## Common commands

Show all commands:

```bash
./relaxtech --help
```

Power control:

```bash
./relaxtech on
./relaxtech off
```

Temperature and status:

```bash
./relaxtech temp --help
./relaxtech status
```

Schedules:

```bash
./relaxtech schedule list
./relaxtech schedule next
./relaxtech schedule create --start 22:00 --level -10 --days 1,2,3,4,5
```

Sleep and metrics:

```bash
./relaxtech sleep --help
./relaxtech metrics --help
```

Available command groups include alarm, audio, autopilot, base, daemon, device, feats, household, presence, schedule, sleep, tempmode, tracks, and travel.

## Output formats

Many commands support:

```bash
--output table
--output json
--output csv
```

Restrict output fields with:

```bash
--fields id,title,type
```

Enable diagnostics with:

```bash
--verbose
```

Diagnostics are written to stderr. The client does not intentionally print passwords or authentication tokens.

## Daemon

The local scheduler can run timed device actions from the RelaxTech configuration.

```bash
./relaxtech daemon --help
```

The daemon uses a PID file, handles SIGINT/SIGTERM, prevents duplicate instances when the recorded process is still alive, and safely removes stale PID files.

## Development

Format:

```bash
gofmt -w .
```

Test:

```bash
go test ./...
```

Static analysis:

```bash
go vet ./...
```

Build:

```bash
go build ./...
```

## Architecture

```text
cmd/relaxtech
      |
      v
internal/cmd       CLI commands (Cobra + Viper)
      |
      v
internal/client    Eight Sleep API client
      |
      +---- internal/tokencache
      |          |
      |          v
      |      OS keyring / encrypted file backend
      |
      +---- internal/config
      +---- internal/output
      +---- internal/daemon
```

## Security improvements in this version

- Removed the hard-coded OAuth client secret from source.
- OAuth client ID and secret are configurable.
- Authentication failures no longer log response bodies or headers.
- API error responses no longer echo arbitrary upstream response bodies.
- Authentication refresh is bounded to prevent infinite retry loops.
- File-keyring fallback no longer uses a predictable built-in password.
- Daemon startup detects and cleans stale PID files.
- The API transport requires TLS 1.2 or newer.
- Added a real Go module manifest with pinned direct dependencies.

## License

MIT

## Author

Krishna Shahane

GitHub: https://github.com/krishnashahane/relaxtech
