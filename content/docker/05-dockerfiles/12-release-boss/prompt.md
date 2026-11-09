# Release 1.1

Boss fight! Ship a new version of the app, using only what you learned in this chapter.

# Assignment

In `~/docker-course`:

1. Change your `Dockerfile` so the image has two default settings: **`APP_VERSION=1.1`** and **`GREETING=Welcome`**.
2. Make sure the test data in `app/test-data` still stays out of the image.
3. Build it as **`cq/webapp:1.1`**.
4. Run **`cq-release`** from `cq/webapp:1.1`, detached, on port **8088**. **Don't** use `-e`: the defaults must come from the image.
5. Open http://localhost:8088. It should say `Welcome from webapp 1.1`.
6. Submit `dk.dockerfile.release-boss`.

Tip: you now have lots of containers running. Clear out the old ones with `docker rm -f cq-web cq-site-a cq-site-b cq-settings cq-keeper` if you like; nothing later needs them.
