# Where Data Can Live

ElixirBytes now keeps an order log. Right now it's written to the container's own filesystem, which is thrown away every time the pod is replaced: every deploy, every crash, every node drain. `emptyDir` from chapter 2 doesn't help either: it lives exactly as long as the pod.

Kubernetes splits durable storage into three objects:

- A **PersistentVolume (PV)** is a piece of real storage: a cloud disk, an NFS share, a folder on a node.
- A **PersistentVolumeClaim (PVC)** is a request for storage: "1 GiB, mountable by one node". Pods mount claims, never volumes directly.
- A **StorageClass** says *how* to create volumes when a claim asks for one. This is **dynamic provisioning**: nobody pre-creates disks.

```sh
kubectl get storageclass          # short name: sc
kubectl describe sc <name>
```

Columns worth knowing:

- **PROVISIONER**: what creates the volumes. On a cloud it's the cloud's disk driver; on kind it's `local-path`, which makes a folder on a node.
- **RECLAIMPOLICY**: what happens to the volume and its data when the claim is deleted (`Delete` or `Retain`).
- **VOLUMEBINDINGMODE**: `Immediate` creates the volume as soon as the claim exists. `WaitForFirstConsumer` waits until a pod uses the claim, so the volume is created on (or near) the node that pod was scheduled to.

The class marked `(default)` is used by any claim that doesn't name one.

## Assignment

Look at your cluster's storage classes and add three lines to `answers.txt`:

```
default_sc=<the name of the default StorageClass>
sc_binding=<its volume binding mode>
sc_reclaim=<its reclaim policy>
```
