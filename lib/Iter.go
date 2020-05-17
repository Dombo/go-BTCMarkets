package BTCMarkets

import (
	"reflect"
)

// Would like a way to set default values
type ListOptions struct {
	Before, After, Limit  uint64
	PrevBefore, PrevAfter uint64
}

func (lo *ListOptions) GetPageOptions() *ListOptions {
	return lo
}

type Iter struct {
	cur         interface{}
	err         error
	pageOptions ListOptions
	query       Query
	values      []interface{}
}

func (it *Iter) Current() interface{} {
	return it.cur
}

func (it *Iter) Err() error {
	return it.err
}

func (it *Iter) getPage() {
	it.values, it.err = it.query(it.pageOptions.GetPageOptions())

	if it.pageOptions.After != 0 {
		// We are paging forwards, but the API only returns things in DESC order
		reverse(it.values)
	}
}

func reverse(a []interface{}) {
	for i := 0; i < len(a)/2; i++ {
		a[i], a[len(a)-i-1] = a[len(a)-i-1], a[i]
	}
}

func (it *Iter) Next() bool {
	if len(it.values) == 0 {
		/*
			// API does not support paging between a range which is a bit lame
			if it.pageOptions.Before != 0 && it.pageOptions.After != 0 { // Paging between a range
				if it.pageOptions.Before > it.pageOptions.After && it.pageOptions.PrevBefore != 0  { // Paging backwards between a range
					it.pageOptions.Before = it.pageOptions.PrevBefore
				} else if it.pageOptions.After > it.pageOptions.Before && it.pageOptions.PrevAfter != 0 { // Paging forwards between a range
					it.pageOptions.After = it.pageOptions.PrevAfter
				}
			}
		*/
		if it.pageOptions.Before != 0 { // Paging backwards
			it.pageOptions.Before = it.pageOptions.PrevBefore
		} else if it.pageOptions.After != 0 { // Paging forwards
			it.pageOptions.After = it.pageOptions.PrevAfter
		}
		//if it.pageOptions.PrevAfter == reflectCurrentItemID(it.cur) {}
		//if it.pageOptions.Before == reflectCurrentItemID(it.cur) {}
		if it.pageOptions.Limit != 0 {

		}
		it.getPage()
	}
	if len(it.values) == 0 { // If the API responds with nothing more, we're done paging
		return false
	}
	it.cur = it.values[0]
	it.values = it.values[1:]
	return true
}

type Query func(*ListOptions) ([]interface{}, error)

func GetIter(pageOptions ListOptions, query Query) *Iter {
	iter := &Iter{
		pageOptions: pageOptions,
		query:       query,
	}

	iter.getPage()

	return iter
}

func reflectCurrentItemID(x interface{}) uint64 {
	return reflect.ValueOf(x).FieldByName("Id").Uint()
}
