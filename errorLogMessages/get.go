package errorLogMessages

import (
	"log"
	"net/http"
	"time"

	"sap/errorlog-rest-dataingestion/cassandra"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

func Get(c echo.Context) error {
	var (
		err                 error
		errorLogMessageList []ErrorLogMessage
	)
	fromTime := c.QueryParam("from")
	toTime := c.QueryParam("to")
	m := map[string]interface{}{}
	tenant := c.Get("tenant").(string)

	fromTimeStamp, err := time.Parse(timeLayout, fromTime)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid timestamps")
	}
	toTimeStamp, err := time.Parse(timeLayout, toTime)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid timestamps")
	}

	fromTimeStampISO8601 := fromTimeStamp.Format("2006-01-02T15:04:05Z07:00")
	toTimeStampISO8601 := toTimeStamp.Format("2006-01-02T15:04:05Z07:00")

	epochDay := fromTimeStamp.Unix() / 86400

	query := "SELECT time, status, message_id FROM errormessage WHERE tenant = ? AND epoch_day = ? AND time >= ? AND  time <= ?"
	iterable := cassandra.Session.Query(query, tenant, epochDay, fromTimeStampISO8601, toTimeStampISO8601).Consistency(gocql.One).Iter()

	for iterable.MapScan(m) {
		errorLogMessageList = append(errorLogMessageList, ErrorLogMessage{
			Tenant:    m["tenant"].(string),
			EpochDay:  m["epoch_day"].(int),
			Time:      m["time"].(string),
			MessageID: m["message_id"].(string),
			Message:   m["message"].(string),
			Reason:    m["reason"].(string),
			Status:    m["status"].(string),
		})
		m = map[string]interface{}{}
	}
	return c.JSON(http.StatusOK, ErrorLogMessageResponse{ErrorLogMessages: errorLogMessageList})
}
