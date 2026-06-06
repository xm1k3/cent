package jobs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

var (
	ReTemplateID = regexp.MustCompile(`^([a-zA-Z0-9]+[-_])*[a-zA-Z0-9]+$`)
)

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*s = []string{single}
	} else {
		*s = multi
	}
	return nil
}

type Info struct {
	Name     string        `yaml:"name" json:"name"`
	Author   StringOrSlice `yaml:"author" json:"author"`
	Severity string        `yaml:"severity" json:"severity"`
}

type TemplateValidation struct {
	ID         string                 `yaml:"id" json:"id"`
	Info       Info                   `yaml:"info" json:"info"`
	Requests   []interface{}          `yaml:"requests,omitempty" json:"requests,omitempty"`
	HTTP       []interface{}          `yaml:"http,omitempty" json:"http,omitempty"`
	DNS        []interface{}          `yaml:"dns,omitempty" json:"dns,omitempty"`
	File       []interface{}          `yaml:"file,omitempty" json:"file,omitempty"`
	Network    []interface{}          `yaml:"network,omitempty" json:"network,omitempty"`
	TCP        []interface{}          `yaml:"tcp,omitempty" json:"tcp,omitempty"`
	Headless   []interface{}          `yaml:"headless,omitempty" json:"headless,omitempty"`
	SSL        []interface{}          `yaml:"ssl,omitempty" json:"ssl,omitempty"`
	Websocket  []interface{}          `yaml:"websocket,omitempty" json:"websocket,omitempty"`
	WHOIS      []interface{}          `yaml:"whois,omitempty" json:"whois,omitempty"`
	Code       []interface{}          `yaml:"code,omitempty" json:"code,omitempty"`
	Javascript []interface{}          `yaml:"javascript,omitempty" json:"javascript,omitempty"`
	Workflows  []interface{}          `yaml:"workflows,omitempty" json:"workflows,omitempty"`
	Flow       string                 `yaml:"flow,omitempty" json:"flow,omitempty"`
	Variables  map[string]interface{} `yaml:"variables,omitempty" json:"variables,omitempty"`
	Constants  map[string]interface{} `yaml:"constants,omitempty" json:"constants,omitempty"`
}

func isBlank(s string) bool {
	if s == "" {
		return true
	}
	for _, c := range s {
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			return false
		}
	}
	return true
}

func validateTemplateMandatoryFields(template *TemplateValidation) error {
	info := template.Info

	var validateErrors []error

	if isBlank(info.Name) {
		validateErrors = append(validateErrors, fmt.Errorf("mandatory 'name' field is missing"))
	}

	if len(info.Author) == 0 {
		validateErrors = append(validateErrors, fmt.Errorf("mandatory 'author' field is missing"))
	}

	if template.ID == "" {
		validateErrors = append(validateErrors, fmt.Errorf("mandatory 'id' field is missing"))
	} else if !ReTemplateID.MatchString(template.ID) {
		validateErrors = append(validateErrors, fmt.Errorf("invalid field format for 'id' (allowed format is %s)", ReTemplateID.String()))
	}

	if len(validateErrors) > 0 {
		return errors.Join(validateErrors...)
	}

	return nil
}

func validateTemplateOptionalFields(template *TemplateValidation) error {
	info := template.Info

	var warnings []error

	hasProtocol := len(template.Workflows) > 0 ||
		len(template.Requests) > 0 ||
		len(template.HTTP) > 0 ||
		len(template.DNS) > 0 ||
		len(template.File) > 0 ||
		len(template.Network) > 0 ||
		len(template.TCP) > 0 ||
		len(template.Headless) > 0 ||
		len(template.SSL) > 0 ||
		len(template.Websocket) > 0 ||
		len(template.WHOIS) > 0 ||
		len(template.Code) > 0 ||
		len(template.Javascript) > 0

	if hasProtocol && len(template.Workflows) == 0 && isBlank(info.Severity) {
		warnings = append(warnings, fmt.Errorf("field 'severity' is missing"))
	}

	if len(warnings) > 0 {
		return errors.Join(warnings...)
	}

	return nil
}

func ValidateTemplateFile(templatePath string) error {
	file, err := os.Open(templatePath)
	if err != nil {
		return fmt.Errorf("could not open template file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read template file: %w", err)
	}

	template := &TemplateValidation{}
	err = yaml.Unmarshal(data, template)
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	if err := validateTemplateMandatoryFields(template); err != nil {
		return fmt.Errorf("mandatory field validation failed: %w", err)
	}

	if err := validateTemplateOptionalFields(template); err != nil {
		fmt.Printf("[WARN] %s: %v\n", templatePath, err)
	}

	return nil
}
