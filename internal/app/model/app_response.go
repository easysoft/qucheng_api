package model

import "gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"

type AppRespStatus struct {
	Status     string                   `json:"status"`
	AccessHost string                   `json:"access_host"`
	Components []AppRespStatusComponent `json:"components"`
}

type AppRespStatusComponent struct {
	Name       string                 `json:"name"`
	Kind       string                 `json:"kind"`
	StatusCode constant.AppStatusType `json:"-"`
	Status     string                 `json:"status"`
	Replicas   int32                  `json:"replicas"`
}
