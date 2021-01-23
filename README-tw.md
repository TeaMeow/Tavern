# Tavern [English](./README.md) [![GoDoc](https://godoc.org/github.com/teacat/tavern?status.svg)](https://godoc.org/github.com/teacat/tavern) [![Coverage Status](https://coveralls.io/repos/github/teacat/tavern/badge.svg?branch=master)](https://coveralls.io/github/teacat/tavern?branch=master) [![Build Status](https://travis-ci.org/teacat/tavern.svg?branch=master)](https://travis-ci.org/teacat/tavern) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/tavern)](https://goreportcard.com/report/github.com/teacat/tavern)

Tavern 是一個 _不需要結構體標籤_ 的驗證函式庫且同時支援自訂驗證規則。

## 為什麼？

我對於所有基於結構體標籤的表單驗證函式庫感到厭倦了。透過 Tavern 你可以不需要建立任何結構體就可以進行驗證。

## 安裝方式

```bash
$ go get github.com/teacat/tavern
```

## 範例

```go
err := tavern.Validate([]tavern.Rule{
    {
        Value: "你好，世界！",
        Validators: []tavern.Validator{
            tavern.WithRequired(),
            tavern.WithMinLength(5),
        },
    },
})
if err != nil {
    panic(err)
}
```

## 驗證器

一個驗證包含了許多規則，而一個規則需要至少一個值跟相關的驗證器。驗證器之間可以相互串連，並且透過傳遞 `context` 來進行溝通。

內建的基礎驗證器有如：`WithRequired`、`WithLength`、`WithRange`、等…。查看 [GoDoc](https://pkg.go.dev/github.com/teacat/tavern) 來了解更多內建的驗證器。

### 自訂驗證器

你能夠建立自己的驗證器。

```go
type Validator func(ctx context.Context, value interface{}) (context.Context, error)
```

`context` 參數提供你一個上下文內容可以讓你在不同的驗證器之間共享來達到溝通的方式。透過 `ctx.Value(tavern.KeyRequired).(bool)` 你可以取得一個布林值，而其狀態是基於目前的值是否為必填（Required）。

## 自訂錯誤

預設的情況下，Tavern 僅會回傳 `ErrRequired`、`ErrLength`、等…內建的錯誤訊息，但這很有可能不是你期望的。因此你可以透過 `WithCustomError(validator Validator, err error)` 函式來替每個驗證器自訂自己的錯誤訊息。

事實上這個函式也是一個驗證器，但會在傳入的驗證器發生錯誤時回傳你自訂的錯誤訊息。

```go
err := tavern.Validate([]tavern.Rule{
    {
        Value: "",
        Validators: []tavern.Validator{
            tavern.WithCustomError(tavern.WithRequired(), errors.New("你在公三小")),
            tavern.WithMinLength(5),
        },
    },
})
if err != nil {
    panic(err) // 輸出：你在公三小
}
```

## 已知錯誤

-   `WithIPv4Address` 允許 `::0` 而這其實是 IPv6 的東西。
-   `WithIPAddress`, `WithIPv4Address`, `WithIPv6Address` 允許帶有通訊埠的 IP 位址。

IPAddress 驗證器使用 `net.ResolveIPAddr` 作為驗證基礎，不知道為什麼這些都能通過驗證。
