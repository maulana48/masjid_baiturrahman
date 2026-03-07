package main

import (
	"encoding/json"
	"html/template"
	"log"
	"masjid_baiturrahman/internal/db"
	"masjid_baiturrahman/internal/models"
	"net/http"
	"time"
)

func main() {
	// Initialize database
	db.InitDB()

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/about", handleAbout)
	http.HandleFunc("/program", handleProgram)
	http.HandleFunc("/gallery", handleGallery)
	http.HandleFunc("/contact", handleContact)
	http.HandleFunc("/dashboard", handleDashboard)

	// HTMX endpoints
	http.HandleFunc("/prayers", handlePrayers) // Returns partial

	log.Println("Server starting on :8080...")
	log.Println("Visit http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Common page data
type PageData struct {
	Title string
	Path  string
}

// Helper to render full pages
func renderPage(w http.ResponseWriter, page string, data interface{}) {
	files := []string{
		"templates/layout.html",
		"templates/" + page,
		"templates/partials/navbar.html",
		"templates/partials/footer.html",
		"templates/partials/prayer_times.html",
		"templates/partials/home_about.html",
		"templates/partials/home_programs.html",
		"templates/partials/home_gallery.html",
		"templates/partials/home_contact.html",
	}

	tmpl := template.New("layout.html").Funcs(template.FuncMap{
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
	})

	var err error
	tmpl, err = tmpl.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	prayers, err := db.GetPrayers()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	programs, err := db.GetHomePrograms()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		PageData
		Prayers  []models.Prayer
		Programs []models.Program
	}{
		PageData: PageData{
			Title: "Masjid Baiturrahman - Home",
			Path:  "/",
		},
		Prayers:  prayers,
		Programs: programs,
	}

	renderPage(w, "index.html", data)
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Masjid Baiturrahman - About",
		Path:  "/about",
	}
	renderPage(w, "about.html", data)
}

func handleProgram(w http.ResponseWriter, r *http.Request) {
	programs, err := db.GetAllPrograms()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		PageData
		Programs []models.Program
	}{
		PageData: PageData{
			Title: "Masjid Baiturrahman - Program",
			Path:  "/program",
		},
		Programs: programs,
	}
	renderPage(w, "program.html", data)
}

func handleGallery(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Masjid Baiturrahman - Gallery",
		Path:  "/gallery",
	}
	renderPage(w, "gallery.html", data)
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Parse error", http.StatusBadRequest)
			return
		}
		msg := models.ContactMessage{
			Name:      r.FormValue("name"),
			Phone:     r.FormValue("phone"),
			Email:     r.FormValue("email"),
			Topic:     r.FormValue("topic"),
			Message:   r.FormValue("message"),
			CreatedAt: time.Now(),
		}
		if err := db.SaveContact(msg); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Return success message for HTMX
		w.Write([]byte("<div class='bg-green-100 text-green-800 p-4 rounded border border-green-200'>Terima kasih! Pesan Anda telah terkirim.</div>"))
		return
	}

	data := PageData{
		Title: "Masjid Baiturrahman - Contact",
		Path:  "/contact",
	}
	renderPage(w, "contact.html", data)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Masjid Baiturrahman - Dashboard",
		Path:  "/dashboard",
	}
	renderPage(w, "dashboard.html", data)
}

func handlePrayers(w http.ResponseWriter, r *http.Request) {
	prayers, err := db.GetPrayers()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/partials/prayer_times.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Prayers []models.Prayer
	}{
		Prayers: prayers,
	}

	tmpl.Execute(w, data)
}
