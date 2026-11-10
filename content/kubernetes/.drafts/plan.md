# Kubernetes course: authoring notes (no story)

Beginner course in boot.dev style: one idea per lesson (~1 minute read), then do it.
No story, no characters, no company. The learner never reads or writes app code and
never builds an image: every workload uses a public image, and every file they need
comes from one copy-paste block.

Audience: has never heard of Kubernetes. Knows what a container is (one-line refresher
in lesson 1; the Docker course is the suggested first step, not a requirement beyond
having Docker installed and running).

Rewrite of the old intermediate course (ElixirBytes story, 101 lessons, last version in
git at f88e68a). Every concept it taught is kept; see "Coverage" below.

## Cluster and workspace
- Workspace: `~/k8s-course` (run `codequest` from there). Manifests the learner writes go
  in `~/k8s-course/manifests/`. Observed values go in `~/k8s-course/answers.txt`
  (`key=value` lines).
- kind cluster `cq` (context `kind-cq`): one control plane + two workers
  (`cq-control-plane`, `cq-worker`, `cq-worker2`). Config `~/k8s-course/kind.yaml`,
  written in the install lesson. Host ports 80/443 -> ingress, 30080 -> NodePort.
- Everything the learner builds lives in namespace `course` (created in chapter 1, made the
  context's default). The troubleshooting chapter uses `broken`, seeded by a copy-paste block.
- Images: `nginx:1.27-alpine` and `nginx:1.26-alpine` (web pages and rollouts),
  `busybox:1.36` (workers, shells), `curlimages/curl:8.10.1` (API calls).
  Web pages show text from a ConfigMap, so "the app's version" is just HTML we control.

## Checks
- Every kubectl check pins the cluster: `kubectl --context kind-cq -n course ...`. A check
  must never read a learner's work cluster.
- Read-only: `get -o jsonpath`, `auth can-i`, `exec ... cat`, `curl` against localhost.
  Never `apply`, `delete`, `run`, `scale` in a check.
- Labelled output with jsonpath: `-o jsonpath='ready={.status.readyReplicas}'`.
- Pod names are random: select by label (`-l app=web`, `{.items[0]...}`) or `exec deploy/web`.
- Never check IPs, uids, timestamps, hashes or ages. Counts: `echo "n=$(... | wc -w | tr -d ' ')"`.
- Values that differ per learner (container runtime, node count) are captured from the
  cluster in one step and compared with `answers.txt` (`grep -x 'key=${cap}' answers.txt`).
- `${name}` is reserved for captures; jsonpath `{...}` is fine.
- Slow things (cluster create, rollouts, controller installs) are the learner's job; checks
  answer within 30 s.
- Known kind behaviour from the old course's recording notes (still true on v1.37):
  OOM kills may say `Error` with exit 137 (grade 137, not the reason); `logs --previous`
  can exit 0 with "unable to retrieve"; busybox `nslookup` exits 1 (grep `Name:`); use
  EndpointSlices, not Endpoints; untolerated taint message is generic; unbound-PVC pods are
  Pending; kindnet enforces NetworkPolicy (scope policies to their own pods); wait on the
  ingress controller pod, not the Deployment; `kubectl autoscale --cpu=50%`; never print a
  Secret (compare `cksum`); RBAC write checks use `--dry-run=server`.

## Bosses
`project_cli`, 5-8 test rows, a short checklist that combines only that chapter's
commands. Nothing new is taught in a boss.

## Coverage of the old course (nothing dropped)
clusters, control plane parts, nodes, contexts, namespaces, api-resources, declarative apply;
pods, exec, multi-container + emptyDir, sidecar, CrashLoopBackOff, labels/selectors;
ReplicaSets, quarantine by relabel, Deployments, pod-template-hash, scale, rollout, undo,
revisions, maxSurge/maxUnavailable; ClusterIP, EndpointSlices, DNS, kube-proxy, NodePort,
LoadBalancer; ingress controller, path and host routing, Gateway API; ConfigMaps (env,
volume, live update, restart), Secrets, base64 is not encryption, Secret volumes;
readiness, stuck rollout, liveness, requests/limits, OOM, QoS; Jobs, backoffLimit, CronJobs,
concurrencyPolicy, DaemonSets; StorageClasses, PVCs, access modes, data outlives the pod,
StatefulSets, stable identity, headless Services; nodeSelector, taints/tolerations, Pending,
anti-affinity, topology spread; ServiceAccounts, Roles, RoleBindings, resourceNames, API
from a pod, NetworkPolicy; metrics-server, HPA, CRDs, validation, operators; triage of
ImagePullBackOff, missing config, no endpoints, missing StorageClass, unschedulable.

New for beginners: a YAML/manifest chapter before workloads, `kubectl describe`/`logs`/
`port-forward` taught on their own, and self-healing shown before Deployments are explained.

## Ids
Chapter ids use the new prefix `kube.` (old ones were `k8s.`), so no new lesson shares an
id with an old one. Old lessons are unpublished by `cq build`; history is kept.
