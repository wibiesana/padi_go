package routes

import (
	"fmt"
	"net/http"

	controllers "padi-template/app/Controllers"
	"github.com/wibiesana/padi_go_core/config"
	"github.com/wibiesana/padi_go_core/middleware"
	"github.com/wibiesana/padi_go_core/response"
	"github.com/wibiesana/padi_go_core/router"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(r *router.Router) {
	// Root / Health check
	r.Mux.Get("/", func(w http.ResponseWriter, req *http.Request) {
		appName := "Padi REST API Framework (Go)"
		if config.AppConfig != nil && config.AppConfig.AppName != "" {
			appName = config.AppConfig.AppName
		}

		response.Success(w, map[string]interface{}{
			"app_name":  appName,
			"version":   "1.0.0",
			"status":    "running",
		}, fmt.Sprintf("Welcome to %s 🌾", appName))
	})

	r.Mux.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		response.Success(w, map[string]string{
			"status": "healthy",
		}, "System is healthy")
	})

	// Static Storage File Server (Serves uploaded assets)
	fileServer := http.StripPrefix("/storage/", http.FileServer(http.Dir("storage")))
	r.Mux.Handle("/storage/*", fileServer)

	authCtrl := controllers.NewAuthController()

	// Public Auth Endpoints
	r.Mux.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", authCtrl.Register)
		auth.Post("/login", authCtrl.Login)
		auth.Post("/refresh", authCtrl.Refresh)
		auth.Post("/logout", authCtrl.Logout)

		pwdCtrl := controllers.NewPasswordResetController()
		// Support both formats for 100% interchangeability
		auth.Post("/forgot-password", pwdCtrl.ForgotPassword)
		auth.Post("/reset-password", pwdCtrl.ResetPassword)
		auth.Post("/password/forgot", pwdCtrl.ForgotPassword)
		auth.Post("/password/reset", pwdCtrl.ResetPassword)

		// Protected Auth Endpoint
		auth.Group(func(protected chi.Router) {
			protected.Use(middleware.AuthRequired)
			protected.Get("/me", authCtrl.Me)
		})
	})

	// Protected API Resources (Requires JWT Bearer Token)
	r.Mux.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthRequired)

		// Users CRUD Resource
		protected.Route("/users", func(r chi.Router) {
			userCtrl := controllers.NewUserController()
			r.Get("/", userCtrl.Index)
			r.Get("/all", userCtrl.All)
			r.Post("/", userCtrl.Store)
			r.Get("/{id}", userCtrl.Show)
			r.Put("/{id}", userCtrl.Update)
			r.Delete("/{id}", userCtrl.Destroy)
		})
	})

	// Realtime SSE Pub/Sub Endpoints
	r.Mux.Route("/realtime", func(rt chi.Router) {
		rtCtrl := controllers.NewExampleRealtimeController()
		rt.Get("/subscribe", rtCtrl.Subscribe)
		rt.Post("/chat", rtCtrl.Broadcast)
	})
}
