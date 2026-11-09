# Ship It Small, 1.6

Boss fight! Release a small image that still does everything the big one did.

# Assignment

In `~/docker-course`:

1. Use a multi-stage `Dockerfile`. The final stage can be `alpine:3.20` or `scratch`, your choice.
2. Build **`cq/webapp:1.6`** with `VERSION=1.6`. It must report `webapp 1.6`, carry the `1.6` version label and keep `GREETING=Welcome`.
3. It must be small (well under 30 MB) with **at most 2 layers**: one from the base and one for the program.
4. Run **`cq-small`** from it, detached, on port **8089**.
5. Submit `dk.multistage.ship-small-boss`.

If `--version` prints `webapp dev`, check which stage your `ARG` is in.
