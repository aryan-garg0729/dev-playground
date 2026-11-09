# A Cluster on Your Laptop

Black Moon weekend took ElixirBytes down twice. brewd runs in a container now (thanks to you), but it's one container on one server. When that server sneezed, the shop went dark, and Dex had to SSH in at 2am and type `docker start` by hand.

Priya's ticket is short: *"Move us to Kubernetes. Don't touch a cloud account until it works on your laptop."*

Kubernetes (k8s for short: k, eight letters, s) runs containers across a **cluster** of machines and keeps them running without you. We'll build a real one on your laptop with **kind** ("Kubernetes in Docker"): each "machine", or **node**, is a Docker container.

## How these lessons work

Same as the Docker course: you work in your terminal, then the `codequest` CLI runs a few **read-only** checks (it shows each command and asks first). Every check names the cluster explicitly with `--context kind-elixirbytes`, so it can never read some other cluster you have access to.

Two habits for the whole course:

- Run `codequest` from **`~/elixirbytes/k8s`**.
- Use the exact names each lesson gives. The checks find your work by name.

## Assignment

1. Install Docker (Docker Desktop, OrbStack or Colima on macOS; Docker Engine on Linux) and give it at least **4 CPUs and 6 GB** of memory. Then install the tools:

   ```sh
   brew install kind kubectl     # Linux: see kind.sigs.k8s.io and kubernetes.io/docs/tasks/tools
   go install github.com/aryangarg/codequest/api/cmd/codequest@latest
   codequest login
   ```

2. Create the course folder and the cluster's config file. Three nodes on purpose: one **control plane** (the brain) and two **workers** (where your containers run). The port mappings let you reach the cluster from your laptop in later chapters.

   ```sh
   mkdir -p ~/elixirbytes/k8s/manifests && cd ~/elixirbytes/k8s
   cat > kind.yaml <<'YAML'
   kind: Cluster
   apiVersion: kind.x-k8s.io/v1alpha4
   name: elixirbytes
   nodes:
     - role: control-plane
       kubeadmConfigPatches:
         - |
           kind: InitConfiguration
           nodeRegistration:
             kubeletExtraArgs:
               node-labels: "ingress-ready=true"
       extraPortMappings:
         - { containerPort: 80, hostPort: 80, protocol: TCP }
         - { containerPort: 443, hostPort: 443, protocol: TCP }
         - { containerPort: 30080, hostPort: 30080, protocol: TCP }
     - role: worker
     - role: worker
   YAML
   ```

3. Create the cluster (it takes about a minute):

   ```sh
   kind create cluster --config kind.yaml
   kubectl get nodes
   # NAME                        STATUS   ROLES           AGE   VERSION
   # elixirbytes-control-plane   Ready    control-plane   60s   v1.3x.x
   # elixirbytes-worker          Ready    <none>          45s   v1.3x.x
   # elixirbytes-worker2         Ready    <none>          45s   v1.3x.x
   ```

   Workers can show `NotReady` for the first 20 seconds or so. Wait for all three to say `Ready`.

4. Run `codequest run k8s.cluster.setup`, then `codequest submit k8s.cluster.setup`.

If `kind create` fails with `address already in use`, something else on your laptop is using port 80, 443 or 30080 (another kind cluster, a local web server). Stop it and try again.

**Pausing for the day** is safe: stopping Docker (or Colima) pauses the cluster and everything comes back when you start it again. `kind delete cluster --name elixirbytes` is **not** safe: it deletes the cluster and everything you built in it.
