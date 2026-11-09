# Access Modes

A claim's `accessModes` say how the volume may be shared:

| mode | short | meaning |
| --- | --- | --- |
| `ReadWriteOnce` | RWO | read-write, by pods on **one node** at a time |
| `ReadOnlyMany` | ROX | read-only, by pods on many nodes |
| `ReadWriteMany` | RWX | read-write, by pods on many nodes (needs storage that supports it, like NFS or a cloud file service) |
| `ReadWriteOncePod` | RWOP | read-write, by exactly **one pod** |

Most block storage (cloud disks, kind's local-path) only supports RWO.

Mo scales `ledger-writer` to `replicas: 3` "so the order log is highly available". The claim is `ReadWriteOnce`, and the scheduler spreads the pods over both workers.

**What happens?**
