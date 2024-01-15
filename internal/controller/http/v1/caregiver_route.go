package v1

import (
	"lemon_be/internal/controller/http/errorWrapper"
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

	err := r.c.NotifyNearestCaregiver(c.Request.Context(), entity.UserLocation{
		DeviceId:    request.DeviceId,
		Lat:         request.Latitude,
		Long:        request.Longitude,
		Username:    request.Username,
		Destination: request.Destination,
	})

	if err != nil {

		clientError, ok := err.(errorWrapper.ClientError)
		if !ok {
			r.l.Error("http - v1- notifyNearestCaregiver")
			ErrorResponse(c, http.StatusInternalServerError, "notifyNearestCaregiver service problems")
			return
		}

		body, err := clientError.ResponseBody()
		if err != nil {
			r.l.Error("http - v1- notifyNearestCaregiver")
			ErrorResponse(c, http.StatusInternalServerError, "notifyNearestCaregiver service problems")
			return
		}

		status, _ := clientError.ResponseHeaders()
		if status != 0 {
			ErrorResponse(c, status, string(body[:]))
			return
		}
		// errorM := errors.Unwrap(err)
		// errM := errors.Unwrap(err)
		// errorType := strings.Fields(errM.Error())
		// if errorType[0] == "BadRequest" {
		// 	ErrorResponse(c, http.StatusBadRequest, "Bad request:  User With same username or email already exists")
		// 	return
		// }

		r.l.Error("http - v1- notifyNearestCaregiver")
		ErrorResponse(c, http.StatusInternalServerError, "notifyNearestCaregiver service problems: "+err.Error())
		return
	}

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
