# Containerise an App

This is what Docker is really for: taking an app and packing it with everything it needs to run.

We've written a small web app for you. It answers a few web addresses (it's written in Go, but **you don't need Go installed and you don't need to read the code**: Docker will compile it for you inside the build).

| Address | Answer |
| --- | --- |
| `/` | `Hello from webapp dev` |
| `/healthz` | `ok` (used to check the app is alive) |

It listens on port **8080** inside its container.

Here's the Dockerfile that packs it. This part *is* worth reading:

```dockerfile
# start from an image that has the Go tools
FROM golang:1.23-alpine
# work in /src inside the image
WORKDIR /src
# copy the app's files into /src
COPY app/ .
# compile them into a program called webapp
RUN go build -o /usr/local/bin/webapp .
# note: the app listens on 8080
EXPOSE 8080
# start the program when a container starts
CMD ["webapp"]
```

Notice the pattern: **pick a base, copy your files in, build them, say how to start.** Almost every Dockerfile you'll ever read follows it.

# Assignment

1. In `~/docker-course`, create the app. Copy and paste this whole block into your terminal. It creates `app/go.mod` and `app/main.go`:

   ```sh
   mkdir -p app
   cat > app/go.mod <<'EOF'
   module webapp

   go 1.23
   EOF
   cat > app/main.go <<'EOF'
   // webapp: a tiny web server for the Docker course. You don't need to read this.
   package main

   import (
   	"context"
   	"fmt"
   	"log"
   	"net/http"
   	"os"
   	"os/signal"
   	"path/filepath"
   	"strings"
   	"syscall"
   	"time"
   )

   func env(key, def string) string {
   	if v := os.Getenv(key); v != "" {
   		return v
   	}
   	return def
   }

   func main() {
   	version := env("APP_VERSION", "dev")
   	if len(os.Args) > 1 && os.Args[1] == "--version" {
   		fmt.Println("webapp", version)
   		return
   	}
   	greeting := env("GREETING", "Hello")
   	port := env("PORT", "8080")
   	notes := filepath.Join(env("DATA_DIR", "/data"), "notes.txt")

   	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
   		fmt.Fprintf(w, "%s from webapp %s\n", greeting, version)
   	})
   	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
   		fmt.Fprintln(w, "ok")
   	})
   	http.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
   		data, _ := os.ReadFile(notes)
   		w.Write(data)
   	})
   	http.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
   		text := strings.TrimSpace(r.URL.Query().Get("text"))
   		if text == "" {
   			http.Error(w, "add ?text=...", http.StatusBadRequest)
   			return
   		}
   		os.MkdirAll(filepath.Dir(notes), 0o755)
   		f, err := os.OpenFile(notes, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
   		if err != nil {
   			http.Error(w, err.Error(), http.StatusInternalServerError)
   			return
   		}
   		defer f.Close()
   		fmt.Fprintln(f, text)
   		w.WriteHeader(http.StatusCreated)
   		fmt.Fprintln(w, "saved:", text)
   	})

   	srv := &http.Server{Addr: ":" + port}
   	go func() {
   		log.Printf("webapp %s listening on :%s", version, port)
   		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
   			log.Fatal(err)
   		}
   	}()
   	stop := make(chan os.Signal, 1)
   	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
   	<-stop
   	log.Println("webapp: shutting down")
   	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   	defer cancel()
   	srv.Shutdown(ctx)
   }
   EOF
   ```

2. Create `Dockerfile` right in `~/docker-course` (not inside `app`). Copy and paste:

   ```sh
   cat > Dockerfile <<'EOF'
   FROM golang:1.23-alpine
   WORKDIR /src
   COPY app/ .
   RUN go build -o /usr/local/bin/webapp .
   EXPOSE 8080
   CMD ["webapp"]
   EOF
   ```

3. Build it as **`cq/webapp:v1`**. The folder is `.`, which means "this folder":

   ```sh
   docker build -t cq/webapp:v1 .
   ```

   The first build downloads the Go tools, so give it a minute.

4. Run a container named **`cq-app`** from `cq/webapp:v1`, detached, on your port **8085** (the app uses 8080 inside).
5. Open http://localhost:8085 and http://localhost:8085/healthz.
6. Run, then submit `dk.dockerfile.webapp`.
