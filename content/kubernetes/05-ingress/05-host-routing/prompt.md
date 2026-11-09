# Routing by Host Name

Mo has been reading about subdomains: *"Can we just have `potions.elixirbytes.com` and `scrolls.elixirbytes.com`?"*

Every HTTP request carries a `Host` header with the name the client typed. Many sites can share one IP address and one port, and the proxy picks the site by that header. An Ingress rule can match on it with `host:`:

```yaml
spec:
  ingressClassName: nginx
  rules:
    - host: potions.elixirbytes.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: potions
                port:
                  number: 80
    - host: scrolls.elixirbytes.local
      http:
        ...
```

A rule with no `host` (like the `shop` Ingress) matches every host. A rule with a host only matches requests for that name, and wins over the no-host rule for that name.

We don't own `elixirbytes.local` in DNS, but we don't need to. `curl` can send any Host header to any address:

```sh
curl -H 'Host: potions.elixirbytes.local' http://localhost/
```

(Or add `127.0.0.1 potions.elixirbytes.local` to `/etc/hosts` if you want to try it in a browser. That's optional, and outside the course folder, so the course never does it for you.)

## Assignment

Create a **second** Ingress named **`shop-hosts`** (leave `shop` alone), class `nginx`, with two host rules:

- **`potions.elixirbytes.local`**, path `/` to Service `potions`
- **`scrolls.elixirbytes.local`**, path `/` to Service `scrolls`

No rewrite is needed this time: the path is already `/`.
