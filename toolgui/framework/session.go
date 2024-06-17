package framework

type Session struct {
	// Value record components' value
	Values map[string]any `json:"values"`
}

func NewSession() *Session {
	return &Session{
		Values: make(map[string]any),
	}
}
