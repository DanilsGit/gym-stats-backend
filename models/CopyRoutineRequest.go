package models

type CopyRoutineRequest struct {
	UserId          string            `json:"userId"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Public          bool              `json:"public"`
	ExerciseRequest []ExerciseRequest `json:"exercises"`
}
