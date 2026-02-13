package model

type SensorResponse struct {
	ID        string `json:"id,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Unit      string `json:"unit,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateSensorRequest struct {
	DeviceID string `json:"device_id" validate:"required,uuid"`
	Name     string `json:"name" validate:"required,max=100"`
	Type     string `json:"type" validate:"required,max=50"`
	Unit     string `json:"unit,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type UpdateSensorRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,max=100"`
	Type     *string `json:"type,omitempty" validate:"omitempty,max=50"`
	Unit     *string `json:"unit,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}