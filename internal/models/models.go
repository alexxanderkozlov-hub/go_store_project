package models

import "time"

type Store struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Logo      string    `json:"logo"`
	CreatedAt time.Time `json:"created_at"`
}

type Supplier struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	CategoryID  int       `json:"category_id"`
	SKU         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type Supply struct {
	ID         int       `json:"id"`
	SupplierID int       `json:"supplier_id"`
	ProductID  int       `json:"product_id"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	Total      float64   `json:"total"`
	Date       time.Time `json:"date"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}

func GetStatusText(status string) string {
	switch status {
	case "pending":
		return "Ожидается"
	case "delivered":
		return "Доставлено"
	case "cancelled":
		return "Отменено"
	default:
		return status
	}
}

func GetStatusClass(status string) string {
	switch status {
	case "pending":
		return "status-pending"
	case "delivered":
		return "status-delivered"
	case "cancelled":
		return "status-cancelled"
	default:
		return "status-default"
	}
}
