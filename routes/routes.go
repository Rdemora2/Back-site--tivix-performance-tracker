package routes

import (
	"github.com/gofiber/fiber/v2"
	"tivix-performance-tracker-backend/handlers"
)

func SetupRoutes(app *fiber.App) {
	// Grupo principal da API
	api := app.Group("/api/v1")

	// Rotas de times
	teams := api.Group("/teams")
	teams.Get("/", handlers.GetAllTeams)
	teams.Get("/:id", handlers.GetTeamByID)
	teams.Post("/", handlers.CreateTeam)
	teams.Put("/:id", handlers.UpdateTeam)
	teams.Delete("/:id", handlers.DeleteTeam)

	// Rotas de desenvolvedores
	developers := api.Group("/developers")
	developers.Get("/", handlers.GetAllDevelopers)
	developers.Get("/archived", handlers.GetArchivedDevelopers)
	developers.Get("/:id", handlers.GetDeveloperByID)
	developers.Post("/", handlers.CreateDeveloper)
	developers.Put("/:id", handlers.UpdateDeveloper)
	developers.Put("/:id/archive", handlers.ArchiveDeveloper)

	// Rotas de desenvolvedores por time
	teams.Get("/:teamId/developers", handlers.GetDevelopersByTeam)

	// Rotas de relatórios de performance
	reports := api.Group("/performance-reports")
	reports.Get("/", handlers.GetAllPerformanceReports)
	reports.Get("/months", handlers.GetAvailableMonths)
	reports.Get("/stats", handlers.GetPerformanceStats)
	reports.Get("/:id", handlers.GetPerformanceReportByID)
	reports.Post("/", handlers.CreatePerformanceReport)

	// Rotas de relatórios por desenvolvedor
	developers.Get("/:developerId/reports", handlers.GetPerformanceReportsByDeveloper)

	// Rotas de relatórios por mês
	reports.Get("/month/:month", handlers.GetPerformanceReportsByMonth)
}
