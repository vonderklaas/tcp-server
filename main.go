package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Film struct {
	Title string
	Director string
}

func main() {
	fmt.Println("Go!")

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		films := map[string][]Film{
			"Films": {
				{ Title: "The Shawshank Redemption", Director: "Frank Darabont" },
				{ Title: "The Godfather", Director: "Francis Ford Coppola" },
				{ Title: "The Dark Knight", Director: "Christopher Nolan" },
			},
		}
		tmpl.Execute(w, films)
	}

	addFilmHandler := func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")

		if title == "" || director == "" {
			http.Error(w, "Title or director cannot be empty", http.StatusBadRequest)
			return
		}

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})
	}

	http.HandleFunc("/", rootHandler);
	http.HandleFunc("/add-film/", addFilmHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}