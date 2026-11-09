# The Build Context

What does the `.` in `docker build -t cq/webapp:v1 .` actually do?

It's the **build context**: the folder Docker packs up and sends to the builder before it starts. `COPY` can only see files inside the context. That's why the Dockerfile says `COPY app/ .` and works.

Two things follow from that:

- `COPY ../somefile .` fails. Anything outside the context doesn't exist as far as the build is concerned.
- **Everything** in the folder gets sent, even files the Dockerfile never uses.

Look closely at the start of a build in `~/docker-course`:

```
=> [internal] load build context
=> => transferring context: 26.4MB
```

The app is two small files.

**Why is the context 26 MB?**
