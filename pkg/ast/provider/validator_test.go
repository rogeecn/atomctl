package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoValidator_Basic(t *testing.T) {
	validator := NewGoValidator()

	// Test with a valid provider
	validProvider := &Provider{
		StructName:   "UserService",
		ReturnType:   "UserService",
		Mode:         ProviderModeBasic,
		PkgName:      "services",
		ProviderFile: "providers.gen.go",
		InjectParams: map[string]InjectParam{
			"DB": {
				Star:         "*",
				Type:         "DB",
				Package:      "database/sql",
				PackageAlias: "sql",
			},
		},
		Imports: map[string]string{
			"sql": "database/sql",
		},
		Location: SourceLocation{
			File:   "user_service.go",
			Line:   10,
			Column: 1,
		},
	}

	err := validator.Validate(validProvider)
	assert.NoError(t, err, "Valid provider should not return error")

	// Test with an invalid provider
	invalidProvider := &Provider{
		StructName:   "", // Empty struct name
		ReturnType:   "",
		Mode:         "invalid-mode",
		PkgName:      "",
		ProviderFile: "",
		InjectParams: map[string]InjectParam{},
		Imports:      map[string]string{},
	}

	err = validator.Validate(invalidProvider)
	assert.Error(t, err, "Invalid provider should return error")
}

func TestGoValidator_ValidateAll(t *testing.T) {
	validator := NewGoValidator()

	providers := []*Provider{
		{
			StructName:   "UserService",
			ReturnType:   "UserService",
			Mode:         ProviderModeBasic,
			PkgName:      "services",
			ProviderFile: "providers.gen.go",
			InjectParams: map[string]InjectParam{
				"DB": {
					Star:         "*",
					Type:         "DB",
					Package:      "database/sql",
					PackageAlias: "sql",
				},
			},
			Imports: map[string]string{
				"sql": "database/sql",
			},
		},
		{
			StructName:   "JobProcessor",
			ReturnType:   "JobProcessor",
			Mode:         ProviderModeJob,
			PkgName:      "jobs",
			ProviderFile: "providers.gen.go",
			InjectParams: map[string]InjectParam{
				"__job": {
					Type: "Job",
				},
			},
			Imports: map[string]string{},
		},
		{
			StructName:   "", // Invalid
			ReturnType:   "",
			Mode:         "invalid-mode",
			PkgName:      "",
			ProviderFile: "",
			InjectParams: map[string]InjectParam{},
			Imports:      map[string]string{},
		},
	}

	report := validator.ValidateAll(providers)
	assert.NotNil(t, report)
	assert.Equal(t, 3, report.TotalProviders)
	assert.Equal(t, 1, report.ValidCount)
	assert.Equal(t, 2, report.InvalidCount)
	assert.False(t, report.IsValid)
	assert.Len(t, report.Errors, 4)
}

func TestGoValidator_ReportGeneration(t *testing.T) {
	validator := NewGoValidator()

	providers := []*Provider{
		{
			StructName:   "UserService",
			ReturnType:   "UserService",
			Mode:         ProviderModeBasic,
			PkgName:      "services",
			ProviderFile: "providers.gen.go",
			InjectParams: map[string]InjectParam{
				"DB": {
					Star:         "*",
					Type:         "DB",
					Package:      "database/sql",
					PackageAlias: "sql",
				},
			},
			Imports: map[string]string{
				"sql": "database/sql",
			},
		},
	}

	// Test text report generation
	report, err := validator.GenerateReport(providers, ReportFormatText)
	assert.NoError(t, err)
	assert.Contains(t, report, "Provider Validation Report")
	assert.Contains(t, report, "Total Providers: 1")
	assert.Contains(t, report, "Valid Providers: 1")

	// Test JSON report generation
	jsonReport, err := validator.GenerateReport(providers, ReportFormatJSON)
	assert.NoError(t, err)
	assert.Contains(t, jsonReport, `"total_providers": 1`)
	assert.Contains(t, jsonReport, `"valid_count": 1`)

	// Test Markdown report generation
	markdownReport, err := validator.GenerateReport(providers, ReportFormatMarkdown)
	assert.NoError(t, err)
	assert.Contains(t, markdownReport, "# Provider Validation Report")
	assert.Contains(t, markdownReport, "**Total Providers:** 1")
}

func TestValidationReport_Statistics(t *testing.T) {
	validator := NewGoValidator()

	providers := []*Provider{
		{
			StructName:   "UserService",
			ReturnType:   "UserService",
			Mode:         ProviderModeBasic,
			PkgName:      "services",
			ProviderFile: "providers.gen.go",
			InjectParams: map[string]InjectParam{
				"DB": {
					Star:         "*",
					Type:         "DB",
					Package:      "database/sql",
					PackageAlias: "sql",
				},
			},
			Imports: map[string]string{
				"sql": "database/sql",
			},
		},
		{
			StructName:   "JobProcessor",
			ReturnType:   "JobProcessor",
			Mode:         ProviderModeJob,
			PkgName:      "jobs",
			ProviderFile: "providers.gen.go",
			InjectParams: map[string]InjectParam{
				"__job": {
					Type: "Job",
				},
			},
			Imports: map[string]string{},
		},
	}

	report := validator.ValidateWithDetails(providers)
	assert.NotNil(t, report.Statistics)
	assert.Len(t, report.Statistics.ProvidersByMode, 2)
	assert.Contains(t, report.Statistics.ProvidersByMode, "basic")
	assert.Contains(t, report.Statistics.ProvidersByMode, "job")
}

func TestStructNameRule(t *testing.T) {
	rule := &StructNameRule{}

	tests := []struct {
		name        string
		provider    *Provider
		expectError bool
	}{
		{
			name: "Valid struct name",
			provider: &Provider{
				StructName: "UserService",
			},
			expectError: false,
		},
		{
			name: "Empty struct name",
			provider: &Provider{
				StructName: "",
			},
			expectError: true,
		},
		{
			name: "Unexported struct name",
			provider: &Provider{
				StructName: "userService",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule.Validate(tt.provider)
			if tt.expectError {
				assert.NotNil(t, err, "Expected validation error but got nil")
			} else {
				assert.Nil(t, err, "Expected no validation error but got one")
			}
		})
	}
}

func TestInjectionParamsRule(t *testing.T) {
	rule := &InjectionParamsRule{}

	tests := []struct {
		name        string
		provider    *Provider
		expectError bool
	}{
		{
			name: "Valid injection params",
			provider: &Provider{
				StructName: "UserService",
				Mode:       ProviderModeBasic,
				InjectParams: map[string]InjectParam{
					"DB": {
						Star:         "*",
						Type:         "DB",
						Package:      "database/sql",
						PackageAlias: "sql",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Empty parameter type",
			provider: &Provider{
				StructName: "UserService",
				Mode:       ProviderModeBasic,
				InjectParams: map[string]InjectParam{
					"DB": {
						Star:         "*",
						Type:         "",
						Package:      "database/sql",
						PackageAlias: "sql",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Missing package alias",
			provider: &Provider{
				StructName: "UserService",
				Mode:       ProviderModeBasic,
				InjectParams: map[string]InjectParam{
					"DB": {
						Star:         "*",
						Type:         "sql.DB",
						Package:      "database/sql",
						PackageAlias: "",
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rule.Validate(tt.provider)
			if tt.expectError {
				assert.NotNil(t, err, "Expected validation error but got nil")
			} else {
				assert.Nil(t, err, "Expected no validation error but got one")
			}
		})
	}
}
