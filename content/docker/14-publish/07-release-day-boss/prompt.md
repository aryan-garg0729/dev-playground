# Release Day

Final boss! Ship version 2.0.0 the way a real team would: build it, tag it properly, push it to the registry and run it from there.

# Assignment

In `~/docker-course`, with `cq-registry` running:

1. Build the production `Dockerfile` (non-root user, healthcheck) with `VERSION=2.0.0`.
2. Push it to your registry as **`localhost:15000/cq/webapp`** with three tags: **`2.0.0`**, **`2.0`** and **`latest`**, all pointing at the same image.
3. Run **`cq-v2`** from **`localhost:15000/cq/webapp:2.0.0`**, detached, and wait until it's `(healthy)`.
4. Submit `dk.publish.release-day-boss`.

You know Docker now: what it is, how images and containers work, how to build, store, network, compose, harden and ship them. Congratulations! 🐳
