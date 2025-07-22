package config

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Timezone string `yaml:"timezone"`
	Calendar string `yaml:"calendar"`
	Email    Email  `yaml:"email"`
}

type Email struct {
	From       string              `yaml:"from"`
	Recepients []string            `yaml:"recipients"`
	Password   string              `yaml:"password"`
	SMTPHost   string              `yaml:"smtpHost"`
	SMTPPort   string              `yaml:"smtpPort"`
	Templates  map[string]Template `yaml:"templates"`
}

func LoadConfig(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Expand environment variables in the file content
	content := os.ExpandEnv(string(file))

	var cfg Config
	err = yaml.Unmarshal([]byte(content), &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}

type Template struct {
	Subject         string `yaml:"subject"`
	TemplateContent string `yaml:"templateContent"`
}

func (template *Template) ResolveContent(templateVars map[string]string) string {
	content := template.TemplateContent
	for key, value := range templateVars {
		placeholder := "{" + key + "}"
		content = strings.ReplaceAll(content, placeholder, value)
	}
	return content
}

func (config *Config) FindTemplateBy(templateType string) Template {
	template, ok := config.Email.Templates[string(templateType)]
	if !ok {
		log.Fatalf("Template %s not found in configuration", templateType)
	}

	return template
}
