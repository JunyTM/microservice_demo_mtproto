package router

import (
	"ms_gmail/controller"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*", "https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		// MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Ping returns
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })

	userController := controller.NewUserController()
	r.Post("/user/resgister", userController.Register)
	r.Post("/user/login", userController.Login)
	r.Get("/user/profile", userController.GetUserProfile)
	r.Get("/user/gen-data", userController.GenerateUsers)
	r.Get("/test/load-gen-data", userController.LoadUsersGenerated)
	return r
}
