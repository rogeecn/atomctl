package provider

// ProviderMode represents the mode of a provider
type ProviderMode string

const (
	// ProviderModeBasic is the default provider mode
	ProviderModeBasic ProviderMode = "basic"

	// ProviderModeGrpc is for gRPC service providers
	ProviderModeGrpc ProviderMode = "grpc"

	// ProviderModeEvent is for event-based providers
	ProviderModeEvent ProviderMode = "event"

	// ProviderModeJob is for job-based providers
	ProviderModeJob ProviderMode = "job"

	// ProviderModeCronJob is for cron job providers
	ProviderModeCronJob ProviderMode = "cronjob"

	// ProviderModeModel is for model-based providers
	ProviderModeModel ProviderMode = "model"
)

// IsValidProviderMode checks if a provider mode is valid
func IsValidProviderMode(mode string) bool {
	switch ProviderMode(mode) {
	case ProviderModeBasic, ProviderModeGrpc, ProviderModeEvent,
		ProviderModeJob, ProviderModeCronJob, ProviderModeModel:
		return true
	default:
		return false
	}
}

// InjectionMode represents the injection mode for provider fields
type InjectionMode string

const (
	// InjectionModeOnly injects only fields marked with inject:"true"
	InjectionModeOnly InjectionMode = "only"

	// InjectionModeExcept injects all fields except those marked with inject:"false"
	InjectionModeExcept InjectionMode = "except"

	// InjectionModeAuto injects all non-scalar fields automatically
	InjectionModeAuto InjectionMode = "auto"
)

// IsValidInjectionMode checks if an injection mode is valid
func IsValidInjectionMode(mode string) bool {
	switch InjectionMode(mode) {
	case InjectionModeOnly, InjectionModeExcept, InjectionModeAuto:
		return true
	default:
		return false
	}
}
