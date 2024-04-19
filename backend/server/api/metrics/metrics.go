package metrics

import (
	"github.com/apache/incubator-devlake/core/errors"
	// "github.com/apache/incubator-devlake/core/models"
	"net/http"

	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services"

	"github.com/gin-gonic/gin"
)

func GetMetrics(c *gin.Context) {
	params := map[string]interface{}{
		"projects": c.Query("projects"),
		"dateFrom": c.Query("dateFrom"),
		"dateTo":   c.Query("dateTo"),
	}
	metricsData, err := services.GetMetrics(params)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project"))
		return
	}

	shared.ApiOutputSuccess(c, metricsData, http.StatusOK)
}

func GetDeploymentFrequency(c *gin.Context) {
	params := map[string]interface{}{
		"projects": c.Query("projects"),
		"dateFrom": c.Query("dateFrom"),
		"dateTo":   c.Query("dateTo"),
	}
	metricsData, err := services.GetDeploymentFrequency(params)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project"))
		return
	}

	shared.ApiOutputSuccess(c, metricsData, http.StatusOK)
}

func GetLeadTimeForChanges(c *gin.Context) {
	params := map[string]interface{}{
		"projects": c.Query("projects"),
		"dateFrom": c.Query("dateFrom"),
		"dateTo":   c.Query("dateTo"),
	}
	metricsData, err := services.GetLeadTimeForChanges(params)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project"))
		return
	}

	shared.ApiOutputSuccess(c, metricsData, http.StatusOK)
}

func GetMedianTimeToRestore(c *gin.Context) {
	params := map[string]interface{}{
		"projects": c.Query("projects"),
		"dateFrom": c.Query("dateFrom"),
		"dateTo":   c.Query("dateTo"),
	}
	metricsData, err := services.GetMedianTimeToRestore(params)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project"))
		return
	}

	shared.ApiOutputSuccess(c, metricsData, http.StatusOK)
}

func GetChangeFailureRate(c *gin.Context) {
	params := map[string]interface{}{
		"projects": c.Query("projects"),
		"dateFrom": c.Query("dateFrom"),
		"dateTo":   c.Query("dateTo"),
	}
	metricsData, err := services.GetChangeFailureRate(params)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project"))
		return
	}

	shared.ApiOutputSuccess(c, metricsData, http.StatusOK)
}
