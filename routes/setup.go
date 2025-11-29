package routes

import (
	"log"

	"backend-peta/config"
	"backend-peta/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetupApp menyatukan database dan routing
func SetupApp() *fiber.App {
	// Load .env (hanya efek di local, di Vercel ini akan diabaikan/error tp aman)
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env tidak ditemukan (aman jika di Vercel)")
	}

	// Koneksi DB
	db := config.ConnectDB()

	app := fiber.New()
	
	// CORS sangat penting agar frontend GitHub Pages bisa akses
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// --- API CRUD ---
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend Peta Berjalan di Vercel!")
	})

	app.Get("/api/locations", func(c *fiber.Ctx) error {
		var locations []model.Location
		cursor, err := db.Find(c.Context(), bson.M{})
		if err != nil { return c.Status(500).SendString(err.Error()) }
		if err = cursor.All(c.Context(), &locations); err != nil { return c.Status(500).SendString(err.Error()) }
		return c.JSON(locations)
	})

	app.Post("/api/locations", func(c *fiber.Ctx) error {
		var loc model.Location
		if err := c.BodyParser(&loc); err != nil { return c.Status(400).SendString(err.Error()) }
		res, err := db.InsertOne(c.Context(), loc)
		if err != nil { return c.Status(500).SendString(err.Error()) }
		loc.ID = res.InsertedID.(primitive.ObjectID).Hex()
		return c.JSON(loc)
	})

	app.Put("/api/locations/:id", func(c *fiber.Ctx) error {
		idHex := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(idHex)
		if err != nil { return c.Status(400).SendString("ID tidak valid") }
		var updateData model.Location
		if err := c.BodyParser(&updateData); err != nil { return c.Status(400).SendString(err.Error()) }
		update := bson.M{"$set": bson.M{"name": updateData.Name, "desc": updateData.Desc, "category": updateData.Category}}
		_, err = db.UpdateOne(c.Context(), bson.M{"_id": objID}, update)
		if err != nil { return c.Status(500).SendString(err.Error()) }
		updateData.ID = idHex
		return c.JSON(updateData)
	})

	app.Delete("/api/locations/:id", func(c *fiber.Ctx) error {
		idHex := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(idHex)
		_, err := db.DeleteOne(c.Context(), bson.M{"_id": objID})
		if err != nil { return c.Status(500).SendString(err.Error()) }
		return c.SendString("Berhasil dihapus")
	})

	return app
}