# Deploying on free plans (Render + Neon)

One Render web service runs the API and serves the built web app on the same
origin (`deploy/Dockerfile`, `WEB_DIR`). A free Render Key Value instance is
Redis. Postgres is on Neon (its free plan does not expire). Sign-in is Google
and GitHub OAuth, so there is no mail server.

What this setup does not run: the grading workers and the browser sandbox.
The current courses (Git, Docker, Kubernetes) are quiz and CLI lessons, which
are graded inside the API. Python, Go, JS and SQL code lessons would need those
services, which free plans cannot host (the Python worker needs a Docker daemon).

## 1. Database (Neon)

1. Create a project at neon.tech (free plan). Pick a region close to Render's
   `singapore` (for example AWS ap-southeast-1). Supabase works too: use its
   Session pooler string (port 5432, IPv4) with `?sslmode=require`.
2. Copy the connection string. It looks like
   `postgresql://user:pass@ep-xxx.us-west-2.aws.neon.tech/neondb?sslmode=require`.
   Use the direct one, not the `-pooler` host.

## 2. Render

1. render.com → New → Blueprint → connect the GitHub repo. Render reads
   `render.yaml` and creates `codequest` (web, free) and `codequest-redis`
   (Key Value, free).
2. When asked for values, fill in:
   - `APP_URL` and `ALLOWED_ORIGINS`: `https://codequest.onrender.com` (if the
     name is taken, Render shows the real URL on the service page after
     creation; update both then).
   - `DATABASE_URL`: the Neon string.
   - The OAuth ids and secrets from step 3. You need at least one provider;
     leave the other pair empty.
3. Each deploy builds the image; each start applies migrations and loads
   `content/` (a few seconds), then serves. `git push` to `main` redeploys.

## 3. OAuth apps

Google (console.cloud.google.com):
1. APIs & Services → OAuth consent screen: External, app name, your email.
   While it is in "Testing", only the test users you add can sign in; publish
   it for everyone (no review is needed for the `openid email profile` scopes).
2. Credentials → Create credentials → OAuth client ID → Web application.
   - Authorized JavaScript origin: `https://codequest.onrender.com`
   - Authorized redirect URI: `https://codequest.onrender.com/api/v1/auth/google/callback`
3. Put the client id and secret in `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET`.

GitHub (github.com/settings/developers → New OAuth App):
- Homepage URL: `https://codequest.onrender.com`
- Authorization callback URL: `https://codequest.onrender.com/api/v1/auth/github/callback`
- Generate a client secret; set `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET`.

## 4. Learner CLI

The site serves the CLI itself, so the repository can stay private. Each deploy
builds `codequest` for macOS, Linux and Windows (amd64 and arm64) with the
site's address baked in, from `APP_URL` (Render passes service environment
variables to Docker builds as build args; the build fails if it is missing or
not https). The image serves:

- `/cli/install`: the install page the setup lessons link to
- `/install.sh`: `curl -fsSL https://<site>/install.sh | sh` installs into `~/.local/bin`
- `/cli/codequest-<os>-<arch>[.exe]`: the binaries

After changing `APP_URL`, redeploy so the CLI is rebuilt for the new address.
Learners who installed an older build can point it at the new one with
`codequest config base_url https://<new-site>`.

## Free-plan limits

- The web service sleeps after 15 minutes without traffic; the next visit
  takes about a minute to wake. League closes and the daily rollover run only
  while it is awake, and catch up when it wakes.
- Free Key Value has 25 MB and no persistence: a restart clears rate-limit
  counters and the weekly league board cache, which are rebuilt as people use the app.
- Neon's free plan has 0.5 GB of storage and suspends when idle (first query
  after a pause takes a moment).
- Render's free plan blocks outbound mail ports, which no longer matters: there is no email.
