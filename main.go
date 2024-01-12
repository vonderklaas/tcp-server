package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Film struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("main()")

	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/add-film/", addFilmHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("healthCheckHandler()")
	fmt.Fprint(w, "OK")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rootHandler()")

	tmpl := template.Must(template.ParseFiles("index.html"))
	films := map[string][]Film{
		"Films": {
			{Title: "The Shawshank Redemption", Director: "Frank Darabont"},
			{Title: "The Godfather", Director: "Francis Ford Coppola"},
			{Title: "The Dark Knight", Director: "Christopher Nolan"},
			{Title: "Oppenheimer", Director: "Christopher Nolan"},
		},
	}
	tmpl.Execute(w, films)
}

func addFilmHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("addFilmHandler()")

	title := r.PostFormValue("title")
	director := r.PostFormValue("director")

	if title == "" || director == "" {
		http.Error(w, "Title or director cannot be empty", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})
}
