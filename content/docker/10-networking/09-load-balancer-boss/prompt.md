# Load Balancer

Boss fight! One copy of an app can only handle so many visitors. Run **two** copies behind one front door, and let the proxy share visitors between them. That's a **load balancer**.

In nginx, an `upstream` block lists several app containers, and nginx takes turns sending requests to each:

```nginx
upstream apps {
    server cq-app-a:8080;
    server cq-app-b:8080;
}
server {
    listen 80;
    location / {
        proxy_pass http://apps;
    }
}
```

# Assignment

In `~/docker-course`:

1. Start two copies of **`cq/webapp:1.6`** on **`cq-net`**, neither published:
   - **`cq-app-a`** with `GREETING="Server A"`
   - **`cq-app-b`** with `GREETING="Server B"`
2. Put the config above in **`lb/default.conf`**.
3. Start **`cq-lb`** from **`nginx:1.27-alpine`** on **`cq-net`**, published on port **8095**, with the **`lb`** folder mounted read-only at **`/etc/nginx/conf.d`**.
4. Refresh http://localhost:8095 a few times. It should switch between A and B.
5. Submit `dk.networking.load-balancer-boss`.
