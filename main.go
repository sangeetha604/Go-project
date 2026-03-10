package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {

	// Serve static files (CSS, JS)
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/editor", editorHandler)

	fmt.Println("Server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "signup.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "dashboard.html", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create.html", nil)
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "join.html", nil)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "editor.html", nil)
}
