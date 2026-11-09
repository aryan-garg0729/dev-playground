# Meet brewd

Time to meet the service. **brewd** is about 80 lines of Go that take potion orders over HTTP:

| Request | Response |
| --- | --- |
| `GET /healthz` | `ok` |
| `GET /version` | `brewd <version>` |
| `POST /orders?potion=mana` | records the order in `$DATA_DIR/orders.log` (default `/data`) |
| `GET /orders` | every order so far, one per line |

It listens on port **8080**, `brewd --version` prints its version, and it shuts down cleanly when it receives SIGTERM.

**You don't need Go installed.** Docker compiles it for you in the next lesson. You also don't need to understand the Go code: we're here to containerise brewd, not to write it.

## Assignment

1. In `~/elixirbytes`, create **`app/go.mod`**:

   ```
   module elixirbytes/brewd

   go 1.23
   ```

2. Create **`app/main.go`**:

   ```go
   // brewd: the ElixirBytes potion-order API.
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

   var version = "dev"

   func main() {
   	if len(os.Args) > 1 && os.Args[1] == "--version" {
   		fmt.Println("brewd", version)
   		return
   	}

   	dataDir := os.Getenv("DATA_DIR")
   	if dataDir == "" {
   		dataDir = "/data"
   	}
   	orders := filepath.Join(dataDir, "orders.log")

   	mux := http.NewServeMux()
   	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
   		fmt.Fprintln(w, "ok")
   	})
   	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
   		fmt.Fprintln(w, "brewd", version)
   	})
   	mux.HandleFunc("GET /orders", func(w http.ResponseWriter, r *http.Request) {
   		data, _ := os.ReadFile(orders)
   		w.Write(data)
   	})
   	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
   		potion := strings.TrimSpace(r.URL.Query().Get("potion"))
   		if potion == "" {
   			http.Error(w, "missing ?potion=", http.StatusBadRequest)
   			return
   		}
   		if err := os.MkdirAll(dataDir, 0o755); err != nil {
   			http.Error(w, err.Error(), http.StatusInternalServerError)
   			return
   		}
   		f, err := os.OpenFile(orders, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
   		if err != nil {
   			http.Error(w, err.Error(), http.StatusInternalServerError)
   			return
   		}
   		defer f.Close()
   		fmt.Fprintln(f, potion)
   		w.WriteHeader(http.StatusCreated)
   		fmt.Fprintln(w, "ordered", potion)
   	})

   	srv := &http.Server{Addr: ":8080", Handler: mux}
   	go func() {
   		log.Printf("brewd %s listening on :8080", version)
   		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
   			log.Fatal(err)
   		}
   	}()

   	stop := make(chan os.Signal, 1)
   	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
   	<-stop
   	log.Println("brewd: shutting down")
   	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   	defer cancel()
   	srv.Shutdown(ctx)
   }
   ```

3. The QA team keeps a 5 MB file of test fixtures next to the code. brewd never reads it, but it's there, the way junk always is in a real repo. Create it:

   ```sh
   mkdir -p app/testdata
   head -c 5000000 /dev/urandom > app/testdata/fixtures.bin
   ```

4. Submit from `~/elixirbytes`. The checks read `go.mod`, look for three landmarks in `main.go` and list `app/testdata`.
