package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Rahmatdev030605/go-fiber-postgres/model"
	"github.com/Rahmatdev030605/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Book adalah struktur data untuk representasi buku.
type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

// Repository menyimpan koneksi ke database.
type Repository struct {
	DB *gorm.DB
}

// CreateBook membuat buku baru.
func (r *Repository) CreateBook(ctx *fiber.Ctx) error {
	book := Book{}

	// Parse body request ke variabel 'book'.
	if err := ctx.BodyParser(&book); err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}

	// Simpan buku ke database.
	if err := r.DB.Create(&book).Error; err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not create book"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been added", "data": book})
	return nil
}

// DeleteBook menghapus buku berdasarkan ID.
func (r *Repository) DeleteBook(ctx *fiber.Ctx) error {
	bookModel := model.Books{}
	id := ctx.Params("id")

	// Periksa jika ID kosong.
	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	// Hapus buku dari database.
	if err := r.DB.Delete(&bookModel, id).Error; err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete book"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "book deleted successfully"})
	return nil
}

// GetBooks mendapatkan semua buku.
func (r *Repository) GetBooks(ctx *fiber.Ctx) error {
	bookModels := &[]model.Books{}

	// Ambil semua buku dari database.
	if err := r.DB.Find(&bookModels).Error; err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get books"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "books fetched successfully", "data": bookModels})
	return nil
}

// GetBookByID mendapatkan buku berdasarkan ID.
func (r *Repository) GetBookByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	bookModel := &model.Books{}

	// Periksa jika ID kosong.
	if id == "" {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	// Ambil buku dari database berdasarkan ID.
	if err := r.DB.First(&bookModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "book not found"})
			return nil
		}

		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "could not get the book"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "book fetched successfully", "data": bookModel})
	return nil
}


// SetupRoutes mengatur rute API.
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books", r.GetBooks)
	api.Get("/get_book/:id", r.GetBookByID)
}

func main() {
	// Load konfigurasi dari file .env.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Konfigurasi database.
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	// Hubungkan ke database.
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Error connecting to the database")
	}

	// Migrate tabel buku.
	if err := model.MigrateBooks(db); err != nil {
		log.Fatal("Could not migrate database")
	}

	// Inisialisasi repository.
	r := Repository{
		DB: db,
	}

	// Inisialisasi aplikasi Fiber.
	app := fiber.New()

	// Setup rute API.
	r.SetupRoutes(app)

	// Jalankan aplikasi pada port 8080.
	app.Listen(":8080")
}
