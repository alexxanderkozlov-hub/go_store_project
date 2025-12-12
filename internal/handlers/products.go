package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort" // Добавили импорт sort
	"strconv"

	"go_store_project/internal/models"
	"go_store_project/internal/storage"

	"github.com/gin-gonic/gin"
)

func ProductsList(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	productList := storage.GetAllProducts()

	// Добавили сортировку по возрастанию ID
	sort.Slice(productList, func(i, j int) bool {
		return productList[i].ID < productList[j].ID
	})

	log.Printf("DEBUG: Displaying %d products", len(productList))

	c.HTML(http.StatusOK, "products.html", gin.H{
		"username":    GetUsername(c),
		"products":    productList,
		"total_count": len(productList),
	})
}

func ProductCreatePage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	categoryList := storage.GetAllCategories()

	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":      "Добавить товар",
		"action":     "/products/create",
		"product":    models.Product{},
		"categories": categoryList,
		"mode":       "create",
	})
}

func ProductCreateHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	description := c.PostForm("description")
	sku := c.PostForm("sku")
	categoryIDStr := c.PostForm("category_id")

	log.Printf("DEBUG: Creating product - Name: %s, Price: %s, Description: %s, SKU: %s, CategoryID: %s",
		name, priceStr, description, sku, categoryIDStr)

	if name == "" || priceStr == "" {
		log.Printf("DEBUG: Validation failed - name or price is empty")
		categoryList := storage.GetAllCategories()

		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Добавить товар",
			"action": "/products/create",
			"error":  "Название и цена обязательны",
			"product": models.Product{
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
		log.Printf("DEBUG: Invalid price format: %s", priceStr)
		categoryList := storage.GetAllCategories()

		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Добавить товар",
			"action": "/products/create",
			"error":  "Неверный формат цены",
			"product": models.Product{
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

	product := models.Product{
		Name:        name,
		Price:       price,
		Description: description,
		SKU:         sku,
		CategoryID:  categoryID,
	}

	log.Printf("DEBUG: Product object before save: %+v", product)

	// Исправлено: используем возвращаемое значение
	productID := storage.CreateProduct(product)
	log.Printf("DEBUG: Product created with ID: %d", productID)

	// Проверяем что сохранилось
	products := storage.GetAllProducts()
	log.Printf("DEBUG: Total products in storage after create: %d", len(products))
	for i, prod := range products {
		log.Printf("DEBUG: Product %d: ID=%d, Name=%s, Price=%.2f, SKU=%s",
			i, prod.ID, prod.Name, prod.Price, prod.SKU)
	}

	c.Redirect(http.StatusFound, "/products")
}

func ProductEditPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("ERROR: Invalid product ID: %s", c.Param("id"))
		c.Redirect(http.StatusFound, "/products")
		return
	}

	log.Printf("DEBUG: Loading product for edit: ID=%d", id)

	product, exists := storage.GetProduct(id)
	if !exists {
		log.Printf("ERROR: Product not found: ID=%d", id)
		c.Redirect(http.StatusFound, "/products")
		return
	}

	log.Printf("DEBUG: Found product: %+v", product)

	categoryList := storage.GetAllCategories()

	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":      "Редактировать товар",
		"action":     fmt.Sprintf("/products/%d/edit", id),
		"product":    product,
		"categories": categoryList,
		"mode":       "edit",
	})
}

func ProductUpdateHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("ERROR: Invalid product ID: %s", c.Param("id"))
		c.Redirect(http.StatusFound, "/products")
		return
	}

	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	description := c.PostForm("description")
	sku := c.PostForm("sku")
	categoryIDStr := c.PostForm("category_id")

	log.Printf("DEBUG: Updating product ID=%d - Name: %s, Price: %s, Description: %s, SKU: %s, CategoryID: %s",
		id, name, priceStr, description, sku, categoryIDStr)

	if name == "" || priceStr == "" {
		log.Printf("DEBUG: Validation failed for update - name or price is empty")
		categoryList := storage.GetAllCategories()

		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Редактировать товар",
			"action": fmt.Sprintf("/products/%d/edit", id),
			"error":  "Название и цена обязательны",
			"product": models.Product{
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
		log.Printf("DEBUG: Invalid price format for update: %s", priceStr)
		categoryList := storage.GetAllCategories()

		c.HTML(http.StatusOK, "product_form.html", gin.H{
			"title":  "Редактировать товар",
			"action": fmt.Sprintf("/products/%d/edit", id),
			"error":  "Неверный формат цены",
			"product": models.Product{
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
		if catID, err := strconv.Atoi(categoryIDStr); err == nil {
			categoryID = catID
		}
	}

	product := models.Product{
		ID:          id,
		Name:        name,
		Price:       price,
		Description: description,
		SKU:         sku,
		CategoryID:  categoryID,
	}

	log.Printf("DEBUG: Updating product: %+v", product)

	success := storage.UpdateProduct(id, product)
	if !success {
		log.Printf("ERROR: Failed to update product ID=%d", id)
		categoryList := storage.GetAllCategories()

		c.HTML(http.StatusInternalServerError, "product_form.html", gin.H{
			"title":      "Редактировать товар",
			"action":     fmt.Sprintf("/products/%d/edit", id),
			"error":      "Ошибка при обновлении товара",
			"product":    product,
			"categories": categoryList,
			"mode":       "edit",
		})
		return
	}

	log.Printf("DEBUG: Product updated successfully: ID=%d", id)

	c.Redirect(http.StatusFound, "/products")
}

func ProductDeleteHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("ERROR: Invalid product ID: %s", c.Param("id"))
		c.Redirect(http.StatusFound, "/products")
		return
	}

	log.Printf("DEBUG: Deleting product: ID=%d", id)

	success := storage.DeleteProduct(id)
	if !success {
		log.Printf("ERROR: Failed to delete product ID=%d", id)
	}

	log.Printf("DEBUG: Product delete attempt completed: ID=%d, Success: %v", id, success)

	c.Redirect(http.StatusFound, "/products")
}
