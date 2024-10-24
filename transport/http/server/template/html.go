package template

import "html/template"

type Option func(*HTMLTemplate)

// Delims with html template delims config.
func Delims(left, right string) Option {
	return func(s *HTMLTemplate) {
		s.Delims = []string{left, right}
	}
}

// Delims with html template delims config.
func FuncMap(fm template.FuncMap) Option {
	return func(s *HTMLTemplate) {
		s.FuncMap = fm
	}
}

type HTMLTemplate struct {
	Templ   *template.Template
	Delims  []string
	FuncMap template.FuncMap
}

func New(opts ...Option) *HTMLTemplate {
	t := &HTMLTemplate{
		Delims:  []string{"{{", "}}"},
		FuncMap: template.FuncMap{},
	}
	return t
}

// LoadHTMLGlob loads HTML files identified by glob pattern
func (templ *HTMLTemplate) LoadHTMLGlob(pattern string) {
	templ.Templ = template.Must(
		template.New("").Delims(templ.Delims[0], templ.Delims[1]).Funcs(templ.FuncMap).ParseGlob(pattern),
	)
}

// LoadHTMLFiles loads a slice of HTML files
func (templ *HTMLTemplate) LoadHTMLFiles(files ...string) {
	templ.Templ = template.Must(
		template.New("").Delims(templ.Delims[0], templ.Delims[1]).Funcs(templ.FuncMap).ParseFiles(files...),
	)
}
