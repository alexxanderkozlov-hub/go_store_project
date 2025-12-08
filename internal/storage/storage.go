package storage

import (
	"log"
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
	log.Println("STORAGE: Initializing storage...")
	// Можно добавить тестовые данные здесь
}

// Store operations
func GetAllStores() []models.Store {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting all stores, count: %d", len(stores))

	storeList := make([]models.Store, 0, len(stores))
	for id, store := range stores {
		log.Printf("STORAGE: Store %d: ID=%d, Name=%s, Address=%s",
			id, store.ID, store.Name, store.Address)
		storeList = append(storeList, store)
	}

	return storeList
}

func GetStore(id int) (models.Store, bool) {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting store ID=%d", id)

	store, exists := stores[id]
	if exists {
		log.Printf("STORAGE: Found store: %+v", store)
	} else {
		log.Printf("STORAGE: Store not found: ID=%d", id)
	}

	return store, exists
}

func CreateStore(store models.Store) int {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Creating store - Name: %s, Address: %s",
		store.Name, store.Address)
	log.Printf("STORAGE: Current stores count before: %d", len(stores))

	// Показать все текущие магазины
	for id, s := range stores {
		log.Printf("STORAGE: Existing store %d: %s", id, s.Name)
	}

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

	log.Printf("STORAGE: Store created with ID: %d", newID)
	log.Printf("STORAGE: Current stores count after: %d", len(stores))

	return newID
}

func UpdateStore(id int, store models.Store) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Updating store ID=%d", id)

	if _, exists := stores[id]; !exists {
		log.Printf("STORAGE ERROR: Store not found for update: ID=%d", id)
		return false
	}

	store.ID = id
	stores[id] = store

	log.Printf("STORAGE: Store updated successfully: ID=%d", id)

	return true
}

func DeleteStore(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Deleting store ID=%d", id)
	log.Printf("STORAGE: Current stores count before delete: %d", len(stores))

	if _, exists := stores[id]; !exists {
		log.Printf("STORAGE ERROR: Store not found for delete: ID=%d", id)
		return false
	}

	delete(stores, id)

	log.Printf("STORAGE: Store deleted: ID=%d", id)
	log.Printf("STORAGE: Current stores count after delete: %d", len(stores))

	return true
}

// Category operations
func GetAllCategories() []models.Category {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting all categories, count: %d", len(categories))

	categoryList := make([]models.Category, 0, len(categories))
	for id, cat := range categories {
		log.Printf("STORAGE: Category %d: ID=%d, Name=%s",
			id, cat.ID, cat.Name)
		categoryList = append(categoryList, cat)
	}

	return categoryList
}

func GetCategory(id int) (models.Category, bool) {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting category ID=%d", id)

	category, exists := categories[id]
	if exists {
		log.Printf("STORAGE: Found category: %+v", category)
	} else {
		log.Printf("STORAGE: Category not found: ID=%d", id)
	}

	return category, exists
}

func CreateCategory(category models.Category) int {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Creating category - Name: %s, Description: %s",
		category.Name, category.Description)
	log.Printf("STORAGE: Current categories count before: %d", len(categories))

	// Показать все текущие категории
	for id, cat := range categories {
		log.Printf("STORAGE: Existing category %d: %s", id, cat.Name)
	}

	newID := 1
	for {
		if _, exists := categories[newID]; !exists {
			break
		}
		newID++
	}

	category.ID = newID
	categories[newID] = category

	log.Printf("STORAGE: Category created with ID: %d", newID)
	log.Printf("STORAGE: Current categories count after: %d", len(categories))

	return newID
}

func UpdateCategory(id int, category models.Category) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Updating category ID=%d", id)

	if _, exists := categories[id]; !exists {
		log.Printf("STORAGE ERROR: Category not found for update: ID=%d", id)
		return false
	}

	category.ID = id
	categories[id] = category

	log.Printf("STORAGE: Category updated successfully: ID=%d", id)

	return true
}

func DeleteCategory(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Deleting category ID=%d", id)
	log.Printf("STORAGE: Current categories count before delete: %d", len(categories))

	if _, exists := categories[id]; !exists {
		log.Printf("STORAGE ERROR: Category not found for delete: ID=%d", id)
		return false
	}

	delete(categories, id)

	log.Printf("STORAGE: Category deleted: ID=%d", id)
	log.Printf("STORAGE: Current categories count after delete: %d", len(categories))

	return true
}

// Product operations
func GetAllProducts() []models.Product {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting all products, count: %d", len(products))

	productList := make([]models.Product, 0, len(products))
	for id, prod := range products {
		log.Printf("STORAGE: Product %d: ID=%d, Name=%s, Price=%.2f, SKU=%s",
			id, prod.ID, prod.Name, prod.Price, prod.SKU)
		productList = append(productList, prod)
	}

	return productList
}

func GetProduct(id int) (models.Product, bool) {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting product ID=%d", id)

	product, exists := products[id]
	if exists {
		log.Printf("STORAGE: Found product: %+v", product)
	} else {
		log.Printf("STORAGE: Product not found: ID=%d", id)
	}

	return product, exists
}

func CreateProduct(product models.Product) int {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Creating product - Name: %s, Price: %.2f, SKU: %s",
		product.Name, product.Price, product.SKU)
	log.Printf("STORAGE: Current products count before: %d", len(products))

	// Показать все текущие товары
	for id, prod := range products {
		log.Printf("STORAGE: Existing product %d: %s (%.2f)", id, prod.Name, prod.Price)
	}

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

	log.Printf("STORAGE: Product created with ID: %d", newID)
	log.Printf("STORAGE: Current products count after: %d", len(products))

	return newID
}

func UpdateProduct(id int, product models.Product) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Updating product ID=%d", id)

	if _, exists := products[id]; !exists {
		log.Printf("STORAGE ERROR: Product not found for update: ID=%d", id)
		return false
	}

	product.ID = id
	products[id] = product

	log.Printf("STORAGE: Product updated successfully: ID=%d", id)

	return true
}

func DeleteProduct(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Deleting product ID=%d", id)
	log.Printf("STORAGE: Current products count before delete: %d", len(products))

	if _, exists := products[id]; !exists {
		log.Printf("STORAGE ERROR: Product not found for delete: ID=%d", id)
		return false
	}

	delete(products, id)

	log.Printf("STORAGE: Product deleted: ID=%d", id)
	log.Printf("STORAGE: Current products count after delete: %d", len(products))

	return true
}

// Supplier operations
func GetAllSuppliers() []models.Supplier {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting all suppliers, count: %d", len(suppliers))

	supplierList := make([]models.Supplier, 0, len(suppliers))
	for id, sup := range suppliers {
		log.Printf("STORAGE: Supplier %d: ID=%d, Name=%s, Phone=%s",
			id, sup.ID, sup.Name, sup.Phone)
		supplierList = append(supplierList, sup)
	}

	return supplierList
}

func GetSupplier(id int) (models.Supplier, bool) {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting supplier ID=%d", id)

	supplier, exists := suppliers[id]
	if exists {
		log.Printf("STORAGE: Found supplier: %+v", supplier)
	} else {
		log.Printf("STORAGE: Supplier not found: ID=%d", id)
	}

	return supplier, exists
}

func CreateSupplier(supplier models.Supplier) int {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Creating supplier - Name: %s, Phone: %s, Email: %s",
		supplier.Name, supplier.Phone, supplier.Email)
	log.Printf("STORAGE: Current suppliers count before: %d", len(suppliers))

	// Показать все текущие поставщики
	for id, sup := range suppliers {
		log.Printf("STORAGE: Existing supplier %d: %s", id, sup.Name)
	}

	newID := 1
	for {
		if _, exists := suppliers[newID]; !exists {
			break
		}
		newID++
	}

	supplier.ID = newID
	suppliers[newID] = supplier

	log.Printf("STORAGE: Supplier created with ID: %d", newID)
	log.Printf("STORAGE: Current suppliers count after: %d", len(suppliers))

	return newID
}

func UpdateSupplier(id int, supplier models.Supplier) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Updating supplier ID=%d", id)

	if _, exists := suppliers[id]; !exists {
		log.Printf("STORAGE ERROR: Supplier not found for update: ID=%d", id)
		return false
	}

	supplier.ID = id
	suppliers[id] = supplier

	log.Printf("STORAGE: Supplier updated successfully: ID=%d", id)

	return true
}

func DeleteSupplier(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Deleting supplier ID=%d", id)
	log.Printf("STORAGE: Current suppliers count before delete: %d", len(suppliers))

	if _, exists := suppliers[id]; !exists {
		log.Printf("STORAGE ERROR: Supplier not found for delete: ID=%d", id)
		return false
	}

	delete(suppliers, id)

	log.Printf("STORAGE: Supplier deleted: ID=%d", id)
	log.Printf("STORAGE: Current suppliers count after delete: %d", len(suppliers))

	return true
}

// Supply operations
func GetAllSupplies() []models.Supply {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting all supplies, count: %d", len(supplies))

	supplyList := make([]models.Supply, 0, len(supplies))
	for id, sup := range supplies {
		log.Printf("STORAGE: Supply %d: ID=%d, Quantity=%d, Total=%.2f, Status=%s",
			id, sup.ID, sup.Quantity, sup.Total, sup.Status)
		supplyList = append(supplyList, sup)
	}

	return supplyList
}

func GetSupply(id int) (models.Supply, bool) {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Getting supply ID=%d", id)

	supply, exists := supplies[id]
	if exists {
		log.Printf("STORAGE: Found supply: %+v", supply)
	} else {
		log.Printf("STORAGE: Supply not found: ID=%d", id)
	}

	return supply, exists
}

func CreateSupply(supply models.Supply) int {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Creating supply - SupplierID: %d, ProductID: %d, Quantity: %d, Total: %.2f",
		supply.SupplierID, supply.ProductID, supply.Quantity, supply.Total)
	log.Printf("STORAGE: Current supplies count before: %d", len(supplies))

	// Показать все текущие поставки
	for id, sup := range supplies {
		log.Printf("STORAGE: Existing supply %d: SupplierID=%d, ProductID=%d",
			id, sup.SupplierID, sup.ProductID)
	}

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

	log.Printf("STORAGE: Supply created with ID: %d", newID)
	log.Printf("STORAGE: Current supplies count after: %d", len(supplies))

	return newID
}

func UpdateSupply(id int, supply models.Supply) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Updating supply ID=%d", id)

	if _, exists := supplies[id]; !exists {
		log.Printf("STORAGE ERROR: Supply not found for update: ID=%d", id)
		return false
	}

	supply.ID = id
	supplies[id] = supply

	log.Printf("STORAGE: Supply updated successfully: ID=%d", id)

	return true
}

func DeleteSupply(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("STORAGE: Deleting supply ID=%d", id)
	log.Printf("STORAGE: Current supplies count before delete: %d", len(supplies))

	if _, exists := supplies[id]; !exists {
		log.Printf("STORAGE ERROR: Supply not found for delete: ID=%d", id)
		return false
	}

	delete(supplies, id)

	log.Printf("STORAGE: Supply deleted: ID=%d", id)
	log.Printf("STORAGE: Current supplies count after delete: %d", len(supplies))

	return true
}

// Get all data for export
func GetAllData() (map[int]models.Store, map[int]models.Supplier, map[int]models.Product, map[int]models.Category, map[int]models.Supply) {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("STORAGE: Exporting all data - Stores: %d, Suppliers: %d, Products: %d, Categories: %d, Supplies: %d",
		len(stores), len(suppliers), len(products), len(categories), len(supplies))

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
