// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package model

import "gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"

type AppRespStatus struct {
	Status     string                   `json:"status"`
	Version    string                   `json:"version"`
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

type NamespaceAppMetric struct {
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Metrics   *AppMetric `json:"metrics"`
	Status    string     `json:"status"`
	Age       int64      `json:"age"`
}

type AppMetric struct {
	Cpu    ResourceCpu    `json:"cpu"`
	Memory ResourceMemory `json:"memory"`
}

type ResourceCpu struct {
	Usage float64 `json:"usage"`
	Limit float64 `json:"limit"`
}

type ResourceMemory struct {
	Usage int64 `json:"usage"`
	Limit int64 `json:"limit"`
}

type NodeMetric struct {
	Cpu    NodeResourceCpu    `json:"cpu"`
	Memory NodeResourceMemory `json:"memory"`
}

type NodeResourceCpu struct {
	Usage       float64 `json:"usage"`
	Capacity    float64 `json:"capacity"`
	Allocatable float64 `json:"allocatable"`
}

type NodeResourceMemory struct {
	Usage       int64 `json:"usage"`
	Capacity    int64 `json:"capacity"`
	Allocatable int64 `json:"allocatable"`
}

type ClusterMetric struct {
	Status    string     `json:"status"`
	NodeCount int        `json:"node_count"`
	Metrics   NodeMetric `json:"metrics"`
}

type Component struct {
	Name string `json:"name"`
}
