# TLS

## Using `generate_cert.go` tool

1. `mkdir tls` within project directory
2. `cd tls`
3. Find the source code path of your Go installation, e.g. `/usr/local/go/` or relative to `$(go env GOPATH)`
4. `go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost`
5. Add `tls` folder to `.gitignore`

**This does the following:**

- Generates a 2048-bit RSA key pair.
- Stores private key in `key.pem` file, and generates a self-signed TLS certificate for the host `localhost` containing the public key, which are stored in a `cert.pem` file.
