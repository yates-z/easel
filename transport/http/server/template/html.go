package template

import "html/template"

type Option func(*HTMLTemplate)

// Delims with html template delims config.
func Delims(left, right string) Option {
	return func(s *HTMLTemplate) {
		s.Delims = []string{left, right}
	}
}

// FuncMap with html template delims config.
func FuncMap(fm template.FuncMap) Option {
	return func(s *HTMLTemplate) {
		s.FuncMap = fm
	}
}

type HTMLTemplate struct {
	Tpl     *template.Template
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
func (tpl *HTMLTemplate) LoadHTMLGlob(pattern string) {
	tpl.Tpl = template.Must(
		template.New("").Delims(tpl.Delims[0], tpl.Delims[1]).Funcs(tpl.FuncMap).ParseGlob(pattern),
	)
}

// LoadHTMLFiles loads a slice of HTML files
func (tpl *HTMLTemplate) LoadHTMLFiles(files ...string) {
	tpl.Tpl = template.Must(
		template.New("").Delims(tpl.Delims[0], tpl.Delims[1]).Funcs(tpl.FuncMap).ParseFiles(files...),
	)
}
