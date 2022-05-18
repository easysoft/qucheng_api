package model

type NamespaceBase struct {
	QueryCluster
	Name string `form:"name" json:"name" binding:"required"`
}
