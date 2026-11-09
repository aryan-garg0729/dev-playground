# Settings with -e

Most apps in containers are configured with **environment variables**: named settings like `PASSWORD=secret` or `MODE=production` that a program can read when it starts.

Pass them with `-e NAME=value`, once per setting:

```sh
docker run -d --name cq-demo-db -e POSTGRES_PASSWORD=secret postgres:16-alpine
```

(That's how the official database images on Docker Hub ask for a password. Each image's Hub page lists the variables it understands.)

See what a container got with `printenv`:

```sh
docker exec cq-demo-db printenv POSTGRES_PASSWORD
# secret
```

The same image can behave differently in different containers, just by changing the settings. No rebuilding needed.

# Assignment

1. Start a detached container named **`cq-settings`** from **`alpine:3.20`** running `sleep infinity`, with two settings: **`MODE=practice`** and **`LEVEL=3`**.
2. Check them with `docker exec cq-settings printenv MODE LEVEL`.
3. Run, then submit `dk.ports.env-vars`.
