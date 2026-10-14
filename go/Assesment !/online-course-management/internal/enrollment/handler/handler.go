package enrollmenthandler

import (
	"errors"
	"net/http"
	enrollmentModel "online-course-management/internal/enrollment/model"
	"online-course-management/internal/enrollment/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Enroll(c *gin.Context) {
	var req enrollmentModel.EnrollRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: studentId and courseId are required",
		})
		return
	}

	enrollment, err := h.service.Enroll(c.Request.Context(), req.StudentID, req.CourseID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStudentNotFound), errors.Is(err, service.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrAlreadyEnrolled):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrStudentLimitMax), errors.Is(err, service.ErrCourseCapacityMax):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enroll student"})
		}
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}

func (h *Handler) GetCoursesByStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	courses, err := h.service.GetCoursesByStudentID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrStudentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (h *Handler) GetStudentsByCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	students, err := h.service.GetStudentsByCourseID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrCourseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, students)
}

func (h *Handler) GetTopCourses(c *gin.Context) {
	courses, err := h.service.GetTopEnrolledCourses(c.Request.Context(), 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch top courses"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func (h *Handler) GetUnenrolledStudents(c *gin.Context) {
	students, err := h.service.GetUnenrolledStudents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch unenrolled students"})
		return
	}
	c.JSON(http.StatusOK, students)
}

