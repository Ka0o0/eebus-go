package ohpcf

import (
	"fmt"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	spinemodel "github.com/enbility/spine-go/model"
)

var _ ucapi.CemOHPCFInterface = (*OHPCF)(nil)

// A) Announcement of an optional power consumption proces
func (e *OHPCF) ReadOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) (*ucapi.OptionalPowerConsumptionAnnouncement, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	sem, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || sem == nil {
		return nil, err
	}

	data, err := sem.GetData()
	if err != nil || data == nil {
		return nil, err
	}

	if len(data.Alternatives) == 0 ||
		len(data.Alternatives[0].PowerSequence) == 0 ||
		len(data.Alternatives[0].PowerSequence[0].PowerTimeSlot) == 0 ||
		len(data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value) == 0 {
		return nil, api.ErrMissingData
	}

	powerSequence := data.Alternatives[0].PowerSequence[0]

	// if *powerSequence.State.State != spinemodel.PowerSequenceStateTypeInactive ||
	// 	!*powerSequence.State.SequenceRemoteControllable {
	// 	return nil, api.ErrMissingData
	// }

	activeDurationMin, err := powerSequence.OperatingConstraintsDuration.ActiveDurationMin.GetTimeDuration()
	if err != nil {
		return nil, api.ErrMissingData
	}

	var pauseDurationMin *time.Duration
	if powerSequence.OperatingConstraintsDuration.PauseDurationMin != nil {
		pauseDurationMinTemp, err := powerSequence.OperatingConstraintsDuration.PauseDurationMin.GetTimeDuration()
		if err != nil {
			return nil, api.ErrMissingData
		}
		pauseDurationMin = &pauseDurationMinTemp
	}

	// TODO: If more than one "value is present, each "value" SHALL have a different valueType.
	return &ucapi.OptionalPowerConsumptionAnnouncement{
		Power:             powerSequence.PowerTimeSlot[0].ValueList.Value[0].Value.GetValue(),
		PowerType:         *powerSequence.PowerTimeSlot[0].ValueList.Value[0].ValueType,
		IsStoppable:       *powerSequence.OperatingConstraintsInterrupt.IsStoppable,
		IsPausable:        *powerSequence.OperatingConstraintsInterrupt.IsPausable,
		ActiveDurationMin: activeDurationMin,
		PauseDurationMin:  pauseDurationMin,
	}, nil
}

// B) Report of the current state of the scheduled or active power consumption process
func (e *OHPCF) ReadScheduledOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) (*ucapi.ScheduledOptionalPowerConsumptionAnnouncement, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e *OHPCF) ScheduleOptionalPowerConsumptionProcess(entity spineapi.EntityRemoteInterface, startTime time.Time) error {

	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	sem, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || sem == nil {
		return err
	}

	// data, err := sem.GetData()
	// if err != nil || data == nil {
	// 	return err
	// }
	schedule := &model.PowerSequenceScheduleDataType{}
	// duration := time.Until(startTime)
	schedule.StartTime = model.NewAbsoluteOrRelativeTimeType("PT0S")

	sem2 := &spinemodel.SmartEnergyManagementPsDataType{
		Alternatives: []spinemodel.SmartEnergyManagementPsAlternativesType{
			{
				// Relation: &spinemodel.SmartEnergyManagementPsAlternativesRelationType{
				// 	AlternativesId: util.Ptr[spinemodel.AlternativesIdType](0),
				// },
				PowerSequence: []spinemodel.SmartEnergyManagementPsPowerSequenceType{
					{
						Description: &spinemodel.PowerSequenceDescriptionDataType{
							SequenceId: util.Ptr[spinemodel.PowerSequenceIdType](0),
						},
						Schedule: schedule,
					},
				}},
		},
	}

	_, err = sem.WriteData(sem2)

	return err
}

func (e *OHPCF) PauseOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) error {

	return fmt.Errorf("Not implemented")
}

func (e *OHPCF) ResumeOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) error {
	return fmt.Errorf("Not implemented")
}

func (e *OHPCF) StopOptionalPowerConsumption(entity spineapi.EntityRemoteInterface) error {

	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	sem, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || sem == nil {
		return err
	}

	// data, err := sem.GetData()
	// if err != nil || data == nil {
	// 	return err
	// }

	// data.Alternatives[0].PowerSequence[0].State.State = util.Ptr(model.PowerSequenceStateTypeInvalid)

	sem2 := &spinemodel.SmartEnergyManagementPsDataType{
		Alternatives: []spinemodel.SmartEnergyManagementPsAlternativesType{
			{
				// Relation: &spinemodel.SmartEnergyManagementPsAlternativesRelationType{
				// 	AlternativesId: util.Ptr[spinemodel.AlternativesIdType](0),
				// },
				PowerSequence: []spinemodel.SmartEnergyManagementPsPowerSequenceType{
					{
						Description: &spinemodel.PowerSequenceDescriptionDataType{
							SequenceId: util.Ptr[spinemodel.PowerSequenceIdType](0),
						},
						State: &spinemodel.PowerSequenceStateDataType{
							State: util.Ptr(model.PowerSequenceStateTypeInvalid),
						},
					},
				}},
		},
	}

	_, err = sem.WriteData(sem2)

	return err
}

// return the last known SoC of the connected EV
//
// only works with a current ISO15118-2 with VAS or ISO15118-20
// communication between EVSE and EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
// func (e *OHPCF) StateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error) {
// 	if !e.IsCompatibleEntityType(entity) {
// 		return 0, api.ErrNoCompatibleEntity
// 	}

// 	evMeasurement, err := client.NewMeasurement(e.LocalEntity, entity)
// 	if err != nil || evMeasurement == nil {
// 		return 0, err
// 	}

// 	filter := model.MeasurementDescriptionDataType{
// 		MeasurementType: util.Ptr(model.MeasurementTypeTypePercentage),
// 		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
// 		ScopeType:       util.Ptr(model.ScopeTypeTypeStateOfCharge),
// 	}
// 	result, err := evMeasurement.GetDataForFilter(filter)
// 	if err != nil || len(result) == 0 || result[0].Value == nil {
// 		return 0, api.ErrDataNotAvailable
// 	}
// 	return result[0].Value.GetValue(), nil
// }
