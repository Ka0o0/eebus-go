package ohpcf

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *OHPCF) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.deviceConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.SmartEnergyManagementPsDataType:
		e.measurementDataUpdate(payload)
	default:
		return
	}
}

// an EV was connected
func (e *OHPCF) deviceConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	// TODO: Don't do these requests for now, only add it once SPINE supports handling filtering identical pending subscription requests
	// Also: these are covered by EVCEM anyway, which is required

	if sem, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity); err == nil {
		if _, err := sem.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}
		if _, err := sem.Bind(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement descriptions
		if _, err := sem.RequestData(); err != nil {
			logging.Log().Debug(err)
		}
	}

}

func (e *OHPCF) measurementDataUpdate(payload spineapi.EventPayload) {
	// // Scenario 1
	e.EventCB(payload.Ski, payload.Device, payload.Entity, ScheduledOptionalPowerConsumptionUpdated)
	// if sem, err := client.NewSmartEnergyManagementPs(e.LocalEntity, payload.Entity); err == nil {
	// 	filter := model.MeasurementDescriptionDataType{
	// 		ScopeType: util.Ptr(model.ScopeTypeTypeStateOfCharge),
	// 	}
	// 	// if evMeasurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
	// 	// }
	// }

	// 	sem, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	// if err != nil || sem == nil {
	// 	return nil, err
	// }

	// data, err := sem.GetData()
	// if err != nil || data == nil {
	// 	return nil, err
	// }

}
