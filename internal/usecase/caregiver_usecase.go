package usecase



type CaregiverUseCase struct {
	repo CaregiverGeoRepo
}


func NewCaregiverUseCase(r CaregiverGeoRepo) *CaregiverUseCase{
	return &CaregiverUseCase{
		repo: r,
	}

}




