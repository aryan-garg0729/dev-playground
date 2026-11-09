# Ingress and the Gateway API

Ingress was designed in 2015 and covers the basics: hosts, paths, TLS. Everything else (rewrites, timeouts, canaries, auth) ended up in **annotations** that each controller invents for itself. That's why your Ingress says `nginx.ingress.kubernetes.io/rewrite-target`, and why it does nothing on a different controller.

The **Gateway API** is its successor. The same routing as this chapter, sketched:

```yaml
kind: Gateway            # owned by the platform team: the listener on port 80
metadata: { name: shop-gw }
spec:
  gatewayClassName: envoy
  listeners: [{ name: http, port: 80, protocol: HTTP }]
---
kind: HTTPRoute          # owned by the app team: the routing
spec:
  parentRefs: [{ name: shop-gw }]
  rules:
    - matches: [{ path: { type: PathPrefix, value: /potions } }]
      filters:
        - type: URLRewrite         # a real field, not an annotation
          urlRewrite: { path: { type: ReplacePrefixMatch, replacePrefixMatch: / } }
      backendRefs: [{ name: potions, port: 80 }]
```

Meanwhile the ingress-nginx project you installed was retired in March 2026. Dex's question in standup: *"So is all our Ingress YAML garbage now?"*

**Which statement is correct?**
