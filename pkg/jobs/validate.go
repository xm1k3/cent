package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const (
	tmBaseUrlDefault = "https://tm.nuclei.sh"
)

type Mark struct {
	Name     string `json:"name,omitempty"`
	Position int    `json:"position,omitempty"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
}

type Error struct {
	Name string `json:"name"`
	Mark Mark   `json:"mark"`
}

type LintError struct {
	Name   string `json:"name,omitempty"`
	Reason string `json:"reason,omitempty"`
	Mark   Mark   `json:"mark,omitempty"`
}

type TemplateLintResp struct {
	Input     string    `json:"template_input,omitempty"`
	Lint      bool      `json:"template_lint,omitempty"`
	LintError LintError `json:"lint_error,omitempty"`
}

type ValidateError struct {
	Location string      `json:"location,omitempty"`
	Message  string      `json:"message,omitempty"`
	Name     string      `json:"name,omitempty"`
	Argument interface{} `json:"argument,omitempty"`
	Stack    string      `json:"stack,omitempty"`
	Mark     struct {
		Line   int `json:"line,omitempty"`
		Column int `json:"column,omitempty"`
		Pos    int `json:"pos,omitempty"`
	} `json:"mark,omitempty"`
}

// TemplateResponse from templateman to be used for enhancing and formatting
type TemplateResp struct {
	Input              string          `json:"template_input,omitempty"`
	Format             bool            `json:"template_format,omitempty"`
	Updated            string          `json:"updated_template,omitempty"`
	Enhance            bool            `json:"template_enhance,omitempty"`
	Enhanced           string          `json:"enhanced_template,omitempty"`
	Lint               bool            `json:"template_lint,omitempty"`
	LintError          LintError       `json:"lint_error,omitempty"`
	Validate           bool            `json:"template_validate,omitempty"`
	ValidateErrorCount int             `json:"validate_error_count,omitempty"`
	ValidateError      []ValidateError `json:"validate_error,omitempty"`
	Error              Error           `json:"error,omitempty"`
}

func ValidateTemplate(data string) (bool, error) {
	client := retryablehttp.NewClient()
	client.Logger = nil
	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/validate", tmBaseUrlDefault), strings.NewReader(data))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-yaml")

	var resp *http.Response

	resp, err = client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("unexpected status code: %v", resp.Status)
	}
	var validateResp TemplateResp
	if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
		return false, err
	}
	if validateResp.Validate {
		return true, nil
	}
	if validateResp.ValidateErrorCount > 0 {
		if len(validateResp.ValidateError) > 0 {
			return false, fmt.Errorf("Validation failed: %s at line %v", validateResp.ValidateError[0].Message, validateResp.ValidateError[0].Mark.Line)
		}
		return false, fmt.Errorf("validation failed")
	}
	if validateResp.Error.Name != "" {
		return false, fmt.Errorf(validateResp.Error.Name)
	}
	return false, fmt.Errorf("template validation failed")
}
