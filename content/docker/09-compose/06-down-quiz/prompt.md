# down vs down -v

| Command | Containers | Network | Named volumes |
| --- | --- | --- | --- |
| `docker compose stop` | stopped, kept | kept | kept |
| `docker compose down` | **removed** | **removed** | kept |
| `docker compose down -v` | removed | removed | **removed** |

Dex once ran a "quick cleanup" on the staging shop and spent the evening restoring orders from backups. Ever since, the team's runbook says: *know which command you're typing.*

**Which of these are true for our stack?** Pick every one that is.
