package handler

import (
	"net/http"
	"service-parser/internal/app/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	fetchService     *service.FetchService
	statisticService *service.StatisticService
}

func New(fetchService *service.FetchService, statisticService *service.StatisticService) *Handler {
	return &Handler{
		fetchService:     fetchService,
		statisticService: statisticService,
	}
}

func (h *Handler) Download(c *gin.Context) {
	err := h.fetchService.StartDownload()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusAccepted)
}

func (h *Handler) Stats(c *gin.Context) {
	stats, err := h.statisticService.Stats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
