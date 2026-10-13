package coursehandler

import (
	"net/http"
	courseModel "online-course-management/internal/course/model"
	"online-course-management/internal/course/service"
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
	var course courseModel.Course

	err := c.ShouldBindJSON(&course)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	e_rror := h.service.Create(c.Request.Context(), &course)

	if e_rror != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create",
		})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *Handler) GetAll(c *gin.Context) {
	//var courses []courseModel.Cocurse

	courses, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed",
		})
		return
	}
	c.JSON(http.StatusOK, courses)

}

func (h *Handler) GetById(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}
	course, err := h.service.GetById(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	c.JSON(http.StatusOK, course)

}

func (h *Handler) Update(c *gin.Context) {
	var course courseModel.Course

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	err_rr := c.ShouldBindJSON(&course)

	if err_rr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	course.ID = uint(id)
	err_r := h.service.Update(c.Request.Context(), &course)

	if err_r != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update course",
		})
		return
	}

	c.JSON(http.StatusOK, course)

}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	if err := h.service.Delete(
		c.Request.Context(),
		uint(id),
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete course",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
