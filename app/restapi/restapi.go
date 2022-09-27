package restapi

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/prabudzak/article/service"
)

type route struct {
	method  string
	path    string
	handler httprouter.Handle
}

type middleware func(route route, fn httprouter.Handle) httprouter.Handle

// API represent REST API application
type API struct {
	articleService    service.ArticleService
	shortenUrlService service.ShortenUrlService
}

// New create a new instance of REST API application
func New(articleService service.ArticleService, shortenUrlService service.ShortenUrlService) *API {
	return &API{
		articleService:    articleService,
		shortenUrlService: shortenUrlService,
	}
}

// Router return registered REST API path
func (a *API) Router() http.Handler {
	router := httprouter.New()

	routes := []route{
		{method: http.MethodPost, path: "/articles", handler: a.createArticle},
		{method: http.MethodGet, path: "/articles", handler: a.listArticle},
		{method: http.MethodPost, path: "/short", handler: a.generateURL},
		{method: http.MethodGet, path: "/s/:shortenID", handler: a.redirect},
		{method: http.MethodGet, path: "/stats/:shortenID", handler: a.getShortUrlStats},

		{method: http.MethodGet, path: "/healthz", handler: a.healthz},
	}

	for _, route := range routes {
		router.Handle(route.method, route.path, a.log(route, route.handler))
	}

	return router
}

type wrapperResponseWriter struct {
	http.ResponseWriter

	status      int
	wroteHeader bool
}

func (w *wrapperResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}

	w.wroteHeader = true
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (a *API) log(route route, fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		start := time.Now()

		writer := &wrapperResponseWriter{ResponseWriter: w}
		fn(writer, r, param)

		log.Printf("%d %s %s in %dms\n", writer.status, route.method, route.path, time.Since(start).Milliseconds())
	}
}
