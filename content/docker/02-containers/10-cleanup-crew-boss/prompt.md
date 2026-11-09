# Container Control Room

Boss fight! Everything from this chapter, all at once: run in the background, read logs, find exit codes and clean up.

# Assignment

In `~/docker-course`:

1. **Keep one running.** Start **`cq-keeper`** from **`alpine:3.20`**, in the background, running `sleep infinity`.
2. **Run a job.** Start **`cq-job`** from **`alpine:3.20`** that prints `job done` and finishes. Don't remove it afterwards.
3. **Solve the mystery.** This container exits with a different exit code every time. Copy and paste it:

   ```sh
   docker run --name cq-mystery alpine:3.20 sh -c 'exit $(( $(date +%s) % 50 + 1 ))'
   ```

   Find its exit code, then add it to `answers.txt` as a line like `mystery_exit=17`:

   ```sh
   echo "mystery_exit=NUMBER" >> answers.txt
   ```

4. **Clean up.** Remove the practice containers from earlier lessons: **`cq-sleeper`**, **`cq-counter`** and **`cq-crash`**.
5. Submit `dk.containers.cleanup-crew-boss`.

If you need another go at the mystery, `docker rm cq-mystery` and run it again. Then fix the line in `answers.txt` too, because the code will have changed.
