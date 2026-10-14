package studenthandler

import (
	"net/http"
	studentModel "online-course-management/internal/student/model"
	"online-course-management/internal/student/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c *gin.Context) {
	var student studentModel.Student

	err := c.ShouldBindJSON(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err = h.service.Create(c.Request.Context(), &student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create student",
		})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (h *Handler) GetAll(c *gin.Context) {
	students, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch students",
		})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *Handler) GetById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	student, err := h.service.GetById(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	var student studentModel.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	student.ID = uint(id)
	if err := h.service.Update(c.Request.Context(), &student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update student",
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete student",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
