package main

import (
	"context"
	"log"
	"os"
    "net/http"

	"phone-book/handler"
	"phone-book/repository"
	"phone-book/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Подключение к БД через пул
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/phonebook?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Ошибка создания пула: ", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("База данных недоступна: ", err)
	}
	log.Println("База данных подключена")

	// 2. Сборка слоёв
	repo := repository.NewContactRepository(pool)
	svc := service.NewContactService(repo)
	h := handler.NewContactHandler(svc)

	// 3. Настройка Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"https://phone-book-fr.onrender.com"}, // В ПРОДЕ АДРЕС 
    AllowMethods: []string{
        http.MethodGet, http.MethodPost, http.MethodPut,
        http.MethodDelete, http.MethodOptions,
    },
    AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
}))

	// 4. Маршруты
	api := e.Group("/api")
	api.GET("/contacts", h.GetAll)
	api.POST("/contacts", h.Create)
	api.GET("/contacts/:id", h.GetByID)
	api.PUT("/contacts/:id", h.Update)
	api.DELETE("/contacts/:id", h.Delete)

	// 5. Запуск
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}