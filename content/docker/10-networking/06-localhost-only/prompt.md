# Just for Me

Look closely at what `docker port` printed for your earlier containers:

```
0.0.0.0:8085
[::]:8085
```

`0.0.0.0` means **every network your computer is connected to**. Anyone on the same Wi-Fi (in a café, at the office) can open `http://YOUR-IP:8085`. For an admin page or a database, that's a bad idea.

Put an address in front of the port to publish on **your computer only**:

```sh
-p 127.0.0.1:8093:8080
#  └── localhost: only this computer can connect
```

`127.0.0.1` (also called `localhost`) always means "this computer".

# Assignment

1. Start **`cq-admin`** from **`cq/webapp:1.6`**, detached, published on **`127.0.0.1`**, port **8093**.
2. Check `docker port cq-admin`: just one line, with `127.0.0.1`.
3. Run, then submit `dk.networking.localhost-only`.
