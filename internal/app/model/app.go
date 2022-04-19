package model

type AppCreateModel struct {
	QueryNamespace
	Name  string `json:"name" binding:"required"`
	Chart string `json:"chart" binding:"required"`
}
