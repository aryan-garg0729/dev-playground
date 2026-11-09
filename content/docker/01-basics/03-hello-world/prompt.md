# Hello, Container

Every programming course starts with "Hello, world". Docker has one too, and it's a whole app in a box.

```sh
docker run hello-world
```

That one command does a surprising amount:

1. Docker looks for something called `hello-world` on your computer.
2. It's not there, so Docker **downloads** it from the internet (you'll see `Unable to find image... Pulling from library/hello-world`).
3. Docker starts a container from it.
4. The tiny app inside prints a message and finishes.

Run it a second time. It's much faster, because the download is already on your computer.

Don't worry about the words "image" and "pull" yet. They get their own lessons very soon.

# Assignment

1. In `~/docker-course`, run:

   ```sh
   docker run hello-world
   ```

2. Read the message. It should start with **Hello from Docker!**
3. Run `codequest run dk.basics.hello-world`, then `codequest submit dk.basics.hello-world`.

The checks confirm `hello-world` was downloaded and that it still says hello.
