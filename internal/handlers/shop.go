package handlers

import (
	"go_store_project/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShopPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	products := storage.GetAllProducts()

	c.HTML(http.StatusOK, "shop.html", gin.H{
		"products": products,
		"role":     GetRole(c),
	})
}
