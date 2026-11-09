# The 10-Second Stop

Our webapp is polite: when `docker stop` asks it to stop, it prints `webapp: shutting down` and exits within a blink. Try it:

```sh
docker stop cq-release      # from the Dockerfiles boss: instant
docker logs cq-release      # ... webapp: shutting down
```

Now a very common situation: a team wants to print a message (or prepare something) before the app starts, so they add a small **start script**:

```sh
#!/bin/sh
echo "starting webapp..."
webapp "$@"
```

(`"$@"` passes along any arguments the container was given.) The image runs the script with `ENTRYPOINT ["sh", "/src/start.sh"]`. Looks harmless. Let's see what happens.

# Assignment

In `~/docker-course`:

1. Create the script. Copy and paste:

   ```sh
   cat > app/start.sh <<'EOF'
   #!/bin/sh
   echo "starting webapp..."
   webapp "$@"
   EOF
   ```

2. In your `Dockerfile`, replace the last line with `ENTRYPOINT ["sh", "/src/start.sh"]`. (The script is copied into `/src` along with the rest of `app/`.)
3. Build it as **`cq/webapp:banner`**.
4. Run **`cq-slow`** from it, detached. Check `docker logs cq-slow` shows the banner.
5. Time the stop: `time docker stop cq-slow`. About 10 seconds...
6. Look at the exit code and logs. No `shutting down` line, and `Exited (137)`. It was forced.
7. Run, then submit `dk.entrypoint.slow-stop`.
