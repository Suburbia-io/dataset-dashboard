package views

import (
	"html/template"
	"io"
	"path/filepath"
	"sync"
)

type Config struct {
	ViewCache bool
	TmplDir   string
}

type Views struct {
	cache   bool
	TmplDir string

	funcMap     template.FuncMap
	funcMapLock *sync.Mutex

	templates     map[string]*template.Template
	templatesLock *sync.Mutex
}

func NewViews(config Config) (v *Views) {
	return &Views{
		cache:   config.ViewCache,
		TmplDir: config.TmplDir,

		funcMap:     template.FuncMap{},
		funcMapLock: &sync.Mutex{},

		templates:     map[string]*template.Template{},
		templatesLock: &sync.Mutex{},
	}
}

func (v Views) RegisterFunc(key string, fn interface{}) {
	v.funcMapLock.Lock()
	defer v.funcMapLock.Unlock()

	v.funcMap[key] = fn
}

func (v Views) ExecuteSST(w io.Writer, path string, data interface{}) (err error) {
	tmpl, err := v.getTemplate(path, "related")

	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

// relatedDir contains related templates such as the base template and any included templates snippets.
func (v *Views) getTemplate(path string, relatedDir string) (tmpl *template.Template, err error) {
	v.templatesLock.Lock()
	defer v.templatesLock.Unlock()

	name := filepath.Base(path)

	if v.cache {
		if tmpl, ok := v.templates[name]; ok {
			return tmpl, nil
		}
	}

	templateFiles := []string{filepath.Join(v.TmplDir, path)}
	if relatedDir != "" {
		filenames, err := filepath.Glob(filepath.Join(v.TmplDir, relatedDir, "*"))
		if err != nil {
			return nil, err
		}
		templateFiles = append(templateFiles, filenames...)
	}
	v.templates[name], err = template.New(name).Funcs(v.funcMap).ParseFiles(templateFiles...)

	return v.templates[name], err
}
