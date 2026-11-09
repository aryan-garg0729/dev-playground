# Down, or Down -v?

`docker compose down` cleans up the project. It's worth knowing exactly what goes:

| | `stop` | `down` | `down -v` |
| --- | --- | --- | --- |
| Containers stopped | Yes | Yes | Yes |
| Containers removed | No | Yes | Yes |
| Project network removed | No | Yes | Yes |
| **Named volumes removed** | No | **No** | **Yes** |
| Images removed | No | No | No |

**You want to remove every container of the stack for tonight, and get all the saved notes back when you run `docker compose up` tomorrow. Which command?**
