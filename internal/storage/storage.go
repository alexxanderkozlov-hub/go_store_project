package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go_store_project/config"
	"go_store_project/internal/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

var (
	storageInstance *Storage
)

// Инициализация хранилища
func Init() {
	cfg := config.Load()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	storageInstance = &Storage{db: db}

	// Инициализация таблиц
	if err := storageInstance.initTables(); err != nil {
		log.Fatalf("Failed to initialize database tables: %v", err)
	}

	log.Println("Connected to PostgreSQL database")
}

func (s *Storage) initTables() error {
	queries := []string{
		// Таблица магазинов
		`CREATE TABLE IF NOT EXISTS stores (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            address TEXT,
            logo TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,

		// Таблица категорий
		`CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            description TEXT
        )`,

		// Таблица поставщиков
		`CREATE TABLE IF NOT EXISTS suppliers (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            phone VARCHAR(20),
            email VARCHAR(100),
            address TEXT
        )`,

		// Таблица продуктов
		`CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            price DECIMAL(10,2) NOT NULL,
            description TEXT,
            photo TEXT,
            category_id INTEGER REFERENCES categories(id),
            sku VARCHAR(50),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,

		// Таблица пользователей
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL
        )`,

		// Таблица поставок
		`CREATE TABLE IF NOT EXISTS supplies (
            id SERIAL PRIMARY KEY,
            supplier_id INTEGER REFERENCES suppliers(id),
            product_id INTEGER REFERENCES products(id),
            quantity INTEGER NOT NULL,
            price DECIMAL(10,2) NOT NULL,
            total DECIMAL(10,2) GENERATED ALWAYS AS (quantity * price) STORED,
            date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            status VARCHAR(20) DEFAULT 'pending',
            notes TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
	}

	for _, query := range queries {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Создаем начального пользователя, если таблица пуста
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err == nil && count == 0 {
		// Пароль: admin123
		_, err = s.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)",
			"admin", "admin123")
		if err != nil {
			log.Printf("Warning: failed to create default user: %v", err)
		}
	}

	return nil
}

// Store operations
func GetAllStores() []models.Store {
	if storageInstance == nil {
		return []models.Store{}
	}

	rows, err := storageInstance.db.Query(`
        SELECT id, name, address, logo, created_at 
        FROM stores 
        ORDER BY created_at DESC
    `)
	if err != nil {
		log.Printf("Error getting stores: %v", err)
		return []models.Store{}
	}
	defer rows.Close()

	var stores []models.Store
	for rows.Next() {
		var store models.Store
		if err := rows.Scan(&store.ID, &store.Name, &store.Address,
			&store.Logo, &store.CreatedAt); err != nil {
			log.Printf("Error scanning store: %v", err)
			continue
		}
		stores = append(stores, store)
	}

	return stores
}

func GetStore(id int) (models.Store, bool) {
	if storageInstance == nil {
		return models.Store{}, false
	}

	var store models.Store
	err := storageInstance.db.QueryRow(`
        SELECT id, name, address, logo, created_at 
        FROM stores WHERE id = $1
    `, id).Scan(&store.ID, &store.Name, &store.Address, &store.Logo, &store.CreatedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error getting store: %v", err)
		}
		return models.Store{}, false
	}

	return store, true
}

func CreateStore(store models.Store) int {
	if storageInstance == nil {
		return 0
	}

	var id int
	err := storageInstance.db.QueryRow(`
        INSERT INTO stores (name, address, logo, created_at) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id
    `, store.Name, store.Address, store.Logo, time.Now()).Scan(&id)

	if err != nil {
		log.Printf("Error creating store: %v", err)
		return 0
	}

	return id
}

func UpdateStore(id int, store models.Store) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec(`
        UPDATE stores 
        SET name = $1, address = $2, logo = $3 
        WHERE id = $4
    `, store.Name, store.Address, store.Logo, id)

	if err != nil {
		log.Printf("Error updating store: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

func DeleteStore(id int) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec("DELETE FROM stores WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting store: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

// Category operations
func GetAllCategories() []models.Category {
	if storageInstance == nil {
		return []models.Category{}
	}

	rows, err := storageInstance.db.Query(`
        SELECT id, name, description 
        FROM categories 
        ORDER BY name
    `)
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		return []models.Category{}
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}
		categories = append(categories, category)
	}

	return categories
}

func GetCategory(id int) (models.Category, bool) {
	if storageInstance == nil {
		return models.Category{}, false
	}

	var category models.Category
	err := storageInstance.db.QueryRow(`
        SELECT id, name, description 
        FROM categories WHERE id = $1
    `, id).Scan(&category.ID, &category.Name, &category.Description)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error getting category: %v", err)
		}
		return models.Category{}, false
	}

	return category, true
}

func CreateCategory(category models.Category) int {
	if storageInstance == nil {
		return 0
	}

	var id int
	err := storageInstance.db.QueryRow(`
        INSERT INTO categories (name, description) 
        VALUES ($1, $2) 
        RETURNING id
    `, category.Name, category.Description).Scan(&id)

	if err != nil {
		log.Printf("Error creating category: %v", err)
		return 0
	}

	return id
}

func UpdateCategory(id int, category models.Category) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec(`
        UPDATE categories 
        SET name = $1, description = $2 
        WHERE id = $3
    `, category.Name, category.Description, id)

	if err != nil {
		log.Printf("Error updating category: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

func DeleteCategory(id int) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

// Product operations
func GetAllProducts() []models.Product {
	if storageInstance == nil {
		return []models.Product{}
	}

	rows, err := storageInstance.db.Query(`
        SELECT id, name, price, description, photo, category_id, sku, created_at 
        FROM products 
        ORDER BY created_at DESC
    `)
	if err != nil {
		log.Printf("Error getting products: %v", err)
		return []models.Product{}
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price,
			&product.Description, &product.Photo, &product.CategoryID,
			&product.SKU, &product.CreatedAt); err != nil {
			log.Printf("Error scanning product: %v", err)
			continue
		}
		products = append(products, product)
	}

	return products
}

func GetProduct(id int) (models.Product, bool) {
	if storageInstance == nil {
		return models.Product{}, false
	}

	var product models.Product
	err := storageInstance.db.QueryRow(`
        SELECT id, name, price, description, photo, category_id, sku, created_at 
        FROM products WHERE id = $1
    `, id).Scan(&product.ID, &product.Name, &product.Price, &product.Description,
		&product.Photo, &product.CategoryID, &product.SKU, &product.CreatedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error getting product: %v", err)
		}
		return models.Product{}, false
	}

	return product, true
}

func CreateProduct(product models.Product) int {
	if storageInstance == nil {
		return 0
	}

	var id int
	err := storageInstance.db.QueryRow(`
        INSERT INTO products (name, price, description, photo, category_id, sku, created_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id
    `, product.Name, product.Price, product.Description, product.Photo,
		product.CategoryID, product.SKU, time.Now()).Scan(&id)

	if err != nil {
		log.Printf("Error creating product: %v", err)
		return 0
	}

	return id
}

func UpdateProduct(id int, product models.Product) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec(`
        UPDATE products 
        SET name = $1, price = $2, description = $3, 
            photo = $4, category_id = $5, sku = $6 
        WHERE id = $7
    `, product.Name, product.Price, product.Description,
		product.Photo, product.CategoryID, product.SKU, id)

	if err != nil {
		log.Printf("Error updating product: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

func DeleteProduct(id int) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

// Supplier operations
func GetAllSuppliers() []models.Supplier {
	if storageInstance == nil {
		return []models.Supplier{}
	}

	rows, err := storageInstance.db.Query(`
        SELECT id, name, phone, email, address 
        FROM suppliers 
        ORDER BY name
    `)
	if err != nil {
		log.Printf("Error getting suppliers: %v", err)
		return []models.Supplier{}
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var supplier models.Supplier
		if err := rows.Scan(&supplier.ID, &supplier.Name, &supplier.Phone,
			&supplier.Email, &supplier.Address); err != nil {
			log.Printf("Error scanning supplier: %v", err)
			continue
		}
		suppliers = append(suppliers, supplier)
	}

	return suppliers
}

func GetSupplier(id int) (models.Supplier, bool) {
	if storageInstance == nil {
		return models.Supplier{}, false
	}

	var supplier models.Supplier
	err := storageInstance.db.QueryRow(`
        SELECT id, name, phone, email, address 
        FROM suppliers WHERE id = $1
    `, id).Scan(&supplier.ID, &supplier.Name, &supplier.Phone,
		&supplier.Email, &supplier.Address)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error getting supplier: %v", err)
		}
		return models.Supplier{}, false
	}

	return supplier, true
}

func CreateSupplier(supplier models.Supplier) int {
	if storageInstance == nil {
		return 0
	}

	var id int
	err := storageInstance.db.QueryRow(`
        INSERT INTO suppliers (name, phone, email, address) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id
    `, supplier.Name, supplier.Phone, supplier.Email, supplier.Address).Scan(&id)

	if err != nil {
		log.Printf("Error creating supplier: %v", err)
		return 0
	}

	return id
}

func UpdateSupplier(id int, supplier models.Supplier) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec(`
        UPDATE suppliers 
        SET name = $1, phone = $2, email = $3, address = $4 
        WHERE id = $5
    `, supplier.Name, supplier.Phone, supplier.Email, supplier.Address, id)

	if err != nil {
		log.Printf("Error updating supplier: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

func DeleteSupplier(id int) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec("DELETE FROM suppliers WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting supplier: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

// Supply operations
func GetAllSupplies() []models.Supply {
	if storageInstance == nil {
		return []models.Supply{}
	}

	rows, err := storageInstance.db.Query(`
        SELECT id, supplier_id, product_id, quantity, price, 
               total, date, status, notes, created_at 
        FROM supplies 
        ORDER BY created_at DESC
    `)
	if err != nil {
		log.Printf("Error getting supplies: %v", err)
		return []models.Supply{}
	}
	defer rows.Close()

	var supplies []models.Supply
	for rows.Next() {
		var supply models.Supply
		if err := rows.Scan(&supply.ID, &supply.SupplierID, &supply.ProductID,
			&supply.Quantity, &supply.Price, &supply.Total, &supply.Date,
			&supply.Status, &supply.Notes, &supply.CreatedAt); err != nil {
			log.Printf("Error scanning supply: %v", err)
			continue
		}
		supplies = append(supplies, supply)
	}

	return supplies
}

func GetSupply(id int) (models.Supply, bool) {
	if storageInstance == nil {
		return models.Supply{}, false
	}

	var supply models.Supply
	err := storageInstance.db.QueryRow(`
        SELECT id, supplier_id, product_id, quantity, price, 
               total, date, status, notes, created_at 
        FROM supplies WHERE id = $1
    `, id).Scan(&supply.ID, &supply.SupplierID, &supply.ProductID,
		&supply.Quantity, &supply.Price, &supply.Total, &supply.Date,
		&supply.Status, &supply.Notes, &supply.CreatedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error getting supply: %v", err)
		}
		return models.Supply{}, false
	}

	return supply, true
}

func CreateSupply(supply models.Supply) int {
	if storageInstance == nil {
		return 0
	}

	var id int
	err := storageInstance.db.QueryRow(`
        INSERT INTO supplies (supplier_id, product_id, quantity, price, 
                              date, status, notes, created_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
        RETURNING id
    `, supply.SupplierID, supply.ProductID, supply.Quantity, supply.Price,
		supply.Date, supply.Status, supply.Notes, time.Now()).Scan(&id)

	if err != nil {
		log.Printf("Error creating supply: %v", err)
		return 0
	}

	return id
}

func UpdateSupply(id int, supply models.Supply) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec(`
        UPDATE supplies 
        SET supplier_id = $1, product_id = $2, quantity = $3, 
            price = $4, date = $5, status = $6, notes = $7 
        WHERE id = $8
    `, supply.SupplierID, supply.ProductID, supply.Quantity,
		supply.Price, supply.Date, supply.Status, supply.Notes, id)

	if err != nil {
		log.Printf("Error updating supply: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

func DeleteSupply(id int) bool {
	if storageInstance == nil {
		return false
	}

	result, err := storageInstance.db.Exec("DELETE FROM supplies WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting supply: %v", err)
		return false
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

// User operations
func CreateUser(user models.User) int {
	if storageInstance == nil {
		return 0
	}

	var id int
	err := storageInstance.db.QueryRow(`
        INSERT INTO users (username, password, role) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `, user.Username, user.Password, user.Role).Scan(&id)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return 0
	}

	return id
}

func GetUserByUsername(username string) (models.User, bool) {
	if storageInstance == nil {
		return models.User{}, false
	}

	var user models.User
	err := storageInstance.db.QueryRow(`
        SELECT id, username, password, role
        FROM users WHERE username = $1
    `, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error getting user: %v", err)
		}
		return models.User{}, false
	}

	return user, true
}

// Get all data for export
func GetAllData() ([]models.Store, []models.Supplier, []models.Product, []models.Category, []models.Supply) {
	return GetAllStores(), GetAllSuppliers(), GetAllProducts(), GetAllCategories(), GetAllSupplies()
}
