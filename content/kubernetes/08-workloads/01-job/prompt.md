# Work That Ends

Not everything is a web server. Every night ElixirBytes restocks five shelves of the warehouse database. Each restock runs for a few seconds and finishes. Run that as a Deployment and you'd be in trouble: the container exits, the Deployment "fixes" it by restarting it, and it restocks forever.

A **Job** runs pods until a number of them **succeed** (exit 0), then stops:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: restock
spec:
  completions: 5        # 5 successful pods in total
  parallelism: 2        # at most 2 running at once
  backoffLimit: 3       # give up after 3 retries
  template:
    spec:
      restartPolicy: OnFailure    # Jobs can't use Always
      containers:
        - name: restock
          image: busybox:1.36
          command: ["sh", "-c", "echo restocked shelf $(hostname); sleep 3"]
```

```sh
kubectl get jobs
# NAME      STATUS     COMPLETIONS   DURATION
# restock   Running    2/5           8s
kubectl get pods -l job-name=restock     # the Job labels its pods for you
```

The finished pods stay around as `Completed`, so you can still read their logs. Deleting the Job deletes them too. (`ttlSecondsAfterFinished` can clean up automatically; we'll keep them so you can look.)

## Assignment

Save the Job above as `manifests/restock.yaml`, apply it and watch the pods: two at a time, five in total. A Job's pod template can't be changed once it exists, so if you need to fix it, delete the Job and apply again.
