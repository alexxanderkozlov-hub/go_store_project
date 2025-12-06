package main

import (
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

// ====================================================================
// Хранилище в памяти
// ====================================================================

var (
	stores     = make(map[int]Store)
	suppliers  = make(map[int]Supplier)
	products   = make(map[int]Product)
	categories = make(map[int]Category)
	users      = make(map[int]User)
	storeID    = 0
	supplierID = 0
	productID  = 0
	categoryID = 0
	mu         sync.RWMutex
)

func init() {
	// Инициализация для корректного поиска первого свободного ID
	maxID := 0
	for id := range stores {
		if id > maxID {
			maxID = id
		}
	}
	storeID = maxID
}

func main() {
	// Настройка Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Загружаем HTML шаблоны
	r.LoadHTMLGlob("templates/*.html")

	// Статические файлы
	r.Static("/static", "./static")

	// Маршруты
	setupRoutes(r)

	// Запуск сервера
	port := ":8080"
	log.Printf("🚀 Сервер запущен на http://localhost%s", port)

	// Пробуем разные порты если 8080 занят
	err := r.Run(port)
	if err != nil {
		log.Printf("Порт 8080 занят, пробуем 8081...")
		r.Run(":8081")
	}
}

func setupRoutes(r *gin.Engine) {
	// Главная страница (редирект на логин)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// Главное меню (stores.html)
	r.GET("/main", mainMenuPage)

	// Авторизация
	r.GET("/login", loginPage)
	r.POST("/login", loginHandler)
	r.GET("/logout", logoutHandler)

	// Страница магазинов (store.html)
	r.GET("/stores", storeListPage)
	r.GET("/stores/create", storeCreatePage)
	r.POST("/stores/create", storeCreateHandler)
	r.GET("/stores/:id/edit", storeEditPage)
	r.POST("/stores/:id/edit", storeUpdateHandler)
	r.GET("/stores/:id/delete", storeDeleteHandler)

	// Категории товаров (ПОЛНЫЙ CRUD)
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

	// ПОСТАВЩИКИ
	r.GET("/suppliers", suppliersPage)
	r.GET("/suppliers/create", supplierCreatePage)
	r.POST("/suppliers/create", supplierCreateHandler)
	r.GET("/suppliers/:id/edit", supplierEditPage)
	r.POST("/suppliers/:id/edit", supplierUpdateHandler)
	r.GET("/suppliers/:id/delete", supplierDeleteHandler)

	// Статические страницы
	r.GET("/about", aboutPage)
	r.GET("/contacts", contactsPage)
}

// ====================================================================
// ОБРАБОТЧИКИ КАТЕГОРИЙ
// ====================================================================

// Список категорий с сортировкой и поиском
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

	// Получаем параметры сортировки и поиска
	sortBy := c.DefaultQuery("sort", "id")
	direction := c.DefaultQuery("direction", "asc")
	searchQuery := c.Query("search")

	// Фильтрация по поиску
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

	// Сортируем категории
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

// Функция сортировки категорий
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

// Страница создания категории
func categoryCreatePage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Создаем пустую категорию для формы
	category := Category{}

	c.HTML(http.StatusOK, "category_form.html", gin.H{
		"title":    "Создать категорию",
		"action":   "/categories/create",
		"category": category,
	})
}

// Обработчик сохранения категории
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

	// Находим первый свободный ID
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

// Редактирование категории
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

// Обновление категории
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

// Удаление категории
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

// ====================================================================
// АВТОРИЗАЦИЯ
// ====================================================================

func loginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Проверка учетных данных
	if username == "admin" && password == "admin123" {
		// Устанавливаем сессию
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

// ====================================================================
// ГЛАВНОЕ МЕНЮ
// ====================================================================

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

// ====================================================================
// ПОСТАВЩИКИ
// ====================================================================

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

	// Получаем параметры сортировки
	sortBy := c.DefaultQuery("sort", "id")
	direction := c.DefaultQuery("direction", "asc")

	// Фильтрация по поиску
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

	// Сортируем поставщиков
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

// Функция сортировки поставщиков
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

	// Находим первый свободный ID
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

	// Обновляем supplierID если новый ID больше
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

	// Обновляем данные
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

// ====================================================================
// СТРАНИЦА МАГАЗИНОВ
// ====================================================================

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

	// Получаем параметры сортировки
	sortBy := c.DefaultQuery("sort", "id") // по умолчанию сортируем по ID
	direction := c.DefaultQuery("direction", "asc")

	// Сортируем список
	sortStores(storeList, sortBy, direction)

	c.HTML(http.StatusOK, "store.html", gin.H{
		"stores":      storeList,
		"total_count": len(storeList),
		"sort_by":     sortBy,
		"direction":   direction,
	})
}

// Функция сортировки магазинов
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

	// Находим первый свободный ID (1, 2, 3...)
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

	// Обновляем storeID если новый ID больше
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

// ====================================================================
// ТОВАРЫ
// ====================================================================

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

	// Рассчитываем среднюю цену
	var totalPrice float64
	for _, p := range productList {
		totalPrice += p.Price
	}
	avgPrice := 0.0
	if len(productList) > 0 {
		avgPrice = totalPrice / float64(len(productList))
	}

	c.HTML(http.StatusOK, "products.html", gin.H{
		"products":      productList,
		"total_count":   len(productList),
		"average_price": fmt.Sprintf("%.2f", avgPrice),
	})
}

func productCreatePage(c *gin.Context) {
	if !checkAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":  "Добавить товар",
		"action": "/products/create",
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

	if name == "" || priceStr == "" {
		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Добавить товар",
			"action": "/products/create",
			"error":  "Название и цена обязательны",
			"product": Product{
				Name:        name,
				Description: description,
			},
		})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Добавить товар",
			"action": "/products/create",
			"error":  "Неверный формат цены",
			"product": Product{
				Name:        name,
				Description: description,
			},
		})
		return
	}

	mu.Lock()

	// Находим первый свободный ID для товаров (1, 2, 3...)
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
		CreatedAt:   time.Now(),
	}

	// Обновляем productID если новый ID больше
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
	mu.RUnlock()

	if !exists {
		c.Redirect(http.StatusFound, "/products")
		return
	}

	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":   "Редактировать товар",
		"action":  fmt.Sprintf("/products/%d/edit", id),
		"product": product,
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

	mu.Lock()
	product, exists := products[id]
	if !exists {
		mu.Unlock()
		c.Redirect(http.StatusFound, "/products")
		return
	}

	product.Name = c.PostForm("name")
	product.Description = c.PostForm("description")

	priceStr := c.PostForm("price")
	if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
		product.Price = price
	}

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

// ====================================================================
// СТАТИЧЕСКИЕ СТРАНИЦЫ
// ====================================================================
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
