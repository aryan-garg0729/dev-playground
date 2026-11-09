# Night Shift Automation

Dex wants his nights back. His ticket:

> **Automate the night shift.**
> 1. **inventory-count**: a one-off Job that counts the six warehouse aisles, three at a time. `busybox:1.36`, any short command that exits 0.
> 2. **hourly-receipts**: a CronJob that emails receipts at the top of every hour. If a run is still going when the next one is due, kill the old one: only the newest receipts matter. Keep just **1** failed Job around for debugging. Its container prints `receipts sent`.
> 3. I want to see receipts go out **now**, not in an hour: one manual run called **receipts-now**, made from the CronJob.
> 4. **log-shipper**: a DaemonSet that reads logs on every **worker** node. Not on the control plane: that node's logs are shipped by a different system.

## Assignment

Build all four in `manifests/night-shift.yaml` (the manual run is a command, not a manifest). Use exactly these names, label the DaemonSet's pods **`app: log-shipper`**, and keep its container alive (`sleep infinity`).

Everything here is from this chapter: completions and parallelism, cron syntax, concurrency policies, `kubectl create job --from`, and what happens to a DaemonSet when you *don't* add a toleration.

Every check is a hit on the boss. Defeat it to unlock **Volumes and StatefulSets**.
