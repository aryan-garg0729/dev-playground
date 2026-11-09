# Killed for Memory

CPU is squeezable: a container over its CPU limit just runs slower. Memory isn't. When a container goes over its memory limit, the Linux kernel's **OOM killer** ("out of memory") ends it on the spot with SIGKILL. No warning, no graceful shutdown, often no log line, because the process never got the chance to write one.

That makes OOM kills easy to miss. The evidence is in the pod's status, not the logs:

```sh
kubectl describe pod oom-canary
#   Last State:  Terminated
#     Reason:    OOMKilled      <- or just "Error", see below
#     Exit Code: 137
#   Restart Count: 3
#   Limits:
#     memory:  32Mi
```

Exit codes above 128 mean "killed by a signal": the code is 128 plus the signal's number.

The `Reason` field is filled in by the container runtime, and it isn't reliable everywhere. Some runtime versions (including the containerd in recent kind releases) often report `Error` for an OOM kill. Don't wait for the magic word. **Exit code 137, a memory limit, and nobody else stopping the container** is an OOM kill until proven otherwise.

## Assignment

Create a pod named **`oom-canary`** (`busybox:1.36`, label `app: oom-canary`) with a memory **limit of `32Mi`** and no probes, running a command that keeps grabbing more memory until something stops it:

```yaml
command: ["sh", "-c", "echo filling memory; exec tail /dev/zero"]
```

(`tail` keeps everything it has read since the last newline, and `/dev/zero` never sends one, so it grows forever. `exec` makes `tail` the container's main process, so it's the process the kernel kills.)

Watch it die a few times, then add two lines to `answers.txt`:

```
oom_exit=<the exit code of its last run>
oom_signal=<the number of the signal that killed it>
```

The logs say `filling memory` and then nothing. That's what to remember for the next time an app "just restarts" with nothing in its logs: look at `describe` for exit code 137 and a memory limit before you look anywhere else.
