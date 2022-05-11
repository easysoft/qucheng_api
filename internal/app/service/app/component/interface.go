package component

import "gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"

type Component interface {
	Name() string
	Kind() string
	Replicas() int32
	Status() constant.AppStatusType
	Age() int64
}
