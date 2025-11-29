package main

import (
	"log"

	// IMPORT PACKAGE LOKAL KITA:
	"backend-peta/config" // Mengambil fungsi koneksi database
	"backend-peta/model"  // Mengambil struct Location

	// IMPORT LIBRARY PIHAK KETIGA:
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv" // Library untuk baca file .env
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// 1. SETUP ENVIRONMENT & DATABASE
	
	// Load file .env dulu agar Config bisa baca MONGO_URI
	if err := godotenv.Load(); err != nil {
		log.Println("Info: Tidak menemukan file .env, sistem menggunakan default setting.")
	}

	// Panggil fungsi ConnectDB dari folder config.
	// Fungsi ini sudah mengembalikan Collection 'locations' siap pakai.
	db := config.ConnectDB()

	// 2. SETUP FIBER
	app := fiber.New()

	// Enable CORS (Agar frontend beda port bisa akses backend)
	app.Use(cors.New())

	// --- API CRUD ---

	// 1. READ (Ambil Semua Data)
	app.Get("/api/locations", func(c *fiber.Ctx) error {
		var locations []model.Location // Menggunakan struct dari folder model
		
		cursor, err := db.Find(c.Context(), bson.M{})
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		if err = cursor.All(c.Context(), &locations); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		return c.JSON(locations)
	})

	// 2. CREATE (Tambah Data Baru)
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

	// 3. UPDATE (Edit Data)
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

		// Kita hanya update field yang dikirim (Name & Desc)
		update := bson.M{
			"$set": bson.M{
				"name": updateData.Name,
				"desc": updateData.Desc,
			},
		}

		_, err = db.UpdateOne(c.Context(), bson.M{"_id": objID}, update)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		updateData.ID = idHex
		return c.JSON(updateData)
	})

	// 4. DELETE (Hapus Data)
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