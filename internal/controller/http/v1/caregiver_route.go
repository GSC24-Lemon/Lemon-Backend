package v1

import (
	"lemon_be/internal/entity"
	"lemon_be/internal/usecase"
	"lemon_be/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type caregiverRoutes struct {
	c usecase.Caregiver
	l logger.Interface
}

func newCaregiverRoutes(handler *gin.RouterGroup, c usecase.Caregiver, l logger.Interface) {
	r := &caregiverRoutes{c, l}

	h := handler.Group("/caregiver")
	{
		h.POST("/help", r.notifyNearestCaregiver)
		h.POST("/test", r.testGeoAdd)
	}
}

type notifyNearestCaregiverRequest struct {
	DeviceId    string  `json:"deviceId"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Username    string  `json:"username"`
	Destination string  `json:"destination"`
}

// @Summary     send notification to the nearest caregiver
// @Description   send notification to the nearest caregiver
// @ID          getNearestCaregiverRequest
// @Tags  	    group
// @Accept      json
// @Produce     json
// @Param       request body notifyNearestCaregiverRequest true "set up deviceId & Coordinate of user"
// @Success     200 {object}
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /v1/groups [post]
// Author: https://github.com/lintang-b-s
func (r *caregiverRoutes) notifyNearestCaregiver(c *gin.Context) {
	var request notifyNearestCaregiverRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - findNearestCaregiverRequest")
		ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	r.c.NotifyNearestCaregiver(c.Request.Context(), entity.UserLocation{
		DeviceId:    request.DeviceId,
		Lat:         request.Latitude,
		Long:        request.Longitude,
		Username:    request.Username,
		Destination: request.Destination,
	})

	c.JSON(http.StatusOK, okResponse{Messsage: "ok"})
}



// tes geoadd
func (r *caregiverRoutes) testGeoAdd(c *gin.Context) {
	var request notifyNearestCaregiverRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - findNearestCaregiverRequest")
		ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	r.c.TestGeoAdd(c.Request.Context(), entity.UserLocation{
		DeviceId:    request.DeviceId,
		Lat:         request.Latitude,
		Long:        request.Longitude,
		Username:    request.Username,
		Destination: request.Destination,
	})

	c.JSON(http.StatusOK, okResponse{Messsage: "ok"})
}
