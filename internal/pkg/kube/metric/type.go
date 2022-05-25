package metric

import "k8s.io/apimachinery/pkg/api/resource"

type Res struct {
	Cpu    *resource.Quantity
	Memory *resource.Quantity
}
