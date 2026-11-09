# Graceful Release 1.4

Boss fight! Ship a release that starts through the script, stops instantly and still accepts options.

# Assignment

In `~/docker-course`:

1. Make sure `app/start.sh` prints the banner and then hands over to the webapp so that it becomes PID 1, passing along any arguments.
2. Make sure the `Dockerfile`'s `ENTRYPOINT` runs the script in **brackets form**.
3. Build **`cq/webapp:1.4`** with the build argument `VERSION=1.4`.
4. Check it yourself: `docker run --rm cq/webapp:1.4 --version` should print the banner and then `webapp 1.4`.
5. Submit `dk.entrypoint.graceful-boss`.

The second check starts a container named `cq-probe` from your image, looks at PID 1, times `docker stop`, then removes it again.
