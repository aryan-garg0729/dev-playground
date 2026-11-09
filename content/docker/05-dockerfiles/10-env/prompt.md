# Default Settings with ENV

In chapter 3 you passed settings with `-e`. An image can also come with **default** settings, using `ENV` in the Dockerfile:

```dockerfile
ENV APP_VERSION=1.0
```

Every container made from the image gets that variable automatically. `-e` still wins when you want something different for one container:

```sh
docker run -e APP_VERSION=2.0 ...   # overrides the image's default
```

Our app reads two settings:

| Variable | Default | Used in |
| --- | --- | --- |
| `APP_VERSION` | `dev` | `Hello from webapp dev` |
| `GREETING` | `Hello` | `Hello from webapp dev` |

# Assignment

1. Add `ENV APP_VERSION=1.0` to your `Dockerfile`. Copy and paste to replace it:

   ```sh
   cat > Dockerfile <<'EOF'
   FROM golang:1.23-alpine
   WORKDIR /src
   COPY app/ .
   RUN go build -o /usr/local/bin/webapp .
   ENV APP_VERSION=1.0
   EXPOSE 8080
   CMD ["webapp"]
   EOF
   ```

2. Build it as **`cq/webapp:1.0`**.
3. Run **`cq-app-v1`** from `cq/webapp:1.0`, detached, on port **8086**. It should say `Hello from webapp 1.0`.
4. Run **`cq-howdy`** from the same image, detached, on port **8087**, with **`GREETING=Howdy`**. It should say `Howdy from webapp 1.0`.
5. Run, then submit `dk.dockerfile.env`.
