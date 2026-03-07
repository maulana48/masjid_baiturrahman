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
	http.HandleFunc("/prayers", handlePrayers)

	log.Println("Server starting on :8080...")
	log.Println("Visit http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Only handle root path for index, otherwise 404 for unknown routes
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
		Title   string
		Prayers []models.Prayer
	}{
		Title:   "Masjid Baiturrahman",
		Prayers: prayers,
	}

	// Parse all templates required for the layout
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/index.html",
		"templates/partials/navbar.html",
		"templates/partials/footer.html",
		"templates/partials/prayer_times.html",
	)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the "layout.html" template (which is the base)
	// Note: layout.html doesn't have a define "base", it's just the file content.
	// But since it includes {{ template "content" . }}, and index.html defines "content",
	// we should execute the layout file.
	// However, template.ParseFiles parses the first file as the name of the template set?
	// Actually, ParseFiles returns a *Template. Execute runs the first file provided?
	// No, Execute applies the template associated with t. 
	// If I use ParseFiles("layout.html", ...), the returned template name is "layout.html".
	
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func handlePrayers(w http.ResponseWriter, r *http.Request) {
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

	// Parse only the partial
	tmpl, err := template.ParseFiles("templates/partials/prayer_times.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}
