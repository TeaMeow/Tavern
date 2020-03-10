# Tavern [![GoDoc](https://godoc.org/github.com/teacat/tavern?status.svg)](https://godoc.org/github.com/teacat/tavern) [![Coverage Status](https://coveralls.io/repos/github/teacat/tavern/badge.svg?branch=master)](https://coveralls.io/github/teacat/tavern?branch=master) [![Build Status](https://travis-ci.org/teacat/tavern.svg?branch=master)](https://travis-ci.org/teacat/tavern) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/tavern)](https://goreportcard.com/report/github.com/teacat/tavern)

不需建構體標籤（Struct Tag）就能夠以串連式函式進行規則檢查的簡易表單驗證套件。

# 這是什麼？

通常來說你需要透過建構體標籤來配置表單驗證規則，但 Tavern 試圖避免此用法，而直接使用直覺的函式串連，並且最終回傳單個、多個錯誤訊息令開發者更好處理錯誤回傳，你亦能自訂錯誤內容。

# 為什麼？

因為透過函式進行驗證，比起建構體標籤的彈性還要來得大、較不複雜。

# 效能如何？

這裡有份簡略化的效能測試報表。

```
測試規格：
1.7 GHz Intel Core i7 (4650U)
8 GB 1600 MHz DDR3
```

# 索引

* [安裝方法](#安裝方法)
* [使用方式](#使用方式)
    * [規則](#規則)
        * [必填](#必填)
        * [最小／最大](#最小最大)
        * [長度](#長度)
        * [範圍](#範圍)
        * [日期格式](#日期格式)
        * [電子郵件地址](#電子郵件地址)
        * [在清單內](#在清單內)
        * [IP 位址](#ip-位址)
        * [網址](#網址)
        * [相等](#相等)
        * [正規表達式](#正規表達式)
    * [檢查](#檢查)
        * [多重檢查](#多重檢查)
    * [自訂錯誤](#自訂錯誤)
        * [多形錯誤](#多形錯誤)

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/teacat/tavern
```

# 使用方式

使用 Tavern 驗證多個欄位的方法十分簡單。

```go
err := tavern.
	Add(Username).Length(6, 32).Required().Error(errors.New("使用者帳號錯誤，至少需要 6 到 32 個字。")).
	Add(Password).Length(8, 128).Required().Error(errors.New("密碼錯誤，至少需要 8 到 128 個字。")).
	Check()
if err != nil {
	panic(err)
}
```

## 規則

### 必填

一個空的字串、僅有空白、`nil` 值都會被因為必填選項而被拒絕。

```go
tavern.Add(Username).Required()
```

### 最小／最大

如果值是字串，這會驗證字串的最小長度與最大長度；如果值是數值型態，則是範圍大小。

```go
tavern.Add(Username).Min(10)
tavern.Add(Number).Max(30)
```

### 長度

驗證數字、字串的最小與最大長度。

```go
tavern.Add(Username).Length(10, 30)
```

### 範圍

驗證數字的最小與最大範圍。

```go
tavern.Add(Number).Range(10, 30)
```

### 日期格式

驗證傳入的日期、時間字串是否符合指定格式，若有多個日期格式只需符合其中一個即可。

```go
tavern.Add(Birthday).Date("2016-01-02")
tavern.Add(MyDate).Date("2016-01-02", "15:04:05")
```

### 電子郵件地址

確認字串是否為一個電子郵件地址，請注意這並不能當作「最終手段」。這個驗證使用最簡單的正規表達式進行驗證，這意味著奇異的電子郵件地址仍能被認可，如果你真的希望確認電子郵件地址是否正確，請試圖發送郵件至該地址以核對。

```go
tavern.Add(EmailAddress).Email()
```

### 在清單內

確認傳入的值是否為指定值。

```go
tavern.Add(Number).In(1, 2, 3)
tavern.Add(Username).In("YamiOdymel", "Karisu", "Iknore")
```

### IP 位址

驗證傳入的值是否為 IP 位址，若無指定 `v4` 或 `v6` 則預設為兩者皆可。

```go
tavern.Add(IPAddress).IP()
tavern.Add(IPAddress).IP("v4")
tavern.Add(IPAddress).IP("v6")
```

### 網址

確認傳入的字串是否為網址，可指定多個必要前輟，若無指定則為只要是網址即可。

```go
tavern.Add(Website).URL()
tavern.Add(Website).URL("https://", "http://")
```

### 相等

驗證兩者的值是否相等。

```go
tavern.Add(ConfirmPassword).Equal(Password)
```

### 正規表達式

驗證字串是否能通過正規表達式。

```go
tavern.Add(Username).RegExp("a-Z0-9")
```

## 檢查

透過 `Check` 進行檢查，並且回傳第一個錯誤。

```go
err := tavern.Add(IPAddress).IP().
	Add(Website).URL().
	Add(Username).Required().
	Check()
if err != nil {
	panic(err)
}
```

### 多重檢查

如果你希望取得每個錯誤訊息並得到一個錯誤切片，請嘗試 `CheckAll`。此函式較 `Check` 慢一點（因為遇到錯誤會繼續驗證以收集所有錯誤）。

```go
errs := tavern.
	Add(IPAddress).IP().
	Add(Website).URL().
	Add(Username).Required().
	CheckAll()
if errs != nil {
	for _, v := range errs {
		fmt.Println(v.Error())
	}
}
```

## 自訂錯誤

你能夠傳入自訂的錯誤用以取代 Tavern 內建的錯誤訊息，令你更好比對是何種錯誤發生。

```go
ErrUsername := errors.New("使用者帳號欄位錯誤。")
err := tavern.
	Add(Username).Length(6, 32).Required().Error(ErrUsername).
	Check()
```

### 多形錯誤

如果一個值有多個條件，你希望精準地確認該值是因為沒有達到何種條件而失敗，試著同樣使用 `Error` 函式，但傳入一個 `tavern.E` 建構體，並將自訂錯誤傳入該建構體。當錯誤發生時 Tavern 會搜尋並回傳該建構體的相對應錯誤訊息。

```go
err := tavern.Add(Username).Length(6, 32).Required().
	Error(tavern.E{
		Length:   errors.New("使用者帳號的長度不對，應是 6 到 32 個字。"),
		Required: errors.New("使用者帳號為必填選項。"),
	}).
	Check()
```