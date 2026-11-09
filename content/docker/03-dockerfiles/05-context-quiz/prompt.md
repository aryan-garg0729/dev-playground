# The Build Context

Look at the first lines of any build:

```
$ docker build -t elixirbytes/brewd:v1 .
[+] Building 0.4s
 => [internal] load build definition from Dockerfile
 => [internal] load .dockerignore
 => [internal] load build context
 => => transferring context: 5.01MB
```

*Transferring context?* The builder doesn't read your disk directly. Docker packs up the **build context**, the folder you name at the end of `docker build`, and sends it to the builder. `COPY` can only see files inside that bundle. That's also why `COPY ../secrets.txt .` fails: `..` is outside the context.

That 5 MB is your fixtures file. Remember it.

Dex, in a hurry, runs this from their home folder:

```sh
cd ~
docker build -t elixirbytes/brewd:v1 -f elixirbytes/Dockerfile ~
#  => => transferring context: 2.31GB   ...still going
```

**Why is the build so slow before it even starts?**
