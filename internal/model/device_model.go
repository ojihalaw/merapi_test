package model

type DeviceResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Location  string `json:"location,omitempty"`
	Status    string `json:"status,omitempty"`
	Sensors   []SensorResponse `json:"sensors,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateDeviceRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Location string `json:"location,omitempty"`
	Status   string `json:"status,omitempty"`
}

type UpdateDeviceRequest struct {
	Name     *string `json:"name,omitempty"`
	Location *string `json:"location,omitempty"`
	Status   *string `json:"status,omitempty"`
}