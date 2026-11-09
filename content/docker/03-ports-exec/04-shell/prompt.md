# A Shell Inside

Typing `docker exec cq-web ...` before every command gets old. Open a **shell** inside the container instead, and you can type commands as if you'd logged into it:

```sh
docker exec -it cq-web sh
/ #
```

- `-i` keeps the input open so you can type.
- `-t` gives you a proper terminal.
- `sh` is the shell program to start.

The prompt changes to `/ #`: you're inside the box now. Look around:

```sh
/ # whoami
root
/ # ls /usr/share/nginx/html
50x.html  index.html
/ # ps
PID   USER     COMMAND
    1 root     nginx: master process nginx -g daemon off;
   ...
```

`ps` shows the processes inside. There's nginx, and almost nothing else. A container only runs what it needs.

Type `exit` (or press Ctrl+D) to leave.

# Assignment

1. Open a shell inside **`cq-web`** with `docker exec -it cq-web sh`.
2. Inside, run `ps` and `ls /`, and have a look around.
3. Still inside, leave a note:

   ```sh
   echo "I was here" > /tmp/note.txt
   ```

4. Type `exit`.
5. Run, then submit `dk.ports.shell`.
