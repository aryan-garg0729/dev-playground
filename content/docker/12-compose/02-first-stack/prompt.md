# Your First compose.yaml

`compose.yaml` is written in **YAML**, a simple format where **indentation matters** (use spaces, never tabs). Each container you want is a **service**:

```yaml
services:
  api:                      # the service's name
    image: cq/webapp:1.6    # same as the image in docker run
    ports:
      - "8097:8080"         # same as -p 8097:8080
```

The main commands, run from the folder that holds `compose.yaml`:

| Command | What it does |
| --- | --- |
| `docker compose up -d` | Create and start everything, in the background |
| `docker compose ps` | List this project's containers |
| `docker compose logs -f` | Follow the logs of every service at once |
| `docker compose down` | Stop and remove the containers and network it made |

Compose names things after your folder. In `~/docker-course` the container becomes `docker-course-api-1`, and Compose creates a network called `docker-course_default` that every service joins automatically.

# Assignment

In `~/docker-course`:

1. Create `compose.yaml`. Copy and paste:

   ```sh
   cat > compose.yaml <<'EOF'
   services:
     api:
       image: cq/webapp:1.6
       ports:
         - "8097:8080"
   EOF
   ```

2. Start it: `docker compose up -d`.
3. Check `docker compose ps`, then open http://localhost:8097.
4. Run, then submit `dk.compose.first-stack`.
