package routes

import (
	"net/http"

	controllers "padi-template/app/Controllers"
	"github.com/wibiesana/padi-core/middleware"
	"github.com/wibiesana/padi-core/response"
	"github.com/wibiesana/padi-core/router"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(r *router.Router) {
	// Root / Health check
	r.Mux.Get("/", func(w http.ResponseWriter, req *http.Request) {
		response.Success(w, map[string]interface{}{
			"framework": "Padi REST API Framework (Go)",
			"version":   "1.0.0",
			"status":    "running",
		}, "Welcome to Padi REST API Go Framework 🌾")
	})

	r.Mux.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		response.Success(w, map[string]string{
			"status": "healthy",
		}, "System is healthy")
	})

	// Static Storage File Server (Serves uploaded assets)
	fileServer := http.StripPrefix("/storage/", http.FileServer(http.Dir("storage")))
	r.Mux.Handle("/storage/*", fileServer)

	// API v1 Routes
	r.Version("v1", func(v1 chi.Router) {
		authCtrl := controllers.NewAuthController()

		// Public Auth Endpoints
		v1.Route("/auth", func(auth chi.Router) {
			auth.Post("/register", authCtrl.Register)
			auth.Post("/login", authCtrl.Login)

			pwdCtrl := controllers.NewPasswordResetController()
			auth.Post("/password/forgot", pwdCtrl.ForgotPassword)
			auth.Post("/password/reset", pwdCtrl.ResetPassword)

			// Protected Auth Endpoint
			auth.Group(func(protected chi.Router) {
				protected.Use(middleware.AuthRequired)
				protected.Get("/me", authCtrl.Me)
			})
		})

		// Users CRUD Resource
		v1.Route("/users", func(r chi.Router) {
			userCtrl := controllers.NewUserController()
			r.Get("/", userCtrl.Index)
			r.Get("/all", userCtrl.All)
			r.Post("/", userCtrl.Store)
			r.Get("/{id}", userCtrl.Show)
			r.Put("/{id}", userCtrl.Update)
			r.Delete("/{id}", userCtrl.Destroy)
		})

		// Protected API Resources
		v1.Group(func(protected chi.Router) {
			protected.Use(middleware.AuthRequired)

			// Example: Protected resource groups or generated CRUDs
		})

		// Realtime SSE Pub/Sub Endpoints
		v1.Route("/realtime", func(rt chi.Router) {
			rtCtrl := controllers.NewExampleRealtimeController()
			rt.Get("/subscribe", rtCtrl.Subscribe)
			rt.Post("/chat", rtCtrl.Broadcast)
		})
	})
}
