package domain

const (
	DefaultPageLimit = 25
	MaxPageLimit     = 100
)

type Filters struct {
	Limit  int32 `json:"limit" query:"limit"`
	Offset int32 `json:"offset" query:"offset"`
}

func (f *Filters) Validate() error {
	if f.Limit > 100 {
		return ErrPageLimit
	}
	if f.Limit == 0 {
		f.Limit = DefaultPageLimit
	}
	return nil
}

type Page struct {
	Total    int64                  `json:"total"`
	Limit    int32                  `json:"limit"`
	Offset   int32                  `json:"offset"`
	Metadata map[string]interface{} `json:"metadata"`
}
