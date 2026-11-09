# Same Name, Same Data

Each StatefulSet pod gets a DNS name of its own, through the headless Service:

```
<pod>.<service>.<namespace>.svc.cluster.local
vault-0.vault.shop.svc.cluster.local
```

And the Service's own name, `vault.shop.svc.cluster.local`, returns **every** ready pod's address instead of one virtual IP, so clients (or the vault members themselves) can see all the members.

```sh
kubectl exec storefront -- nslookup vault-0.vault
kubectl exec storefront -- nslookup vault.shop.svc.cluster.local
```

## Assignment

1. Read `vault-0`'s identity file: `kubectl exec vault-0 -- cat /data/id`.
2. **Delete the pod `vault-0`** and watch it come back with the same name. Read the file again: the new pod got the old pod's claim, so the file wasn't rewritten.
3. From `storefront`, look up `vault-1.vault.shop.svc.cluster.local`, and then the headless name `vault.shop.svc.cluster.local`.
4. Add one line to `answers.txt`:

   ```
   vault_addresses=<how many Address lines the headless name returns>
   ```

   (Count only the `Address:` lines under the answer, not the DNS server's own `Address:` line at the top, which has a `#53` or `:53` after it.)
