package types

type Ehr struct {
	Id      string `json:"id"`
	Patient string `json:"patient"`

	Details []EhrDetail `json:"details"`
}

type EhrDetail struct {
	Id string `json:"id"`

	Date        string `json:"date"`
	Description string `json:"description"`

	Type  string `json:"type"`
	Value string `json:"value"`
}
