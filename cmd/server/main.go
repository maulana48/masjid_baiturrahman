package main

import (
	"html/template"
	"log"
	"masjid_baiturrahman/internal/db"
	"masjid_baiturrahman/internal/models"
	"net/http"
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

	tmpl, err := template.ParseFiles(files...)
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

	data := struct {
		PageData
		Prayers []models.Prayer
	}{
		PageData: PageData{
			Title: "Masjid Baiturrahman - Home",
			Path:  "/",
		},
		Prayers: prayers,
	}

	renderPage(w, "index.html", data)
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	data := struct {
		PageData
	}{
		PageData: PageData{
			Title: "About - Masjid Baiturrahman",
			Path:  "/about",
		},
	}
	renderPage(w, "about.html", data)
}

func handleProgram(w http.ResponseWriter, r *http.Request) {
	data := struct {
		PageData
	}{
		PageData: PageData{
			Title: "Program - Masjid Baiturrahman",
			Path:  "/program",
		},
	}
	renderPage(w, "program.html", data)
}

func handleGallery(w http.ResponseWriter, r *http.Request) {
	data := struct {
		PageData
	}{
		PageData: PageData{
			Title: "Gallery - Masjid Baiturrahman",
			Path:  "/gallery",
		},
	}
	renderPage(w, "gallery.html", data)
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		contact := models.ContactMessage{
			Name:    r.FormValue("name"),
			Phone:   r.FormValue("phone"),
			Email:   r.FormValue("email"),
			Topic:   r.FormValue("topic"),
			Message: r.FormValue("message"),
		}

		if err := db.SaveContact(contact); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return success message partial
		tmpl := `
		<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative" role="alert">
			<strong class="font-bold">Terima kasih!</strong>
			<span class="block sm:inline">Pesan Anda telah kami terima. Kami akan segera menghubungi Anda.</span>
		</div>
		`
		t, _ := template.New("success").Parse(tmpl)
		t.Execute(w, nil)
		return
	}

	// GET request
	topic := r.URL.Query().Get("topic")
	data := struct {
		PageData
		Topic string
	}{
		PageData: PageData{
			Title: "Contact - Masjid Baiturrahman",
			Path:  "/contact",
		},
		Topic: topic,
	}
	renderPage(w, "contact.html", data)
}

func handlePrayers(w http.ResponseWriter, r *http.Request) {
	// This endpoint returns ONLY the partial for HTMX
	prayers, err := db.GetPrayers()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Prayers []models.Prayer
	}{
		Prayers: prayers,
	}

	tmpl, err := template.ParseFiles("templates/partials/prayer_times.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}
