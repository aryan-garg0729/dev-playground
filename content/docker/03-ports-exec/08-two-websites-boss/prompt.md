# Two Websites, One Laptop

Boss fight! Run two separate websites side by side from the same image, each with its own page and its own port.

# Assignment

In `~/docker-course`:

1. Start **`cq-site-a`** from **`nginx:1.27-alpine`**, detached, reachable at **http://localhost:8082**.
2. Start **`cq-site-b`** from **`nginx:1.27-alpine`**, detached, reachable at **http://localhost:8083**, with the setting **`SITE_NAME=b`**.
3. Change each site's `index.html` (it lives in `/usr/share/nginx/html` inside the container):
   - Site A's page must contain **`Site A is live`**.
   - Site B's page must contain **`Site B is live`**.
4. Open both in your browser to check.
5. Submit `dk.ports.two-websites-boss`.

Stuck? Look back at "Run a Website" for the ports and "Run Commands Inside" for changing a page.
