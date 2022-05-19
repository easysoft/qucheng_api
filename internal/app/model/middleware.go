package model

type Middleware struct {
	Name      string `json:"name,omitempty" binding:"required"`
	Namespace string `json:"namespace,omitempty" binding:"required"`
	Type      string `json:"type,omitempty" binding:"required"` // 类型
}
