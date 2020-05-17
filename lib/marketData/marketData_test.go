package marketData

import (
	"BTCMarkets/lib/internal/test"
	"fmt"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {
	test.EnableServer(m)
}

func TestList(t *testing.T) {
	test.WillReturnTestdata(t, "marketDataListObject.json", http.StatusOK)
	client := test.Client(t)

	list, err := List(client)
	if err != nil {
		t.Fatalf("unexpected error retrieving MarketData list: %s", err)
	}

	fmt.Println("Test got this back: ", list)

	if list.List[0].MarketId != "BTC-AUD" {
		t.Fatalf("expected BTC-AUD, got %s", list.List[0].MarketId)
	}

	// TODO Test the pagination somehow..

	test.AssertEndpointCalled(t, http.MethodGet, "markets")
}