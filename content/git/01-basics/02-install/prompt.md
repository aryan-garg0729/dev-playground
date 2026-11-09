# Install Git

Let's get Git on your computer.

## How these lessons work

You do each lesson in your own terminal. Then the `codequest` command runs a few small checks (it shows each command and asks first) and sends the result to be graded.

## Assignment

1. Install Git 2.28 or newer:
   - **macOS**: run `xcode-select --install`, or `brew install git` if you use Homebrew.
   - **Linux**: use your package manager, for example `sudo apt install git`.

   Then check it:

   ```sh
   git version
   # git version 2.47.1
   ```

2. Install the CodeQuest CLI from the [install page](/cli/install) and sign in. (If you did the Docker course, you already have it.)

   ```sh
   codequest login
   ```

3. Make a folder for this course and move into it:

   ```sh
   mkdir -p ~/git-course && cd ~/git-course
   ```

4. Run `codequest run gt.basics.install` to test, then `codequest submit gt.basics.install`.
