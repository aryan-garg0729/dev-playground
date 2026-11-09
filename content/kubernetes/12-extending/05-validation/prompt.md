# The Schema Says No

A custom resource looks like any other object:

```yaml
apiVersion: elixirbytes.io/v1
kind: Recipe
metadata:
  name: dragonbrew
spec:
  potion: Dragonbrew
  difficulty: legendary
  brewMinutes: 90
```

```sh
kubectl get recipes
# NAME         POTION       DIFFICULTY
# dragonbrew   Dragonbrew   legendary
kubectl get rcp dragonbrew -o yaml
```

Even with no code reacting to recipes, the schema is already worth a lot: nobody can store a recipe with a typo'd difficulty or a 9000-minute brew, because the API server refuses it before it ever reaches etcd, with a message saying which field and why.

## Assignment

1. Create two Recipes in `shop`, in `manifests/recipes.yaml`:
   - **`dragonbrew`**: potion `Dragonbrew`, difficulty `legendary`, brewMinutes `90`
   - **`moonmilk`**: potion `Moon Milk`, difficulty `easy`, brewMinutes `5`
2. Try to create a third with `difficulty: maybe`, and another with `brewMinutes: 500`. Read both refusals.
3. Add one line to `answers.txt`:

   ```
   rejected_field=<the field named in the brewMinutes refusal, as the message writes it>
   ```

The checks test the schema with `--dry-run=server`: the API server validates the object exactly as for a real create, then throws it away. It's a handy way to test a manifest against a cluster without changing anything.
