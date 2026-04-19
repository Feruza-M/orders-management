package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"order-management-app/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := h.DB.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) ListOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	rows, err := h.DB.Query(ctx, `SELECT id, customer, email, amount, status, description, created_at, updated_at FROM orders ORDER BY id DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.Customer, &o.Email, &o.Amount, &o.Status, &o.Description, &o.CreatedAt, &o.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = append(orders, o)
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	var o models.Order
	err = h.DB.QueryRow(ctx, `SELECT id, customer, email, amount, status, description, created_at, updated_at FROM orders WHERE id = $1`, id).Scan(&o.ID, &o.Customer, &o.Email, &o.Amount, &o.Status, &o.Description, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	var id int64
	err := h.DB.QueryRow(ctx, `INSERT INTO orders (customer, email, amount, status, description) VALUES ($1, $2, $3, $4, $5) RETURNING id`, req.Customer, req.Email, req.Amount, req.Status, req.Description).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) UpdateOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	cmd, err := h.DB.Exec(ctx, `UPDATE orders SET customer = $1, email = $2, amount = $3, status = $4, description = $5, updated_at = NOW() WHERE id = $6`, req.Customer, req.Email, req.Amount, req.Status, req.Description, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": true})
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	cmd, err := h.DB.Exec(ctx, `DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
