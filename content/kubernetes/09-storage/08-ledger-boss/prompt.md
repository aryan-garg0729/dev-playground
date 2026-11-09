# The Order Vault

Priya's ticket:

> **ledger: three replicas, three sets of books.** Each ledger member keeps its own journal on its own disk, survives restarts with its journal intact, and is reachable by name so the members can talk to each other. Prove the journals survive by restarting one member.

## Assignment

1. A **headless** Service named **`ledger`** selecting `app: ledger`, port 80.
2. A StatefulSet named **`ledger`**: 3 replicas, governed by that Service, pods labelled `app: ledger`, `busybox:1.36`, with a `volumeClaimTemplates` entry named **`books`** (1Gi, ReadWriteOnce) mounted at **`/books`**.
3. Every time a member starts, it **appends** one line to `/books/journal` and then stays alive:

   ```yaml
   command: ["sh", "-c", "echo \"$(hostname) opened the books at $(date)\" >> /books/journal; exec sleep infinity"]
   ```

4. Delete the pod **`ledger-1`**, let it come back, and check that its journal now has two lines.

Every check is a hit on the boss. Defeat it to unlock **Scheduling**.
