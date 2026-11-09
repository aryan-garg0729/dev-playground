# Release Night for the Orders API

The orders service is moving to the cluster tonight, and Priya wants it done properly:

> **orders: deploy, then upgrade, with a rollback ready.** Start it on the old version (`nginx:1.27-alpine`) so we have something to roll back to, then roll it forward to `nginx:1.28-alpine`. Four pods. During rollouts, it may go **2 pods over** and **1 pod under** the count: it's faster, and orders can take a short dip. Keep only the last **5** ReplicaSets around.
>
> Dex saw one pod logging weird timeouts after the upgrade. Pull one 1.28 pod out of rotation for him, alive, and let the Deployment replace it.

## Assignment

1. Write **`manifests/orders.yaml`**: a Deployment named **`orders`**, label and selector **`app: orders`**, 4 replicas, container `orders` on **`nginx:1.27-alpine`**, port 80, a `RollingUpdate` strategy with **`maxSurge: 2`** and **`maxUnavailable: 1`**, and **`revisionHistoryLimit: 5`**. Apply it and wait for it to be ready.
2. Roll it forward to **`nginx:1.28-alpine`**, and update the file to match.
3. Quarantine **one** of the new pods by relabelling it **`app=orders-quarantine`**.

Use `kubectl explain deployment.spec` if you've forgotten where a field goes.

Every check is a hit on the boss. Defeat it to unlock **Services and DNS**.
