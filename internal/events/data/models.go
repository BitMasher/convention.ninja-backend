package data

import (
	"convention.ninja/internal/data"
	data2 "convention.ninja/internal/organizations/data"
	"time"
)

type Event struct {
	data.SnowflakeModel
	Name           string `json:"name"`
	Recurring      bool   `json:"recurring"`
	OrganizationId int64  `json:"organizationId,string"`
}

type EventSchedule struct {
	data.SnowflakeModel
	EventId        int64                `json:"eventId,string"`
	OrganizationId int64                `json:"organizationId,string"`
	StartDate      time.Time            `json:"startDate"`
	EndDate        time.Time            `json:"endDate"`
	Schedule       []EventScheduleEntry `json:"schedule"`
}

type EventScheduleEntry struct {
	data.SnowflakeModel
	Name            string          `json:"name"`
	StartTime       time.Time       `json:"startDate"`
	EndTime         time.Time       `json:"endDate"`
	VenueRoomId     int64           `json:"venueRoomId,string"`
	VenueRoom       data2.VenueRoom `json:"venueRoom"`
	EventScheduleId int64           `json:"eventScheduleId,string"`
	OrganizationId  int64           `json:"organizationId,string"`
}
