package sensors

import "time"

type baseSensor struct {
	UpdatedAt *time.Time `json:"updated_at"`
}

func (bs *baseSensor) PostUpdate() error {
	now := time.Now()
	bs.UpdatedAt = &now
	return nil
}
