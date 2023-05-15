package model

type Superhero struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	FullName     string `json:"full_name"`
	GenderId     int    `json:"gender_id"`
	EyeColourId  int    `json:"eye_colour_id"`
	HairColourId int    `json:"hair_colour_id"`
	SkinColourId int    `json:"skin_colour_id"`
	RaceId       int    `json:"race_id"`
	PublisherId  int    `json:"publisher_id"`
	AlignmentId  int    `json:"alignment_id"`
	HeightCm     int    `json:"height_cm"`
	WeightKg     int    `json:"weight_kg"`
}
