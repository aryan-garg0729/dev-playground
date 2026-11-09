# The File Changes, the Variable Doesn't

Monday morning. Marketing wants the v2 banner, and the database team has moved brewd's database to port 6432. Let's change the ConfigMap and see what the running pods notice, **without restarting anything**.

When you update a ConfigMap:

- **Mounted files are updated in place.** The kubelet checks for changes periodically, so it takes up to a minute or two, not an instant. (Except files mounted with `subPath`, which never update.)
- **Environment variables never change.** They're set once, when the container process starts. There's no mechanism in Linux to change another running process's environment.

So a pod that read `DB_PORT` at startup keeps the old value until the pod is replaced. Whether the *app* notices a changed file is the app's business: some watch their config files, most read them once at startup.

## Assignment

1. Change `brewd-config` in `manifests/brewd-config.yaml`:
   - `db-port`: **`6432`**
   - `greeting`: **`Fresh potions daily v2`**
2. Apply **only that file**. Don't restart or roll out `brewd`.
3. Wait, then check in a running pod, until the file shows the new port:

   ```sh
   kubectl exec deploy/brewd -- sh -c 'echo "env=$DB_PORT file=$(cat /etc/brewd/db-port)"'
   ```

4. Record what that pod sees in `answers.txt`:

   ```
   env_db_port=<what $DB_PORT says>
   file_db_port=<what /etc/brewd/db-port says>
   ```

If you accidentally restarted brewd, the variable already shows 6432 too and the check can't see the difference. Set `db-port` back to `5432`, apply, run `kubectl rollout restart deploy/brewd` and wait for it, then set `6432` again and apply **only** the ConfigMap.
