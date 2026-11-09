# Labels and Selectors

Kubernetes doesn't group things by folder or by name. It groups them by **labels**: short `key: value` tags on any object. A **selector** is a query over labels, and it's how almost everything finds everything else: Services find their pods by selector, controllers count their pods by selector, and so do you.

```sh
kubectl label pod storefront tier=frontend            # add a label
kubectl label pod storefront tier=backend --overwrite # change one
kubectl label pod storefront tier-                    # remove one (trailing minus)
kubectl get pods --show-labels
kubectl get pods -l tier=frontend                     # equality
kubectl get pods -l 'tier=frontend,app!=storefront'   # AND, not-equal
kubectl get pods -l 'status in (quarantined,broken)'  # set-based
```

Labels are for **selecting**, so keep them short and meaningful. For everything else (who owns it, a ticket link, a git commit), use **annotations**. They look the same but can hold long values, and nothing selects on them:

```sh
kubectl annotate pod storefront elixirbytes.io/owner="Your Name"
```

The `elixirbytes.io/` prefix is a convention: it marks the key as ours, so it can't clash with keys Kubernetes or other tools add.

## Assignment

Dex wants to be able to find "everything customer-facing" with one command.

1. Label `storefront` and `menu-board` with **`tier=frontend`**.
2. Label `cauldron` with **`tier=backend`** and **`status=quarantined`**.
3. Annotate `storefront` with **`elixirbytes.io/owner`** set to your name.
4. Try the selectors from above until you can list each group on its own.

Labels you add with `kubectl label` live only in the cluster. Add them to your manifests too, or the next `apply` won't know about them.
