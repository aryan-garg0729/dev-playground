# Teaching the Cluster a New Noun

Every potion recipe lives in a different wiki, spreadsheet or someone's head. Priya wants them in the cluster, next to the services that brew them. Kubernetes has no `Recipe` kind. But you can add one.

A **CustomResourceDefinition** (CRD) teaches the API server a new kind. Afterwards a `Recipe` is as real as a `Pod`: stored in etcd, validated, covered by RBAC, and works with `kubectl get`, `describe`, `-o yaml`, `explain` and watch.

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: recipes.elixirbytes.io        # must be <plural>.<group>
spec:
  group: elixirbytes.io
  scope: Namespaced
  names:
    kind: Recipe
    plural: recipes
    singular: recipe
    shortNames: [rcp]
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              required: [potion, difficulty]
              properties:
                potion:      { type: string }
                difficulty:  { type: string, enum: [easy, tricky, legendary] }
                brewMinutes: { type: integer, minimum: 1, maximum: 120 }
      additionalPrinterColumns:        # extra columns for kubectl get
        - { name: Potion, type: string, jsonPath: .spec.potion }
        - { name: Difficulty, type: string, jsonPath: .spec.difficulty }
```

The schema is **OpenAPI**: types, required fields, enums, ranges. The API server enforces it on every create and update.

## Assignment

Save the CRD above as `manifests/recipe-crd.yaml` and apply it. Then look at what you got for free:

```sh
kubectl api-resources --api-group=elixirbytes.io
kubectl explain recipe.spec
kubectl get recipes        # "No resources found", and that's not an error
```

CRDs are cluster-scoped (the *kind* exists everywhere), even though each `Recipe` will live in a namespace.
