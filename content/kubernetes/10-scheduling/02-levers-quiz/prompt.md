# Pulling and Pushing Pods

Five levers decide where pods go. Three are the **pod** choosing; one is the **node** refusing:

| lever | who decides | says |
| --- | --- | --- |
| `nodeSelector` | pod | "only nodes with these labels" |
| node affinity | pod | the same, with `In`/`NotIn`/`Exists`, and soft `preferred...` rules |
| pod (anti-)affinity | pod | "near / away from pods with these labels" |
| taints + tolerations | node | "no pods here unless they tolerate my taint" |
| `topologySpreadConstraints` | pod | "spread my replicas evenly across nodes or zones" |

A **toleration only permits**. It lets a pod onto a tainted node; it doesn't send it there.

ElixirBytes buys one node with a GPU for the potion-recommendation model. It's expensive. Priya wants the ML pods on it, and **nothing else**.

**What's the right setup?**
