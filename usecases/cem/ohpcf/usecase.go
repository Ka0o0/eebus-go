package ohpcf

import (
	"github.com/enbility/eebus-go/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type OHPCF struct {
	*usecase.UseCaseBase
}

// Add support for the EV State of Charge (EVSOC) use case
// as a CEM actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewOHPCF(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *OHPCF {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeCompressor,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeCompressor,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
			// TODO: do we need to set this?
			// ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeMeasurement},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: true,
			// ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeMeasurement},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeOptimizationOfSelfConsumptionByHeatPumpCompressorFlexibility,
		"1.0.0",
		model.UseCaseDocumentSubRevisionRelease,
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &OHPCF{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *OHPCF) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeSmartEnergyManagementPs,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}

func (e *OHPCF) UpdateUseCaseAvailability(available bool) {
	e.LocalEntity.SetUseCaseAvailability(model.UseCaseFilterType{
		Actor:       model.UseCaseActorTypeCEM,
		UseCaseName: e.UseCaseName,
	}, available)
}
