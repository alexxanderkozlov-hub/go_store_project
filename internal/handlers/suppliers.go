package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"go_store_project/internal/models"
	"go_store_project/internal/storage"

	"github.com/gin-gonic/gin"
)

func SuppliersPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	supplierList := storage.GetAllSuppliers()

	c.HTML(http.StatusOK, "suppliers.html", gin.H{
		"username":    GetUsername(c),
		"suppliers":   supplierList,
		"total_count": len(supplierList),
	})
}

func SupplierCreatePage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "supplier_form.html", gin.H{
		"title":  "Добавить поставщика",
		"action": "/suppliers/create",
	})
}

func SupplierCreateHandler(c *gin.Context) {
	if !CheckAuth(c) {
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
			"supplier": models.Supplier{
				Name:    name,
				Phone:   phone,
				Email:   email,
				Address: address,
			},
		})
		return
	}

	supplier := models.Supplier{
		Name:    name,
		Phone:   phone,
		Email:   email,
		Address: address,
	}

	storage.CreateSupplier(supplier)

	c.Redirect(http.StatusFound, "/suppliers")
}

func SupplierEditPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	supplier, exists := storage.GetSupplier(id)
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

func SupplierUpdateHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	supplier, exists := storage.GetSupplier(id)
	if !exists {
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	supplier.Name = c.PostForm("name")
	supplier.Phone = c.PostForm("phone")
	supplier.Email = c.PostForm("email")
	supplier.Address = c.PostForm("address")

	storage.UpdateSupplier(id, supplier)

	c.Redirect(http.StatusFound, "/suppliers")
}

func SupplierDeleteHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/suppliers")
		return
	}

	storage.DeleteSupplier(id)
	c.Redirect(http.StatusFound, "/suppliers")
}
