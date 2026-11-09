# When a Run Overruns

`concurrencyPolicy` decides what a CronJob does when it's time for the next run and the previous one hasn't finished:

| policy | if the last run is still going... |
| --- | --- |
| `Allow` (default) | start another one anyway |
| `Forbid` | skip this run |
| `Replace` | stop the old run and start the new one |

The ledger CronJob runs every hour with `concurrencyPolicy: Forbid`. On Black Moon night there are so many orders that the 00:00 run takes 80 minutes.

**What happens at 01:00?**
