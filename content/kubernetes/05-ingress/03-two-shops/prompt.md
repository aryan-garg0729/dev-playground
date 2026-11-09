# Two Shops Behind One Door

Before we route to two services, we need two services that look different, so you can tell which one answered. Every service so far has served the same nginx welcome page.

busybox has a tiny web server built in (`httpd`), so a pod can write its own page and serve it:

```yaml
containers:
  - name: web
    image: busybox:1.36
    command: ["sh", "-c", "mkdir -p /www && echo 'potions shop' > /www/index.html && exec httpd -f -v -p 80 -h /www"]
    ports:
      - containerPort: 80
```

- `-f` keeps httpd in the foreground (a container's main process must not exit), and `-v` logs every request.
- `-p 80` is the port, and `-h /www` is the folder it serves.
- `exec` replaces the shell with httpd, so httpd gets the signals when the pod stops (the PID 1 lesson from the Docker course).

## Assignment

For each of **`potions`** and **`scrolls`**:

1. A Deployment with that name, 2 replicas, label/selector `app: <name>`, running the container above, with its page saying **`potions shop`** or **`scrolls shop`**.
2. A ClusterIP Service with the same name, port 80 to target port 80.

Put both in `manifests/shops.yaml`. Check each one from the storefront pod: `kubectl exec storefront -- wget -qO- http://potions`.
