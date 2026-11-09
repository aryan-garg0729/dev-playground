# What a Controller Adds

Mo, delighted: *"So now if I create a Recipe, a brewer pod starts up and brews it?"*

Try it. Create a Recipe and watch `kubectl get pods`. Nothing happens. The API server stored your object, and that's all.

Everything that *does* things in Kubernetes is a **controller**, a program that watches objects and reacts. The ReplicaSet controller watches ReplicaSets. The Deployment controller watches Deployments. Nobody is watching Recipes.

**What would you have to add so that creating a Recipe starts a brewer pod, and deleting it stops the pod?**
