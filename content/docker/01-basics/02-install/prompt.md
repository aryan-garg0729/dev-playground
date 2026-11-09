# Install Docker

Time to get Docker on your computer. Docker comes in two halves, and you need both:

- the **`docker` command** you type in the terminal (the *client*)
- the **Docker engine**, a program that runs quietly in the background and does the real work

## How these lessons work

You do each lesson in your own terminal. Then the `codequest` command runs a few small checks (it shows you each command and asks first) and sends the result to be graded.

# Assignment

1. Install Docker:
   - **macOS**: install [Docker Desktop](https://www.docker.com/products/docker-desktop/) and open it once. (OrbStack or Colima also work.)
   - **Linux**: install [Docker Engine](https://docs.docker.com/engine/install/).
2. Check both halves are there:

   ```sh
   docker version
   ```

   You should see a **Client** section *and* a **Server** section. If there's no Server section, the engine isn't running: open Docker Desktop and wait for it to say it's running.

3. Install the CodeQuest CLI from the [install page](/cli/install) and sign in:

   ```sh
   codequest login
   ```

4. Make a folder for this course. You'll run every lesson from here:

   ```sh
   mkdir -p ~/docker-course && cd ~/docker-course
   ```

5. Run `codequest run dk.basics.install` to test, then `codequest submit dk.basics.install`.

On Linux, `permission denied` means your user isn't allowed to talk to Docker yet. Run `sudo usermod -aG docker $USER`, then log out and back in. That's the only `sudo` in this course.
