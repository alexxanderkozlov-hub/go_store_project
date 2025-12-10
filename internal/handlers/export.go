package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go_store_project/internal/models"
	"go_store_project/internal/storage"

	"github.com/gin-gonic/gin"
)

// Экспорт в TXT
func ExportTXT(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	dataType := c.Query("type")
	filename := "export.txt"

	stores, suppliers, products, categories, supplies := storage.GetAllData()

	// Создаем мапы для быстрого доступа по ID
	categoryMap := make(map[int]models.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	supplierMap := make(map[int]models.Supplier)
	for _, supplier := range suppliers {
		supplierMap[supplier.ID] = supplier
	}

	productMap := make(map[int]models.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	var content strings.Builder
	content.WriteString("Экспорт данных системы\n")
	content.WriteString("========================\n")
	content.WriteString(fmt.Sprintf("Дата экспорта: %s\n\n", time.Now().Format("02.01.2006 15:04")))

	switch dataType {
	case "stores":
		filename = "stores.txt"
		content.WriteString("МАГАЗИНЫ\n")
		content.WriteString("========\n")
		for _, store := range stores {
			content.WriteString(fmt.Sprintf("ID: %d\n", store.ID))
			content.WriteString(fmt.Sprintf("Название: %s\n", store.Name))
			content.WriteString(fmt.Sprintf("Адрес: %s\n", store.Address))
			if store.Logo != "" {
				content.WriteString(fmt.Sprintf("Логотип: %s\n", store.Logo))
			}
			content.WriteString(fmt.Sprintf("Дата создания: %s\n", store.CreatedAt.Format("02.01.2006")))
			content.WriteString("---\n")
		}

	case "products":
		filename = "products.txt"
		content.WriteString("ТОВАРЫ\n")
		content.WriteString("======\n")
		for _, product := range products {
			categoryName := "Без категории"
			if cat, exists := categoryMap[product.CategoryID]; exists {
				categoryName = cat.Name
			}

			content.WriteString(fmt.Sprintf("ID: %d\n", product.ID))
			content.WriteString(fmt.Sprintf("Название: %s\n", product.Name))
			content.WriteString(fmt.Sprintf("Артикул: %s\n", product.SKU))
			content.WriteString(fmt.Sprintf("Цена: %.2f руб.\n", product.Price))
			content.WriteString(fmt.Sprintf("Описание: %s\n", product.Description))
			if product.Photo != "" {
				content.WriteString(fmt.Sprintf("Фото: %s\n", product.Photo))
			}
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
			if supplier.Phone != "" {
				content.WriteString(fmt.Sprintf("Телефон: %s\n", supplier.Phone))
			}
			if supplier.Email != "" {
				content.WriteString(fmt.Sprintf("Email: %s\n", supplier.Email))
			}
			if supplier.Address != "" {
				content.WriteString(fmt.Sprintf("Адрес: %s\n", supplier.Address))
			}
			content.WriteString("---\n")
		}

	case "supplies":
		filename = "supplies.txt"
		content.WriteString("ПОСТАВКИ\n")
		content.WriteString("========\n")
		for _, supply := range supplies {
			supplierName := "Неизвестный"
			if sup, exists := supplierMap[supply.SupplierID]; exists {
				supplierName = sup.Name
			}

			productName := "Неизвестный"
			if prod, exists := productMap[supply.ProductID]; exists {
				productName = prod.Name
			}

			content.WriteString(fmt.Sprintf("ID поставки: %d\n", supply.ID))
			content.WriteString(fmt.Sprintf("Поставщик: %s\n", supplierName))
			content.WriteString(fmt.Sprintf("Товар: %s\n", productName))
			content.WriteString(fmt.Sprintf("Количество: %d\n", supply.Quantity))
			content.WriteString(fmt.Sprintf("Цена за ед.: %.2f руб.\n", supply.Price))
			content.WriteString(fmt.Sprintf("Итого: %.2f руб.\n", supply.Total))
			content.WriteString(fmt.Sprintf("Дата: %s\n", supply.Date.Format("02.01.2006")))
			content.WriteString(fmt.Sprintf("Статус: %s\n", models.GetStatusText(supply.Status)))
			if supply.Notes != "" {
				content.WriteString(fmt.Sprintf("Примечания: %s\n", supply.Notes))
			}
			content.WriteString(fmt.Sprintf("Дата создания: %s\n", supply.CreatedAt.Format("02.01.2006")))
			content.WriteString("---\n")
		}

	case "categories":
		filename = "categories.txt"
		content.WriteString("КАТЕГОРИИ\n")
		content.WriteString("=========\n")
		for _, category := range categories {
			content.WriteString(fmt.Sprintf("ID: %d\n", category.ID))
			content.WriteString(fmt.Sprintf("Название: %s\n", category.Name))
			if category.Description != "" {
				content.WriteString(fmt.Sprintf("Описание: %s\n", category.Description))
			}
			content.WriteString("---\n")
		}

	default:
		filename = "full_export.txt"
		// Экспорт всего
		content.WriteString("МАГАЗИНЫ\n")
		content.WriteString("========\n")
		for _, store := range stores {
			content.WriteString(fmt.Sprintf("ID: %d | Название: %s | Адрес: %s | Дата: %s\n",
				store.ID, store.Name, store.Address, store.CreatedAt.Format("02.01.2006")))
		}
		content.WriteString("\n")

		content.WriteString("КАТЕГОРИИ\n")
		content.WriteString("=========\n")
		for _, category := range categories {
			content.WriteString(fmt.Sprintf("ID: %d | Название: %s | Описание: %s\n",
				category.ID, category.Name, category.Description))
		}
		content.WriteString("\n")

		content.WriteString("ТОВАРЫ\n")
		content.WriteString("======\n")
		for _, product := range products {
			categoryName := "Без категории"
			if cat, exists := categoryMap[product.CategoryID]; exists {
				categoryName = cat.Name
			}
			content.WriteString(fmt.Sprintf("ID: %d | Название: %s | Цена: %.2f | Артикул: %s | Категория: %s\n",
				product.ID, product.Name, product.Price, product.SKU, categoryName))
		}
		content.WriteString("\n")

		content.WriteString("ПОСТАВЩИКИ\n")
		content.WriteString("==========\n")
		for _, supplier := range suppliers {
			content.WriteString(fmt.Sprintf("ID: %d | Название: %s | Телефон: %s | Email: %s | Адрес: %s\n",
				supplier.ID, supplier.Name, supplier.Phone, supplier.Email, supplier.Address))
		}
		content.WriteString("\n")

		content.WriteString("ПОСТАВКИ\n")
		content.WriteString("========\n")
		for _, supply := range supplies {
			supplierName := "Неизвестный"
			if sup, exists := supplierMap[supply.SupplierID]; exists {
				supplierName = sup.Name
			}

			productName := "Неизвестный"
			if prod, exists := productMap[supply.ProductID]; exists {
				productName = prod.Name
			}

			content.WriteString(fmt.Sprintf("ID: %d | Поставщик: %s | Товар: %s | Количество: %d | Цена: %.2f | Итого: %.2f | Статус: %s\n",
				supply.ID, supplierName, productName, supply.Quantity, supply.Price, supply.Total, models.GetStatusText(supply.Status)))
		}

		content.WriteString("\n=== СВОДКА ===\n")
		content.WriteString(fmt.Sprintf("Всего магазинов: %d\n", len(stores)))
		content.WriteString(fmt.Sprintf("Всего категорий: %d\n", len(categories)))
		content.WriteString(fmt.Sprintf("Всего товаров: %d\n", len(products)))
		content.WriteString(fmt.Sprintf("Всего поставщиков: %d\n", len(suppliers)))
		content.WriteString(fmt.Sprintf("Всего поставок: %d\n", len(supplies)))
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, content.String())
}

// Экспорт в CSV
func ExportCSV(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	dataType := c.Query("type")
	filename := "export.csv"

	stores, suppliers, products, categories, supplies := storage.GetAllData()

	// Создаем мапы для быстрого доступа по ID
	categoryMap := make(map[int]models.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	supplierMap := make(map[int]models.Supplier)
	for _, supplier := range suppliers {
		supplierMap[supplier.ID] = supplier
	}

	productMap := make(map[int]models.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

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

	var content bytes.Buffer
	content.Write([]byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM

	switch dataType {
	case "stores":
		content.WriteString("ID;Название;Адрес;Логотип;Дата создания\r\n")
		for _, store := range stores {
			content.WriteString(fmt.Sprintf("%d;%s;%s;%s;%s\r\n",
				store.ID,
				escapeCSV(store.Name),
				escapeCSV(store.Address),
				escapeCSV(store.Logo),
				store.CreatedAt.Format("02.01.2006"),
			))
		}

	case "products":
		content.WriteString("ID;Название;Артикул;Цена;Описание;Фото;Категория;Дата создания\r\n")
		for _, product := range products {
			categoryName := "Без категории"
			if cat, exists := categoryMap[product.CategoryID]; exists {
				categoryName = cat.Name
			}

			content.WriteString(fmt.Sprintf("%d;%s;%s;%.2f;%s;%s;%s;%s\r\n",
				product.ID,
				escapeCSV(product.Name),
				escapeCSV(product.SKU),
				product.Price,
				escapeCSV(product.Description),
				escapeCSV(product.Photo),
				escapeCSV(categoryName),
				product.CreatedAt.Format("02.01.2006"),
			))
		}

	case "suppliers":
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
		content.WriteString("ID;Поставщик;Товар;Количество;Цена за ед.;Итого;Дата;Статус;Примечания;Дата создания\r\n")
		for _, supply := range supplies {
			supplierName := "Неизвестный"
			if sup, exists := supplierMap[supply.SupplierID]; exists {
				supplierName = sup.Name
			}

			productName := "Неизвестный"
			if prod, exists := productMap[supply.ProductID]; exists {
				productName = prod.Name
			}

			content.WriteString(fmt.Sprintf("%d;%s;%s;%d;%.2f;%.2f;%s;%s;%s;%s\r\n",
				supply.ID,
				escapeCSV(supplierName),
				escapeCSV(productName),
				supply.Quantity,
				supply.Price,
				supply.Total,
				supply.Date.Format("02.01.2006"),
				escapeCSV(models.GetStatusText(supply.Status)),
				escapeCSV(supply.Notes),
				supply.CreatedAt.Format("02.01.2006"),
			))
		}

	case "categories":
		content.WriteString("ID;Название;Описание\r\n")
		for _, category := range categories {
			content.WriteString(fmt.Sprintf("%d;%s;%s\r\n",
				category.ID,
				escapeCSV(category.Name),
				escapeCSV(category.Description),
			))
		}

	default:
		// Полный экспорт
		content.WriteString("=== МАГАЗИНЫ ===\r\n")
		content.WriteString("ID;Название;Адрес;Логотип;Дата создания\r\n")
		for _, store := range stores {
			content.WriteString(fmt.Sprintf("%d;%s;%s;%s;%s\r\n",
				store.ID,
				escapeCSV(store.Name),
				escapeCSV(store.Address),
				escapeCSV(store.Logo),
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
		content.WriteString("ID;Название;Артикул;Цена;Описание;Фото;Категория;Дата создания\r\n")
		for _, product := range products {
			categoryName := "Без категории"
			if cat, exists := categoryMap[product.CategoryID]; exists {
				categoryName = cat.Name
			}

			content.WriteString(fmt.Sprintf("%d;%s;%s;%.2f;%s;%s;%s;%s\r\n",
				product.ID,
				escapeCSV(product.Name),
				escapeCSV(product.SKU),
				product.Price,
				escapeCSV(product.Description),
				escapeCSV(product.Photo),
				escapeCSV(categoryName),
				product.CreatedAt.Format("02.01.2006"),
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
		content.WriteString("ID;Поставщик;Товар;Количество;Цена за ед.;Итого;Дата;Статус;Примечания;Дата создания\r\n")
		for _, supply := range supplies {
			supplierName := "Неизвестный"
			if sup, exists := supplierMap[supply.SupplierID]; exists {
				supplierName = sup.Name
			}

			productName := "Неизвестный"
			if prod, exists := productMap[supply.ProductID]; exists {
				productName = prod.Name
			}

			content.WriteString(fmt.Sprintf("%d;%s;%s;%d;%.2f;%.2f;%s;%s;%s;%s\r\n",
				supply.ID,
				escapeCSV(supplierName),
				escapeCSV(productName),
				supply.Quantity,
				supply.Price,
				supply.Total,
				supply.Date.Format("02.01.2006"),
				escapeCSV(models.GetStatusText(supply.Status)),
				escapeCSV(supply.Notes),
				supply.CreatedAt.Format("02.01.2006"),
			))
		}

		content.WriteString("\r\n=== СВОДКА ===\r\n")
		content.WriteString("Тип;Количество\r\n")
		content.WriteString(fmt.Sprintf("Магазины;%d\r\n", len(stores)))
		content.WriteString(fmt.Sprintf("Категории;%d\r\n", len(categories)))
		content.WriteString(fmt.Sprintf("Товары;%d\r\n", len(products)))
		content.WriteString(fmt.Sprintf("Поставщики;%d\r\n", len(suppliers)))
		content.WriteString(fmt.Sprintf("Поставки;%d\r\n", len(supplies)))
		content.WriteString(fmt.Sprintf("Дата экспорта;%s\r\n", time.Now().Format("02.01.2006 15:04:05")))
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", content.Bytes())
}

// Экспорт в Excel
func ExportExcel(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	dataType := c.Query("type")
	filename := "export.xlsx"

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Redirect(http.StatusFound, fmt.Sprintf("/export/csv?type=%s", dataType))
}

// Функция для экранирования CSV значений
func escapeCSV(value string) string {
	if value == "" {
		return ""
	}
	// Заменяем разделители и переносы строк
	result := strings.ReplaceAll(value, ";", ",")
	result = strings.ReplaceAll(result, "\r\n", " ")
	result = strings.ReplaceAll(result, "\n", " ")
	result = strings.ReplaceAll(result, "\r", " ")
	result = strings.ReplaceAll(result, "\"", "'")
	return result
}
