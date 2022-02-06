package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

var dataRoot = "data/"

func (p *Page) save() error {
	filename := dataRoot + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

// this should match validPath parsing
var linkFormat = regexp.MustCompile(`\[[a-zA-Z0-9]+\]`)

func (p Page) linkPages() *Page {
	formatted := linkFormat.ReplaceAllFunc(p.Body, func(link []byte) []byte {
		title := string(link[1 : len(link)-1])
		newLink := `<a href="/view/` + title + `">` + title + "</a>"
		return []byte(newLink)
	})

	return &Page{Title: p.Title, Body: formatted}
}

func loadPage(title string) (*Page, error) {
	filename := dataRoot + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func frontPageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p.linkPages())
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// https://stackoverflow.com/questions/18175630/go-template-executetemplate-include-html
//
// the `safeHTML` template function is provided so that HTML tags can be rendered in the view
var templates = template.Must(template.New("main").Funcs(template.FuncMap{
	"safeHTML": func(b []byte) template.HTML {
		return template.HTML(b)
	}}).ParseGlob("tmpl/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("static"))
	// the favicon handler is added to avoid falling back to the front page on every page load
	http.Handle("/favicon.ico", http.StripPrefix("/", fileServer))
	http.Handle("/styles.css", http.StripPrefix("/", fileServer))

	http.HandleFunc("/", frontPageHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
