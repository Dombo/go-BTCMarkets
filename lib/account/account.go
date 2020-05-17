package account

import (
	"BTCMarkets/lib"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	pathBase = "accounts"
)


type BalanceList struct {
	List []Balance
}

func (bl *BalanceList) UnmarshalJSON(input []byte) error {
	var keys []Balance
	err := json.Unmarshal(input, &keys)
	if err != nil {
		return err
	}
	bl.List = keys
	return nil
}

type Balance struct {
	AssetName string
	Balance float64 `json:",string"`
	Available float64 `json:",string"`
	Locked float64 `json:",string"`
}

func ListBalances(c *BTCMarkets.Client) (*BalanceList, error) {
	path := fmt.Sprintf("%s/%s", pathBase, "me/balances")
	balanceList := &BalanceList{}
	//if err := c.Request(balanceList, http.MethodGet, path, nil); err != nil {
	//	return nil, err
	//}
	if err := c.Do(balanceList, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return balanceList, nil
}

type TransactionList struct {
	List []Transaction
}

func (tl *TransactionList) UnmarshalJSON(input []byte) error {
	var keys []Transaction
	err := json.Unmarshal(input, &keys)
	if err != nil {
		return err
	}
	tl.List = keys
	return nil
}

type Transaction struct {
	Id uint64 `json:",string"`
	CreationTime BTCMarkets.SpecialDatetime
	Description string
	AssetName string
	Amount float64 `json:",string"`
	Balance float64 `json:",string"`
	Type string
	RecordType string
	ReferenceId int64 `json:",string"`
}

type TransactionsOptions struct {
	AssetName string
}

// Can hook the endpoint specific query param building in here
func (i TransactionsOptions) buildEndpointSpecificOptions() url.Values {
	query := url.Values{}
	query.Set("assetName", i.AssetName)
	return query
}

// TODO Ideally whatever type is passed to paginationQuery() would be some sort of interface which guarantees:
//		The buildEndpointSpecificOptions method exists on the "proxied"(?) type/struct/whatever

func ListTransactionsPaged(c *BTCMarkets.Client, transOptions TransactionsOptions, listOptions *BTCMarkets.ListOptions) *Iter {
	return &Iter{BTCMarkets.GetIter(*listOptions, func(listOptions *BTCMarkets.ListOptions) ([]interface{}, error) {
		endpointSpecificOptions := transOptions.buildEndpointSpecificOptions()
		path := fmt.Sprintf("%s/%s", pathBase, "me/transactions")
		transactionList := &TransactionList{}

		if err := c.DoWithPagination(transactionList, http.MethodGet, path, endpointSpecificOptions, listOptions); err != nil {
			return nil, err
		}

		ret := make([]interface{}, len(transactionList.List))
		for i, v := range transactionList.List {
			ret[i] = v
		}


		return ret, nil
	})}
}


// Iter is an iterator for customers.
type Iter struct {
	*BTCMarkets.Iter
}

// Customer returns the customer which the iterator is currently pointing to.
func (i *Iter) Transaction() Transaction {
	return i.Current().(Transaction)
}

//func reflectCurrentItemID(x interface{}) uint64 {
//	return reflect.ValueOf(x).FieldByName("Id").Uint()
//}