package main

import (
	"bytes" // ДОБАВЬТЕ ЭТОТ ИМПОРТ
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ====================================================================
// Структуры данных (в памяти)
// ====================================================================

type Store struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Logo      string    `json:"logo"`
	CreatedAt time.Time `json:"created_at"`
}

type Supplier struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	CategoryID  int       `json:"category_id"`
	SKU         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Supply struct {
	ID         int       `json:"id"`
	SupplierID int       `json:"supplier_id"`
	ProductID  int       `json:"product_id"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	Total      float64   `json:"total"`
	Date       time.Time `json:"date"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}

// ====================================================================
// Хранилище в памяти
// ====================================================================

var (
	stores     = make(map[int]Store)
	suppliers  = make(map[int]Supplier)
	products   = make(map[int]Product)
	categories = make(map[int]Category)
	supplies   = make(map[int]Supply)
	users      = make(map[int]User)
	storeID    = 0
	supplierID = 0
	productID  = 0
	categoryID = 0
	supplyID   = 0
	mu         sync.RWMutex
)

func init() {
	maxID := 0
	for id := range stores {
		if id > maxID {
			maxID = id
		}
	}
	storeID = maxID
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	setupRoutes(r)

	port := ":8080"
	log.Printf("🚀 Сервер запущен на http://localhost%s", port)

	err := r.Run(port)
	if err != nil {
		log.Printf("Порт 8080 занят, пробуем 8081...")
		r.Run(":8081")
	}
}

func setupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	r.GET("/main", mainMenuPage)

	// Авторизация
	r.GET("/login", loginPage)
	r.POST("/login", loginHandler)
	r.GET("/logout", logoutHandler)

	// Магазины
	r.GET("/stores", storeListPage)
	r.GET("/stores/create", storeCreatePage)
	r.POST("/stores/create", storeCreateHandler)
	r.GET("/stores/:id/edit", storeEditPage)
	r.POST("/stores/:id/edit", storeUpdateHandler)
	r.GET("/stores/:id/delete", storeDeleteHandler)

	// Категории
	r.GET("/categories", categoriesListPage)
	r.GET("/categories/create", categoryCreatePage)
	r.POST("/categories/create", categoryCreateHandler)
	r.GET("/categories/:id/edit", categoryEditPage)
	r.POST("/categories/:id/edit", categoryUpdateHandler)
	r.GET("/categories/:id/delete", categoryDeleteHandler)

	// Товары
	r.GET("/products", productsList)
	r.GET("/products/create", productCreatePage)
	r.POST("/products/create", productCreateHandler)
	r.GET("/products/:id/edit", productEditPage)
	r.POST("/products/:id/edit", productUpdateHandler)
	r.GET("/products/:id/delete", productDeleteHandler)

	// Поставщики
	r.GET("/suppliers", suppliersPage)
	r.GET("/suppliers/create", supplierCreatePage)
	r.POST("/suppliers/create", supplierCreateHandler)
	r.GET("/suppliers/:id/edit", supplierEditPage)
	r.POST("/suppliers/:id/edit", supplierUpdateHandler)
	r.GET("/suppliers/:id/delete", supplierDeleteHandler)

	// Поставки
	r.GET("/supplies", suppliesPage)
	r.GET("/supplies/create", supplyCreatePage)
	r.POST("/supplies/create", supplyCreateHandler)
	r.GET("/supplies/:id/edit", supplyEditPage)
	r.POST("/supplies/:id/edit", supplyUpdateHandler)
	r.GET("/supplies/:id/delete", supplyDeleteHandler)

	// Экспорт данных - НОВЫЕ ОБРАБОТЧИКИ
	r.GET("/export/txt", exportTXT)
	r.GET("/export/csv", exportCSV)
	r.GET("/export/excel", exportExcel)

	// Статические страницы
	r.GET("/about", aboutPage)
	r.GET("/contacts", contactsPage)
}

// ====================================================================
// ЭКСПОРТ ДАННЫХ - НОВЫЕ ФУНКЦИИ
// ====================================================================

// Экспорт в TXT
func exportTXT(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Получаем тип данных для экспорта
	dataType := c.Query("type")
	filename := "export.txt"

	var content strings.Builder
	content.WriteString("Экспорт данных системы\n")
	content.WriteString("========================\n")
	content.WriteString(fmt.Sprintf("Дата экспорта: %s\n\n", time.Now().Format("02.01.2006 15:04")))

	mu.RLock()
	defer mu.RUnlock()

	switch dataType {
	case "stores":
		filename = "stores.txt"
		content.WriteString("МАГАЗИНЫ\n")
		content.WriteString("========\n")
		for _, store := range stores {
			content.WriteString(fmt.Sprintf("ID: %d\n", store.ID))
			content.WriteString(fmt.Sprintf("Название: %s\n", store.Name))
			content.WriteString(fmt.Sprintf("Адрес: %s\n", store.Address))
			content.WriteString(fmt.Sprintf("Дата создания: %s\n", store.CreatedAt.Format("02.01.2006")))
			content.WriteString("---\n")
		}

	case "products":
		filename = "products.txt"
		content.WriteString("ТОВАРЫ\n")
		content.WriteString("======\n")
		for _, product := range products {
			categoryName := "Без категории"
			if cat, exists := categories[product.CategoryID]; exists {
				categoryName = cat.Name
			}

			content.WriteString(fmt.Sprintf("ID: %d\n", product.ID))
			content.WriteString(fmt.Sprintf("Название: %s\n", product.Name))
			content.WriteString(fmt.Sprintf("Артикул: %s\n", product.SKU))
			content.WriteString(fmt.Sprintf("Цена: %.2f руб.\n", product.Price))
			content.WriteString(fmt.Sprintf("Описание: %s\n", product.Description))
			content.WriteString(fmt.Sprintf("Категория: %s\n", categoryName))
			content.WriteString(fmt.Sprintf("Дата создания: %s\n", product.CreatedAt.Format("02.01.2006")))
			content.WriteString("---\n")
		}

	case "suppliers":
		filename = "suppliers.txt"
		content.WriteString("ПОСТАВЩИКИ\n")
		content.WriteString("==========\n")
		for _, supplier := range suppliers {
			content.WriteString(fmt.Sprintf("ID: %d\n", supplier.ID))
			content.WriteString(fmt.Sprintf("Название: %s\n", supplier.Name))
			content.WriteString(fmt.Sprintf("Телефон: %s\n", supplier.Phone))
			content.WriteString(fmt.Sprintf("Email: %s\n", supplier.Email))
			content.WriteString(fmt.Sprintf("Адрес: %s\n", supplier.Address))
			content.WriteString("---\n")
		}

	case "supplies":
		filename = "supplies.txt"
		content.WriteString("ПОСТАВКИ\n")
		content.WriteString("========\n")
		for _, supply := range supplies {
			supplierName := "Неизвестный"
			if sup, exists := suppliers[supply.SupplierID]; exists {
				supplierName = sup.Name
			}

			productName := "Неизвестный"
			if prod, exists := products[supply.ProductID]; exists {
				productName = prod.Name
			}

			content.WriteString(fmt.Sprintf("ID поставки: %d\n", supply.ID))
			content.WriteString(fmt.Sprintf("Поставщик: %s\n", supplierName))
			content.WriteString(fmt.Sprintf("Товар: %s\n", productName))
			content.WriteString(fmt.Sprintf("Количество: %d\n", supply.Quantity))
			content.WriteString(fmt.Sprintf("Цена за ед.: %.2f руб.\n", supply.Price))
			content.WriteString(fmt.Sprintf("Итого: %.2f руб.\n", supply.Total))
			content.WriteString(fmt.Sprintf("Дата: %s\n", supply.Date.Format("02.01.2006")))
			content.WriteString(fmt.Sprintf("Статус: %s\n", getStatusText(supply.Status)))
			content.WriteString(fmt.Sprintf("Примечания: %s\n", supply.Notes))
			content.WriteString("---\n")
		}

	case "categories":
		filename = "categories.txt"
		content.WriteString("КАТЕГОРИИ\n")
		content.WriteString("=========\n")
		for _, category := range categories {
			content.WriteString(fmt.Sprintf("ID: %d\n", category.ID))
			content.WriteString(fmt.Sprintf("Название: %s\n", category.Name))
			content.WriteString(fmt.Sprintf("Описание: %s\n", category.Description))
			content.WriteString("---\n")
		}

	default:
		// Экспорт всего
		filename = "full_export.txt"
		// Экспорт магазинов
		content.WriteString("МАГАЗИНЫ\n")
		content.WriteString("========\n")
		for _, store := range stores {
			content.WriteString(fmt.Sprintf("ID: %d | Название: %s | Адрес: %s\n",
				store.ID, store.Name, store.Address))
		}
		content.WriteString("\n")

		// Экспорт категорий
		content.WriteString("КАТЕГОРИИ\n")
		content.WriteString("=========\n")
		for _, category := range categories {
			content.WriteString(fmt.Sprintf("ID: %d | Название: %s | Описание: %s\n",
				category.ID, category.Name, category.Description))
		}
		content.WriteString("\n")

		// Экспорт товаров
		content.WriteString("ТОВАРЫ\n")
		content.WriteString("======\n")
		for _, product := range products {
			categoryName := "Без категории"
			if cat, exists := categories[product.CategoryID]; exists {
				categoryName = cat.Name
			}
			content.WriteString(fmt.Sprintf("ID: %d | Название: %s | Цена: %.2f | Артикул: %s | Категория: %s\n",
				product.ID, product.Name, product.Price, product.SKU, categoryName))
		}
		content.WriteString("\n")

		// Экспорт поставщиков
		content.WriteString("ПОСТАВЩИКИ\n")
		content.WriteString("==========\n")
		for _, supplier := range suppliers {
			content.WriteString(fmt.Sprintf("ID: %d | Название: %s | Телефон: %s | Email: %s\n",
				supplier.ID, supplier.Name, supplier.Phone, supplier.Email))
		}
		content.WriteString("\n")

		// Экспорт поставок
		content.WriteString("ПОСТАВКИ\n")
		content.WriteString("========\n")
		for _, supply := range supplies {
			supplierName := "Неизвестный"
			if sup, exists := suppliers[supply.SupplierID]; exists {
				supplierName = sup.Name
			}

			productName := "Неизвестный"
			if prod, exists := products[supply.ProductID]; exists {
				productName = prod.Name
			}

			content.WriteString(fmt.Sprintf("ID: %d | Поставщик: %s | Товар: %s | Количество: %d | Итого: %.2f | Статус: %s\n",
				supply.ID, supplierName, productName, supply.Quantity, supply.Total, getStatusText(supply.Status)))
		}
	}

	// Устанавливаем заголовки для скачивания файла
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, content.String())
}

// Экспорт в CSV (Excel) - ИСПРАВЛЕННАЯ ВЕРСИЯ
func exportCSV(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Получаем тип данных для экспорта
	dataType := c.Query("type")
	filename := "export.csv"

	mu.RLock()
	defer mu.RUnlock()

	// Определяем имя файла
	switch dataType {
	case "stores":
		filename = "stores.csv"
	case "products":
		filename = "products.csv"
	case "suppliers":
		filename = "suppliers.csv"
	case "supplies":
		filename = "supplies.csv"
	case "categories":
		filename = "categories.csv"
	default:
		filename = "full_export.csv"
	}

	// Создаем буфер для данных
	var content bytes.Buffer

	// Записываем BOM (Byte Order Mark) для UTF-8 - ВАЖНО для Excel!
	content.Write([]byte{0xEF, 0xBB, 0xBF})

	switch dataType {
	case "stores":
		// Заголовки для магазинов
		content.WriteString("ID;Название;Адрес;Дата создания\r\n")
		for _, store := range stores {
			content.WriteString(fmt.Sprintf("%d;%s;%s;%s\r\n",
				store.ID,
				escapeCSV(store.Name),
				escapeCSV(store.Address),
				store.CreatedAt.Format("02.01.2006"),
			))
		}

	case "products":
		// Заголовки для товаров
		content.WriteString("ID;Название;Артикул;Цена;Описание;Категория;Дата создания\r\n")
		for _, product := range products {
			categoryName := "Без категории"
			if cat, exists := categories[product.CategoryID]; exists {
				categoryName = cat.Name
			}

			content.WriteString(fmt.Sprintf("%d;%s;%s;%.2f;%s;%s;%s\r\n",
				product.ID,
				escapeCSV(product.Name),
				escapeCSV(product.SKU),
				product.Price,
				escapeCSV(product.Description),
				escapeCSV(categoryName),
				product.CreatedAt.Format("02.01.2006"),
			))
		}

	case "suppliers":
		// Заголовки для поставщиков
		content.WriteString("ID;Название;Телефон;Email;Адрес\r\n")
		for _, supplier := range suppliers {
			content.WriteString(fmt.Sprintf("%d;%s;%s;%s;%s\r\n",
				supplier.ID,
				escapeCSV(supplier.Name),
				escapeCSV(supplier.Phone),
				escapeCSV(supplier.Email),
				escapeCSV(supplier.Address),
			))
		}

	case "supplies":
		// Заголовки для поставок
		content.WriteString("ID;Поставщик;Товар;Количество;Цена за ед.;Итого;Дата;Статус;Примечания\r\n")
		for _, supply := range supplies {
			supplierName := "Неизвестный"
			if sup, exists := suppliers[supply.SupplierID]; exists {
				supplierName = sup.Name
			}

			productName := "Неизвестный"
			if prod, exists := products[supply.ProductID]; exists {
				productName = prod.Name
			}

			content.WriteString(fmt.Sprintf("%d;%s;%s;%d;%.2f;%.2f;%s;%s;%s\r\n",
				supply.ID,
				escapeCSV(supplierName),
				escapeCSV(productName),
				supply.Quantity,
				supply.Price,
				supply.Total,
				supply.Date.Format("02.01.2006"),
				escapeCSV(getStatusText(supply.Status)),
				escapeCSV(supply.Notes),
			))
		}

	case "categories":
		// Заголовки для категорий
		content.WriteString("ID;Название;Описание\r\n")
		for _, category := range categories {
			content.WriteString(fmt.Sprintf("%d;%s;%s\r\n",
				category.ID,
				escapeCSV(category.Name),
				escapeCSV(category.Description),
			))
		}

	default:
		// Экспорт всего в один файл
		content.WriteString("=== МАГАЗИНЫ ===\r\n")
		content.WriteString("ID;Название;Адрес;Дата создания\r\n")
		for _, store := range stores {
			content.WriteString(fmt.Sprintf("%d;%s;%s;%s\r\n",
				store.ID,
				escapeCSV(store.Name),
				escapeCSV(store.Address),
				store.CreatedAt.Format("02.01.2006"),
			))
		}

		content.WriteString("\r\n=== КАТЕГОРИИ ===\r\n")
		content.WriteString("ID;Название;Описание\r\n")
		for _, category := range categories {
			content.WriteString(fmt.Sprintf("%d;%s;%s\r\n",
				category.ID,
				escapeCSV(category.Name),
				escapeCSV(category.Description),
			))
		}

		content.WriteString("\r\n=== ТОВАРЫ ===\r\n")
		content.WriteString("ID;Название;Артикул;Цена;Описание;Категория\r\n")
		for _, product := range products {
			categoryName := "Без категории"
			if cat, exists := categories[product.CategoryID]; exists {
				categoryName = cat.Name
			}

			content.WriteString(fmt.Sprintf("%d;%s;%s;%.2f;%s;%s\r\n",
				product.ID,
				escapeCSV(product.Name),
				escapeCSV(product.SKU),
				product.Price,
				escapeCSV(product.Description),
				escapeCSV(categoryName),
			))
		}

		content.WriteString("\r\n=== ПОСТАВЩИКИ ===\r\n")
		content.WriteString("ID;Название;Телефон;Email;Адрес\r\n")
		for _, supplier := range suppliers {
			content.WriteString(fmt.Sprintf("%d;%s;%s;%s;%s\r\n",
				supplier.ID,
				escapeCSV(supplier.Name),
				escapeCSV(supplier.Phone),
				escapeCSV(supplier.Email),
				escapeCSV(supplier.Address),
			))
		}

		content.WriteString("\r\n=== ПОСТАВКИ ===\r\n")
		content.WriteString("ID;Поставщик;Товар;Количество;Цена за ед.;Итого;Дата;Статус\r\n")
		for _, supply := range supplies {
			supplierName := "Неизвестный"
			if sup, exists := suppliers[supply.SupplierID]; exists {
				supplierName = sup.Name
			}

			productName := "Неизвестный"
			if prod, exists := products[supply.ProductID]; exists {
				productName = prod.Name
			}

			content.WriteString(fmt.Sprintf("%d;%s;%s;%d;%.2f;%.2f;%s;%s\r\n",
				supply.ID,
				escapeCSV(supplierName),
				escapeCSV(productName),
				supply.Quantity,
				supply.Price,
				supply.Total,
				supply.Date.Format("02.01.2006"),
				escapeCSV(getStatusText(supply.Status)),
			))
		}
	}

	// Устанавливаем заголовки
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Отправляем данные
	c.Data(http.StatusOK, "text/csv; charset=utf-8", content.Bytes())
}

// Функция для экранирования CSV значений
func escapeCSV(value string) string {
	// Заменяем точку с запятой на запятую и удаляем переносы строк
	result := strings.ReplaceAll(value, ";", ",")
	result = strings.ReplaceAll(result, "\r\n", " ")
	result = strings.ReplaceAll(result, "\n", " ")
	result = strings.ReplaceAll(result, "\r", " ")
	return result
}

// Экспорт в Excel - используем реальный Excel файл
func exportExcel(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Просто вызываем тот же CSV, но с другим именем файла
	// Excel откроет CSV если он правильно сформирован
	dataType := c.Query("type")
	filename := "export.xlsx"

	// Меняем заголовки, чтобы Excel знал что это CSV
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Редирект на CSV функцию
	c.Redirect(http.StatusFound, fmt.Sprintf("/export/csv?type=%s", dataType))
}

// ====================================================================
// ОСТАВШИЙСЯ КОД (без изменений)
// ====================================================================

// ОБРАБОТЧИКИ КАТЕГОРИЙ
func categoriesListPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	mu.RLock()
	categoryList := make([]Category, 0, len(categories))
	for _, cat := range categories {
		categoryList = append(categoryList, cat)
	}
	mu.RUnlock()

	sortBy := c.DefaultQuery("sort", "id")
	direction := c.DefaultQuery("direction", "asc")
	searchQuery := c.Query("search")

	if searchQuery != "" {
		filteredList := make([]Category, 0)
		searchLower := strings.ToLower(searchQuery)
		for _, category := range categoryList {
			if strings.Contains(strings.ToLower(category.Name), searchLower) ||
				strings.Contains(strings.ToLower(category.Description), searchLower) {
				filteredList = append(filteredList, category)
			}
		}
		categoryList = filteredList
	}

	sortCategories(categoryList, sortBy, direction)

	c.HTML(http.StatusOK, "categories.html", gin.H{
		"username":     getUsername(c),
		"categories":   categoryList,
		"total_count":  len(categoryList),
		"search_query": searchQuery,
		"sort_by":      sortBy,
		"direction":    direction,
	})
}

func sortCategories(categories []Category, sortBy, direction string) {
	switch sortBy {
	case "id":
		if direction == "asc" {
			sort.Slice(categories, func(i, j int) bool {
				return categories[i].ID < categories[j].ID
			})
		} else {
			sort.Slice(categories, func(i, j int) bool {
				return categories[i].ID > categories[j].ID
			})
		}
	case "name":
		if direction == "asc" {
			sort.Slice(categories, func(i, j int) bool {
				return strings.ToLower(categories[i].Name) < strings.ToLower(categories[j].Name)
			})
		} else {
			sort.Slice(categories, func(i, j int) bool {
				return strings.ToLower(categories[i].Name) > strings.ToLower(categories[j].Name)
			})
		}
	}
}

func categoryCreatePage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	category := Category{}

	c.HTML(http.StatusOK, "category_form.html", gin.H{
		"title":    "Создать категорию",
		"action":   "/categories/create",
		"category": category,
	})
}

func categoryCreateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")

	if name == "" {
		c.HTML(http.StatusOK, "category_form.html", gin.H{
			"title":  "Создать категорию",
			"action": "/categories/create",
			"error":  "Название категории обязательно",
			"category": Category{
				Name:        name,
				Description: description,
			},
		})
		return
	}

	mu.Lock()

	newID := 1
	for {
		if _, exists := categories[newID]; !exists {
			break
		}
		newID++
	}

	categories[newID] = Category{
		ID:          newID,
		Name:        name,
		Description: description,
	}

	if newID > categoryID {
		categoryID = newID
	}

	mu.Unlock()

	c.Redirect(http.StatusFound, "/categories")
}

func categoryEditPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	mu.RLock()
	category, exists := categories[id]
	mu.RUnlock()

	if !exists {
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	c.HTML(http.StatusOK, "category_form.html", gin.H{
		"title":    "Редактировать категорию",
		"action":   fmt.Sprintf("/categories/%d/edit", id),
		"category": category,
	})
}

func categoryUpdateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")

	if name == "" {
		c.HTML(http.StatusOK, "category_form.html", gin.H{
			"title":  "Редактировать категорию",
			"action": fmt.Sprintf("/categories/%d/edit", id),
			"error":  "Название категории обязательно",
			"category": Category{
				ID:          id,
				Name:        name,
				Description: description,
			},
		})
		return
	}

	mu.Lock()
	category, exists := categories[id]
	if !exists {
		mu.Unlock()
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	category.Name = name
	category.Description = description
	categories[id] = category
	mu.Unlock()

	c.Redirect(http.StatusFound, "/categories")
}

func categoryDeleteHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	mu.Lock()
	delete(categories, id)
	mu.Unlock()

	c.Redirect(http.StatusFound, "/categories")
}

// АВТОРИЗАЦИЯ
func loginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "admin" && password == "admin123" {
		c.SetCookie("auth", "true", 3600, "/", "", false, true)
		c.SetCookie("username", username, 3600, "/", "", false, false)
		c.Redirect(http.StatusFound, "/main")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"error": "Неверный логин или пароль",
	})
}

func logoutHandler(c *gin.Context) {
	c.SetCookie("auth", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func checkAuth(c *gin.Context) bool {
	auth, err := c.Cookie("auth")
	return err == nil && auth == "true"
}

// ГЛАВНОЕ МЕНЮ
func mainMenuPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "stores.html", gin.H{
		"username": getUsername(c),
	})
}

func getUsername(c *gin.Context) string {
	username, err := c.Cookie("username")
	if err != nil {
		return "Гость"
	}
	return username
}

// ПОСТАВЩИКИ
func suppliersPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	mu.RLock()
	supplierList := make([]Supplier, 0, len(suppliers))
	for _, supplier := range suppliers {
		supplierList = append(supplierList, supplier)
	}
	mu.RUnlock()

	sortBy := c.DefaultQuery("sort", "id")
	direction := c.DefaultQuery("direction", "asc")
	searchQuery := c.Query("search")
	if searchQuery != "" {
		filteredList := make([]Supplier, 0)
		searchLower := strings.ToLower(searchQuery)
		for _, supplier := range supplierList {
			if strings.Contains(strings.ToLower(supplier.Name), searchLower) ||
				strings.Contains(strings.ToLower(supplier.Phone), searchLower) ||
				strings.Contains(strings.ToLower(supplier.Email), searchLower) ||
				strings.Contains(strings.ToLower(supplier.Address), searchLower) {
				filteredList = append(filteredList, supplier)
			}
		}
		supplierList = filteredList
	}

	sortSuppliers(supplierList, sortBy, direction)

	c.HTML(http.StatusOK, "suppliers.html", gin.H{
		"username":     getUsername(c),
		"suppliers":    supplierList,
		"total_count":  len(supplierList),
		"search_query": searchQuery,
		"sort_by":      sortBy,
		"direction":    direction,
	})
}

func sortSuppliers(suppliers []Supplier, sortBy, direction string) {
	switch sortBy {
	case "id":
		if direction == "asc" {
			sort.Slice(suppliers, func(i, j int) bool {
				return suppliers[i].ID < suppliers[j].ID
			})
		} else {
			sort.Slice(suppliers, func(i, j int) bool {
				return suppliers[i].ID > suppliers[j].ID
			})
		}
	case "name":
		if direction == "asc" {
			sort.Slice(suppliers, func(i, j int) bool {
				return strings.ToLower(suppliers[i].Name) < strings.ToLower(suppliers[j].Name)
			})
		} else {
			sort.Slice(suppliers, func(i, j int) bool {
				return strings.ToLower(suppliers[i].Name) > strings.ToLower(suppliers[j].Name)
			})
		}
	case "phone":
		if direction == "asc" {
			sort.Slice(suppliers, func(i, j int) bool {
				return suppliers[i].Phone < suppliers[j].Phone
			})
		} else {
			sort.Slice(suppliers, func(i, j int) bool {
				return suppliers[i].Phone > suppliers[j].Phone
			})
		}
	}
}

func supplierCreatePage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "supplier_form.html", gin.H{
		"title":  "Добавить поставщика",
		"action": "/suppliers/create",
	})
}

func supplierCreateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	name := c.PostForm("name")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	address := c.PostForm("address")

	if name == "" {
		c.HTML(http.StatusOK, "supplier_form.html", gin.H{
			"title":  "Добавить поставщика",
			"action": "/suppliers/create",
			"error":  "Название обязательно",
			"supplier": Supplier{
				Name:    name,
				Phone:   phone,
				Email:   email,
				Address: address,
			},
		})
		return
	}

	mu.Lock()

	newID := 1
	for {
		if _, exists := suppliers[newID]; !exists {
			break
		}
		newID++
	}

	suppliers[newID] = Supplier{
		ID:      newID,
		Name:    name,
		Phone:   phone,
		Email:   email,
		Address: address,
	}

	if newID > supplierID {
		supplierID = newID
	}

	mu.Unlock()

	c.Redirect(http.StatusFound, "/suppliers")
}

func supplierEditPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	mu.RLock()
	supplier, exists := suppliers[id]
	mu.RUnlock()

	if !exists {
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	c.HTML(http.StatusOK, "supplier_form.html", gin.H{
		"title":    "Редактировать поставщика",
		"action":   fmt.Sprintf("/suppliers/%d/edit", id),
		"supplier": supplier,
	})
}

func supplierUpdateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	mu.Lock()
	supplier, exists := suppliers[id]
	if !exists {
		mu.Unlock()
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	supplier.Name = c.PostForm("name")
	supplier.Phone = c.PostForm("phone")
	supplier.Email = c.PostForm("email")
	supplier.Address = c.PostForm("address")

	suppliers[id] = supplier
	mu.Unlock()

	c.Redirect(http.StatusFound, "/suppliers")
}

func supplierDeleteHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	mu.Lock()
	delete(suppliers, id)
	mu.Unlock()

	c.Redirect(http.StatusFound, "/suppliers")
}

// СТРАНИЦА МАГАЗИНОВ
func storeListPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	mu.RLock()
	storeList := make([]Store, 0, len(stores))
	for _, store := range stores {
		storeList = append(storeList, store)
	}
	mu.RUnlock()

	sortBy := c.DefaultQuery("sort", "id")
	direction := c.DefaultQuery("direction", "asc")

	sortStores(storeList, sortBy, direction)

	c.HTML(http.StatusOK, "store.html", gin.H{
		"stores":      storeList,
		"total_count": len(storeList),
		"sort_by":     sortBy,
		"direction":   direction,
	})
}

func sortStores(stores []Store, sortBy, direction string) {
	switch sortBy {
	case "id":
		if direction == "asc" {
			sort.Slice(stores, func(i, j int) bool {
				return stores[i].ID < stores[j].ID
			})
		} else {
			sort.Slice(stores, func(i, j int) bool {
				return stores[i].ID > stores[j].ID
			})
		}
	case "name":
		if direction == "asc" {
			sort.Slice(stores, func(i, j int) bool {
				return strings.ToLower(stores[i].Name) < strings.ToLower(stores[j].Name)
			})
		} else {
			sort.Slice(stores, func(i, j int) bool {
				return strings.ToLower(stores[i].Name) > strings.ToLower(stores[j].Name)
			})
		}
	case "created_at":
		if direction == "asc" {
			sort.Slice(stores, func(i, j int) bool {
				return stores[i].CreatedAt.Before(stores[j].CreatedAt)
			})
		} else {
			sort.Slice(stores, func(i, j int) bool {
				return stores[i].CreatedAt.After(stores[j].CreatedAt)
			})
		}
	}
}

func storeCreatePage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "store_form.html", gin.H{
		"title":  "Добавить магазин",
		"action": "/stores/create",
	})
}

func storeCreateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	name := c.PostForm("name")
	address := c.PostForm("address")

	if name == "" {
		c.HTML(http.StatusOK, "store_form.html", gin.H{
			"title":  "Добавить магазин",
			"action": "/stores/create",
			"error":  "Название обязательно",
			"store":  Store{Name: name, Address: address},
		})
		return
	}

	mu.Lock()

	newID := 1
	for {
		if _, exists := stores[newID]; !exists {
			break
		}
		newID++
	}

	stores[newID] = Store{
		ID:        newID,
		Name:      name,
		Address:   address,
		CreatedAt: time.Now(),
	}

	if newID > storeID {
		storeID = newID
	}

	mu.Unlock()

	c.Redirect(http.StatusFound, "/stores")
}

func storeEditPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	mu.RLock()
	store, exists := stores[id]
	mu.RUnlock()

	if !exists {
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	c.HTML(http.StatusOK, "store_form.html", gin.H{
		"title":  "Редактировать магазин",
		"action": fmt.Sprintf("/stores/%d/edit", id),
		"store":  store,
	})
}

func storeUpdateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	mu.Lock()
	store, exists := stores[id]
	if !exists {
		mu.Unlock()
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	store.Name = c.PostForm("name")
	store.Address = c.PostForm("address")
	stores[id] = store
	mu.Unlock()

	c.Redirect(http.StatusFound, "/stores")
}

func storeDeleteHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	mu.Lock()
	delete(stores, id)
	mu.Unlock()

	c.Redirect(http.StatusFound, "/stores")
}

// ТОВАРЫ
func productsList(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	mu.RLock()
	productList := make([]Product, 0, len(products))
	for _, product := range products {
		productList = append(productList, product)
	}
	mu.RUnlock()

	sortBy := c.DefaultQuery("sort", "id")
	direction := c.DefaultQuery("direction", "asc")
	searchQuery := c.Query("search")

	if searchQuery != "" {
		filteredList := make([]Product, 0)
		searchLower := strings.ToLower(searchQuery)
		for _, product := range productList {
			if strings.Contains(strings.ToLower(product.Name), searchLower) ||
				strings.Contains(strings.ToLower(product.Description), searchLower) ||
				strings.Contains(strings.ToLower(product.SKU), searchLower) {
				filteredList = append(filteredList, product)
			}
		}
		productList = filteredList
	}

	sortProducts(productList, sortBy, direction)

	var totalPrice float64
	var minPrice, maxPrice float64
	if len(productList) > 0 {
		minPrice = productList[0].Price
		maxPrice = productList[0].Price
	}

	for _, p := range productList {
		totalPrice += p.Price
		if p.Price < minPrice {
			minPrice = p.Price
		}
		if p.Price > maxPrice {
			maxPrice = p.Price
		}
	}

	avgPrice := 0.0
	if len(productList) > 0 {
		avgPrice = totalPrice / float64(len(productList))
	}

	c.HTML(http.StatusOK, "products.html", gin.H{
		"username":      getUsername(c),
		"products":      productList,
		"total_count":   len(productList),
		"average_price": fmt.Sprintf("%.2f", avgPrice),
		"min_price":     fmt.Sprintf("%.2f", minPrice),
		"max_price":     fmt.Sprintf("%.2f", maxPrice),
		"total_value":   fmt.Sprintf("%.2f", totalPrice),
		"search_query":  searchQuery,
		"sort_by":       sortBy,
		"direction":     direction,
	})
}

func sortProducts(products []Product, sortBy, direction string) {
	switch sortBy {
	case "id":
		if direction == "asc" {
			sort.Slice(products, func(i, j int) bool {
				return products[i].ID < products[j].ID
			})
		} else {
			sort.Slice(products, func(i, j int) bool {
				return products[i].ID > products[j].ID
			})
		}
	case "name":
		if direction == "asc" {
			sort.Slice(products, func(i, j int) bool {
				return strings.ToLower(products[i].Name) < strings.ToLower(products[j].Name)
			})
		} else {
			sort.Slice(products, func(i, j int) bool {
				return strings.ToLower(products[i].Name) > strings.ToLower(products[j].Name)
			})
		}
	case "price":
		if direction == "asc" {
			sort.Slice(products, func(i, j int) bool {
				return products[i].Price < products[j].Price
			})
		} else {
			sort.Slice(products, func(i, j int) bool {
				return products[i].Price > products[j].Price
			})
		}
	case "sku":
		if direction == "asc" {
			sort.Slice(products, func(i, j int) bool {
				return strings.ToLower(products[i].SKU) < strings.ToLower(products[j].SKU)
			})
		} else {
			sort.Slice(products, func(i, j int) bool {
				return strings.ToLower(products[i].SKU) > strings.ToLower(products[j].SKU)
			})
		}
	}
}

func productCreatePage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	mu.RLock()
	categoryList := make([]Category, 0, len(categories))
	for _, cat := range categories {
		categoryList = append(categoryList, cat)
	}
	mu.RUnlock()

	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":      "Добавить товар",
		"action":     "/products/create",
		"product":    Product{},
		"categories": categoryList,
		"mode":       "create",
	})
}

func productCreateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	description := c.PostForm("description")
	sku := c.PostForm("sku")
	categoryIDStr := c.PostForm("category_id")

	if name == "" || priceStr == "" {
		mu.RLock()
		categoryList := make([]Category, 0, len(categories))
		for _, cat := range categories {
			categoryList = append(categoryList, cat)
		}
		mu.RUnlock()

		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Добавить товар",
			"action": "/products/create",
			"error":  "Название и цена обязательны",
			"product": Product{
				Name:        name,
				Description: description,
				SKU:         sku,
			},
			"categories": categoryList,
			"mode":       "create",
		})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		mu.RLock()
		categoryList := make([]Category, 0, len(categories))
		for _, cat := range categories {
			categoryList = append(categoryList, cat)
		}
		mu.RUnlock()

		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Добавить товар",
			"action": "/products/create",
			"error":  "Неверный формат цены",
			"product": Product{
				Name:        name,
				Description: description,
				SKU:         sku,
			},
			"categories": categoryList,
			"mode":       "create",
		})
		return
	}

	var categoryID int
	if categoryIDStr != "" {
		if id, err := strconv.Atoi(categoryIDStr); err == nil {
			categoryID = id
		}
	}

	mu.Lock()

	newID := 1
	for {
		if _, exists := products[newID]; !exists {
			break
		}
		newID++
	}

	products[newID] = Product{
		ID:          newID,
		Name:        name,
		Price:       price,
		Description: description,
		SKU:         sku,
		CategoryID:  categoryID,
		CreatedAt:   time.Now(),
	}

	if newID > productID {
		productID = newID
	}

	mu.Unlock()

	c.Redirect(http.StatusFound, "/products")
}

func productEditPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/products")
		return
	}

	mu.RLock()
	product, exists := products[id]

	categoryList := make([]Category, 0, len(categories))
	for _, cat := range categories {
		categoryList = append(categoryList, cat)
	}
	mu.RUnlock()

	if !exists {
		c.Redirect(http.StatusFound, "/products")
		return
	}

	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":      "Редактировать товар",
		"action":     fmt.Sprintf("/products/%d/edit", id),
		"product":    product,
		"categories": categoryList,
		"mode":       "edit",
	})
}

func productUpdateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/products")
		return
	}

	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	description := c.PostForm("description")
	sku := c.PostForm("sku")
	categoryIDStr := c.PostForm("category_id")

	if name == "" || priceStr == "" {
		mu.RLock()
		categoryList := make([]Category, 0, len(categories))
		for _, cat := range categories {
			categoryList = append(categoryList, cat)
		}
		mu.RUnlock()

		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Редактировать товар",
			"action": fmt.Sprintf("/products/%d/edit", id),
			"error":  "Название и цена обязательны",
			"product": Product{
				ID:          id,
				Name:        name,
				Description: description,
				SKU:         sku,
			},
			"categories": categoryList,
			"mode":       "edit",
		})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		mu.RLock()
		categoryList := make([]Category, 0, len(categories))
		for _, cat := range categories {
			categoryList = append(categoryList, cat)
		}
		mu.RUnlock()

		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Редактировать товар",
			"action": fmt.Sprintf("/products/%d/edit", id),
			"error":  "Неверный формат цены",
			"product": Product{
				ID:          id,
				Name:        name,
				Description: description,
				SKU:         sku,
			},
			"categories": categoryList,
			"mode":       "edit",
		})
		return
	}

	var categoryID int
	if categoryIDStr != "" {
		if id, err := strconv.Atoi(categoryIDStr); err == nil {
			categoryID = id
		}
	}

	mu.Lock()
	product, exists := products[id]
	if !exists {
		mu.Unlock()
		c.Redirect(http.StatusFound, "/products")
		return
	}

	product.Name = name
	product.Price = price
	product.Description = description
	product.SKU = sku
	product.CategoryID = categoryID

	products[id] = product
	mu.Unlock()

	c.Redirect(http.StatusFound, "/products")
}

func productDeleteHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/products")
		return
	}

	mu.Lock()
	delete(products, id)
	mu.Unlock()

	c.Redirect(http.StatusFound, "/products")
}

// ПОСТАВКИ
func suppliesPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	mu.RLock()
	supplyList := make([]Supply, 0, len(supplies))
	for _, supply := range supplies {
		supplyList = append(supplyList, supply)
	}

	supplierMap := make(map[int]Supplier)
	for id, supplier := range suppliers {
		supplierMap[id] = supplier
	}

	productMap := make(map[int]Product)
	for id, product := range products {
		productMap[id] = product
	}
	mu.RUnlock()

	sortBy := c.DefaultQuery("sort", "id")
	direction := c.DefaultQuery("direction", "asc")
	searchQuery := c.Query("search")
	statusFilter := c.Query("status")

	if searchQuery != "" || statusFilter != "" {
		filteredList := make([]Supply, 0)
		searchLower := strings.ToLower(searchQuery)

		for _, supply := range supplyList {
			if statusFilter != "" && supply.Status != statusFilter {
				continue
			}

			if searchQuery != "" {
				supplier := supplierMap[supply.SupplierID]
				product := productMap[supply.ProductID]

				if strings.Contains(strings.ToLower(supplier.Name), searchLower) ||
					strings.Contains(strings.ToLower(product.Name), searchLower) ||
					strings.Contains(strings.ToLower(supply.Notes), searchLower) {
					filteredList = append(filteredList, supply)
				}
			} else {
				filteredList = append(filteredList, supply)
			}
		}
		supplyList = filteredList
	}

	sortSupplies(supplyList, sortBy, direction)

	var totalQuantity int
	var totalCost float64
	for _, s := range supplyList {
		totalQuantity += s.Quantity
		totalCost += s.Total
	}

	supplyData := make([]gin.H, 0, len(supplyList))
	for _, supply := range supplyList {
		supplier := supplierMap[supply.SupplierID]
		product := productMap[supply.ProductID]

		supplyData = append(supplyData, gin.H{
			"ID":          supply.ID,
			"SupplierID":  supply.SupplierID,
			"Supplier":    supplier.Name,
			"ProductID":   supply.ProductID,
			"Product":     product.Name,
			"Quantity":    supply.Quantity,
			"Price":       supply.Price,
			"Total":       supply.Total,
			"Date":        supply.Date.Format("02.01.2006"),
			"Status":      supply.Status,
			"StatusText":  getStatusText(supply.Status),
			"StatusClass": getStatusClass(supply.Status),
			"Notes":       supply.Notes,
			"CreatedAt":   supply.CreatedAt.Format("02.01.2006 15:04"),
		})
	}

	c.HTML(http.StatusOK, "supplies.html", gin.H{
		"username":       getUsername(c),
		"supplies":       supplyData,
		"total_count":    len(supplyList),
		"total_quantity": totalQuantity,
		"total_cost":     fmt.Sprintf("%.2f", totalCost),
		"search_query":   searchQuery,
		"status_filter":  statusFilter,
		"sort_by":        sortBy,
		"direction":      direction,
	})
}

func sortSupplies(supplies []Supply, sortBy, direction string) {
	switch sortBy {
	case "id":
		if direction == "asc" {
			sort.Slice(supplies, func(i, j int) bool {
				return supplies[i].ID < supplies[j].ID
			})
		} else {
			sort.Slice(supplies, func(i, j int) bool {
				return supplies[i].ID > supplies[j].ID
			})
		}
	case "date":
		if direction == "asc" {
			sort.Slice(supplies, func(i, j int) bool {
				return supplies[i].Date.Before(supplies[j].Date)
			})
		} else {
			sort.Slice(supplies, func(i, j int) bool {
				return supplies[i].Date.After(supplies[j].Date)
			})
		}
	case "total":
		if direction == "asc" {
			sort.Slice(supplies, func(i, j int) bool {
				return supplies[i].Total < supplies[j].Total
			})
		} else {
			sort.Slice(supplies, func(i, j int) bool {
				return supplies[i].Total > supplies[j].Total
			})
		}
	case "quantity":
		if direction == "asc" {
			sort.Slice(supplies, func(i, j int) bool {
				return supplies[i].Quantity < supplies[j].Quantity
			})
		} else {
			sort.Slice(supplies, func(i, j int) bool {
				return supplies[i].Quantity > supplies[j].Quantity
			})
		}
	}
}

func getStatusText(status string) string {
	switch status {
	case "pending":
		return "Ожидается"
	case "delivered":
		return "Доставлено"
	case "cancelled":
		return "Отменено"
	default:
		return status
	}
}

func getStatusClass(status string) string {
	switch status {
	case "pending":
		return "status-pending"
	case "delivered":
		return "status-delivered"
	case "cancelled":
		return "status-cancelled"
	default:
		return "status-default"
	}
}

func supplyCreatePage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	mu.RLock()
	supplierList := make([]Supplier, 0, len(suppliers))
	for _, supplier := range suppliers {
		supplierList = append(supplierList, supplier)
	}

	productList := make([]Product, 0, len(products))
	for _, product := range products {
		productList = append(productList, product)
	}
	mu.RUnlock()

	c.HTML(http.StatusOK, "supplies_form.html", gin.H{
		"title":     "Создать поставку",
		"action":    "/supplies/create",
		"supply":    Supply{Date: time.Now(), Status: "pending"},
		"suppliers": supplierList,
		"products":  productList,
		"mode":      "create",
		"now":       time.Now(),
	})
}

func supplyCreateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	supplierIDStr := c.PostForm("supplier_id")
	productIDStr := c.PostForm("product_id")
	quantityStr := c.PostForm("quantity")
	priceStr := c.PostForm("price")
	dateStr := c.PostForm("date")
	status := c.PostForm("status")
	notes := c.PostForm("notes")

	if supplierIDStr == "" || productIDStr == "" || quantityStr == "" || priceStr == "" {
		mu.RLock()
		supplierList := make([]Supplier, 0, len(suppliers))
		for _, supplier := range suppliers {
			supplierList = append(supplierList, supplier)
		}

		productList := make([]Product, 0, len(products))
		for _, product := range products {
			productList = append(productList, product)
		}
		mu.RUnlock()

		c.HTML(http.StatusOK, "supplies_form.html", gin.H{
			"title":  "Создать поставку",
			"action": "/supplies/create",
			"error":  "Заполните обязательные поля",
			"supply": Supply{
				Status: status,
				Notes:  notes,
			},
			"suppliers": supplierList,
			"products":  productList,
			"mode":      "create",
			"now":       time.Now(),
		})
		return
	}

	supplierID, err1 := strconv.Atoi(supplierIDStr)
	productID, err2 := strconv.Atoi(productIDStr)
	quantity, err3 := strconv.Atoi(quantityStr)
	price, err4 := strconv.ParseFloat(priceStr, 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || quantity <= 0 || price < 0 {
		mu.RLock()
		supplierList := make([]Supplier, 0, len(suppliers))
		for _, supplier := range suppliers {
			supplierList = append(supplierList, supplier)
		}

		productList := make([]Product, 0, len(products))
		for _, product := range products {
			productList = append(productList, product)
		}
		mu.RUnlock()

		c.HTML(http.StatusOK, "supplies_form.html", gin.H{
			"title":  "Создать поставку",
			"action": "/supplies/create",
			"error":  "Неверный формат данных",
			"supply": Supply{
				Status: status,
				Notes:  notes,
			},
			"suppliers": supplierList,
			"products":  productList,
			"mode":      "create",
			"now":       time.Now(),
		})
		return
	}

	var date time.Time
	if dateStr != "" {
		date, _ = time.Parse("2006-01-02", dateStr)
	} else {
		date = time.Now()
	}

	total := float64(quantity) * price

	mu.Lock()

	newID := 1
	for {
		if _, exists := supplies[newID]; !exists {
			break
		}
		newID++
	}

	supplies[newID] = Supply{
		ID:         newID,
		SupplierID: supplierID,
		ProductID:  productID,
		Quantity:   quantity,
		Price:      price,
		Total:      total,
		Date:       date,
		Status:     status,
		Notes:      notes,
		CreatedAt:  time.Now(),
	}

	if newID > supplyID {
		supplyID = newID
	}

	mu.Unlock()

	c.Redirect(http.StatusFound, "/supplies")
}

func supplyEditPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/supplies")
		return
	}

	mu.RLock()
	supply, exists := supplies[id]

	supplierList := make([]Supplier, 0, len(suppliers))
	for _, supplier := range suppliers {
		supplierList = append(supplierList, supplier)
	}

	productList := make([]Product, 0, len(products))
	for _, product := range products {
		productList = append(productList, product)
	}
	mu.RUnlock()

	if !exists {
		c.Redirect(http.StatusFound, "/supplies")
		return
	}

	c.HTML(http.StatusOK, "supplies_form.html", gin.H{
		"title":     "Редактировать поставку",
		"action":    fmt.Sprintf("/supplies/%d/edit", id),
		"supply":    supply,
		"suppliers": supplierList,
		"products":  productList,
		"mode":      "edit",
		"now":       time.Now(),
	})
}

func supplyUpdateHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/supplies")
		return
	}

	supplierIDStr := c.PostForm("supplier_id")
	productIDStr := c.PostForm("product_id")
	quantityStr := c.PostForm("quantity")
	priceStr := c.PostForm("price")
	dateStr := c.PostForm("date")
	status := c.PostForm("status")
	notes := c.PostForm("notes")

	if supplierIDStr == "" || productIDStr == "" || quantityStr == "" || priceStr == "" {
		mu.RLock()
		supplierList := make([]Supplier, 0, len(suppliers))
		for _, supplier := range suppliers {
			supplierList = append(supplierList, supplier)
		}

		productList := make([]Product, 0, len(products))
		for _, product := range products {
			productList = append(productList, product)
		}
		mu.RUnlock()

		c.HTML(http.StatusOK, "supplies_form.html", gin.H{
			"title":  "Редактировать поставку",
			"action": fmt.Sprintf("/supplies/%d/edit", id),
			"error":  "Заполните обязательные поля",
			"supply": Supply{
				ID:     id,
				Status: status,
				Notes:  notes,
			},
			"suppliers": supplierList,
			"products":  productList,
			"mode":      "edit",
			"now":       time.Now(),
		})
		return
	}

	supplierID, err1 := strconv.Atoi(supplierIDStr)
	productID, err2 := strconv.Atoi(productIDStr)
	quantity, err3 := strconv.Atoi(quantityStr)
	price, err4 := strconv.ParseFloat(priceStr, 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || quantity <= 0 || price < 0 {
		mu.RLock()
		supplierList := make([]Supplier, 0, len(suppliers))
		for _, supplier := range suppliers {
			supplierList = append(supplierList, supplier)
		}

		productList := make([]Product, 0, len(products))
		for _, product := range products {
			productList = append(productList, product)
		}
		mu.RUnlock()

		c.HTML(http.StatusOK, "supplies_form.html", gin.H{
			"title":  "Редактировать поставку",
			"action": fmt.Sprintf("/supplies/%d/edit", id),
			"error":  "Неверный формат данных",
			"supply": Supply{
				ID:     id,
				Status: status,
				Notes:  notes,
			},
			"suppliers": supplierList,
			"products":  productList,
			"mode":      "edit",
			"now":       time.Now(),
		})
		return
	}

	var date time.Time
	if dateStr != "" {
		date, _ = time.Parse("2006-01-02", dateStr)
	} else {
		date = time.Now()
	}

	total := float64(quantity) * price

	mu.Lock()
	supply, exists := supplies[id]
	if !exists {
		mu.Unlock()
		c.Redirect(http.StatusFound, "/supplies")
		return
	}

	supply.SupplierID = supplierID
	supply.ProductID = productID
	supply.Quantity = quantity
	supply.Price = price
	supply.Total = total
	supply.Date = date
	supply.Status = status
	supply.Notes = notes

	supplies[id] = supply
	mu.Unlock()

	c.Redirect(http.StatusFound, "/supplies")
}

func supplyDeleteHandler(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/supplies")
		return
	}

	mu.Lock()
	delete(supplies, id)
	mu.Unlock()

	c.Redirect(http.StatusFound, "/supplies")
}

// СТАТИЧЕСКИЕ СТРАНИЦЫ
func aboutPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "about.html", gin.H{})
}

func contactsPage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "contacts.html", gin.H{})
}
