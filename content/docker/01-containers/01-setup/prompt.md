# Day One at ElixirBytes

Welcome to **ElixirBytes**, the internet's third-most-popular potion shop. Every order goes through one small Go service called **brewd**. Right now brewd runs because somebody SSHed into a server, started it inside `tmux` and never logged out. When that box reboots, nobody can buy a potion.

Priya, our platform lead, wants brewd in containers. Before we touch brewd, we need a machine where Docker works and a way to grade what you do.

## How these lessons work

You work in your own terminal. The `codequest` CLI then runs a few read-only checks (it shows each command and asks first) and sends the output to be graded. Two habits for the whole course:

- Run `codequest` from one folder, `~/elixirbytes`.
- Use the exact names each lesson gives: containers start with `eb-`, images with `elixirbytes/`. That's how the checks find your work.

## Assignment

1. Install Docker: Docker Desktop, OrbStack or Colima on macOS, or Docker Engine on Linux. Make sure `docker version` prints a **Server** section. If it only shows a Client, the Docker daemon (the background service that actually runs containers) isn't running.
2. Install the CodeQuest CLI from the [install page](/cli/install) (one command on macOS and Linux), then sign in:

   ```sh
   codequest login
   ```

3. Create the course folder and download the small image we'll use all chapter:

   ```sh
   mkdir -p ~/elixirbytes && cd ~/elixirbytes
   docker pull alpine:3.20
   ```

4. Run `codequest run dk.containers.setup`, then `codequest submit dk.containers.setup`.

On Linux, `permission denied` means you're not in the `docker` group: `sudo usermod -aG docker $USER`, then log out and back in. That's the only `sudo` in this course.
