# Pin an Exact Image

Even a version tag like `alpine:3.20` can move: when Alpine releases 3.20.4 with a security fix, the publishers point `3.20` at the new build. Usually that's what you want.

But sometimes you need **exactly the same image, byte for byte**, forever. For that, every image has a **digest**: a long fingerprint calculated from its contents.

```sh
docker images --digests alpine
# REPOSITORY   TAG    DIGEST                                                                    IMAGE ID
# alpine       3.20   sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d   0a9a5dfd008f
```

Change a single byte of the image and the digest changes completely, so a digest can never point anywhere else. Use it instead of a tag with `@`:

```sh
docker run --rm alpine@sha256:beefdbd8...   # the full digest
```

Tags are friendly nicknames that can move. Digests are fingerprints that never do.

# Assignment

1. Find the digest of your `alpine:3.20` with `docker images --digests alpine`.
2. Start a container named **`cq-pinned`** from `alpine@sha256:<that digest>` (copy the full digest) that runs `cat /etc/alpine-release`.
3. Run, then submit `dk.images.digests`.
