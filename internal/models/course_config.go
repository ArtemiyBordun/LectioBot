package models

type CourseConfig struct {
	MaxPoints           int `json:"max_point"`
	MinPointsOPForPass  int `json:"min_points_op_for_pass"`
	MinPointsOOPForPass int `json:"min_points_oop_for_pass"`
	TotalLectures       int `json:"total_lectures"`
}

func (c *CourseConfig) Update(maxPoints, minPointsOPForPass, minPointsOOPForPass, totalLectures int) {
	if maxPoints != -1 {
		c.MaxPoints = maxPoints
	}
	if minPointsOPForPass != -1 {
		c.MinPointsOPForPass = minPointsOPForPass
	}
	if minPointsOOPForPass != -1 {
		c.MinPointsOOPForPass = minPointsOOPForPass
	}
	if totalLectures != -1 {
		c.TotalLectures = totalLectures
	}
}
