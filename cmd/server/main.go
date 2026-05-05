package main

import (
	"go_store_project/internal/handlers"
	"go_store_project/internal/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// Инициализация хранилища (подключается к PostgreSQL)
	storage.Init()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	setupRoutes(r)

	port := ":8082"
	log.Printf("🚀 Сервер запущен на http://localhost%s", port)
	log.Printf("📊 PostgreSQL подключена")

	err := r.Run(port)
	if err != nil {
		log.Printf("Порт 8082 занят, пробуем 8081...")
		r.Run(":8081")
	}
}

func setupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// Главное меню и авторизация
	r.GET("/main", handlers.MainMenuPage)
	r.GET("/login", handlers.LoginPage)
	r.POST("/login", handlers.LoginHandler)
	r.GET("/logout", handlers.LogoutHandler)

	// 👇 ДОБАВИЛИ (страница покупателя)
	r.GET("/shop", handlers.ShopPage)
	r.GET("/cart/add/:id", handlers.AddToCartHandler)
	r.GET("/cart/delete/:id", handlers.RemoveFromCart)
	r.GET("/cart/buy/:id", handlers.BuyFromCart)
	r.GET("/cart", handlers.CartPage)
	r.GET("/orders", handlers.OrdersPage)

	r.GET("/orders/delete/:id", handlers.DeleteOrder)
	r.GET("/orders/status/:id/:status", handlers.UpdateOrderStatusHandler)

	// Магазины
	r.GET("/stores", handlers.StoreListPage)
	r.GET("/stores/create", handlers.StoreCreatePage)
	r.POST("/stores/create", handlers.StoreCreateHandler)
	r.GET("/stores/:id/edit", handlers.StoreEditPage)
	r.POST("/stores/:id/edit", handlers.StoreUpdateHandler)
	r.GET("/stores/:id/delete", handlers.StoreDeleteHandler)

	// Категории
	r.GET("/categories", handlers.CategoriesListPage)
	r.GET("/categories/create", handlers.CategoryCreatePage)
	r.POST("/categories/create", handlers.CategoryCreateHandler)
	r.GET("/categories/:id/edit", handlers.CategoryEditPage)
	r.POST("/categories/:id/edit", handlers.CategoryUpdateHandler)
	r.GET("/categories/:id/delete", handlers.CategoryDeleteHandler)

	// Товары
	r.GET("/products", handlers.ProductsList)
	r.GET("/products/create", handlers.ProductCreatePage)
	r.POST("/products/create", handlers.ProductCreateHandler)
	r.GET("/products/:id/edit", handlers.ProductEditPage)
	r.POST("/products/:id/edit", handlers.ProductUpdateHandler)
	r.GET("/products/:id/delete", handlers.ProductDeleteHandler)

	// Поставщики
	r.GET("/suppliers", handlers.SuppliersPage)
	r.GET("/suppliers/create", handlers.SupplierCreatePage)
	r.POST("/suppliers/create", handlers.SupplierCreateHandler)
	r.GET("/suppliers/:id/edit", handlers.SupplierEditPage)
	r.POST("/suppliers/:id/edit", handlers.SupplierUpdateHandler)
	r.GET("/suppliers/:id/delete", handlers.SupplierDeleteHandler)

	// Поставки
	r.GET("/supplies", handlers.SuppliesPage)
	r.GET("/supplies/create", handlers.SupplyCreatePage)
	r.POST("/supplies/create", handlers.SupplyCreateHandler)
	r.GET("/supplies/:id/edit", handlers.SupplyEditPage)
	r.POST("/supplies/:id/edit", handlers.SupplyUpdateHandler)
	r.GET("/supplies/:id/delete", handlers.SupplyDeleteHandler)

	// Экспорт данных
	r.GET("/export/txt", handlers.ExportTXT)
	r.GET("/export/csv", handlers.ExportCSV)
	r.GET("/export/excel", handlers.ExportExcel)

	// Статические страницы
	r.GET("/about", handlers.AboutPage)
	r.GET("/contacts", handlers.ContactsPage)
}
