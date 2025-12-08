package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_store_project/internal/models"
	"go_store_project/internal/storage"

	"github.com/gin-gonic/gin"
)

func SuppliesPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	supplyList := storage.GetAllSupplies()
	supplierList := storage.GetAllSuppliers()
	productList := storage.GetAllProducts()

	supplierMap := make(map[int]models.Supplier)
	for _, supplier := range supplierList {
		supplierMap[supplier.ID] = supplier
	}

	productMap := make(map[int]models.Product)
	for _, product := range productList {
		productMap[product.ID] = product
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
			"StatusText":  models.GetStatusText(supply.Status),
			"StatusClass": models.GetStatusClass(supply.Status),
			"Notes":       supply.Notes,
			"CreatedAt":   supply.CreatedAt.Format("02.01.2006 15:04"),
		})
	}

	c.HTML(http.StatusOK, "supplies.html", gin.H{
		"username":    GetUsername(c),
		"supplies":    supplyData,
		"total_count": len(supplyList),
	})
}

func SupplyCreatePage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	supplierList := storage.GetAllSuppliers()
	productList := storage.GetAllProducts()

	c.HTML(http.StatusOK, "supplies_form.html", gin.H{
		"title":     "Создать поставку",
		"action":    "/supplies/create",
		"supply":    models.Supply{Date: time.Now(), Status: "pending"},
		"suppliers": supplierList,
		"products":  productList,
		"mode":      "create",
		"now":       time.Now().Format("2006-01-02"),
	})
}

func SupplyCreateHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	supplierIDStr := c.PostForm("supplier_id")
	productIDStr := c.PostForm("product_id")
	quantityStr := c.PostForm("quantity")
	priceStr := c.PostForm("price")
	dateStr := c.PostForm("date") // <-- получаем дату из формы
	status := c.PostForm("status")
	notes := c.PostForm("notes")

	if supplierIDStr == "" || productIDStr == "" || quantityStr == "" || priceStr == "" {
		supplierList := storage.GetAllSuppliers()
		productList := storage.GetAllProducts()

		c.HTML(http.StatusOK, "supplies_form.html", gin.H{
			"title":  "Создать поставку",
			"action": "/supplies/create",
			"error":  "Заполните обязательные поля",
			"supply": models.Supply{
				Status: status,
				Notes:  notes,
			},
			"suppliers": supplierList,
			"products":  productList,
			"mode":      "create",
			"now":       time.Now().Format("2006-01-02"),
		})
		return
	}

	supplierID, err1 := strconv.Atoi(supplierIDStr)
	productID, err2 := strconv.Atoi(productIDStr)
	quantity, err3 := strconv.Atoi(quantityStr)
	price, err4 := strconv.ParseFloat(priceStr, 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || quantity <= 0 || price < 0 {
		supplierList := storage.GetAllSuppliers()
		productList := storage.GetAllProducts()

		c.HTML(http.StatusOK, "supplies_form.html", gin.H{
			"title":  "Создать поставку",
			"action": "/supplies/create",
			"error":  "Неверный формат данных",
			"supply": models.Supply{
				Status: status,
				Notes:  notes,
			},
			"suppliers": supplierList,
			"products":  productList,
			"mode":      "create",
			"now":       time.Now().Format("2006-01-02"),
		})
		return
	}

	// Парсим дату из строки
	var date time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			date = time.Now() // Если ошибка парсинга, используем текущую дату
		} else {
			date = parsedDate
		}
	} else {
		date = time.Now() // Если дата не указана, используем текущую
	}

	total := float64(quantity) * price

	supply := models.Supply{
		SupplierID: supplierID,
		ProductID:  productID,
		Quantity:   quantity,
		Price:      price,
		Total:      total,
		Date:       date, // <-- используем дату
		Status:     status,
		Notes:      notes,
		CreatedAt:  time.Now(),
	}

	storage.CreateSupply(supply)

	c.Redirect(http.StatusFound, "/supplies")
}

func SupplyEditPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/supplies")
		return
	}

	supply, exists := storage.GetSupply(id)
	if !exists {
		c.Redirect(http.StatusFound, "/supplies")
		return
	}

	supplierList := storage.GetAllSuppliers()
	productList := storage.GetAllProducts()

	c.HTML(http.StatusOK, "supplies_form.html", gin.H{
		"title":     "Редактировать поставку",
		"action":    fmt.Sprintf("/supplies/%d/edit", id),
		"supply":    supply,
		"suppliers": supplierList,
		"products":  productList,
		"mode":      "edit",
		"now":       time.Now().Format("2006-01-02"),
	})
}

func SupplyUpdateHandler(c *gin.Context) {
	if !CheckAuth(c) {
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
	dateStr := c.PostForm("date") // <-- получаем дату из формы
	status := c.PostForm("status")
	notes := c.PostForm("notes")

	if supplierIDStr == "" || productIDStr == "" || quantityStr == "" || priceStr == "" {
		supplierList := storage.GetAllSuppliers()
		productList := storage.GetAllProducts()

		c.HTML(http.StatusOK, "supplies_form.html", gin.H{
			"title":  "Редактировать поставку",
			"action": fmt.Sprintf("/supplies/%d/edit", id),
			"error":  "Заполните обязательные поля",
			"supply": models.Supply{
				ID:     id,
				Status: status,
				Notes:  notes,
			},
			"suppliers": supplierList,
			"products":  productList,
			"mode":      "edit",
			"now":       time.Now().Format("2006-01-02"),
		})
		return
	}

	supplierID, err1 := strconv.Atoi(supplierIDStr)
	productID, err2 := strconv.Atoi(productIDStr)
	quantity, err3 := strconv.Atoi(quantityStr)
	price, err4 := strconv.ParseFloat(priceStr, 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || quantity <= 0 || price < 0 {
		supplierList := storage.GetAllSuppliers()
		productList := storage.GetAllProducts()

		c.HTML(http.StatusOK, "supplies_form.html", gin.H{
			"title":  "Редактировать поставку",
			"action": fmt.Sprintf("/supplies/%d/edit", id),
			"error":  "Неверный формат данных",
			"supply": models.Supply{
				ID:     id,
				Status: status,
				Notes:  notes,
			},
			"suppliers": supplierList,
			"products":  productList,
			"mode":      "edit",
			"now":       time.Now().Format("2006-01-02"),
		})
		return
	}

	// Парсим дату из строки
	var date time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			date = time.Now() // Если ошибка парсинга, используем текущую дату
		} else {
			date = parsedDate
		}
	} else {
		date = time.Now() // Если дата не указана, используем текущую
	}

	total := float64(quantity) * price

	supply := models.Supply{
		ID:         id,
		SupplierID: supplierID,
		ProductID:  productID,
		Quantity:   quantity,
		Price:      price,
		Total:      total,
		Date:       date, // <-- используем дату
		Status:     status,
		Notes:      notes,
	}

	storage.UpdateSupply(id, supply)

	c.Redirect(http.StatusFound, "/supplies")
}

func SupplyDeleteHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/supplies")
		return
	}

	storage.DeleteSupply(id)
	c.Redirect(http.StatusFound, "/supplies")
}
