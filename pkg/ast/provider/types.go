package provider

// SourceLocation represents a location in source code
type SourceLocation struct {
	File   string // File path
	Line   int    // Line number
	Column int    // Column number
}

// InjectParam represents a parameter to be injected
type InjectParam struct {
	Star         string // "*" for pointer types, empty for value types
	Type         string // The type name
	Package      string // The package path
	PackageAlias string // The package alias used in the file
}

// Provider represents a provider struct with metadata
type Provider struct {
	StructName       string                 // Name of the struct
	ReturnType       string                 // Return type of the provider
	Mode             ProviderMode           // Provider mode (basic, grpc, event, job, cronjob, model)
	ProviderGroup    string                 // Provider group for dependency injection
	GrpcRegisterFunc string                 // gRPC register function name
	NeedPrepareFunc  bool                   // Whether prepare function is needed
	InjectParams     map[string]InjectParam // Parameters to inject
	Imports          map[string]string      // Required imports
	PkgName          string                 // Package name
	ProviderFile     string                 // Output file path
	Location         SourceLocation         // Location in source code
	Comment          string                 // Provider comment/documentation
}

// ProviderDescribe represents the parsed provider annotation
type ProviderDescribe struct {
	IsOnly     bool   // Whether only mode is enabled
	Mode       string // Provider mode (job, grpc, event, etc.)
	ReturnType string // Return type specification
	Group      string // Provider group
}
