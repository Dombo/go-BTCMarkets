package BTCMarkets

import (
	"strings"
	"time"
)

type SpecialDatetime struct {
	time.Time
}

func (sd *SpecialDatetime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", strInput)
	if err != nil {
		return err
	}

	sd.Time = newTime
	return nil
}
