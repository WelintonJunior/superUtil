package sutests

import (
	"testing"

	sucontrollers "github.com/WelintonJunior/superUtil/controllers"
	surepository "github.com/WelintonJunior/superUtil/repository"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type customerEntity struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}

func TestStructureTest_UsesRealControllerRoutes(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in memory: %v", err)
	}

	if err := db.AutoMigrate(&customerEntity{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	repo := surepository.NewSuperUtilRepository[customerEntity](db)
	controller := sucontrollers.NewSuperUtilController[customerEntity](repo)

	app := fiber.New()
	app.Post("/customers", controller.Create())
	app.Get("/customers", controller.GetAll())
	app.Get("/customers/:id", controller.GetByID())
	app.Put("/customers/:id", controller.Update())
	app.Delete("/customers/:id", controller.DeleteByID())

	if err := StructureTest[customerEntity](app, "/customers"); err != nil {
		t.Fatalf("controller integration test failed: %v", err)
	}
}

func TestStructureTest_InvalidIDRouteReturnsBadRequest(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in memory: %v", err)
	}

	if err := db.AutoMigrate(&customerEntity{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	repo := surepository.NewSuperUtilRepository[customerEntity](db)
	controller := sucontrollers.NewSuperUtilController[customerEntity](repo)

	app := fiber.New()
	app.Get("/customers/:id", controller.GetByID())

	resp, err := app.Test(newRequest("GET", "/customers/invalid-id"))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}
