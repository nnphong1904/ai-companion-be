package router

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ai-companion-be/internal/config"
	"ai-companion-be/internal/handler"
	"ai-companion-be/internal/middleware"
	"ai-companion-be/internal/response"
)

// New creates a fully configured chi router with all routes.
func New(
	cfg *config.Config,
	authH *handler.AuthHandler,
	companionH *handler.CompanionHandler,
	storyH *handler.StoryHandler,
	messageH *handler.MessageHandler,
	relationshipH *handler.RelationshipHandler,
	memoryH *handler.MemoryHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware.
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   getAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		// Health check.
		r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		// Public auth routes.
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", authH.Signup)
			r.Post("/login", authH.Login)
		})

		// Protected routes.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWT))

			r.Get("/auth/me", authH.Me)

			// Companions.
			r.Get("/companions", companionH.GetAll)
			r.Get("/companions/{id}", companionH.GetByID)

			// Stories.
			r.Get("/stories", storyH.GetActiveStories)
			r.Get("/companions/{id}/stories", storyH.GetByCompanion)
			r.Post("/stories/{id}/react", storyH.React)

			// Messages (chat).
			r.Get("/companions/{id}/messages", messageH.GetHistory)
			r.Post("/companions/{id}/messages", messageH.Send)

			// Relationships.
			r.Get("/relationships", relationshipH.GetAllRelationships)
			r.Get("/companions/{id}/relationship", relationshipH.GetRelationship)

			// Onboarding.
			r.Post("/onboarding/select-companion", relationshipH.SelectCompanion)

			// Memories.
			r.Get("/companions/{id}/memories", memoryH.GetByCompanion)
			r.Post("/companions/{id}/memories", memoryH.Create)
			r.Delete("/memories/{id}", memoryH.Delete)
			r.Patch("/memories/{id}/pin", memoryH.TogglePin)
		})
	})

	return r
}

func getAllowedOrigins() []string {
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		return []string{origins}
	}
	return []string{"http://localhost:3000"}
}
