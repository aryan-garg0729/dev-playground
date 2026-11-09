# Recipe Book and Rush Hour

Two tickets landed on the same morning. Do both.

> **From Dex: menu-api falls over every Monday.** Autoscale it: never fewer than **2** pods, never more than **6**, aiming for **60%** CPU. (Hint from Dex: "last time someone tried, the HPA said `<unknown>` for a week.")

> **From Priya: recipes need ingredients.** A recipe with no ingredients is useless. Make `ingredients` a **required** list of strings with **at least one** item, and add a **`Brew`** column to `kubectl get recipes` showing `brewMinutes`. Then add our first real recipe: **`elixir-of-focus`** (potion `Elixir of Focus`, difficulty `tricky`, brewMinutes `25`, ingredients `sage`, `moonwater`, `owl feather`).

## Assignment

1. Give menu-api what it needs for CPU autoscaling, then create an HPA named **`menu-api`**.
2. Update the CRD. In OpenAPI, a list of strings is `{ type: array, items: { type: string } }`, and `minItems` sets a minimum length. A field is required when it's listed in its parent's `required`.
3. Existing recipes (`dragonbrew`, `moonmilk`) have no ingredients. They stay stored, but the next update to either would be rejected, so add ingredients to them too (any you like) and re-apply.
4. Create `elixir-of-focus`.

Every check is a hit on the boss. Defeat it to unlock **Troubleshooting**, the last chapter.
