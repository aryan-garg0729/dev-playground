# Every Two Minutes

The ledger has to be totalled every night. Someone has been doing it with a `crontab` on the old server, which is now the single most fragile thing at ElixirBytes.

A **CronJob** creates a Job on a schedule. The schedule is classic cron syntax: minute, hour, day of month, month, day of week.

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: nightly-ledger
spec:
  schedule: "*/2 * * * *"         # every 2 minutes (so we don't wait until night)
  concurrencyPolicy: Forbid       # don't start a new run while one is still going
  successfulJobsHistoryLimit: 3   # keep the last 3 successful Jobs
  jobTemplate:                    # a Job spec, as in lesson 1
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: ledger
              image: busybox:1.36
              command: ["sh", "-c", "echo totalling the ledger; sleep 5; echo done"]
```

The schedule uses the controller's time zone (UTC on most clusters) unless you set `spec.timeZone`. The Jobs it creates are named `nightly-ledger-<number>`.

Waiting for a schedule is boring when you're debugging. You can start a run right now from the CronJob's template:

```sh
kubectl create job ledger-manual --from=cronjob/nightly-ledger
```

## Assignment

1. Save the CronJob above as `manifests/nightly-ledger.yaml` and apply it.
2. Wait for it to fire at least once (`kubectl get jobs -w`; up to 2 minutes).
3. Trigger one run by hand, named **`ledger-manual`**, with the command above.
