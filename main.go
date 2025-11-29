package main

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

func main() {
	// 1. LOAD .env
	if err := godotenv.Load(); err != nil {
		log.Println("Info: Tidak menemukan file .env, sistem menggunakan default setting.")
	}

	// 2. KONEKSI DB
	db := config.ConnectDB()

	app := fiber.New()

	// 3. MIDDLEWARE
	app.Use(cors.New())

	// --- API CRUD ---

	// 1. READ 
	app.Get("/api/locations", func(c *fiber.Ctx) error {
		var locations []model.Location
		
		cursor, err := db.Find(c.Context(), bson.M{})
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		if err = cursor.All(c.Context(), &locations); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		return c.JSON(locations)
	})

	// 2. CREATE 
	app.Post("/api/locations", func(c *fiber.Ctx) error {
		var loc model.Location
		if err := c.BodyParser(&loc); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		res, err := db.InsertOne(c.Context(), loc)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// Ambil ID yang baru dibuat untuk dikembalikan ke frontend
		loc.ID = res.InsertedID.(primitive.ObjectID).Hex()
		return c.JSON(loc)
	})

	// 3. UPDATE 
	app.Put("/api/locations/:id", func(c *fiber.Ctx) error {
		idHex := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			return c.Status(400).SendString("ID tidak valid")
		}

		var updateData model.Location
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		// PENTING: Update field 'category' juga di sini
		update := bson.M{
			"$set": bson.M{
				"name":     updateData.Name,
				"desc":     updateData.Desc,
				"category": updateData.Category, 
			},
		}

		_, err = db.UpdateOne(c.Context(), bson.M{"_id": objID}, update)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		// Kembalikan data yang sudah diupdate ke frontend
		updateData.ID = idHex
		return c.JSON(updateData)
	})

	// 4. DELETE 
	app.Delete("/api/locations/:id", func(c *fiber.Ctx) error {
		idHex := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(idHex)
		
		_, err := db.DeleteOne(c.Context(), bson.M{"_id": objID})
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		return c.SendString("Berhasil dihapus")
	})

	log.Println("Backend berjalan di port 3000")
	app.Listen(":3000")
}