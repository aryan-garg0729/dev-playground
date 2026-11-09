# The Notebook That Survives

Boss fight! Prove that your data outlives your containers.

# Assignment

In `~/docker-course`:

1. Create a volume named **`cq-book`**.
2. Run the webapp (**`cq/webapp:1.6`**) in a container on port **8092**, with `cq-book` at **`/data`** and your **`config`** folder mounted **read-only** at **`/config`**. Call it whatever you like for now.
3. Save two notes, `first note` and `second note`:

   ```sh
   curl -X POST "http://localhost:8092/notes?text=first+note"
   curl -X POST "http://localhost:8092/notes?text=second+note"
   ```

4. Remove that container.
5. Start a **new** container named **`cq-notebook`**, set up exactly the same way.
6. Check `curl http://localhost:8092/notes` still lists both notes.
7. Submit `dk.storage.notebook-boss`. **Don't add more notes after step 5**: one check confirms `cq-notebook` was created after the notes were written.
