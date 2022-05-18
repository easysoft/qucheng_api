// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package model

import "gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"

type AppRespStatus struct {
	Status     string                   `json:"status"`
	Age        int64                    `json:"age"`
	AccessHost string                   `json:"access_host"`
	Components []AppRespStatusComponent `json:"components"`
}

type AppRespStatusComponent struct {
	Name       string                 `json:"name"`
	Kind       string                 `json:"kind"`
	StatusCode constant.AppStatusType `json:"-"`
	Status     string                 `json:"status"`
	Replicas   int32                  `json:"replicas"`
	Age        int64                  `json:"age"`
}
