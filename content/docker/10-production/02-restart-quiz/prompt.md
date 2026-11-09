# Restart Policies

`docker run --restart POLICY` tells Docker what to do when the main process exits:

| Policy | Restarts after a crash? | After the Docker daemon or machine restarts? | After `docker stop`? |
| --- | --- | --- | --- |
| `no` (default) | no | no | no |
| `on-failure[:N]` | only on a non-zero exit, up to N times | yes, if it was still meant to be running | no |
| `always` | yes, any exit | **yes, even if you'd stopped it** | not until the daemon restarts |
| `unless-stopped` | yes, any exit | yes, **unless you'd stopped it** | no |

Restarts back off: Docker waits 100 ms, then 200 ms, 400 ms... so a crash loop doesn't eat the CPU.

Dex wants brewd to come back after crashes and reboots, but if someone deliberately runs `docker stop` during an incident, it must **stay stopped**, even across a reboot.

**Which policy?**
