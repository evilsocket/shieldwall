package api

import (
	"compress/flate"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/mailer"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"strings"
	"time"
)

type API struct {
	config      Config
	mail        EmailConfig
	sendmail    *mailer.Mailer
	router      *chi.Mux
	certManager *autocert.Manager
}

func Setup(config Config, email EmailConfig, sendmail *mailer.Mailer) *API {
	api := &API{
		config:   config,
		mail:     email,
		sendmail: sendmail,
		router:   chi.NewRouter(),
	}

	if config.SSL {
		log.Info("ssl enabled for %s (caching on %s)", strings.Join(config.Domains, ", "), config.CertsCache)
		api.certManager = &autocert.Manager{
			Cache:      autocert.DirCache(config.CertsCache),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.Domains...),
		}
	} else {
		log.Warning("ssl disabled")
	}

	// use response compression
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	api.router.Use(compressor.Handler)

	// set CORS rules
	api.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API routes
	api.router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/rules", api.GetRules)

			r.Route("/subnets", func(r chi.Router) {
				r.Use(httprate.LimitByIP(1, 1*time.Second))
				r.Get("/cloudflare", api.GetCloudflareSubnets)
			})

			r.Route("/user", func(r chi.Router) {
				r.Use(httprate.LimitByIP(2, 1*time.Second))

				r.Post("/register", api.UserRegister)
				r.Get("/verify/{verification:[A-Fa-f0-9]{64}}", api.UserVerify)
				r.Post("/login", api.UserLogin)
				r.Post("/2step", api.UserSecondStep)

				r.Post("/", api.UserUpdate)

				r.Route("/agents", func(r chi.Router) {
					r.Put("/new", api.UserCreateAgent)

					r.Get("/", api.UserGetAgents)
					r.Get("/{id:[0-9]+}", api.UserGetAgent)

					r.Put("/{id:[0-9]+}", api.UserUpdateAgent)
					r.Delete("/{id:[0-9]+}", api.UserDeleteAgent)
				})
			})
		})
	})

	// frontend
	api.router.Handle("/*", http.FileServer(MockFS()))

	return api
}

func (api *API) Run() {
	log.Info("api starting on %s", api.config.Address)

	go api.alertingLoop()

	if api.config.SSL {
		server := &http.Server{
			Addr:      api.config.Address,
			TLSConfig: api.certManager.TLSConfig(),
			Handler:   api.router,
		}
		log.Fatal("%v", server.ListenAndServeTLS("", ""))
	} else {
		log.Fatal("%v", http.ListenAndServe(api.config.Address, api.router))
	}
}
