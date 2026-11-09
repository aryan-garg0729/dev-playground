# On Call

Boss fight! You're on call and three containers are broken. Each one has a different problem. Use the debugging checklist to find out what's wrong, then fix it.

# Assignment

In `~/docker-course`:

1. Start all three broken containers. Copy and paste:

   ```sh
   docker run -d --name cq-bug-port -p 8096:80 cq/webapp:1.6
   docker run -d --name cq-bug-env alpine:3.20 sh -c '[ -n "$WORKER_TOKEN" ] || { echo "fatal: WORKER_TOKEN missing"; exit 3; }; echo "worker ready"; sleep infinity'
   docker run -d --name cq-bug-net alpine:3.20 sh -c 'until wget -qO- -T 2 http://cq-api:8080/healthz >/dev/null 2>&1; do echo "cannot reach cq-api"; sleep 2; done; echo "connected to cq-api"; sleep infinity'
   ```

2. Investigate and fix each one. Keep each container's **name** and **command**; change only what's needed around it:
   - **`cq-bug-port`**: http://localhost:8096 doesn't load.
   - **`cq-bug-env`**: isn't running.
   - **`cq-bug-net`**: keeps saying it can't reach `cq-api`. (`cq-api` from the networking chapter must still be running.)
3. Submit `dk.debug.fix-it-boss`.
