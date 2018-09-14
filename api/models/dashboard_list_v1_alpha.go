package models

import "encoding/json"

type DashboardListV1Alpha struct {
	Dashboards []DashboardV1Alpha `json:"dashboards" yaml:"dashboards"`
}

func NewDashboardListV1AlphaFromJson(data []byte) (*DashboardListV1Alpha, error) {
	list := DashboardListV1Alpha{}

	err := json.Unmarshal(data, &list)

	if err != nil {
		return nil, err
	}

	for _, s := range list.Dashboards {
		if s.ApiVersion == "" {
			s.ApiVersion = "v1alpha"
		}

		if s.Kind == "" {
			s.Kind = "Dashboard"
		}
	}

	return &list, nil
}
