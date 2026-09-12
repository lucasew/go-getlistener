# go-getlistener

Get a TCP `net.Listener` that works with plain `HOST`/`PORT` env vars and with
systemd socket activation.

Any TCP-based protocol can sit on top of the returned listener. See
`cmd/go-listener-demo` for a minimal HTTP example.

## Usage

```go
ln, err := getlistener.GetListener()
if err != nil {
    return err
}
defer ln.Close()
// serve on ln (http.Serve, custom accept loop, …)
```

## Configuration

`GetListener` reads process environment only. There is no options struct to
pass from code.

| Variable | Meaning | Default |
|----------|---------|---------|
| `PORT` | TCP port to bind. `0` or unset picks an ephemeral port. Non-numeric values and values outside `0`–`65535` are an error (the bad value is included in the message). | unset → ephemeral |
| `HOST` | Address to bind. Non-loopback values log a security warning via `log/slog`. | `127.0.0.1` |

Bound address after listen: `ln.Addr()`.

## Systemd socket activation (unix)

When `LISTEN_PID` is set, the library prefers the passed file descriptor over
`HOST`/`PORT`:

1. `LISTEN_PID` must equal this process’s PID, otherwise `GetListener` returns
   an error (it does **not** fall back to TCP).
2. `LISTEN_FDS` must be `1`. Zero or multiple sockets are unsupported.
3. The socket FD is the usual systemd start FD (`3`).

If `LISTEN_PID` is unset, systemd activation is skipped and TCP listen uses
`HOST`/`PORT` as above.

Example unit sketches (adjust paths and user):

```ini
# app.socket
[Socket]
ListenStream=127.0.0.1:8080
# or: ListenStream=8080

[Install]
WantedBy=sockets.target
```

```ini
# app.service
[Service]
ExecStart=/usr/local/bin/your-app
# inherits the socket from app.socket
```

## Platforms

- **unix:** TCP + systemd socket activation
- **windows:** TCP only (`HOST`/`PORT`)

## Develop

Tooling is pinned in `mise.toml`:

```bash
mise install
mise run test
mise run lint
mise run ci
```
