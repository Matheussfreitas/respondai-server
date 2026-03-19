package handler

import (
	"html/template"
	"net/http"
)

type scalarReferencePageData struct {
	Title   string
	SpecURL string
}

var scalarReferencePageTemplate = template.Must(template.New("scalar-reference").Parse(`<!doctype html>
<html lang="pt-BR">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{.Title}}</title>
    <style>
      html, body { height: 100%; margin: 0; }
      body { font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif; }
      .fallback { padding: 16px; }
      code { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace; }
    </style>
  </head>
  <body>
    <div id="scalar-fallback" class="fallback">
      Carregando documentação interativa…
      <noscript>
        <div style="margin-top: 8px;">Habilite JavaScript para visualizar a documentação interativa.</div>
      </noscript>
    </div>

    <script id="api-reference" data-url="{{.SpecURL}}"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>

    <script>
        var fallback = document.getElementById('scalar-fallback')
        if (!fallback) return

        var attempts = 0
        var timer = window.setInterval(function () {
          attempts++

          var rendered = !!document.querySelector('scalar-api-reference')
          if (!rendered) return

          fallback.style.display = 'none'
          window.clearInterval(timer)
        }, 250)

        window.setTimeout(function () {
          if (fallback.style.display === 'none') return
          fallback.innerHTML = 'Se a UI não carregar, verifique se o endpoint do spec está acessível em <code>{{.SpecURL}}</code>.'
        }, 3000)
    </script>
  </body>
</html>`))

func ScalarReferenceHandler(defaultSpecURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reference" && r.URL.Path != "/reference/" {
			http.NotFound(w, r)
			return
		}

		specURL := r.URL.Query().Get("spec")
		if specURL == "" {
			specURL = defaultSpecURL
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = scalarReferencePageTemplate.Execute(w, scalarReferencePageData{
			Title:   "API Reference",
			SpecURL: specURL,
		})
	}
}
