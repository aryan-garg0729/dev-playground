# Liveness or Readiness

There are three kinds of probe. They look the same in YAML; what differs is what the kubelet *does* when they fail:

| probe | question | on failure |
| --- | --- | --- |
| `readinessProbe` | Should I get traffic? | removed from Service endpoints; **not restarted** |
| `livenessProbe` | Am I stuck? | container is **killed and restarted** |
| `startupProbe` | Have I finished starting? | other probes wait until it passes; restarted if it never does |

`startupProbe` exists for slow starters: without it, a liveness probe can kill an app that's still loading.

Mo has an idea for brewd: *"Let's make the liveness probe call `/healthz?check=db`, which fails if brewd can't reach the database. Then a broken pod restarts itself."*

**What do you tell Mo?**
