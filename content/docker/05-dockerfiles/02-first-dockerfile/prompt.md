# Your First Image

Time to build an image of your own.

`docker build` takes a folder that contains a Dockerfile:

```sh
docker build -t cq/demo:v1 somefolder
#            │             └── the folder with the Dockerfile in it
#            └── -t: the name (tag) to give the new image
```

Docker runs each step in the Dockerfile and prints its progress. At the end, the new image is in `docker images`, just like the ones you pulled.

# Assignment

1. In `~/docker-course`, make a folder and a Dockerfile in it. Copy and paste:

   ```sh
   mkdir -p hello
   cat > hello/Dockerfile <<'EOF'
   FROM alpine:3.20
   CMD ["echo", "Hello from my own image!"]
   EOF
   ```

   (That `cat > file <<'EOF' ... EOF` trick writes everything between the two `EOF` lines into the file. You can also create the file in your editor.)

2. Build it as **`cq/hello:v1`**:

   ```sh
   docker build -t cq/hello:v1 hello
   ```

3. Run it with `docker run --rm cq/hello:v1`. You made that!
4. Run, then submit `dk.dockerfile.first-dockerfile`.
