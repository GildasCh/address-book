package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GildasCh/fermentation-notebook/model"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

var db *bolt.DB

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "path/to/addresses.db")
		return
	}

	db, err := bolt.Open(os.Args[1], 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	AddFromDB()

	if len(os.Args) >= 2 {
		AddFromYaml(os.Args[2])
	}

	err = serve()
	fmt.Println(err)
}

func serve() error {
	batchesHandler := func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("address-book.html").Funcs(template.FuncMap{
			"nl2br": func(s string) template.HTML {
				return template.HTML(strings.Replace(s, "\n", "<br />\n", -1))
			},
			"date": func(t time.Time) string {
				return t.Format("2006-01-02 15:04")
			},
			"until": func(t time.Time) string {
				return model.DurationToString(time.Until(t))
			}}).ParseFiles("tmpl/address-book.html")
		if err != nil {
			fmt.Println(err)
		}
		err = t.Execute(w, struct {
			AddrByTown
		}{ByTown()})
		if err != nil {
			fmt.Println(err)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/", batchesHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)
	port := "8080"
	fmt.Printf("Listening on %s...\n", port)
	return http.ListenAndServe(":"+port, nil)
}
