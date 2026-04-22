package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	iofs "io/fs"
	"log"
	"math"
	"net/http"
	"os"

	"wynnbreeder/database"
	"wynnbreeder/handlers"
	"wynnbreeder/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/bcrypt"
)

//go:embed web
var webFS embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <serve|seed-admin|promote> [flags]\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		runServe(os.Args[2:])
	case "seed-admin":
		runSeedAdmin(os.Args[2:])
	case "promote":
		runPromote(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.String("port", "", "Port to listen on (overrides WYNNBREEDER_PORT)")
	dbPath := fs.String("db", "", "SQLite DB path (overrides WYNNBREEDER_DB)")
	_ = fs.Parse(args)

	cfg := loadConfig()
	if *port != "" {
		cfg.Port = *port
	}
	if *dbPath != "" {
		cfg.DBPath = *dbPath
	}

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	tmpl, err := loadTemplates()
	if err != nil {
		log.Fatalf("load templates: %v", err)
	}

	h := &handlers.Handler{
		DB:         db,
		Templates:  tmpl,
		SessionTTL: cfg.SessionTTL,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files — strip "web/" prefix from embedded FS
	staticFS, err := iofs.Sub(webFS, "web")
	if err != nil {
		log.Fatalf("static fs: %v", err)
	}
	r.Handle("/static/*", http.FileServer(http.FS(staticFS)))

	// Public routes
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.Login)
	r.Post("/api/login", h.APILogin)

	// Authenticated page routes
	r.Group(func(r chi.Router) {
		r.Use(h.RequireAuth)
		r.Get("/", h.DashboardPage)
		r.Get("/submit", h.SubmitPage)
		r.Get("/submissions", h.SubmissionsPage)
		r.Get("/submissions/{id}", h.SubmissionDetailPage)
		r.Get("/analytics", h.AnalyticsPage)
		r.Get("/account", h.AccountPage)
		r.Get("/logout", h.Logout)
		r.Post("/logout", h.Logout)
	})

	// Admin page routes
	r.Group(func(r chi.Router) {
		r.Use(h.RequireAdmin)
		r.Get("/admin", h.AdminPage)
	})

	// API routes (JSON)
	r.Group(func(r chi.Router) {
		r.Use(h.RequireAuthAPI)
		r.Post("/api/logout", h.APILogout)
		r.Post("/api/account/password", h.APIChangePassword)
		r.Get("/api/submissions", h.APIListSubmissions)
		r.Post("/api/submissions", h.APICreateSubmission)
		r.Get("/api/submissions/{id}", h.APIGetSubmission)
		r.Delete("/api/submissions/{id}", h.APIDeleteSubmission)
		r.Patch("/api/submissions/{id}/offspring", h.APIAddOffspring)
		r.Get("/api/analytics/stats", h.APIAnalyticsStats)
	})

	// Admin API routes
	r.Group(func(r chi.Router) {
		r.Use(h.RequireAdminAPI)
		r.Get("/api/admin/users", h.APIListUsers)
		r.Post("/api/admin/users", h.APICreateUser)
		r.Delete("/api/admin/users/{id}", h.APIDeleteUser)
		r.Patch("/api/admin/users/{id}/role", h.APISetUserRole)
		r.Get("/api/admin/export", h.APIExportSubmissions)
	})

	log.Printf("WynnBreeder running on http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func runSeedAdmin(args []string) {
	fs := flag.NewFlagSet("seed-admin", flag.ExitOnError)
	username := fs.String("username", "", "Admin username (required)")
	password := fs.String("password", "", "Admin password (required)")
	dbPath := fs.String("db", "", "SQLite DB path (overrides WYNNBREEDER_DB)")
	_ = fs.Parse(args)

	if *username == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "Usage: seed-admin --username <name> --password <pass>")
		os.Exit(1)
	}
	if len(*password) < 8 {
		fmt.Fprintln(os.Stderr, "Password must be at least 8 characters")
		os.Exit(1)
	}

	cfg := loadConfig()
	if *dbPath != "" {
		cfg.DBPath = *dbPath
	}

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	user, err := db.CreateSuperAdmin(*username, string(hash))
	if err != nil {
		log.Fatalf("create admin: %v", err)
	}
	fmt.Printf("Admin created: id=%d username=%s\n", user.ID, user.Username)
}

func runPromote(args []string) {
	fs := flag.NewFlagSet("promote", flag.ExitOnError)
	username := fs.String("username", "", "Username to promote to superadmin (required)")
	dbPath := fs.String("db", "", "SQLite DB path (overrides WYNNBREEDER_DB)")
	_ = fs.Parse(args)

	if *username == "" {
		fmt.Fprintln(os.Stderr, "Usage: promote --username <name>")
		os.Exit(1)
	}

	cfg := loadConfig()
	if *dbPath != "" {
		cfg.DBPath = *dbPath
	}

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	user, err := db.PromoteToSuperAdmin(*username)
	if err != nil {
		log.Fatalf("promote: %v", err)
	}
	fmt.Printf("Promoted: id=%d username=%s is_superadmin=true\n", user.ID, user.Username)
}

// pages lists every page template file that needs a layout+content pair.
var pages = []string{
	"login.html",
	"dashboard.html",
	"submit.html",
	"submissions.html",
	"submission_detail.html",
	"analytics.html",
	"admin.html",
	"account.html",
}

func loadTemplates() (map[string]*template.Template, error) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"pct": func(val, max int) int {
			if max == 0 {
				return 0
			}
			return int(float64(val) / float64(max) * 100)
		},
		"delta": func(off, pa, pb int) float64 {
			avg := float64(pa+pb) / 2.0
			return float64(off) - avg
		},
		"fmtDelta": func(d float64) string {
			if d > 0 {
				return fmt.Sprintf("+%.1f", d)
			}
			return fmt.Sprintf("%.1f", d)
		},
		"roundf": func(f float64) string {
			return fmt.Sprintf("%.2f", f)
		},
		"absDelta": func(d float64) float64 {
			return math.Abs(d)
		},
		"roleName": func(role models.MountRole) string {
			switch role {
			case models.RoleParentA:
				return "Parent A"
			case models.RoleParentB:
				return "Parent B"
			case models.RoleOffspring:
				return "Offspring"
			}
			return string(role)
		},
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i + 1
			}
			return s
		},
		"args": func(vals ...any) []any { return vals },
	}

	// Parse each page together with the layout so their "content" defines
	// don't collide across pages.
	tmpls := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		t, err := template.New("").Funcs(funcMap).ParseFS(
			webFS,
			"web/templates/layout.html",
			"web/templates/"+page,
		)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", page, err)
		}
		tmpls[page] = t
	}
	return tmpls, nil
}
