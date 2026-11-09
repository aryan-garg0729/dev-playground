# Fast Release 1.3

Boss fight! Make a release build that's quick to rebuild and gets its version from the command line.

# Assignment

In `~/docker-course`:

1. Make sure your `Dockerfile`:
   - copies `app/go.mod` and downloads libraries **before** copying the rest of the code,
   - builds with a **cache mount** for `/root/.cache/go-build`,
   - takes the version from a **build argument** named `VERSION`, passes it to the app as `APP_VERSION` and puts it in the `org.opencontainers.image.version` label.
2. Build once with `VERSION=1.3` as **`cq/webapp:1.3`**.
3. Change the code: `echo "// release 1.3" >> app/main.go`
4. Build **`cq/webapp:1.3`** again with `VERSION=1.3`, in plain mode, saving the output to **`release.log`**.
5. Submit `dk.cache.fast-release-boss`.

If your Dockerfile from "A Cache That Survives" is still in place, you're most of the way there already.
