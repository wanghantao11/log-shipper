package log

// swagger:model
type Log struct {
	ID string `json:"id"`
	L  string `json:"@l"`
	M  string `json:"@m"`
	T  string `json:"@t"`
}
