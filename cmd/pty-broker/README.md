# warpbuild-pty-broker

Tiny WebSocket → PTY broker meant to live inside an isolated environment
(e2b sandbox, runner, etc.) and front a single login shell per WS connection.
Keystrokes from a browser flow straight to the broker — no proxy hop in the
byte path.

## Layout

- `main.go` — broker source (Go, ~250 lines). Single binary, single port,
  single concurrent PTY per WS connection.

## Auth

HMAC-SHA256 short-lived tokens in the `?token=` query string. The HMAC key
is provided to the broker via `BROKER_HMAC_KEY` env (or `--hmac-key` flag)
at startup. Token format is a compact `<base64url(payload)>.<base64url(sig)>`
with a JSON payload carrying:

```
{
  "sid": "<sandbox-id>",     // scope (audited)
  "uid": "<user-id>",        // optional, audited
  "cwd": "/path",            // working dir for the spawned shell
  "c":   <cols>,             // optional initial cols
  "r":   <rows>,             // optional initial rows
  "e":   <unix-seconds>      // expiry
}
```

Tokens that don't verify or are past `e` get a 401 before the WS upgrade.

## Wire protocol

- `BinaryMessage` from client → raw stdin bytes (UTF-8 keystrokes).
- `BinaryMessage` from broker → raw stdout bytes (PTY output).
- `TextMessage` from client → JSON control envelope:
  - `{"type":"resize","cols":N,"rows":N}`
  - legacy `{"type":"input","data":"..."}`
- WS close with reason `"pty exited"` when the shell terminates.

## Flags

| Flag | Default | Notes |
|---|---|---|
| `--port` | `7681` | TCP port to listen on. |
| `--hmac-key` | `$BROKER_HMAC_KEY` | Required; verifier secret. |
| `--cwd` | `/home/user` | Default shell cwd if the token doesn't carry one. |
| `--shell` | `$SHELL` → `/bin/bash` → `/bin/sh` | Shell to spawn (always with `-l`). |

## Endpoints

- `GET /term?token=<jwt-ish>` — WS upgrade; spawns one PTY per connection.
- `GET /health` — `200 ok`.

## Build

Local:

```sh
make build-pty-broker      # cross-compiles linux/amd64 into bin/warpbuild-pty-broker
```

CI: published as a goreleaser archive (`warpbuild-pty-broker_Linux_x86_64.tar.gz`,
`..._Linux_arm64.tar.gz`) by the `release-testing.yaml` workflow. Tarballs are
synced to `r2:warpbuild-packages-dev/warpbuild-agentd/<branch>/` and reachable
via the public R2 CDN.
