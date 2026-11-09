# Black Moon Sale Switch

Friday, 08:55. The sale starts at 09:00. Priya's ticket:

> **Flip brewd into sale mode.**
> - The banner says `Black Moon sale is live`.
> - brewd reads a `SALE_MODE` environment variable at startup. It must be `on` in **every running pod**, from a new `sale-mode` key in `brewd-config`.
> - The payments partner issued an API token for the sale. Add it to the `brewd-db` Secret as `api-token`. brewd reads it from `/etc/creds/api-token`.
> - No dropped orders.
>
> Mo asked if we could "just put the token in the ConfigMap to save time". We can't.

## Assignment

1. Update `brewd-config`: `greeting` is **`Black Moon sale is live`**, and add **`sale-mode: "on"`**. (Quote it: unquoted `on` is a boolean in some YAML parsers.)
2. Add an env var **`SALE_MODE`** to brewd from that key.
3. Add the key **`api-token`** to the `brewd-db` Secret, with any value you like (`openssl rand -hex 20` will do). Keep `username` and `password`. To add one key without retyping the others, patch it in: `kubectl patch secret brewd-db --type merge -p '{"stringData":{"api-token":"..."}}'`.
4. Make sure every brewd pod is running with all of it, using a rolling update.

Every check is a hit on the boss. Defeat it to unlock **Probes and resources**.
