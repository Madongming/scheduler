package api

import "errors"

var (
	ErrorPluginNotFound = errors.New("Plugin is not found")
	ErrorNotScheduleJob = errors.New("Not a scheduled job")
)
