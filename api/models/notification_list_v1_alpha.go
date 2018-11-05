package models

import (
	"encoding/json"
)

type NotificationListV1Alpha struct {
	Notifications []NotificationV1Alpha `json:"notifications" yaml:"notifications"`
}

func NewNotificationListV1AlphaFromJson(data []byte) (*NotificationListV1Alpha, error) {
	list := NotificationListV1Alpha{}

	err := json.Unmarshal(data, &list)

	if err != nil {
		return nil, err
	}

	for _, s := range list.Notifications {
		if s.ApiVersion == "" {
			s.ApiVersion = "v1alpha"
		}

		if s.Kind == "" {
			s.Kind = "Notification"
		}
	}

	return &list, nil
}
