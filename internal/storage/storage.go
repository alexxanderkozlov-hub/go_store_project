package storage

import (
	"sync"
	"time"

	"go_store_project/internal/models"
)

var (
	stores     = make(map[int]models.Store)
	suppliers  = make(map[int]models.Supplier)
	products   = make(map[int]models.Product)
	categories = make(map[int]models.Category)
	supplies   = make(map[int]models.Supply)
	users      = make(map[int]models.User)
	mu         sync.RWMutex
)

func Init() {
	mu.Lock()
	defer mu.Unlock()
	// Хранилище инициализируется пустым
}

// Store operations
func GetAllStores() []models.Store {
	mu.RLock()
	defer mu.RUnlock()

	storeList := make([]models.Store, 0, len(stores))
	for _, store := range stores {
		storeList = append(storeList, store)
	}

	return storeList
}

func GetStore(id int) (models.Store, bool) {
	mu.RLock()
	defer mu.RUnlock()

	store, exists := stores[id]
	return store, exists
}

func CreateStore(store models.Store) int {
	mu.Lock()
	defer mu.Unlock()

	newID := 1
	for {
		if _, exists := stores[newID]; !exists {
			break
		}
		newID++
	}

	store.ID = newID
	store.CreatedAt = time.Now()
	stores[newID] = store

	return newID
}

func UpdateStore(id int, store models.Store) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := stores[id]; !exists {
		return false
	}

	store.ID = id
	stores[id] = store
	return true
}

func DeleteStore(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := stores[id]; !exists {
		return false
	}

	delete(stores, id)
	return true
}

// Category operations
func GetAllCategories() []models.Category {
	mu.RLock()
	defer mu.RUnlock()

	categoryList := make([]models.Category, 0, len(categories))
	for _, cat := range categories {
		categoryList = append(categoryList, cat)
	}

	return categoryList
}

func GetCategory(id int) (models.Category, bool) {
	mu.RLock()
	defer mu.RUnlock()

	category, exists := categories[id]
	return category, exists
}

func CreateCategory(category models.Category) int {
	mu.Lock()
	defer mu.Unlock()

	newID := 1
	for {
		if _, exists := categories[newID]; !exists {
			break
		}
		newID++
	}

	category.ID = newID
	categories[newID] = category

	return newID
}

func UpdateCategory(id int, category models.Category) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := categories[id]; !exists {
		return false
	}

	category.ID = id
	categories[id] = category
	return true
}

func DeleteCategory(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := categories[id]; !exists {
		return false
	}

	delete(categories, id)
	return true
}

// Product operations
func GetAllProducts() []models.Product {
	mu.RLock()
	defer mu.RUnlock()

	productList := make([]models.Product, 0, len(products))
	for _, prod := range products {
		productList = append(productList, prod)
	}

	return productList
}

func GetProduct(id int) (models.Product, bool) {
	mu.RLock()
	defer mu.RUnlock()

	product, exists := products[id]
	return product, exists
}

func CreateProduct(product models.Product) int {
	mu.Lock()
	defer mu.Unlock()

	newID := 1
	for {
		if _, exists := products[newID]; !exists {
			break
		}
		newID++
	}

	product.ID = newID
	product.CreatedAt = time.Now()
	products[newID] = product

	return newID
}

func UpdateProduct(id int, product models.Product) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := products[id]; !exists {
		return false
	}

	product.ID = id
	products[id] = product
	return true
}

func DeleteProduct(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := products[id]; !exists {
		return false
	}

	delete(products, id)
	return true
}

// Supplier operations
func GetAllSuppliers() []models.Supplier {
	mu.RLock()
	defer mu.RUnlock()

	supplierList := make([]models.Supplier, 0, len(suppliers))
	for _, sup := range suppliers {
		supplierList = append(supplierList, sup)
	}

	return supplierList
}

func GetSupplier(id int) (models.Supplier, bool) {
	mu.RLock()
	defer mu.RUnlock()

	supplier, exists := suppliers[id]
	return supplier, exists
}

func CreateSupplier(supplier models.Supplier) int {
	mu.Lock()
	defer mu.Unlock()

	newID := 1
	for {
		if _, exists := suppliers[newID]; !exists {
			break
		}
		newID++
	}

	supplier.ID = newID
	suppliers[newID] = supplier

	return newID
}

func UpdateSupplier(id int, supplier models.Supplier) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := suppliers[id]; !exists {
		return false
	}

	supplier.ID = id
	suppliers[id] = supplier
	return true
}

func DeleteSupplier(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := suppliers[id]; !exists {
		return false
	}

	delete(suppliers, id)
	return true
}

// Supply operations
func GetAllSupplies() []models.Supply {
	mu.RLock()
	defer mu.RUnlock()

	supplyList := make([]models.Supply, 0, len(supplies))
	for _, sup := range supplies {
		supplyList = append(supplyList, sup)
	}

	return supplyList
}

func GetSupply(id int) (models.Supply, bool) {
	mu.RLock()
	defer mu.RUnlock()

	supply, exists := supplies[id]
	return supply, exists
}

func CreateSupply(supply models.Supply) int {
	mu.Lock()
	defer mu.Unlock()

	newID := 1
	for {
		if _, exists := supplies[newID]; !exists {
			break
		}
		newID++
	}

	supply.ID = newID
	supply.CreatedAt = time.Now()
	supplies[newID] = supply

	return newID
}

func UpdateSupply(id int, supply models.Supply) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := supplies[id]; !exists {
		return false
	}

	supply.ID = id
	supplies[id] = supply
	return true
}

func DeleteSupply(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := supplies[id]; !exists {
		return false
	}

	delete(supplies, id)
	return true
}

// User operations
func CreateUser(user models.User) int {
	mu.Lock()
	defer mu.Unlock()

	newID := 1
	for {
		if _, exists := users[newID]; !exists {
			break
		}
		newID++
	}

	user.ID = newID
	users[newID] = user

	return newID
}

func GetUserByUsername(username string) (models.User, bool) {
	mu.RLock()
	defer mu.RUnlock()

	for _, user := range users {
		if user.Username == username {
			return user, true
		}
	}
	return models.User{}, false
}

// Get all data for export
func GetAllData() (map[int]models.Store, map[int]models.Supplier, map[int]models.Product, map[int]models.Category, map[int]models.Supply) {
	mu.RLock()
	defer mu.RUnlock()

	storesCopy := make(map[int]models.Store)
	suppliersCopy := make(map[int]models.Supplier)
	productsCopy := make(map[int]models.Product)
	categoriesCopy := make(map[int]models.Category)
	suppliesCopy := make(map[int]models.Supply)

	for k, v := range stores {
		storesCopy[k] = v
	}
	for k, v := range suppliers {
		suppliersCopy[k] = v
	}
	for k, v := range products {
		productsCopy[k] = v
	}
	for k, v := range categories {
		categoriesCopy[k] = v
	}
	for k, v := range supplies {
		suppliesCopy[k] = v
	}

	return storesCopy, suppliersCopy, productsCopy, categoriesCopy, suppliesCopy
}
