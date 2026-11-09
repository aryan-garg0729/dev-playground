# Locked Doors

RBAC controls who can talk to the **API server**. It says nothing about pod-to-pod traffic: by default, every pod in the cluster can reach every other pod, in every namespace. The vault proxy holds secrets, and right now the storefront could connect to it.

A **NetworkPolicy** is a firewall rule for pods, written with labels:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: vault-proxy-ingress
spec:
  podSelector:                 # which pods this policy protects
    matchLabels: { app: vault-proxy }
  policyTypes: [Ingress]       # it restricts incoming traffic
  ingress:
    - from:
        - podSelector:         # allowed: pods in this namespace with this label
            matchLabels: { role: vault-client }
```

The rules, which mirror RBAC:

- A pod that **no** policy selects accepts all traffic.
- As soon as **any** policy selects a pod for Ingress, that pod accepts **only** what some policy allows. Policies only ever add allowances, so there's no "deny" rule here either.
- A policy with an empty `ingress:` list (or none) allows nothing: that's "default deny".

**Who enforces this?** Not Kubernetes itself. The **CNI plugin** does. kind's kindnet enforces NetworkPolicies, as do Calico and Cilium. Some plugins don't, and on those your policy is accepted and silently does nothing. Always test a policy from both sides.

## Assignment

1. Create a pod named **`vault-client`** (`busybox:1.36`, `sleep infinity`) with the label **`role: vault-client`**.
2. Apply the NetworkPolicy **`vault-proxy-ingress`** above.
3. Test both sides, using a vault-proxy pod's IP (`kubectl get pods -l app=vault-proxy -o wide`):

   ```sh
   kubectl exec storefront -- wget -qO- -T 3 http://<pod-ip>      # should time out
   kubectl exec vault-client -- wget -qO- -T 3 http://<pod-ip>    # should get the page
   ```

We only lock down `vault-proxy`. A default-deny for the whole `shop` namespace is a good idea in production, but here it would break the other chapters' checks.
