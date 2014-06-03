# tempconsul

tempconsul makes it easy to start and stop temporary consul agent server
processes. This is useful for things like integration tests (and not much else).

See `tempconsul_test.go` for an usage example.

[API documentation](http://godoc.org/github.com/stvp/tempconsul)

## Should I use this?

Probably not.

## Why not?

It only supports spinning up a single consul agent server, bound to the default
ports locally. It doesn't support anything else, really.

