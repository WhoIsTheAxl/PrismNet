# PrismNet

Minecraft: Bedrock Edition networking abstractions in Go, powered by Gophertunnel.

## Experimental AI branch

This branch is **work in progress**, not a production-ready release. It adds:
- Authentication with cached Microsoft tokens.
- Direct server and Realm connection helpers.
- Session packet callbacks and a single packet reader.
- Bidirectional proxy relay with packet interception hooks.

The proxy relay accepts two **already connected** Gophertunnel connections. It does not yet set up a server listener or handle authentication/protocol translation. Player actions and commands are still placeholders.

### Development checks

Run `go mod tidy`, `go test ./...`, and `go vet ./...` locally.
No CI or runtime validation has been performed for these changes.

Never commit `token.json` or other credentials.
