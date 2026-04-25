package shared

import (
	"fmt"
	"strings"
)

type TemplateVariableOption struct {
	Key         string
	Label       string
	Description string
	Sample      string
}

var templateVariableOptions = []TemplateVariableOption{
	{
		Key:         "customer_name",
		Label:       "Customer Name",
		Description: "Nama customer/pelanggan.",
		Sample:      "Pelanggan",
	},
	{
		Key:         "days_since_last_visit",
		Label:       "Days Since Last Visit",
		Description: "Jumlah hari sejak kunjungan terakhir.",
		Sample:      "30",
	},
	{
		Key:         "last_visit_date",
		Label:       "Last Visit Date",
		Description: "Tanggal kunjungan terakhir customer.",
		Sample:      "10 Maret 2026",
	},
	{
		Key:         "service_name",
		Label:       "Service Name",
		Description: "Nama layanan atau treatment.",
		Sample:      "Facial Treatment",
	},
	{
		Key:         "business_name",
		Label:       "Business Name",
		Description: "Nama bisnis outlet.",
		Sample:      "Celestial Spa & Wellness",
	},
}

var templateVariableOptionByKey = func() map[string]TemplateVariableOption {
	m := make(map[string]TemplateVariableOption, len(templateVariableOptions))
	for _, opt := range templateVariableOptions {
		m[opt.Key] = opt
	}
	return m
}()

func ListTemplateVariableOptions() []TemplateVariableOption {
	out := make([]TemplateVariableOption, len(templateVariableOptions))
	copy(out, templateVariableOptions)
	return out
}

func NormalizeTemplateVariableKeys(keys []string) ([]string, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("bodyExample is required")
	}
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		k := strings.TrimSpace(key)
		if k == "" {
			return nil, fmt.Errorf("bodyExample cannot contain empty value")
		}
		if _, ok := templateVariableOptionByKey[k]; !ok {
			return nil, fmt.Errorf("invalid bodyExample key: %s", k)
		}
		out = append(out, k)
	}
	return out, nil
}

func TemplateVariableSampleForKey(key string) (string, bool) {
	opt, ok := templateVariableOptionByKey[strings.TrimSpace(key)]
	if !ok {
		return "", false
	}
	return opt.Sample, true
}
