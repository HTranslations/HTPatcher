package rpgmaker

import (
	"encoding/json"
	"htpatcher/internal/util"
)

// CommonEvent represents a common event that can be triggered
type CommonEvent struct {
	ID       int                        `json:"id"`
	List     []*EventCommand            `json:"list"`
	Name     string                     `json:"name"`
	SwitchId int                        `json:"switchId"`
	Trigger  int                        `json:"trigger"`
	Extras   map[string]json.RawMessage `json:"-"`
}

func (c *CommonEvent) UnmarshalJSON(data []byte) error {
	type Alias CommonEvent
	aux := (*Alias)(c)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	c.Extras, _ = util.UnmarshalExtras(data, util.GetJSONFieldNames(c))
	return nil
}

func (c CommonEvent) MarshalJSON() ([]byte, error) {
	type Alias CommonEvent
	return util.MarshalWithExtras((*Alias)(&c), c.Extras)
}

// CommonEventsData is an array of common events
type CommonEventsData []*CommonEvent
