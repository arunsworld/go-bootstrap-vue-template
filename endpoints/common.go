package endpoints

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

// Theme is either bootstrap of vuetify
var Theme = "bootstrap"

type htmlFromTemplate struct {
	t *template.Template
}

func newHTMLFromTemplate(name string, dir string) (htmlFromTemplate, error) {
	tmpl, err := template.New(name).Delims("[[", "]]").ParseGlob(path.Join(dir, "*"))
	if err != nil {
		return htmlFromTemplate{}, fmt.Errorf("problem parsing template: %v", err)
	}
	if tmpl.Tree == nil {
		return htmlFromTemplate{}, fmt.Errorf("template for %s is empty, cannot be used", name)
	}
	return htmlFromTemplate{
		t: tmpl,
	}, nil
}

func newHTMLFromTemplateFromMinfiedTemplates(templates []string, name string) htmlFromTemplate {
	tmpl := template.New(name).Delims("[[[", "]]]")
	for _, tmplData := range templates {
		tmpl.Parse(tmplData)
	}
	return htmlFromTemplate{
		t: tmpl,
	}
}

// MinifiedTemplates reads the given template dir and loads them into memory after minimfying them
func MinifiedTemplates(dir string) ([]string, error) {
	filenames, err := filepath.Glob(path.Join(dir, "*"))
	if err != nil {
		return nil, err
	}
	result := []string{}
	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
	})
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, err
		}
		result = append(result, string(mb))
		log.Println("read template with:", filename)
	}
	return result, nil
}

func (tmpl htmlFromTemplate) Execute(w http.ResponseWriter, params interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	err := tmpl.t.Execute(w, params)
	if err != nil {
		return fmt.Errorf("problem executing template with params: %v: %v", params, err)
	}

	return nil
}

// NavBar defines the parameterized requirements of the navbar
type NavBar struct {
	NavItems   []NavItem
	ActiveItem int
}

// NavItem are elements within the NavBar
type NavItem struct {
	Title    string
	Link     string
	Icon     string
	IsActive bool
}
