package entity

type Queue struct {
	RequestID   string  `json:"request_id"`
	UniqID      string  `json:"uniq_id,omitempty"`
	Description string  `json:"description,omitempty"`
	Condition   string  `json:"condition,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Color       string  `json:"color,omitempty"`
	Size        string  `json:"size,omitempty"`
	AgeGroup    string  `json:"age_group,omitempty"`
	Material    string  `json:"material,omitempty"`
	WeightKG    float64 `json:"weight_kg,omitempty"`
}
