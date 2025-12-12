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

func CategoriesListPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	categoryList := storage.GetAllCategories()

	// Добавили сортировку по возрастанию ID
	sort.Slice(categoryList, func(i, j int) bool {
		return categoryList[i].ID < categoryList[j].ID
	})

	log.Printf("DEBUG: Displaying %d categories", len(categoryList))

	c.HTML(http.StatusOK, "categories.html", gin.H{
		"username":    GetUsername(c),
		"categories":  categoryList,
		"total_count": len(categoryList),
	})
}

func CategoryCreatePage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	category := models.Category{}

	c.HTML(http.StatusOK, "category_form.html", gin.H{
		"title":    "Создать категорию",
		"action":   "/categories/create",
		"category": category,
	})
}

func CategoryCreateHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")

	log.Printf("DEBUG: Received form data - Name: %s, Description: %s", name, description)

	if name == "" {
		log.Printf("DEBUG: Validation failed - name is empty")
		c.HTML(http.StatusOK, "category_form.html", gin.H{
			"title":  "Создать категорию",
			"action": "/categories/create",
			"error":  "Название категории обязательно",
			"category": models.Category{
				Name:        name,
				Description: description,
			},
		})
		return
	}

	category := models.Category{
		Name:        name,
		Description: description,
	}

	log.Printf("DEBUG: Creating category object: %+v", category)

	// Исправлено: используем возвращаемое значение
	categoryID := storage.CreateCategory(category)
	log.Printf("DEBUG: Category created with ID: %d", categoryID)

	// Проверяем что сохранилось
	categories := storage.GetAllCategories()
	log.Printf("DEBUG: Total categories in storage after create: %d", len(categories))
	for i, cat := range categories {
		log.Printf("DEBUG: Category %d: ID=%d, Name=%s, Description=%s",
			i, cat.ID, cat.Name, cat.Description)
	}

	c.Redirect(http.StatusFound, "/categories")
}

func CategoryEditPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("ERROR: Invalid category ID: %s", c.Param("id"))
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	log.Printf("DEBUG: Loading category for edit: ID=%d", id)

	category, exists := storage.GetCategory(id)
	if !exists {
		log.Printf("ERROR: Category not found: ID=%d", id)
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	log.Printf("DEBUG: Found category: %+v", category)

	c.HTML(http.StatusOK, "category_form.html", gin.H{
		"title":    "Редактировать категорию",
		"action":   fmt.Sprintf("/categories/%d/edit", id),
		"category": category,
	})
}

func CategoryUpdateHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("ERROR: Invalid category ID: %s", c.Param("id"))
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")

	log.Printf("DEBUG: Updating category ID=%d - Name: %s, Description: %s",
		id, name, description)

	if name == "" {
		log.Printf("DEBUG: Validation failed for update - name is empty")
		c.HTML(http.StatusOK, "category_form.html", gin.H{
			"title":  "Редактировать категорию",
			"action": fmt.Sprintf("/categories/%d/edit", id),
			"error":  "Название категории обязательно",
			"category": models.Category{
				ID:          id,
				Name:        name,
				Description: description,
			},
		})
		return
	}

	category := models.Category{
		ID:          id,
		Name:        name,
		Description: description,
	}

	log.Printf("DEBUG: Updating category: %+v", category)

	success := storage.UpdateCategory(id, category)
	if !success {
		log.Printf("ERROR: Failed to update category ID=%d", id)
		c.HTML(http.StatusInternalServerError, "category_form.html", gin.H{
			"title":    "Редактировать категорию",
			"action":   fmt.Sprintf("/categories/%d/edit", id),
			"error":    "Ошибка при обновлении категории",
			"category": category,
		})
		return
	}

	log.Printf("DEBUG: Category updated successfully: ID=%d", id)

	c.Redirect(http.StatusFound, "/categories")
}

func CategoryDeleteHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("ERROR: Invalid category ID: %s", c.Param("id"))
		c.Redirect(http.StatusFound, "/categories")
		return
	}

	log.Printf("DEBUG: Deleting category: ID=%d", id)

	success := storage.DeleteCategory(id)
	if !success {
		log.Printf("ERROR: Failed to delete category ID=%d", id)
	}

	log.Printf("DEBUG: Category delete attempt completed: ID=%d, Success: %v", id, success)

	c.Redirect(http.StatusFound, "/categories")
}
