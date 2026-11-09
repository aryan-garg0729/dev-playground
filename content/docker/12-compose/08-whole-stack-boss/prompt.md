# The Whole Stack in One File

Boss fight! Make the stack production-shaped and prove its data survives a full `down` and `up`.

# Assignment

In `~/docker-course`, change `compose.yaml` so that:

1. **`api`** keeps its build, settings, notes volume and healthcheck, but is **not published** on any port. Only `front` should be reachable from your computer.
2. **`front`** is published on port **8099** instead of 8098, and still waits for a healthy `api`.
3. The `toolbox` profile can stay.

Then:

4. `docker compose up -d` and wait for `api` to be healthy.
5. Save a note **through the front door**:

   ```sh
   curl -X POST "http://localhost:8099/notes?text=saved+through+the+front+door"
   ```

6. Take everything down (keeping the data!), then bring it back up.
7. Check `curl http://localhost:8099/notes` still shows the note.
8. Submit `dk.compose.whole-stack-boss`. Don't save more notes after step 6.
