package entity

type Gender int

// const (
// 	Male Gender  = iota +1
// 	Female
// )

// func (g Gender) String() string{
// 	return [...]string{"male", "female"}[g-1];
// }

type CreateCaregiverRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	Job      string `json:"job"`
	Age      uint   `json:"age"`
}

type Caregiver struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Age            uint   `json:"age"`
	Gender         string `json:"gender"`
	Job            string `json:"job"`
	HashedPassword string `json:"password"`
}

type CaregiverLocation struct {
	Latitude  string    `json:"latitude"`
	Longitude string    `json:"longitude"`
	Caregiver Caregiver `json:"caregiver"`
}
