package constants

import "time"

const (
	//ApiPrefix = "/api/v1"
	ApiPrefix  = ""
	ApiAddress = ":8000"
)

const (
	SessionName       = "JSESSIONID"
	SessionLivingTime = 5 * time.Hour
)
