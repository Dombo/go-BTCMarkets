package BTCMarkets

import "fmt"

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("The API returned %s: %s", e.Code, e.Message)
}
