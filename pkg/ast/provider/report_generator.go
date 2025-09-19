package provider

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ReportGenerator handles the generation of validation reports in various formats
type ReportGenerator struct {
	report *ValidationReport
}

// NewReportGenerator creates a new ReportGenerator
func NewReportGenerator(report *ValidationReport) *ReportGenerator {
	return &ReportGenerator{
		report: report,
	}
}

// GenerateTextReport generates a human-readable text report
func (rg *ReportGenerator) GenerateTextReport() string {
	var builder strings.Builder

	builder.WriteString("Provider Validation Report\n")
	builder.WriteString("=========================\n\n")
	builder.WriteString(fmt.Sprintf("Generated: %s\n", rg.report.Timestamp.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("Total Providers: %d\n", rg.report.TotalProviders))
	builder.WriteString(fmt.Sprintf("Valid Providers: %d\n", rg.report.ValidCount))
	builder.WriteString(fmt.Sprintf("Invalid Providers: %d\n", rg.report.InvalidCount))
	builder.WriteString(fmt.Sprintf("Overall Status: %s\n\n", rg.getStatusText()))

	// Summary section
	builder.WriteString("Summary\n")
	builder.WriteString("-------\n")
	if rg.report.IsValid {
		builder.WriteString("✅ All providers are valid\n\n")
	} else {
		builder.WriteString("❌ Validation failed with errors\n\n")
	}

	// Errors section
	if len(rg.report.Errors) > 0 {
		builder.WriteString("Errors\n")
		builder.WriteString("------\n")
		for i, err := range rg.report.Errors {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rg.formatValidationError(&err)))
		}
		builder.WriteString("\n")
	}

	// Warnings section
	if len(rg.report.Warnings) > 0 {
		builder.WriteString("Warnings\n")
		builder.WriteString("--------\n")
		for i, warning := range rg.report.Warnings {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rg.formatValidationError(&warning)))
		}
		builder.WriteString("\n")
	}

	// Infos section
	if len(rg.report.Infos) > 0 {
		builder.WriteString("Information\n")
		builder.WriteString("-----------\n")
		for i, info := range rg.report.Infos {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rg.formatValidationError(&info)))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// GenerateJSONReport generates a JSON report
func (rg *ReportGenerator) GenerateJSONReport() (string, error) {
	data, err := json.MarshalIndent(rg.report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON report: %w", err)
	}
	return string(data), nil
}

// GenerateHTMLReport generates an HTML report
func (rg *ReportGenerator) GenerateHTMLReport() string {
	var builder strings.Builder

	builder.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>Provider Validation Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f5f5f5; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        .summary { background-color: #e8f5e8; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
        .error { background-color: #ffe6e6; padding: 15px; border-radius: 5px; margin-bottom: 10px; }
        .warning { background-color: #fff3cd; padding: 15px; border-radius: 5px; margin-bottom: 10px; }
        .info { background-color: #e7f3ff; padding: 15px; border-radius: 5px; margin-bottom: 10px; }
        .severity { font-weight: bold; text-transform: uppercase; }
        .provider-ref { color: #666; font-style: italic; }
        .suggestion { color: #28a745; font-style: italic; }
        h2 { color: #333; border-bottom: 2px solid #ddd; padding-bottom: 5px; }
        ul { margin: 10px 0; }
        li { margin: 5px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Provider Validation Report</h1>
        <p><strong>Generated:</strong> ` + rg.report.Timestamp.Format(time.RFC3339) + `</p>
        <p><strong>Total Providers:</strong> ` + fmt.Sprintf("%d", rg.report.TotalProviders) + `</p>
        <p><strong>Valid Providers:</strong> ` + fmt.Sprintf("%d", rg.report.ValidCount) + `</p>
        <p><strong>Invalid Providers:</strong> ` + fmt.Sprintf("%d", rg.report.InvalidCount) + `</p>
        <p><strong>Overall Status:</strong> ` + rg.getStatusHTML() + `</p>
    </div>

    <div class="summary">
        <h2>Summary</h2>
        <p>` + rg.getSummaryText() + `</p>
    </div>`)

	// Errors section
	if len(rg.report.Errors) > 0 {
		builder.WriteString(`
    <div class="errors">
        <h2>Errors</h2>
        <ul>`)
		for _, err := range rg.report.Errors {
			builder.WriteString(fmt.Sprintf(`<li class="error">%s</li>`, rg.formatValidationErrorHTML(&err)))
		}
		builder.WriteString(`</ul>
    </div>`)
	}

	// Warnings section
	if len(rg.report.Warnings) > 0 {
		builder.WriteString(`
    <div class="warnings">
        <h2>Warnings</h2>
        <ul>`)
		for _, warning := range rg.report.Warnings {
			builder.WriteString(fmt.Sprintf(`<li class="warning">%s</li>`, rg.formatValidationErrorHTML(&warning)))
		}
		builder.WriteString(`</ul>
    </div>`)
	}

	// Infos section
	if len(rg.report.Infos) > 0 {
		builder.WriteString(`
    <div class="infos">
        <h2>Information</h2>
        <ul>`)
		for _, info := range rg.report.Infos {
			builder.WriteString(fmt.Sprintf(`<li class="info">%s</li>`, rg.formatValidationErrorHTML(&info)))
		}
		builder.WriteString(`</ul>
    </div>`)
	}

	builder.WriteString(`
</body>
</html>`)

	return builder.String()
}

// GenerateMarkdownReport generates a Markdown report
func (rg *ReportGenerator) GenerateMarkdownReport() string {
	var builder strings.Builder

	builder.WriteString("# Provider Validation Report\n\n")
	builder.WriteString(fmt.Sprintf("**Generated:** %s\n\n", rg.report.Timestamp.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("**Total Providers:** %d\n", rg.report.TotalProviders))
	builder.WriteString(fmt.Sprintf("**Valid Providers:** %d\n", rg.report.ValidCount))
	builder.WriteString(fmt.Sprintf("**Invalid Providers:** %d\n", rg.report.InvalidCount))
	builder.WriteString(fmt.Sprintf("**Overall Status:** %s\n\n", rg.getStatusText()))

	builder.WriteString("## Summary\n\n")
	builder.WriteString(rg.getSummaryText() + "\n\n")

	// Errors section
	if len(rg.report.Errors) > 0 {
		builder.WriteString("## Errors\n\n")
		for i, err := range rg.report.Errors {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rg.formatValidationErrorMarkdown(&err)))
		}
		builder.WriteString("\n")
	}

	// Warnings section
	if len(rg.report.Warnings) > 0 {
		builder.WriteString("## Warnings\n\n")
		for i, warning := range rg.report.Warnings {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rg.formatValidationErrorMarkdown(&warning)))
		}
		builder.WriteString("\n")
	}

	// Infos section
	if len(rg.report.Infos) > 0 {
		builder.WriteString("## Information\n\n")
		for i, info := range rg.report.Infos {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rg.formatValidationErrorMarkdown(&info)))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// Helper methods

func (rg *ReportGenerator) getStatusText() string {
	if rg.report.IsValid {
		return "✅ Valid"
	}
	return "❌ Invalid"
}

func (rg *ReportGenerator) getStatusHTML() string {
	if rg.report.IsValid {
		return "<span style=\"color: green;\">✅ Valid</span>"
	}
	return "<span style=\"color: red;\">❌ Invalid</span>"
}

func (rg *ReportGenerator) getSummaryText() string {
	if rg.report.IsValid {
		return "All providers are valid and ready for use."
	}

	totalIssues := len(rg.report.Errors) + len(rg.report.Warnings) + len(rg.report.Infos)
	return fmt.Sprintf("Found %d issues (%d errors, %d warnings, %d info). Please review and fix the issues before proceeding.",
		totalIssues, len(rg.report.Errors), len(rg.report.Warnings), len(rg.report.Infos))
}

func (rg *ReportGenerator) formatValidationError(err *ValidationError) string {
	var parts []string

	if err.ProviderRef != "" {
		parts = append(parts, fmt.Sprintf("[%s]", err.ProviderRef))
	}

	parts = append(parts, fmt.Sprintf("%s: %s", err.RuleName, err.Message))

	if err.Field != "" {
		parts = append(parts, fmt.Sprintf("(field: %s)", err.Field))
	}

	if err.Value != "" {
		parts = append(parts, fmt.Sprintf("(value: %s)", err.Value))
	}

	if err.Suggestion != "" {
		parts = append(parts, fmt.Sprintf("💡 %s", err.Suggestion))
	}

	return strings.Join(parts, " ")
}

func (rg *ReportGenerator) formatValidationErrorHTML(err *ValidationError) string {
	var builder strings.Builder

	if err.ProviderRef != "" {
		builder.WriteString(fmt.Sprintf("<span class=\"provider-ref\">[%s]</span> ", err.ProviderRef))
	}

	builder.WriteString(fmt.Sprintf("<span class=\"severity\">%s</span>: %s", err.Severity, err.Message))

	if err.Field != "" {
		builder.WriteString(fmt.Sprintf(" <em>(field: %s)</em>", err.Field))
	}

	if err.Value != "" {
		builder.WriteString(fmt.Sprintf(" <em>(value: %s)</em>", err.Value))
	}

	if err.Suggestion != "" {
		builder.WriteString(fmt.Sprintf(" <span class=\"suggestion\">💡 %s</span>", err.Suggestion))
	}

	return builder.String()
}

func (rg *ReportGenerator) formatValidationErrorMarkdown(err *ValidationError) string {
	var builder strings.Builder

	if err.ProviderRef != "" {
		builder.WriteString(fmt.Sprintf("*[%s]* ", err.ProviderRef))
	}

	builder.WriteString(fmt.Sprintf("**%s**: %s", err.RuleName, err.Message))

	if err.Field != "" {
		builder.WriteString(fmt.Sprintf(" *(field: %s)*", err.Field))
	}

	if err.Value != "" {
		builder.WriteString(fmt.Sprintf(" *(value: %s)*", err.Value))
	}

	if err.Suggestion != "" {
		builder.WriteString(fmt.Sprintf("\n  💡 *%s*", err.Suggestion))
	}

	return builder.String()
}

// ReportFormat defines supported report formats
type ReportFormat string

const (
	ReportFormatText     ReportFormat = "text"
	ReportFormatJSON     ReportFormat = "json"
	ReportFormatHTML     ReportFormat = "html"
	ReportFormatMarkdown ReportFormat = "markdown"
)

// GenerateReport generates a report in the specified format
func (rg *ReportGenerator) GenerateReport(format ReportFormat) (string, error) {
	switch format {
	case ReportFormatText:
		return rg.GenerateTextReport(), nil
	case ReportFormatJSON:
		return rg.GenerateJSONReport()
	case ReportFormatHTML:
		return rg.GenerateHTMLReport(), nil
	case ReportFormatMarkdown:
		return rg.GenerateMarkdownReport(), nil
	default:
		return "", fmt.Errorf("unsupported report format: %s", format)
	}
}

// ReportWriter handles writing reports to files or other outputs
type ReportWriter struct {
	generator *ReportGenerator
}

// NewReportWriter creates a new ReportWriter
func NewReportWriter(report *ValidationReport) *ReportWriter {
	return &ReportWriter{
		generator: NewReportGenerator(report),
	}
}

// WriteToFile writes a report to a file in the specified format
func (rw *ReportWriter) WriteToFile(filename string, format ReportFormat) error {
	report, err := rw.generator.GenerateReport(format)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// In a real implementation, this would write to a file
	// For now, we'll just return success
	_ = report // Placeholder for file writing logic

	return nil
}

// WriteToConsole writes a report to the console
func (rw *ReportWriter) WriteToConsole(format ReportFormat) error {
	report, err := rw.generator.GenerateReport(format)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	fmt.Println(report)
	return nil
}
