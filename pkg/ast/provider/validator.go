package provider

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Validator defines the interface for validating provider configurations
type Validator interface {
	// Validate validates a provider configuration
	Validate(provider *Provider) error

	// ValidateComment validates a provider comment annotation
	ValidateComment(comment string) error

	// AddRule adds a validation rule to the validator
	AddRule(rule ValidationRule)

	// RemoveRule removes a validation rule from the validator
	RemoveRule(name string)

	// GetRules returns all validation rules
	GetRules() []ValidationRule

	// ValidateAll validates multiple providers and returns a comprehensive report
	ValidateAll(providers []*Provider) *ValidationReport
}

// ValidationRule defines the interface for validation rules
type ValidationRule interface {
	// Name returns the name of the validation rule
	Name() string

	// Validate validates a provider against this rule
	Validate(provider *Provider) *ValidationError

	// Description returns a human-readable description of the rule
	Description() string
}

// ValidationError represents a validation error
type ValidationError struct {
	RuleName    string `json:"rule_name"`
	Message     string `json:"message"`
	Field       string `json:"field,omitempty"`
	Value       string `json:"value,omitempty"`
	Severity    string `json:"severity"` // "error", "warning", "info"
	Suggestion  string `json:"suggestion,omitempty"`
	ProviderRef string `json:"provider_ref,omitempty"`
	Cause       error  `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.ProviderRef != "" {
		return fmt.Sprintf("[%s] %s: %s", e.ProviderRef, e.RuleName, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.RuleName, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *ValidationError) Unwrap() error {
	return e.Cause
}

// ValidationReport represents the result of validating multiple providers
type ValidationReport struct {
	Timestamp      time.Time             `json:"timestamp"`
	TotalProviders int                   `json:"total_providers"`
	ValidCount     int                   `json:"valid_count"`
	InvalidCount   int                   `json:"invalid_count"`
	Errors         []ValidationError     `json:"errors"`
	Warnings       []ValidationError     `json:"warnings"`
	Infos          []ValidationError     `json:"infos"`
	IsValid        bool                  `json:"is_valid"`
	Statistics     *ValidationStatistics `json:"statistics,omitempty"`
}

// ValidationStatistics contains detailed statistics about validation results
type ValidationStatistics struct {
	ProvidersByMode   map[string]int          `json:"providers_by_mode"`
	ProvidersByStatus map[string]int          `json:"providers_by_status"`
	RuleViolations    map[string]int          `json:"rule_violations"`
	ErrorByField      map[string]int          `json:"error_by_field"`
	CommonErrors      []ValidationErrorDetail `json:"common_errors"`
}

// ValidationErrorDetail represents a detailed validation error with count
type ValidationErrorDetail struct {
	ErrorKey string `json:"error_key"`
	Count    int    `json:"count"`
}

// GoValidator implements the Validator interface
type GoValidator struct {
	rules []ValidationRule
}

// NewGoValidator creates a new GoValidator with default rules
func NewGoValidator() *GoValidator {
	validator := &GoValidator{
		rules: make([]ValidationRule, 0),
	}

	// Add default validation rules
	validator.AddDefaultRules()

	return validator
}

// Validate implements Validator.Validate
func (v *GoValidator) Validate(provider *Provider) error {
	var errors []ValidationError

	for _, rule := range v.rules {
		if err := rule.Validate(provider); err != nil {
			errors = append(errors, *err)
		}
	}

	if len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}

	return nil
}

// ValidateComment implements Validator.ValidateComment
func (v *GoValidator) ValidateComment(comment string) error {
	if !strings.HasPrefix(comment, "@provider") {
		return &ValidationError{
			RuleName: "CommentFormat",
			Message:  "provider comment must start with @provider",
			Severity: "error",
		}
	}

	// Parse the comment to check its structure
	providerDoc := parseProvider(comment)

	// Validate the parsed structure
	if providerDoc.Mode != "" && !IsValidProviderMode(providerDoc.Mode) {
		return &ValidationError{
			RuleName: "ProviderMode",
			Message:  fmt.Sprintf("invalid provider mode: %s", providerDoc.Mode),
			Severity: "error",
		}
	}

	return nil
}

// AddRule implements Validator.AddRule
func (v *GoValidator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// RemoveRule implements Validator.RemoveRule
func (v *GoValidator) RemoveRule(name string) {
	for i, rule := range v.rules {
		if rule.Name() == name {
			v.rules = append(v.rules[:i], v.rules[i+1:]...)
			break
		}
	}
}

// GetRules implements Validator.GetRules
func (v *GoValidator) GetRules() []ValidationRule {
	return v.rules
}

// ValidateAll implements Validator.ValidateAll
func (v *GoValidator) ValidateAll(providers []*Provider) *ValidationReport {
	report := &ValidationReport{
		Timestamp:      time.Now(),
		TotalProviders: len(providers),
		ValidCount:     0,
		InvalidCount:   0,
		Errors:         make([]ValidationError, 0),
		Warnings:       make([]ValidationError, 0),
		Infos:          make([]ValidationError, 0),
	}

	for _, provider := range providers {
		providerRef := fmt.Sprintf("%s:%s", provider.PkgName, provider.StructName)

		providerValid := true
		for _, rule := range v.rules {
			if err := rule.Validate(provider); err != nil {
				err.ProviderRef = providerRef

				switch err.Severity {
				case "error":
					report.Errors = append(report.Errors, *err)
					providerValid = false
				case "warning":
					report.Warnings = append(report.Warnings, *err)
				case "info":
					report.Infos = append(report.Infos, *err)
				}
			}
		}

		if providerValid {
			report.ValidCount++
		} else {
			report.InvalidCount++
		}
	}

	report.IsValid = len(report.Errors) == 0

	return report
}

// ValidateWithDetails validates all providers and returns a detailed report with statistics
func (v *GoValidator) ValidateWithDetails(providers []*Provider) *ValidationReport {
	report := v.ValidateAll(providers)

	// Add detailed statistics
	stats := calculateValidationStatistics(report, providers)
	report.Statistics = &stats

	return report
}

// calculateValidationStatistics calculates detailed validation statistics
func calculateValidationStatistics(report *ValidationReport, providers []*Provider) ValidationStatistics {
	stats := ValidationStatistics{
		ProvidersByMode:   make(map[string]int),
		ProvidersByStatus: make(map[string]int),
		RuleViolations:    make(map[string]int),
		ErrorByField:      make(map[string]int),
		CommonErrors:      make([]ValidationErrorDetail, 0),
	}

	// Count providers by mode
	for _, provider := range providers {
		mode := string(provider.Mode)
		stats.ProvidersByMode[mode]++
	}

	// Count violations by rule
	allIssues := append(append(report.Errors, report.Warnings...), report.Infos...)
	for _, issue := range allIssues {
		stats.RuleViolations[issue.RuleName]++
		stats.ErrorByField[issue.Field]++
	}

	// Find common errors
	errorCounts := make(map[string]int)
	for _, err := range report.Errors {
		key := fmt.Sprintf("%s:%s", err.RuleName, err.Message)
		errorCounts[key]++
	}

	// Get top 5 common errors
	for key, count := range errorCounts {
		stats.CommonErrors = append(stats.CommonErrors, ValidationErrorDetail{
			ErrorKey: key,
			Count:    count,
		})
	}

	// Sort by count (descending)
	for i := 0; i < len(stats.CommonErrors)-1; i++ {
		for j := i + 1; j < len(stats.CommonErrors); j++ {
			if stats.CommonErrors[i].Count < stats.CommonErrors[j].Count {
				stats.CommonErrors[i], stats.CommonErrors[j] = stats.CommonErrors[j], stats.CommonErrors[i]
			}
		}
	}

	// Keep only top 5
	if len(stats.CommonErrors) > 5 {
		stats.CommonErrors = stats.CommonErrors[:5]
	}

	return stats
}

// GenerateReport generates a formatted validation report
func (v *GoValidator) GenerateReport(providers []*Provider, format ReportFormat) (string, error) {
	report := v.ValidateWithDetails(providers)
	generator := NewReportGenerator(report)
	return generator.GenerateReport(format)
}

// SaveReport saves a validation report to a file
func (v *GoValidator) SaveReport(providers []*Provider, filename string, format ReportFormat) error {
	report := v.ValidateWithDetails(providers)
	writer := NewReportWriter(report)
	return writer.WriteToFile(filename, format)
}

// PrintReport prints a validation report to the console
func (v *GoValidator) PrintReport(providers []*Provider, format ReportFormat) error {
	report := v.ValidateWithDetails(providers)
	writer := NewReportWriter(report)
	return writer.WriteToConsole(format)
}

// AddDefaultRules adds the default validation rules
func (v *GoValidator) AddDefaultRules() {
	v.AddRule(&StructNameRule{})
	v.AddRule(&ReturnTypeRule{})
	v.AddRule(&ProviderModeRule{})
	v.AddRule(&InjectionParamsRule{})
	v.AddRule(&PackageAliasRule{})
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError
}

// Error implements the error interface
func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "no validation errors"
	}

	if len(e.Errors) == 1 {
		return e.Errors[0].Message
	}

	return fmt.Sprintf("%d validation errors occurred, first: %s", len(e.Errors), e.Errors[0].Message)
}

// StructNameRule validates struct names
type StructNameRule struct{}

func (r *StructNameRule) Name() string { return "StructName" }

func (r *StructNameRule) Description() string {
	return "Validates that struct names follow Go naming conventions"
}

func (r *StructNameRule) Validate(provider *Provider) *ValidationError {
	if provider.StructName == "" {
		return &ValidationError{
			RuleName: r.Name(),
			Message:  "struct name cannot be empty",
			Field:    "StructName",
			Severity: "error",
		}
	}

	if !isExportedName(provider.StructName) {
		return &ValidationError{
			RuleName: r.Name(),
			Message:  "struct name must be exported (start with uppercase letter)",
			Field:    "StructName",
			Value:    provider.StructName,
			Severity: "error",
		}
	}

	// Check for Go naming conventions
	if !isValidGoIdentifier(provider.StructName) {
		return &ValidationError{
			RuleName: r.Name(),
			Message:  "struct name must be a valid Go identifier",
			Field:    "StructName",
			Value:    provider.StructName,
			Severity: "error",
		}
	}

	// Check for common naming conventions (CamelCase for structs)
	if !isCamelCase(provider.StructName) {
		return &ValidationError{
			RuleName:   r.Name(),
			Message:    "struct name should use CamelCase convention",
			Field:      "StructName",
			Value:      provider.StructName,
			Severity:   "warning",
			Suggestion: "consider using CamelCase (e.g., UserService instead of userService or USER_SERVICE)",
		}
	}

	// Check for reserved words or problematic names
	if isReservedWord(provider.StructName) {
		return &ValidationError{
			RuleName:   r.Name(),
			Message:    "struct name should not use reserved words or common Go types",
			Field:      "StructName",
			Value:      provider.StructName,
			Severity:   "warning",
			Suggestion: "choose a more descriptive name that doesn't conflict with Go built-ins",
		}
	}

	return nil
}

// ReturnTypeRule validates return types
type ReturnTypeRule struct{}

func (r *ReturnTypeRule) Name() string { return "ReturnType" }

func (r *ReturnTypeRule) Description() string {
	return "Validates that return types are properly formatted and consistent"
}

func (r *ReturnTypeRule) Validate(provider *Provider) *ValidationError {
	if provider.ReturnType == "" {
		return &ValidationError{
			RuleName: r.Name(),
			Message:  "return type cannot be empty",
			Field:    "ReturnType",
			Severity: "error",
		}
	}

	// Check if return type is a valid Go type identifier
	if !isValidGoType(provider.ReturnType) {
		return &ValidationError{
			RuleName: r.Name(),
			Message:  fmt.Sprintf("invalid return type format: %s", provider.ReturnType),
			Field:    "ReturnType",
			Value:    provider.ReturnType,
			Severity: "error",
		}
	}

	// Check for mode-specific return type requirements
	if err := validateModeSpecificReturnType(provider); err != nil {
		return err
	}

	// Check for interface types (should generally be avoided for providers)
	if strings.HasPrefix(provider.ReturnType, "*") &&
		(strings.HasSuffix(provider.ReturnType, "Interface") ||
			strings.HasSuffix(provider.ReturnType, "IFace")) {
		return &ValidationError{
			RuleName:   r.Name(),
			Message:    "pointer to interface types should generally be avoided",
			Field:      "ReturnType",
			Value:      provider.ReturnType,
			Severity:   "warning",
			Suggestion: "consider using the interface type directly or a concrete type",
		}
	}

	// Check for overly complex generic types
	if strings.Count(provider.ReturnType, "[") > 2 || strings.Count(provider.ReturnType, "]") > 2 {
		return &ValidationError{
			RuleName:   r.Name(),
			Message:    "return type appears overly complex with too many generic parameters",
			Field:      "ReturnType",
			Value:      provider.ReturnType,
			Severity:   "warning",
			Suggestion: "consider simplifying the type or using type aliases",
		}
	}

	// Check for pointer vs value type consistency
	if strings.HasPrefix(provider.ReturnType, "*") {
		// Pointer type validations
		pointedType := strings.TrimPrefix(provider.ReturnType, "*")
		if !isExportedName(pointedType) {
			return &ValidationError{
				RuleName: r.Name(),
				Message:  "pointed type in pointer return type must be exported",
				Field:    "ReturnType",
				Value:    provider.ReturnType,
				Severity: "error",
			}
		}
	} else {
		// Value type validations
		if !isExportedName(provider.ReturnType) && !isBuiltInType(provider.ReturnType) {
			return &ValidationError{
				RuleName: r.Name(),
				Message:  "non-pointer return type must be either exported or a built-in type",
				Field:    "ReturnType",
				Value:    provider.ReturnType,
				Severity: "error",
			}
		}
	}

	return nil
}

// validateModeSpecificReturnType checks return type requirements for specific provider modes
func validateModeSpecificReturnType(provider *Provider) *ValidationError {
	switch provider.Mode {
	case ProviderModeGrpc:
		// gRPC providers should typically return interface types
		if !strings.HasSuffix(provider.ReturnType, "Client") &&
			!strings.HasSuffix(provider.ReturnType, "Service") {
			return &ValidationError{
				RuleName:   "ReturnType",
				Message:    "gRPC providers should typically return client or service interface types",
				Field:      "ReturnType",
				Value:      provider.ReturnType,
				Severity:   "warning",
				Suggestion: "consider naming the return type ending with Client or Service",
			}
		}

	case ProviderModeModel:
		// Model providers should return struct types
		if strings.Contains(provider.ReturnType, "interface") ||
			strings.Contains(provider.ReturnType, "Interface") {
			return &ValidationError{
				RuleName: "ReturnType",
				Message:  "model providers should return concrete struct types, not interfaces",
				Field:    "ReturnType",
				Value:    provider.ReturnType,
				Severity: "error",
			}
		}

	case ProviderModeJob, ProviderModeCronJob:
		// Job providers should return types that implement job interfaces
		if !strings.Contains(provider.ReturnType, "Job") {
			return &ValidationError{
				RuleName:   "ReturnType",
				Message:    "job providers should return types that implement job interfaces",
				Field:      "ReturnType",
				Value:      provider.ReturnType,
				Severity:   "warning",
				Suggestion: "consider naming the return type to indicate it's a job handler",
			}
		}
	}

	return nil
}

// isBuiltInType checks if a type is a Go built-in type
func isBuiltInType(typeStr string) bool {
	builtInTypes := map[string]bool{
		"string": true, "int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"float32": true, "float64": true, "complex64": true, "complex128": true,
		"bool": true, "byte": true, "rune": true, "error": true,
	}

	return builtInTypes[typeStr]
}

// ProviderModeRule validates provider modes
type ProviderModeRule struct{}

func (r *ProviderModeRule) Name() string { return "ProviderMode" }

func (r *ProviderModeRule) Description() string {
	return "Validates that provider modes are valid and appropriate for the provider configuration"
}

func (r *ProviderModeRule) Validate(provider *Provider) *ValidationError {
	if !IsValidProviderMode(string(provider.Mode)) {
		return &ValidationError{
			RuleName: r.Name(),
			Message:  fmt.Sprintf("invalid provider mode: %s", provider.Mode),
			Field:    "Mode",
			Value:    string(provider.Mode),
			Severity: "error",
		}
	}

	// Validate mode-specific configurations
	if err := validateModeSpecificConfiguration(provider); err != nil {
		return err
	}

	// Check for deprecated or discouraged mode combinations
	if err := validateModeCombinations(provider); err != nil {
		return err
	}

	return nil
}

// validateModeSpecificConfiguration checks mode-specific configuration requirements
func validateModeSpecificConfiguration(provider *Provider) *ValidationError {
	switch provider.Mode {
	case ProviderModeGrpc:
		// gRPC providers need gRPC register function
		if provider.GrpcRegisterFunc == "" {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "gRPC providers must specify a gRPC register function",
				Field:      "GrpcRegisterFunc",
				Value:      provider.GrpcRegisterFunc,
				Severity:   "error",
				Suggestion: "add gRPC register function to the provider annotation",
			}
		}

	case ProviderModeJob, ProviderModeCronJob:
		// Job providers should have prepare function
		if !provider.NeedPrepareFunc {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "job providers should typically need a prepare function",
				Field:      "NeedPrepareFunc",
				Value:      fmt.Sprintf("%t", provider.NeedPrepareFunc),
				Severity:   "warning",
				Suggestion: "consider setting prepare function for job initialization",
			}
		}

		// Check if job provider has proper injection parameters
		hasJobParam := false
		for paramName := range provider.InjectParams {
			if paramName == "__job" {
				hasJobParam = true
				break
			}
		}
		if !hasJobParam {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "job providers should inject __job parameter",
				Field:      "InjectParams",
				Severity:   "warning",
				Suggestion: "add __job *job.Job parameter to injection parameters",
			}
		}

	case ProviderModeEvent:
		// Event providers should have __event parameter
		hasEventParam := false
		for paramName := range provider.InjectParams {
			if paramName == "__event" {
				hasEventParam = true
				break
			}
		}
		if !hasEventParam {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "event providers should inject __event parameter",
				Field:      "InjectParams",
				Severity:   "warning",
				Suggestion: "add __event parameter to injection parameters",
			}
		}

	case ProviderModeBasic:
		// Basic providers should not have special parameters
		for paramName := range provider.InjectParams {
			if paramName == "__job" || paramName == "__event" || paramName == "__grpc" {
				return &ValidationError{
					RuleName: "ProviderMode",
					Message:  "basic providers should not inject special parameters (__job, __event, __grpc)",
					Field:    "InjectParams",
					Value:    paramName,
					Severity: "error",
				}
			}
		}
	}

	return nil
}

// validateModeCombinations checks for problematic mode combinations
func validateModeCombinations(provider *Provider) *ValidationError {
	// Check for incompatible mode combinations
	if provider.ProviderGroup != "" {
		groupLower := strings.ToLower(provider.ProviderGroup)
		modeStr := string(provider.Mode)

		// Basic checks for common issues
		if (modeStr == "grpc" && groupLower == "job") ||
			(modeStr == "job" && groupLower == "grpc") {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "incompatible combination of provider mode and group",
				Field:      "ProviderGroup",
				Value:      provider.ProviderGroup,
				Severity:   "warning",
				Suggestion: "ensure provider mode and group are compatible",
			}
		}
	}

	// Check for mode-specific naming suggestions
	if err := validateModeNamingConventions(provider); err != nil {
		return err
	}

	return nil
}

// validateModeNamingConventions provides naming suggestions for different modes
func validateModeNamingConventions(provider *Provider) *ValidationError {
	structName := provider.StructName

	switch provider.Mode {
	case ProviderModeGrpc:
		if !strings.HasSuffix(structName, "Client") &&
			!strings.HasSuffix(structName, "Service") &&
			!strings.HasSuffix(structName, "Provider") {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "gRPC provider struct names should typically end with Client, Service, or Provider",
				Field:      "StructName",
				Value:      structName,
				Severity:   "warning",
				Suggestion: "consider naming the struct to indicate it's a gRPC provider",
			}
		}

	case ProviderModeJob, ProviderModeCronJob:
		if !strings.HasSuffix(structName, "Job") &&
			!strings.HasSuffix(structName, "Worker") &&
			!strings.HasSuffix(structName, "Handler") {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "job provider struct names should typically end with Job, Worker, or Handler",
				Field:      "StructName",
				Value:      structName,
				Severity:   "warning",
				Suggestion: "consider naming the struct to indicate it's a job provider",
			}
		}

	case ProviderModeEvent:
		if !strings.HasSuffix(structName, "Handler") &&
			!strings.HasSuffix(structName, "Listener") &&
			!strings.HasSuffix(structName, "Subscriber") {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "event provider struct names should typically end with Handler, Listener, or Subscriber",
				Field:      "StructName",
				Value:      structName,
				Severity:   "warning",
				Suggestion: "consider naming the struct to indicate it's an event provider",
			}
		}

	case ProviderModeModel:
		if !strings.HasSuffix(structName, "Model") &&
			!strings.HasSuffix(structName, "Entity") &&
			!strings.HasSuffix(structName, "Repo") &&
			!strings.HasSuffix(structName, "Repository") {
			return &ValidationError{
				RuleName:   "ProviderMode",
				Message:    "model provider struct names should typically end with Model, Entity, Repo, or Repository",
				Field:      "StructName",
				Value:      structName,
				Severity:   "warning",
				Suggestion: "consider naming the struct to indicate it's a model provider",
			}
		}
	}

	return nil
}

// InjectionParamsRule validates injection parameters
type InjectionParamsRule struct{}

func (r *InjectionParamsRule) Name() string { return "InjectionParams" }

func (r *InjectionParamsRule) Description() string {
	return "Validates injection parameters for consistency and correctness"
}

func (r *InjectionParamsRule) Validate(provider *Provider) *ValidationError {
	if len(provider.InjectParams) == 0 {
		// Some providers might not need injection parameters
		if provider.Mode == ProviderModeBasic {
			return nil // Basic providers can have no injection params
		}

		return &ValidationError{
			RuleName:   r.Name(),
			Message:    "providers should have at least one injection parameter",
			Field:      "InjectParams",
			Severity:   "warning",
			Suggestion: "consider adding injection parameters for dependency injection",
		}
	}

	// Check for duplicate parameter names
	paramNames := make(map[string]bool)
	for paramName := range provider.InjectParams {
		if paramNames[paramName] {
			return &ValidationError{
				RuleName: r.Name(),
				Message:  "duplicate injection parameter name",
				Field:    "InjectParams",
				Value:    paramName,
				Severity: "error",
			}
		}
		paramNames[paramName] = true
	}

	// Validate each injection parameter
	for paramName, param := range provider.InjectParams {
		if err := validateInjectionParameter(paramName, param, provider.Mode); err != nil {
			return err
		}
	}

	// Check for parameter consistency
	if err := validateParameterConsistency(provider); err != nil {
		return err
	}

	return nil
}

// validateInjectionParameter validates a single injection parameter
func validateInjectionParameter(paramName string, param InjectParam, mode ProviderMode) *ValidationError {
	// Validate parameter name
	if paramName == "" {
		return &ValidationError{
			RuleName: "InjectionParams",
			Message:  "injection parameter name cannot be empty",
			Field:    "InjectParams",
			Severity: "error",
		}
	}

	// Special parameters are allowed to be unexported
	if !isSpecialParameter(paramName) && !isExportedName(paramName) {
		return &ValidationError{
			RuleName:   "InjectionParams",
			Message:    "injection parameter names must be exported",
			Field:      "InjectParams",
			Value:      paramName,
			Severity:   "warning",
			Suggestion: "use exported names for injection parameters",
		}
	}

	// Validate parameter type
	if param.Type == "" {
		return &ValidationError{
			RuleName: "InjectionParams",
			Message:  "injection parameter type cannot be empty",
			Field:    "InjectParams",
			Value:    paramName,
			Severity: "error",
		}
	}

	// Validate type format
	if !isValidGoType(param.Type) {
		return &ValidationError{
			RuleName: "InjectionParams",
			Message:  fmt.Sprintf("invalid injection parameter type format: %s", param.Type),
			Field:    "InjectParams",
			Value:    paramName,
			Severity: "error",
		}
	}

	// Validate package and alias consistency
	if param.Package != "" && param.PackageAlias == "" {
		return &ValidationError{
			RuleName: "InjectionParams",
			Message:  "package alias is required when package is specified",
			Field:    "InjectParams",
			Value:    paramName,
			Severity: "error",
		}
	}

	if param.Package == "" && param.PackageAlias != "" {
		return &ValidationError{
			RuleName: "InjectionParams",
			Message:  "package cannot be empty when package alias is specified",
			Field:    "InjectParams",
			Value:    paramName,
			Severity: "error",
		}
	}

	// Validate special parameters
	if isSpecialParameter(paramName) {
		if err := validateSpecialParameter(paramName, param, mode); err != nil {
			return err
		}
	}

	return nil
}

// validateParameterConsistency checks for consistency across all parameters
func validateParameterConsistency(provider *Provider) *ValidationError {
	// Check for conflicting parameter types
	paramTypes := make(map[string]string)
	for paramName, param := range provider.InjectParams {
		// Check if same type is used with different aliases
		if existingAlias, exists := paramTypes[param.Type]; exists {
			// Extract package alias from existing parameter
			existingParam := provider.InjectParams[existingAlias]
			if existingParam.PackageAlias != param.PackageAlias {
				return &ValidationError{
					RuleName:   "InjectionParams",
					Message:    "same type used with different package aliases",
					Field:      "InjectParams",
					Value:      fmt.Sprintf("%s (%s vs %s)", param.Type, existingParam.PackageAlias, param.PackageAlias),
					Severity:   "warning",
					Suggestion: "use consistent package aliases for the same type",
				}
			}
		}
		paramTypes[param.Type] = paramName
	}

	// Check for circular dependencies (simplified check)
	for paramName, param := range provider.InjectParams {
		if param.Type == provider.StructName {
			return &ValidationError{
				RuleName: "InjectionParams",
				Message:  "provider cannot inject itself",
				Field:    "InjectParams",
				Value:    paramName,
				Severity: "error",
			}
		}
	}

	// Check for mode-specific parameter requirements
	if err := validateModeSpecificParameters(provider); err != nil {
		return err
	}

	return nil
}

// validateModeSpecificParameters checks mode-specific parameter requirements
func validateModeSpecificParameters(provider *Provider) *ValidationError {
	switch provider.Mode {
	case ProviderModeJob, ProviderModeCronJob:
		// Job providers should have __job parameter
		if _, hasJobParam := provider.InjectParams["__job"]; !hasJobParam {
			return &ValidationError{
				RuleName:   "InjectionParams",
				Message:    "job providers should inject __job parameter",
				Field:      "InjectParams",
				Severity:   "warning",
				Suggestion: "add __job *job.Job parameter for job context",
			}
		}

	case ProviderModeEvent:
		// Event providers should have __event parameter
		if _, hasEventParam := provider.InjectParams["__event"]; !hasEventParam {
			return &ValidationError{
				RuleName:   "InjectionParams",
				Message:    "event providers should inject __event parameter",
				Field:      "InjectParams",
				Severity:   "warning",
				Suggestion: "add __event parameter for event context",
			}
		}

	case ProviderModeGrpc:
		// gRPC providers should have __grpc parameter
		if _, hasGrpcParam := provider.InjectParams["__grpc"]; !hasGrpcParam {
			return &ValidationError{
				RuleName:   "InjectionParams",
				Message:    "gRPC providers should inject __grpc parameter",
				Field:      "InjectParams",
				Severity:   "warning",
				Suggestion: "add __grpc parameter for gRPC context",
			}
		}
	}

	return nil
}

// validateSpecialParameter validates special parameters like __job, __event, etc.
func validateSpecialParameter(paramName string, param InjectParam, mode ProviderMode) *ValidationError {
	switch paramName {
	case "__job":
		if mode != ProviderModeJob && mode != ProviderModeCronJob {
			return &ValidationError{
				RuleName: "InjectionParams",
				Message:  "__job parameter should only be used in job providers",
				Field:    "InjectParams",
				Value:    paramName,
				Severity: "error",
			}
		}
		if param.Type != "*job.Job" {
			return &ValidationError{
				RuleName: "InjectionParams",
				Message:  "__job parameter should have type *job.Job",
				Field:    "InjectParams",
				Value:    param.Type,
				Severity: "error",
			}
		}

	case "__event":
		if mode != ProviderModeEvent {
			return &ValidationError{
				RuleName: "InjectionParams",
				Message:  "__event parameter should only be used in event providers",
				Field:    "InjectParams",
				Value:    paramName,
				Severity: "error",
			}
		}
		if !strings.Contains(param.Type, "Event") {
			return &ValidationError{
				RuleName:   "InjectionParams",
				Message:    "__event parameter should have an event type",
				Field:      "InjectParams",
				Value:      param.Type,
				Severity:   "warning",
				Suggestion: "use a type that indicates it's an event",
			}
		}

	case "__grpc":
		if mode != ProviderModeGrpc {
			return &ValidationError{
				RuleName: "InjectionParams",
				Message:  "__grpc parameter should only be used in gRPC providers",
				Field:    "InjectParams",
				Value:    paramName,
				Severity: "error",
			}
		}
		if !strings.Contains(param.Type, "grpc") {
			return &ValidationError{
				RuleName:   "InjectionParams",
				Message:    "__grpc parameter should have a gRPC-related type",
				Field:      "InjectParams",
				Value:      param.Type,
				Severity:   "warning",
				Suggestion: "use a type that indicates it's gRPC-related",
			}
		}
	}

	return nil
}

// isSpecialParameter checks if a parameter name is a special parameter
func isSpecialParameter(paramName string) bool {
	specialParams := map[string]bool{
		"__job":   true,
		"__event": true,
		"__grpc":  true,
	}
	return specialParams[paramName]
}

// PackageAliasRule validates package aliases
type PackageAliasRule struct{}

func (r *PackageAliasRule) Name() string { return "PackageAlias" }

func (r *PackageAliasRule) Description() string {
	return "Validates package aliases for consistency"
}

func (r *PackageAliasRule) Validate(provider *Provider) *ValidationError {
	for alias, path := range provider.Imports {
		if alias == "" {
			return &ValidationError{
				RuleName: r.Name(),
				Message:  "package alias cannot be empty",
				Field:    "Imports",
				Value:    path,
				Severity: "error",
			}
		}

		if path == "" {
			return &ValidationError{
				RuleName: r.Name(),
				Message:  "package path cannot be empty",
				Field:    "Imports",
				Value:    alias,
				Severity: "error",
			}
		}

		if !isValidGoIdentifier(alias) {
			return &ValidationError{
				RuleName: r.Name(),
				Message:  "package alias must be a valid Go identifier",
				Field:    "Imports",
				Value:    alias,
				Severity: "error",
			}
		}
	}

	return nil
}

// Helper functions

func isExportedName(name string) bool {
	if name == "" {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}

func isValidGoType(typeStr string) bool {
	// Enhanced validation for Go type identifiers including pointers
	return regexp.MustCompile(`^(\*[a-zA-Z_][a-zA-Z0-9_.]*|[a-zA-Z_][a-zA-Z0-9_.]*(\[\])?|[a-zA-Z_][a-zA-Z0-9_.]*(\[\])?\*?)$`).MatchString(typeStr)
}

func isValidGoIdentifier(name string) bool {
	if name == "" {
		return false
	}

	// Check if it's a valid Go identifier
	for i, r := range name {
		if i == 0 && !unicode.IsLetter(r) && r != '_' {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return true
}

func isCamelCase(name string) bool {
	if !isValidGoIdentifier(name) {
		return false
	}

	// Check if first character is uppercase (exported)
	if !unicode.IsUpper(rune(name[0])) {
		return false
	}

	// Check for snake_case or kebab-case patterns
	if strings.Contains(name, "_") || strings.Contains(name, "-") {
		return false
	}

	// Check for ALL_CAPS (usually used for constants)
	if strings.ToUpper(name) == name && len(name) > 1 {
		return false
	}

	return true
}

func isReservedWord(name string) bool {
	reservedWords := map[string]bool{
		// Go keywords
		"break": true, "case": true, "chan": true, "const": true, "continue": true,
		"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
		"func": true, "go": true, "goto": true, "if": true, "import": true,
		"interface": true, "map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true, "var": true,

		// Predeclared identifiers
		"bool": true, "byte": true, "complex64": true, "complex128": true,
		"error": true, "float32": true, "float64": true, "int": true, "int8": true,
		"int16": true, "int32": true, "int64": true, "rune": true, "string": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"uintptr": true, "true": true, "false": true, "nil": true, "iota": true,

		// Common problematic names for providers
		"Provider": true, "Service": true, "Handler": true, "Manager": true,
		"Controller": true, "Repository": true, "Config": true, "Client": true,
	}

	return reservedWords[name]
}

func isValidPackageName(name string) bool {
	if name == "" {
		return false
	}

	// Package names should be lowercase, short, and descriptive
	if !isValidGoIdentifier(name) {
		return false
	}

	// Package names should typically be lowercase
	if name != strings.ToLower(name) {
		return false
	}

	// Avoid common problematic package names
	problematicNames := map[string]bool{
		"main": true, "testing": true, "fmt": true, "strings": true,
		"container": true, "service": true, "utils": true, "common": true,
	}

	return !problematicNames[name]
}
