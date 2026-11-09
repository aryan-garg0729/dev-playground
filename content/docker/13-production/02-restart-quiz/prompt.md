# Restart Policies

Apps crash. Servers reboot. In production, someone has to start things again, and it shouldn't be you at 3 a.m. A **restart policy** tells Docker what to do when a container stops:

```sh
docker run -d --restart unless-stopped myapp
```

| `--restart` | Restarts after a crash? | Restarts when Docker or the computer restarts? | After you `docker stop` it |
| --- | --- | --- | --- |
| `no` (default) | No | No | Stays stopped |
| `on-failure[:N]` | Yes, if the exit code isn't 0 (at most N times) | Only if it had failed | Stays stopped |
| `always` | Yes | Yes | Starts again next time Docker starts |
| `unless-stopped` | Yes | Yes | **Stays stopped** |

Docker waits a little longer between each restart, so a container that crashes instantly doesn't hammer your machine.

**You want a web server to come back after crashes and after the server reboots, but to stay down when you stop it on purpose for maintenance. Which policy?**
