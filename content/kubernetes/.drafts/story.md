# Kubernetes course: story bible (ElixirBytes)

## World
Same company as the Git and Docker courses: **ElixirBytes**, the online potion shop.
brewd is containerised now, and one box is no longer enough: Black Moon sales took the
shop down twice last month. Priya has approved a move to Kubernetes. You are the engineer
who gets the cluster, the pager and the blame.

To keep the course standalone, workloads use public stand-in images
(`nginx:1.27-alpine` for anything that serves HTTP, `busybox:1.36` for workers and
throwaway shells, `curlimages/curl:8.10.1` for API calls). Prompts say once that in
production these would be `elixirbytes/brewd`.

## Characters (same as Git and Docker)
- **Priya**: platform lead, writes the tickets. Precise, dry.
- **Dex**: on-call SRE. Grumpy, right. Owns the triage chapter.
- **Mo**: product manager, "can we just..." questions (just one more LoadBalancer, just put
  the password in the ConfigMap).

## Cluster and workspace
- A local kind cluster named **`elixirbytes`** (context **`kind-elixirbytes`**): 1 control
  plane + 2 workers (`elixirbytes-control-plane`, `elixirbytes-worker`,
  `elixirbytes-worker2`). Host ports 80/443 -> ingress, 30080 -> NodePort lessons.
  Config written in lesson 1 to `~/elixirbytes/k8s/kind.yaml`.
- Learner works in **`~/elixirbytes/k8s`**, all course: `manifests/` (their YAML, by
  chapter) and `answers.txt` (`key=value` lines of things they observed).
- Everything the learner builds lives in namespace **`shop`** (created ch1). The triage
  chapter uses **`triage`**, seeded from a manifest in its own prompt.

## Check style
- Every kubectl check pins the cluster: `kubectl --context kind-elixirbytes -n shop ...`.
  This is a safety rule: a check must never read a learner's work cluster.
- Read-only only: `get -o jsonpath`, `auth can-i`, `exec ... cat` (reading a file),
  `curl` against localhost. Never `apply`, `delete`, `run`.
- Print labelled values with jsonpath: `-o jsonpath='ready={.status.readyReplicas} image={.spec.template.spec.containers[0].image}'`.
- Pod names are random: select by label (`-l app=storefront`) and use `{.items[0]...}`, or
  `kubectl exec deploy/<name>`.
- Never check IPs, uids, timestamps, hashes or ages. Counts: `echo "n=$(... | wc -w | tr -d ' ')"`.
- Observed values (runtime, QoS, pending reason) are captured from the cluster in one step
  and compared with `answers.txt` with `grep -x 'key=${cap}' answers.txt` in the next.
- `${name}` is reserved for captures; jsonpath `{...}` is fine.
- The learner runs slow things (cluster create, rollouts, helm). Checks are fast (30 s).

## Bosses
`project_cli`, 5-8 test rows, a realistic ticket from Priya or Dex combining only that
chapter's skills. Objects are named after the ticket.

## Workspace state per chapter (namespace shop unless noted)
- ch1: cluster; ns `shop` (default ns of the context); answers nodes/runtime/version.
- ch2: pods `storefront`, `menu-board` (2 containers, emptyDir `shared`), `cauldron`
  (crashlooping, then deleted), labels.
- ch3: ReplicaSet `brewd-rs` + quarantined pod; Deployment `brewd` (3 replicas,
  maxSurge 1 / maxUnavailable 0, revision >= 3 on nginx:1.27-alpine).
- ch4: Services `brewd` (ClusterIP), `brewd-np` (NodePort 30080), `brewd-lb` (pending).
- ch5: ingress controller; apps `potions`, `scrolls`; Ingress `shop`.
- ch6: ConfigMap `brewd-config`, Secret `brewd-db`, both wired into `brewd`.
- ch7: probes on `brewd`; Deployment `steady` (Guaranteed); pod `oom-canary`.
- ch8: Job `restock`, CronJob `nightly-ledger`, DaemonSet `node-scout`.
- ch9: PVC `orders`, Deployment `ledger-writer`, StatefulSet + headless Service `vault`.
- ch10: node label `tier=fast` on worker, taint on worker2; Deployments `pinned`,
  `tolerant`, `spread`; pod `stubborn`.
- ch11: SA `menu-reader`, Role `pod-reader`, RoleBinding `menu-reader-pods`, pod `api-caller`.
- ch12: metrics-server; Deployment + HPA `rush`; CRD `recipes.elixirbytes.io`, Recipes.
- ch13: ns `triage`, five broken workloads fixed.

## Verified on kind (Kubernetes v1.37.0, containerd 2.3.4, kindnet 2026-08), 2026-11-09
Every solution.json is a recorded run against a real cluster built from lesson 1's kind.yaml,
except k8s.storage.claim (hand-built from real output: it must be recorded while the PVC is
still Pending, before lesson 3 binds it). Things that surprised us, and how lessons handle them:
- OOM kills often report `Reason: Error` (not `OOMKilled`) with exit 137 on this containerd.
  The OOM lesson grades exit 137 + memory limit + no probes, never the reason string.
- `kubectl logs --previous` exits 0 with "unable to retrieve container logs" right after a
  crash. Checks read both `--previous` and current logs and grep the line; the cauldron pod
  stores its random numbers in an emptyDir so every run prints the same line.
- busybox `sh` as PID 1 ignores SIGTERM, so liveness restarts take ~65 s (30 s grace) and
  exit 137. The liveness prompt explains this.
- busybox `nslookup` exits 1 and prints NXDOMAIN lines from the search list; checks look
  for the `Name:` line only.
- Endpoints API is deprecated: checks use EndpointSlices
  (`-l kubernetes.io/service-name=<svc>`, `endpoints[?(@.conditions.ready==true)]`).
- The scheduler message no longer names the taint: "N node(s) had untolerated taint(s)".
- A pod whose PVC is unbound is `Pending` ("unbound immediate PersistentVolumeClaims"),
  not ContainerCreating.
- kindnet enforces NetworkPolicy, so ch11 has a hands-on policy lesson. Policies are scoped
  to vault-proxy and status-page only; a namespace-wide default-deny would break later checks.
- ingress-nginx: wait on the controller pod (`condition=ready`), not the Deployment's
  `available`, which reports too early. Project retired March 2026; Gateway API is a later add.
- `kubectl autoscale --cpu-percent` is deprecated; use `--cpu=50%`.
- Secret checks never print a secret: they compare `cksum` of the mounted file and the decoded
  Secret, and report only lengths.
- Destructive-looking API checks use `?dryRun=All` / `--dry-run=server`, so a wrong RBAC grant
  can't make a check delete anything.

## Lesson ordering that matters
- ch9 lesson 2 (`orders_phase=Pending`) must be submitted before lesson 3 binds the claim.
- ch10 leaves worker2 tainted until its boss removes the taint; ch10 boss frees `stubborn`.
- ch6 lesson 3 needs a running brewd pod started before the ConfigMap change.
