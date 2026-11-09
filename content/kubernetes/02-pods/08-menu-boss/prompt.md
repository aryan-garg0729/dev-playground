# Potion of the Day Board

Priya's ticket for the end of week one:

> **Potion of the Day board, as a proper manifest.** One pod, two containers. The web server serves whatever the writer last wrote; the writer picks the potion. Label it so `-l tier=frontend` finds it, and put your name on it so I know who to ask. While you're there, get rid of the cauldron pod: it's been crashing all week and it's spamming Dex's dashboard.

## Assignment

Write **`manifests/potion-menu.yaml`** and apply it:

- A pod named **`potion-menu`**, labels **`app: potion-menu`** and **`tier: frontend`**, annotation **`elixirbytes.io/owner`** with your name.
- An `emptyDir` volume named **`pages`**.
- Container **`web`** (`nginx:1.27-alpine`, port 80) mounting `pages` at **`/usr/share/nginx/html`**.
- Container **`writer`** (`busybox:1.36`) mounting `pages` at **`/out`**, writing an `index.html` there that contains the text **`Potion of the Day`** and then staying alive. You write this command yourself; the one from lesson 4 is a good model.

Then **delete the `cauldron` pod**.

Check it the way you'd check anything: `kubectl get pods` for `2/2`, `describe` if not, and a port-forward to see the page.

Every check is a hit on the boss. Defeat it to unlock **ReplicaSets and Deployments**.
