package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GildasCh/fermentation-notebook/model"
	"github.com/gorilla/mux"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "path/to/addresses.yaml")
		return
	}

	as, err := readAddresses(os.Args[1])
	if err != nil {
		fmt.Println("Error reading addresses:", err)
		return
	}

	err = serve(as)
	fmt.Println(err)
}

func readAddresses(input string) (as Addresses, err error) {
	af, err := os.Open(input)
	if err != nil {
		return
	}
	ab, err := ioutil.ReadAll(af)
	if err != nil {
		return
	}

	as, err = ParseAddresses(ab)
	if err != nil {
		return
	}
	return
}

func serve(as Addresses) error {
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
			Addresses
		}{as})
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
