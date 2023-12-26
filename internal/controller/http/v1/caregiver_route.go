package v1

// import (
// 	"lemon_be/internal/entity"
// 	"lemon_be/internal/usecase"
// 	"lemon_be/pkg/logger"

// 	"github.com/gin-gonic/gin"
// )

// type CaregiverRoutes struct {
// 	c usecase.Caregiver
// 	l logger.Interface
// }

// func newCaregiverRoutes(handler *gin.RouterGroup, c usecase.Caregiver, l logger.Interface){
// 	 r := &CaregiverRoutes{c, l};

// 	 h := handler.Group("/caregiver")
// 	 {
// 			h.GET("/getNearestCaregiver", r.getNearestCaregiver);
// 			h.POST("/register", r.register);
// 	 }
// }

// type CaregiverResponse  struct {
// 	Caregiver entity.Caregiver `json:"caregiver"`
// }
