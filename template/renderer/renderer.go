package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	templateCache = make(map[string]*template.Template)
	templates     *template.Template
)

// Caches templates recursively within a base templates path
func SetupTemplates(templatesPath string) {
	templateFiles := make(map[string][]string)
	filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			// Get path relative to templates path
			relPath := strings.Replace(path, templatesPath+"/", "", 1)
			// Get key for templates, relPath without base filename
			templatesKey := filepath.Dir(relPath)
			if _, hasKey := templateFiles[templatesKey]; !hasKey {
				templateFiles[templatesKey] = []string{}
			}
			templateFiles[templatesKey] = append(templateFiles[templatesKey], path)
		}

		return nil
	})

	// Store a template obj per template
	for key, files := range templateFiles {
		tmpl := template.New(key)
		templateCache[key] = template.Must(tmpl.ParseFiles(files...))
	}
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	templatesKey := filepath.Dir(name)

	tmpl, hasTmpl := templateCache[templatesKey]
	if !hasTmpl {
		errMessage := fmt.Sprintln("No cached template at key: ", templatesKey)
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}

	templateName := filepath.Base(name)
	templatePath := fmt.Sprintf("%s.html", templateName)
	if err := tmpl.ExecuteTemplate(w, templatePath, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
