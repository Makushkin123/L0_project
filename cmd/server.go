package main

import (
	"html/template"
	"net/http"
)

type Application struct {
	cache *Cache
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	//fmt.Println("hello")
	id := r.URL.Query().Get("name")
	ts, err := template.ParseFiles("main.html")

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if id != "" {
		value, isin := app.cache.Get(id)
		if isin {
			err = ts.Execute(w, string(value))
		} else {
			err = ts.Execute(w, "нет данных")
		}
	} else {
		err = ts.Execute(w, nil)
	}

}

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	return mux
}
