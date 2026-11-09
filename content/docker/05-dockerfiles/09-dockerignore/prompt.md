# Leave the Junk Out

A **`.dockerignore`** file lists what to leave out of the build context. It sits next to the Dockerfile, one pattern per line:

```
*.tar
secrets.env
node_modules
```

Files it matches are never sent to the builder, so they can't end up in your image by accident. That makes builds faster, images smaller and, most importantly, it stops secrets (passwords, keys) from being baked in.

# Assignment

In `~/docker-course`:

1. Pretend someone left 5 MB of test files inside the app folder. Copy and paste:

   ```sh
   mkdir -p app/test-data
   head -c 5000000 /dev/urandom > app/test-data/big.bin
   ```

2. Build as **`cq/webapp:fat`**, then look inside: `docker run --rm cq/webapp:fat ls /src`. The junk came along for the ride.
3. Create `.dockerignore`. Copy and paste:

   ```sh
   cat > .dockerignore <<'EOF'
   *.tar
   answers.txt
   app/test-data
   EOF
   ```

4. Build again as **`cq/webapp:lean`**. Notice the much smaller "transferring context" number. Then check `ls /src` again.
5. Run, then submit `dk.dockerfile.dockerignore`.

Patterns are relative to the context folder, so it's `app/test-data`, not just `test-data`.
