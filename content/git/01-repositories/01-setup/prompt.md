# The Handbook Problem

**ElixirBytes**, the potion shop, has an ops handbook: who to page, how to deploy brewd, what's on the menu. It lives in one shared online doc. Last Tuesday Dex pasted over the deploy steps by mistake, and nobody could say what they used to be.

Undo buttons in a doc only go back so far, and they never tell you *who* changed *what* and *why*. Priya wants the handbook in **Git**, a version control system: a tool that records every change to a folder as a snapshot you can inspect, compare and go back to. You're its new maintainer.

## How these lessons work

You work in your own terminal. The `codequest` CLI then runs a few **read-only** checks in your folder (it shows each command and asks first) and sends the output to be graded. They never change your files.

Two habits for the whole course:

- Run `codequest` from one folder: `~/elixirbytes/handbook`.
- Use the exact names, labels and lines each lesson gives. That's how the checks find your work.

## Assignment

1. Install Git 2.28 or newer. On macOS, `xcode-select --install` or `brew install git`. On Linux, your package manager (`sudo apt install git`, `sudo dnf install git`). Check it:

   ```sh
   git version
   # git version 2.47.1
   ```

2. Install the CodeQuest CLI and sign in. If you did the Docker course, you already have it.

   ```sh
   go install github.com/aryangarg/codequest/api/cmd/codequest@latest
   codequest login
   ```

3. Create the handbook folder and move into it:

   ```sh
   mkdir -p ~/elixirbytes/handbook && cd ~/elixirbytes/handbook
   ```

4. Run `codequest run gt.repos.setup`, then `codequest submit gt.repos.setup`.

The checks confirm your Git version is at least 2.28 (older versions lack commands we use, like `git switch`) and that you're in `~/elixirbytes/handbook`.
