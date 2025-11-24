package api

import (
	"time"

	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	spinemodel "github.com/enbility/spine-go/model"
)

type OptionalPowerConsumptionAnnouncement struct {
	Power             float64                               // Power in Watts
	PowerType         spinemodel.PowerTimeSlotValueTypeType // "power", "powerMax"
	IsStoppable       bool
	IsPausable        bool
	ActiveDurationMin time.Duration
	PauseDurationMin  *time.Duration
}

type ScheduledOptionalPowerConsumptionAnnouncement struct {
	StartTime time.Time
	State     string // scheduled, running, paused, stopped
	// TODO: what is point 3? Current options for the Actor CEM to interrupt or resume the consumption [OHPCF-012/3].
}

// Actor: CEM
// UseCase: Optimization of Self Consumption by Heat Pump Compressor Flexibility
type CemOHPCFInterface interface {
	api.UseCaseInterface

	// Scenario 1
	// Monitor heat pump compressor's power consumption flexibility

	// A) Announcement of an optional power consumption proces
	ReadOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) (*OptionalPowerConsumptionAnnouncement, error)

	// B) Report of the current state of the scheduled or active power consumption process
	ReadScheduledOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) (*ScheduledOptionalPowerConsumptionAnnouncement, error)

	// // Scenario 2
	// // Control heat pump compressor's power consumption flexibility

	// // A) Control heat pump compressor's power consumption flexibility

	ScheduleOptionalPowerConsumptionProcess(entity spineapi.EntityRemoteInterface, startTime time.Time) error

	// // B) Interrupt or resume a scheduled power consumption process

	PauseOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) error

	ResumeOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) error

	StopOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) error
}
