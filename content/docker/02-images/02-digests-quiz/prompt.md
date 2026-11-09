# Tags Move, Digests Don't

Every image also has a **digest**: a SHA-256 hash of its content.

```sh
docker image inspect -f '{{range .RepoDigests}}{{println .}}{{end}}' alpine:3.20
# alpine@sha256:d9e853e87e55526f6b2917df91a2115c36dd7c696a35be12163d44e6e2a4b6bc
```

| | Tag (`alpine:3.20`) | Digest (`alpine@sha256:d9e8…`) |
| --- | --- | --- |
| Chosen by | A human | Computed from the content |
| Can it change what it points at? | **Yes**, at any time | **No**. Different bytes give a different hash. |
| Good for | Humans and "give me the newest 3.20.x" | Reproducible deploys |

Hashes like this are how you check that bytes haven't changed. We won't cover how SHA-256 works here.

Mo reports something odd. A Dockerfile says `FROM potionlib:2.4`. Built in March on one server, it works. Built in April on another server, the same Dockerfile fails, and nobody changed it.

**What's the most likely explanation?**
