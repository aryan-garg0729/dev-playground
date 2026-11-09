# Release 1.5.0

Priya's ticket:

> **REL-150: cut brewd 1.5.0.** Product added the line `// 1.5.0: faster potion lookups` to `main.go` (it's a potion shop; the bar is low). Build the release so that it **says** what it is, inside and out, and so that this rebuild is a fast one. Attach the build log as proof.

## Assignment

Using only this chapter, in `~/elixirbytes`:

1. Make sure your `Dockerfile` downloads dependencies before copying the source, stamps `BUILD_VERSION` into the binary and the `org.opencontainers.image.version` label, and builds with the Go compiler cache mounted.
2. Build it once as anything, so the cache is warm.
3. Add the comment line `// 1.5.0: faster potion lookups` to `app/main.go`.
4. Build **`elixirbytes/brewd:1.5.0`** with `BUILD_VERSION=1.5.0`, using `--progress=plain`, and save the output to **`release.log`**.

The log is checked step by step: the dependency download must be **reused**, the source copy must **run fresh** (you changed a file, so a cached copy means the log is from the wrong build), and the build step must show the cache mount and the release version.

Submit from `~/elixirbytes` to attack. Every check is a hit on the boss. Defeat it to unlock **ENTRYPOINT, CMD and PID 1**.
