package model

type QueryCluster struct {
	Cluster string `form:"cluster" json:"cluster" binding:"required"`
}

type QueryNamespace struct {
	QueryCluster
	Namespace string `form:"namespace" json:"namespace" binding:"required"`
}
