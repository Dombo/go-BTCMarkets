package marketData

import (
	"BTCMarkets/lib"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	pathBase = "markets"
)

type MarketList struct {
	List []Market
}

// TODO This was annoying, can it be addressed better? If not, can we hoist it
func (ml *MarketList) UnmarshalJSON(input []byte) error {
	var keys []Market
	err := json.Unmarshal(input, &keys)
	if err != nil {
		return err
	}
	ml.List = keys
	return nil
}

// TODO Notify BTCMarkets doc maintainer that their docs mismatch their response format
type Market struct {
	MarketId string
	BaseAssetName string
	QuoteAssetName string
	MinOrderAmount float64 `json:",string"`
	MaxOrderAmount float64 `json:",string"`
	AmountDecimals uint32 `json:",string"`
	PriceDecimals uint32 `json:",string"`
}

func List(c *BTCMarkets.Client) (*MarketList, error) {
	marketList := &MarketList{}
	if err := c.Do(marketList, http.MethodGet, pathBase, nil); err != nil {
		return nil, err
	}

	return marketList, nil
}

type MarketTicker struct {
	MarketId string
	BestBid float64 `json:",string"`
	BestAsk float64 `json:",string"`
	LastPrice float64 `json:",string"`
	Volume24h float64 `json:",string"`
	Price24h float64 `json:",string"`
	Low24h float64 `json:",string"`
	High24h float64 `json:",string"`
	Timestamp BTCMarkets.SpecialDatetime
}

func ReadTicker(c *BTCMarkets.Client, marketId string) (*MarketTicker, error) {
	path := fmt.Sprintf("%s/%s/ticker", pathBase, marketId)
	marketTicker := &MarketTicker{}
	if err := c.Do(marketTicker, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return marketTicker, nil
}

type MarketCandleList struct {
	List []MarketCandle
}

func (mcl *MarketCandleList) UnmarshalJSON(input []byte) error {
	var keys []MarketCandle
	err := json.Unmarshal(input, &keys)
	if err != nil {
		return err
	}
	mcl.List = keys
	return nil
}

// https://stackoverflow.com/questions/49415573/golang-json-how-do-i-unmarshal-array-of-strings-into-int64
type Float64Str float64

func (f *Float64Str) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil { // Try string first
		value, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*f = Float64Str(value)
		return nil
	}

	return json.Unmarshal(b, (*float64)(f)) // Fallback to number
}

type MarketCandle struct {
	Time BTCMarkets.SpecialDatetime
	Open  Float64Str
	High  Float64Str
	Low  Float64Str
	Close  Float64Str
	Volume  Float64Str
}

// https://eagain.net/articles/go-json-array-to-struct/
func (mc *MarketCandle) UnmarshalJSON(input []byte) error {
	tmp := []interface{}{&mc.Time,&mc.Open,&mc.High,&mc.Low,&mc.Close,&mc.Volume}
	desiredLen := len(tmp)
	if err := json.Unmarshal(input, &tmp); err != nil {
		return err
	}
	if received, expected := len(tmp), desiredLen; received != expected {
		return fmt.Errorf("wrong number of fields in MarketCandle: %d != %d", received, expected)
	}
	return nil
}

// Default time window is 1 day, enforced by server API
// Retrieve candles either by pagination (before, after, limit) or by specifying timestamp parameters (from and/or to).
// Pagination parameters can't be combined with timestamp parameters.
// Default behavior is pagination when no query param is specified.
func ReadCandles(c *BTCMarkets.Client, marketId string) (*MarketCandleList, error) {
	path := fmt.Sprintf("%s/%s/candles", pathBase, marketId)
	marketCandleList := &MarketCandleList{}
	if err := c.Do(marketCandleList, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return marketCandleList, nil
}