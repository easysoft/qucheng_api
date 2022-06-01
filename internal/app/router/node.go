package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"
)

func ClusterStatistics(c *gin.Context) {
	var (
		ctx   = c.Request.Context()
		err   error
		query model.QueryCluster
	)

	if err = c.ShouldBindQuery(&query); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	metricData, err := service.Nodes(ctx, query.Cluster).Statistic()
	if err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	data := model.ClusterMetric{
		Status:    constant.ClusterStatusMap[constant.ClusterStatusNormal],
		NodeCount: len(service.Nodes(ctx, query.Cluster).GetNodes()),
		Metrics:   metricData,
	}
	renderJson(c, http.StatusOK, data)
}
