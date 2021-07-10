package service

import (
	"math"
	"math/rand"
	"net/http"
	"nsaop/model"
	"nsaop/router/resp"
	"nsaop/utils/constant"
	"time"

	"github.com/gin-gonic/gin"
)

type TrafficHistory struct {
	Traffic []float64 `traffic`
}

// @Tags		Service
// @Summary		Get Traffic in the past 30 days of a service
// @Description	with jwt auth, get Traffic of a certain service
// @Accept		json
// @Product		json
// @Param 		id path string true "service id"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Success		200 {object} resp.Response{data=TrafficHistory}
// @Failure		400 {object} resp.Response "format error<br>id not available"
// @Router		/service/{id}/traffic [GET]
func GetTraffic(c *gin.Context) {
	service := c.MustGet(constant.TableService).(model.Service)
	if service.OnAt == time.Unix(0, 0) {
		resp.OK(c, http.StatusOK, gin.H{})
		return
	}
	// fake traffic
	durationDay := int(time.Now().Sub(service.OnAt).Hours() / 24)
	if durationDay > 30 {
		durationDay = 30
	} else if durationDay == 0 {
		durationDay++
	}
	var history TrafficHistory
	for day := 0; day < durationDay; day++ {
		history.Traffic = append(history.Traffic, 24*3600*(math.Abs(rand.NormFloat64()*0.2+10)))
	}
	resp.OK(c, http.StatusOK, history)
}

/*
return traffic (MB) during some time
*/
func GetTrafficByID(start time.Time, end time.Time) (traffic float64) {
	traffic = end.Sub(start).Seconds()*math.Abs(rand.NormFloat64()) + 1
	return
}
