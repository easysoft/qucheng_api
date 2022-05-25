// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package model

type AppModel struct {
	QueryNamespace
	Name string `form:"name" json:"name" binding:"required"`
}

type AppCreateModel struct {
	AppModel
	Channel  string          `json:"channel"`
	Chart    string          `json:"chart" binding:"required"`
	Settings []stringSetting `json:"settings"`
}

type stringSetting struct {
	Key string `json:"key"`
	Val string `json:"value"`
}

type AppManageModel struct {
	AppModel
	Channel string `json:"channel"`
	Chart   string `json:"chart" binding:"required"`
}

type AppSettingMode struct {
	AppModel
	Mode string `form:"mode" binding:"oneof=list map"`
}

type AppListModel struct {
	QueryCluster
	Apps []NamespacedApp `json:"apps" binding:"required"`
}

type NamespacedApp struct {
	Namespace string `json:"namespace" binding:"required"`
	Name      string `json:"name" binding:"required"`
}
