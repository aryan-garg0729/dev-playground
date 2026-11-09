# Endpoints Exist, Traffic Fails

A new kind of broken. The pods are `1/1 Running`, and this time the Service **has** endpoints:

```sh
kubectl -n triage get endpointslices -l kubernetes.io/service-name=coupons
# NAME            ADDRESSTYPE   PORTS   ENDPOINTS
# coupons-x8f2k   IPv4          8080    10.244.1.31,10.244.2.17
kubectl exec storefront -- wget -qO- -T 3 http://coupons.triage
# wget: can't connect to remote host (10.96.77.4): Connection refused
```

The pods run `nginx:1.27-alpine`, which listens on port 80.

**What's wrong?**
