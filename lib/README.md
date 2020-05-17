```go
    package main
    
    import (
        "BTCMarkets/lib"
        "BTCMarkets/lib/account"
        "BTCMarkets/lib/marketData"
        "BTCMarkets/lib/report"
        "fmt"
    )
    
    const apiKey = "YOUR-API-KEY"
    const privateKey = "YOUR-PRIVATE-KEY"
    
    func main() {
        client := BTCMarkets.New(apiKey, privateKey)
    
        balances, err := account.ListBalances(client)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("Resp: ", balances)
    
    
        listTransOptions := account.TransactionsOptions{
            AssetName: "BTC",
        }
        // Backwards paging
        listOptions := BTCMarkets.ListOptions{
            Limit: 2,
            Before: 5713565494,
        }
        i := account.ListTransactionsPaged(client, listTransOptions, &listOptions)
        fmt.Println("Backwards paging")
        for i.Next() {
            fmt.Println("Resp: ", i.Transaction().Id)
            //fmt.Println("Resp: ", i)
        }
    
        // Forwards paging
        listOptions = BTCMarkets.ListOptions{
            Limit: 2,
            After: 5713565494,
        }
        i = account.ListTransactionsPaged(client, listTransOptions, &listOptions)
        fmt.Println("Forwards paging")
        for i.Next() {
            fmt.Println("Resp: ", i.Transaction().Id)
        }
    
        if err := i.Err(); err != nil {
            fmt.Println(err)
        }
    
    
        markets, err := marketData.List(client)
        if err != nil {
            fmt.Println(err)
            return
        }
    
        fmt.Println("Resp: ", markets)
    
        marketTicker, err := marketData.ReadTicker(client, markets.List[0].MarketId)
        if err != nil {
            fmt.Println(err)
            return
        }
    
        fmt.Println("Resp: ", marketTicker)
    
        marketCandles, err := marketData.ReadCandles(client, markets.List[0].MarketId)
        if err != nil {
            fmt.Println(err)
            return
        }
    
        fmt.Println("Resp: ", marketCandles)
    
        // Would defining the allowed values on the struct make sense here? It should offer a better DX
        reportReq, err := report.Create(client, "json")
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("Resp: ", reportReq)
        report, err := report.Get(client, "mjeuclhpi2l0***fp94ldp4rsl")
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("Resp: ", report)
    }
```