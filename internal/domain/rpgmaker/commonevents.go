package rpgmaker

// CommonEvent represents a common event that can be triggered
type CommonEvent struct {
	ID       int             `json:"id"`
	List     []*EventCommand `json:"list"`
	Name     string          `json:"name"`
	SwitchId int             `json:"switchId"`
	Trigger  int             `json:"trigger"`
}

// CommonEventsData is an array of common events
type CommonEventsData []*CommonEvent




