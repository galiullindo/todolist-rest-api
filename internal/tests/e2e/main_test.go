package e2e

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/galiullindo/todolist-rest-api/internal/config"
	"github.com/galiullindo/todolist-rest-api/internal/handler"
	"github.com/galiullindo/todolist-rest-api/internal/repository"
	"github.com/galiullindo/todolist-rest-api/internal/service"
	_ "github.com/mattn/go-sqlite3"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(run(m))
}

func run(m *testing.M) (code int) {
	var cliArguments struct {
		configPath string
	}

	args := flag.Args()
	fs := flag.NewFlagSet("e2e_test", flag.ContinueOnError)
	fs.StringVar(&cliArguments.configPath, "config", "config.yaml", "path to config file")
	if err := fs.Parse(args); err != nil {
		log.Printf("parse args: %v\n", err)
		return 2
	}

	cfg, err := config.Load(cliArguments.configPath)
	if err != nil {
		log.Printf("load config: %v\n", err)
		return 1
	}

	db, err := sql.Open("sqlite3", cfg.Database.URI)
	if err != nil {
		log.Printf("open db: %v\n", err)
		return 1
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		log.Printf("ping db: %v\n", err)
		return 1
	}
	if err := setSavepoint(db); err != nil {
		log.Println(err)
		return 1
	}

	defer func() {
		if err := releaseSavepoint(db); err != nil {
			log.Println(err)
			if code == 0 {
				code = 1
			}
		}
	}()

	defer func() {
		if err := rollbackToSavepoint(db); err != nil {
			log.Println(err)
			if code == 0 {
				code = 1
			}
		}
	}()

	testDB = db
	testServer = newTestServer(db)
	defer testServer.Close()

	testClient = newTestClient()

	return m.Run()
}

var (
	testDB     *sql.DB
	testServer *httptest.Server
	testClient *http.Client
)

func newTestServer(db *sql.DB) *httptest.Server {
	taskRepository := repository.NewRepository(db)
	taskService := service.NewService(taskRepository)
	taskHandler := handler.NewHandler(taskService)

	mux := http.NewServeMux()
	taskHandler.RegisterRoutes(mux)

	return httptest.NewServer(mux)
}

func newTestClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func setSavepoint(db *sql.DB) error {
	if _, err := db.Exec("SAVEPOINT e2e_test"); err != nil {
		return fmt.Errorf("set savepoint: %w", err)
	}
	return nil
}

func rollbackToSavepoint(db *sql.DB) error {
	if _, err := db.Exec("ROLLBACK TO e2e_test"); err != nil {
		return fmt.Errorf("rollback to savepoint: %w", err)
	}
	return nil
}

func releaseSavepoint(db *sql.DB) error {
	if _, err := db.Exec("RELEASE SAVEPOINT e2e_test"); err != nil {
		return fmt.Errorf("release savepoint: %w", err)
	}
	return nil
}
