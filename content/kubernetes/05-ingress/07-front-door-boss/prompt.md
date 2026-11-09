# The Shop Front Door

Priya's ticket for launch day:

> **One public name: `shop.elixirbytes.local`.** `/` is brewd, `/potions` is the potions shop, `/scrolls` is the scrolls shop, and anything under those paths goes to the same shop (so `/potions/today` works too). The subdomains from last week stay. The anonymous `shop` Ingress that answers for *any* host has to go: I don't want random host names reaching our backends.

## Assignment

1. Create an Ingress named **`front-door`**, class `nginx`, with **one** rule for host **`shop.elixirbytes.local`** and three `Prefix` paths:
   - `/` to Service **`brewd`**, port 80
   - `/potions` to Service **`potions`**, port 80
   - `/scrolls` to Service **`scrolls`**, port 80
2. Make the shop pods receive `/` no matter which of those paths was requested (you've done this before).
3. Delete the **`shop`** Ingress. Keep **`shop-hosts`**.

When several `Prefix` paths match, the longest one wins, so `/potions/today` goes to potions even though `/` matches too.

Test it with `curl -H 'Host: shop.elixirbytes.local' http://localhost/...`, exactly as the boss will.

Every check is a hit on the boss. Defeat it to unlock **ConfigMaps and Secrets**.
