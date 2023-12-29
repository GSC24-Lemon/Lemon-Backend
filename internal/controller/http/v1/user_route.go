package v1

import (
	"lemon_be/internal/entity"
	"lemon_be/internal/usecase"
	"lemon_be/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	c usecase.UserUseCaseI
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, c usecase.UserUseCaseI, l logger.Interface) {
	r := &userRoutes{c, l}

	h := handler.Group("/user")
	{
		h.POST("/registerName", r.registerUsername)
	}
}

type userRegisterRequest struct {
	Username  string `json:"username"`
	DeviceId  string `json:"deviceId"`
	Telephone string `json:"telephone"`
}

type okResponse struct {
	Messsage string `json:"message" example:"message"`
}

// @Summary     save (deviceId, username) to redis
// @Description   save (deviceId, username) to redis
// @ID          getNearestCaregiverRequest
// @Tags  	    group
// @Accept      json
// @Produce     json
// @Param       request body userRegisterRequest true "set up user name"
// @Success     200 {object}
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /v1/groups [post]
// Author: https://github.com/lintang-b-s
func (r *userRoutes) registerUsername(c *gin.Context) {
	var request userRegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - findNearestCaregiverRequest")
		ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	r.c.SaveUsernameAndDeviceId(c.Request.Context(), entity.SaveUsername{
		Username:  request.Username,
		DeviceId:  request.DeviceId,
		Telephone: request.Telephone,
	})
	c.JSON(http.StatusOK, okResponse{Messsage: "ok"})
}
