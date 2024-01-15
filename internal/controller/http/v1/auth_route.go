package v1

import (
	"lemon_be/internal/controller/http/errorWrapper"
	"lemon_be/internal/entity"
	"lemon_be/internal/usecase"
	"lemon_be/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	a usecase.Auth
	l logger.Interface
}

func newAuthRoutes(handler *gin.RouterGroup, a usecase.Auth, l logger.Interface) {
	r := &authRoutes{a, l}

	h := handler.Group("/auth")
	{
		h.POST("/register", r.registerUser)
		h.POST("/login", r.loginUser)
		//h.POST("/token", r.renewAccessToken)
		h.DELETE("/logout", r.deleteRefreshToken)
	}

}

type createUserRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Gender   string `json:"gender" binding:"required"`
	Job      string `json:"job" binding:"required"`
	Age      uint   `json:"age" binding:"required"`
}

type userResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    uint   `json:"age"`
	Gender string `json:"gender"`
	Job    string `json:"job"`
}

// @Summary     Register User in Db
// @Description   Register User in Db
// @ID          registerUser
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body createUserRequest true "Set up user"
// @Success     200 {object} userResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /v1/auth/register [post]
// Author: https://github.com/lintang-b-s
func (r *authRoutes) registerUser(c *gin.Context) {
	var request createUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - registerUser")
		ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := r.a.Register(
		c.Request.Context(),
		entity.CreateCaregiverRequest{
			Name:     request.Name,
			Password: request.Password,
			Email:    request.Email,
			Gender:   request.Gender,
			Job:      request.Job,
			Age:      request.Age,
		},
	)

	if err != nil {

		clientError, ok := err.(errorWrapper.ClientError)
		if !ok {
			r.l.Error("http - v1- registerUser")
			ErrorResponse(c, http.StatusInternalServerError, "register service problems")
			return
		}

		body, err := clientError.ResponseBody()
		if err != nil {
			r.l.Error("http - v1- registerUser")
			ErrorResponse(c, http.StatusInternalServerError, "register service problems")
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

		r.l.Error("http - v1- registerUser")
		ErrorResponse(c, http.StatusInternalServerError, "register service problems: "+err.Error())
		return
	}

	resp := userResponse{Name: user.Name, Email: user.Email, Age: user.Age, Gender: user.Gender, Job: user.Job}
	c.JSON(http.StatusCreated, resp)
}

// ---- loginUser ----
type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=3"`
}

type loginUserResponse struct {
	SessionId             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	Email                 string    `json:"email"`
}

// @Summary     Login User
// @Description   Login User
// @ID          loginUser
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body loginUserRequest true "Login  user"
// @Success     200 {object} loginUserResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /v1/auth/login [post]
// Author: https://github.com/lintang-b-s
func (r *authRoutes) loginUser(c *gin.Context) {

	var request loginUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - loginUser")
		ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	loginResponse, err := r.a.Login(
		c.Request.Context(),
		entity.LoginUserRequest{
			Email:    request.Email,
			Password: request.Password,
		},
	)
	if err != nil {

		clientError, ok := err.(errorWrapper.ClientError)
		if !ok {
			r.l.Error("http - v1- loginUser")
			ErrorResponse(c, http.StatusInternalServerError, "loginUser service problems")
			return
		}

		body, err := clientError.ResponseBody()
		if err != nil {
			r.l.Error("http - v1- loginUser")
			ErrorResponse(c, http.StatusInternalServerError, "loginUser service problems")
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

		r.l.Error("http - v1- loginUser")
		ErrorResponse(c, http.StatusInternalServerError, "loginUser service problems: "+err.Error())
		return
	}

	resp := loginUserResponse{
		SessionId:             loginResponse.SessionId,
		AccessToken:           loginResponse.AccessToken,
		AccessTokenExpiresAt:  loginResponse.AccessTokenExpiresAt,
		RefreshToken:          loginResponse.RefreshToken,
		RefreshTokenExpiresAt: loginResponse.RefreshTokenExpiresAt,
		Email:                 loginResponse.User.Email,
	}
	c.JSON(http.StatusCreated, resp)
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

//
//// @Summary     renew Access Token using user refreshToken
//// @Description    renew Access Token using user refreshToken
//// @ID          renewAccessToken
//// @Tags  	    user
//// @Accept      json
//// @Produce     json
//// @Param       request body renewAccessTokenRequest true "Login  user"
//// @Success     200 {object} renewAccessTokenResponse
//// @Failure     400 {object} response
//// @Failure     500 {object} response
//// @Router      /v1/auth/token [post]
//// Author: https://github.com/lintang-b-s
//func (r *authRoutes) renewAccessToken(c *gin.Context) {
//	var request renewAccessTokenRequest
//
//	if err := c.ShouldBindJSON(&request); err != nil {
//		r.l.Error(err, "http - v1 - loginUser")
//		ErrorResponse(c, http.StatusBadRequest, "invalid request body")
//		return
//	}
//
//	renewResponse, err := r.a.RenewAccessToken(
//		c.Request.Context(),
//		entity.RenewAccessTokenRequest{
//			RefreshToken: request.RefreshToken,
//		},
//	)
//	if err != nil {
//
//		unwrapedErr := errors.Unwrap(err)
//		// jika refresh token invalid / expired
//		//errRepo := errors.Unwrap(unwrapedErr)
//		if unwrapedErr == jwt.ErrInvalidToken || unwrapedErr == jwt.ErrExpiredToken {
//			ErrorResponse(c, http.StatusUnauthorized, "Token invalid or token already expired")
//			return
//		}
//
//		if err.Error() == "Invalid session" {
//			ErrorResponse(c, http.StatusUnauthorized, "Refresh Token mismatch with refresh token in database")
//			return
//		}
//
//		r.l.Error("http - v1- renewAccessToken")
//		ErrorResponse(c, http.StatusInternalServerError, "renewAccessToken service problems: "+err)
//		return
//	}
//
//	resp := renewAccessTokenResponse{
//		AccessToken:          renewResponse.AccessToken,
//		AccessTokenExpiresAt: renewResponse.AccessTokenExpiresAt,
//	}
//
//	c.JSON(http.StatusCreated, resp)
//}

type deleteRefreshTokenRequest struct {
	refreshToken string `json:"refresh_token" binding:"required"`
}

type deleteRefreshTokenResponse struct {
	responseMessage string `json:"response_message"`
}

// @Summary     delete refresh token
// @Description   delete refresh token
// @ID          deleteRefreshToken
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body deleteRefreshTokenRequest true "Login  user"
// @Success     200 {object} deleteRefreshTokenResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /v1/auth/logout [delete]
// Author: https://github.com/lintang-b-s
func (r *authRoutes) deleteRefreshToken(c *gin.Context) {
	var request renewAccessTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - loginUser")
		ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := r.a.DeleteRefreshToken(
		c.Request.Context(),
		entity.DeleteRefreshTokenRequest{
			RefreshToken: request.RefreshToken,
		},
	)
	if err != nil {

		clientError, ok := err.(errorWrapper.ClientError)
		if !ok {
			r.l.Error("http - v1- deleteRefreshToken")
			ErrorResponse(c, http.StatusInternalServerError, "deleteRefreshToken service problems")
			return
		}

		body, err := clientError.ResponseBody()
		if err != nil {
			r.l.Error("http - v1- deleteRefreshToken")
			ErrorResponse(c, http.StatusInternalServerError, "deleteRefreshToken service problems")
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

		r.l.Error("http - v1- deleteRefreshToken")
		ErrorResponse(c, http.StatusInternalServerError, "deleteRefreshToken service problems: "+err.Error())
		return
	}

	resp := deleteRefreshTokenResponse{
		responseMessage: "refresh token successfully deleted!",
	}
	c.JSON(http.StatusOK, resp)
}
