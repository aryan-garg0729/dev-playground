# Three Strikes

`on-failure:N` is the policy for things that *should* finish: a migration, a nightly report. Retry a few times in case it was a blip, then stop and let a human look.

All attempts share one container, so they share its logs too, and `RestartCount` counts the retries.

## Assignment

1. Run **`eb-flaky`** from `alpine:3.20`, detached, with the restart policy **on-failure, at most 3 retries**, running:

   ```sh
   sh -c 'echo "cauldron overflow at $(date +%T)"; exit 1'
   ```

2. Wait a few seconds. Docker starts it four times in total (once, plus three retries), then leaves it `exited`.
3. Submit from `~/elixirbytes`. The checks read its restart count and policy and count its log lines.
