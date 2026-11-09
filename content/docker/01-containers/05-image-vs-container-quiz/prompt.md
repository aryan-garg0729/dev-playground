# Image vs Container

Two words get mixed up all the time:

| | Image | Container |
| --- | --- | --- |
| What it is | A read-only template: files plus some settings | A running (or stopped) process made from an image |
| Analogy | A recipe | The potion you brewed from it |
| How many | One `alpine:3.20` | As many containers as you like, from that one image |
| Can you change it? | No | Yes, but only its own private writable layer |

You just created `/etc/eb-badge` inside `eb-hello`. Now, with `eb-hello` still running, someone runs:

```sh
docker run --rm alpine:3.20 cat /etc/eb-badge
```

**What happens?**
