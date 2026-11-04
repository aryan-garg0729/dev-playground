Python ships with a tiny web server. Use it to serve two files from a folder.

1. In an empty folder, create `hello.txt` containing `Hello from my server`.
2. In the same folder, create `info.json`:

   ```json
   { "course": "command-line", "lessons": 2 }
   ```

3. Start the server in that folder and leave it running:

   ```sh
   python3 -m http.server 8000
   ```

4. In a second terminal, run `codequest run cl.local-server`, then
   `codequest submit cl.local-server`.

If port 8000 is busy, start the server on another port and point the CLI at it
with `codequest config base_url http://localhost:<port>`.
