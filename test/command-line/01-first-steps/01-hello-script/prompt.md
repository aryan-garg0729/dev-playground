Some lessons run on your own computer instead of in the browser. You write the
code in your own editor, and the `codequest` CLI runs it and sends the output to
CodeQuest to be graded.

1. Install the CLI and sign in:

   ```sh
   go install github.com/aryangarg/codequest/api/cmd/codequest@latest
   codequest login
   ```

2. In an empty folder, create `hello.py` that prints exactly:

   ```
   Hello, CodeQuest!
   ```

3. From that folder, run `codequest run cl.hello-script` to see the output, then
   `codequest submit cl.hello-script` to grade it.
