package models

type CourseConfig struct {
	MaxPoints        int `json:"max_point"`
	MinPointsForPass int `json:"min_points_for_pass"`
}

func (c *CourseConfig) Update(maxPoints, minPointsForPass int) {
	c.MaxPoints = maxPoints
	c.MinPointsForPass = minPointsForPass
}
