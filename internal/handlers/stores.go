package handlers

import (
	"fmt"
	"net/http"
	"sort" // Добавили импорт sort
	"strconv"

	"go_store_project/internal/models"
	"go_store_project/internal/storage"

	"github.com/gin-gonic/gin"
)

func StoreListPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	storeList := storage.GetAllStores()

	// Добавили сортировку по возрастанию ID
	sort.Slice(storeList, func(i, j int) bool {
		return storeList[i].ID < storeList[j].ID
	})

	c.HTML(http.StatusOK, "store.html", gin.H{
		"stores":      storeList,
		"total_count": len(storeList),
	})
}

func StoreCreatePage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "store_form.html", gin.H{
		"title":  "Добавить магазин",
		"action": "/stores/create",
	})
}

func StoreCreateHandler(c *gin.Context) {
	if !CheckAuth(c) {
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
			"store":  models.Store{Name: name, Address: address},
		})
		return
	}

	store := models.Store{
		Name:    name,
		Address: address,
	}

	storage.CreateStore(store)

	c.Redirect(http.StatusFound, "/stores")
}

func StoreEditPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	store, exists := storage.GetStore(id)
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

func StoreUpdateHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	store, exists := storage.GetStore(id)
	if !exists {
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	store.Name = c.PostForm("name")
	store.Address = c.PostForm("address")

	storage.UpdateStore(id, store)

	c.Redirect(http.StatusFound, "/stores")
}

func StoreDeleteHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/stores")
		return
	}

	storage.DeleteStore(id)
	c.Redirect(http.StatusFound, "/stores")
}
