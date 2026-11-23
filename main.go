package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Product struct {
	ID, Name, Category, Description, Color string
	Price                                  float64
}

func (p Product) PriceString() string { return fmt.Sprintf("£%.2f", p.Price) }

var products = []Product{
	{"ember-mug", "Ember Mug", "Mugs", "Hand-thrown stoneware mug with a warm, speckled glaze. Holds 350ml.", "#c8643b", 24},
	{"dune-bowl", "Dune Bowl", "Bowls", "A wide serving bowl with a soft sand finish. Dishwasher safe.", "#d9b982", 38},
	{"moss-vase", "Moss Vase", "Vases", "Tall, narrow vase in deep green. Watertight and glazed inside.", "#5f7a4f", 46},
	{"tide-plate", "Tide Plate", "Plates", "Dinner plate with a gentle blue rim, made to be stacked.", "#4a7c8c", 22},
	{"clay-jug", "Clay Jug", "Jugs", "Generous 1L jug for water, wine or wildflowers.", "#a8553f", 34},
	{"ash-planter", "Ash Planter", "Planters", "Matte grey planter with a drainage hole and saucer.", "#7b7b7b", 29},
}

var pages = map[string]*template.Template{}

func init() {
	for _, n := range []string{"landing", "catalog", "product", "story", "enquiries"} {
		pages[n] = template.Must(template.ParseFiles("templates/base.html", "templates/"+n+".html"))
	}
}

func render(w http.ResponseWriter, name, title string, data any) {
	err := pages[name].ExecuteTemplate(w, "base", map[string]any{"Title": title, "Data": data})
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		log.Println(err)
	}
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	render(w, "landing", "Handmade ceramics", products[:3])
}

func catalogPage(w http.ResponseWriter, r *http.Request) {
	render(w, "catalog", "Shop", products)
}

func productPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for _, p := range products {
		if p.ID == id {
			render(w, "product", p.Name, p)
			return
		}
	}
	http.NotFound(w, r)
}

func storyPage(w http.ResponseWriter, r *http.Request) {
	render(w, "story", "Our story", nil)
}


func enquiriesPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to process enquiry form", http.StatusBadRequest)
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))

		log.Printf("enquiry from %q <%s>: %s",
			name,
			r.FormValue("email"),
			r.FormValue("message"),
		)

		render(w, "enquiries", "Enquiries", map[string]any{
			"Sent": true,
			"Name": name,
		})
		return
	}

	render(w, "enquiries", "Enquiries", map[string]any{
		"Sent": false,
	})
}


func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /{$}", landingPage)
	mux.HandleFunc("GET /shop", catalogPage)
	mux.HandleFunc("GET /shop/{id}", productPage)
	mux.HandleFunc("GET /story", storyPage)
	mux.HandleFunc("/enquiries", enquiriesPage)

	log.Println("Kiln & Co. running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
