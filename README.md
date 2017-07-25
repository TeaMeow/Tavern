# Tavern [![GoDoc](https://godoc.org/github.com/TeaMeow/Tavern?status.svg)](https://godoc.org/github.com/TeaMeow/Tavern)

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

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/TeaMeow/Tavern
```

# 使用方式

## 規則

### 必填

```go
tavern.Add(Username).Required()
```

### 最小／最大

```go
tavern.Add(Username).Min(10).Max(30)
```

### 長度

```go
tavern.Add(Username).Length(10, 30)
```

### 範圍

```go
tavern.Add(Number).Range(10, 30)
```

### 日期格式

```go
tavern.Add(Birthday).Date("2016-01-02")
```

### 電子郵件地址

```go
tavern.Add(Birthday).Email()
```

### 在清單內

```go
tavern.Add(Gender).In(1, 2, 3)
```

### IP 位址

```go
tavern.Add(IPAddress).IP()
tavern.Add(IPAddress).IP("v4")
tavern.Add(IPAddress).IP("v6")
```

### 網址

```go
tavern.Add(Website).URL()
tavern.Add(Website).URL("https://", "http://")
```

### 相等

```go
tavern.Add(ConfirmPassword).Equal(Password)
```

### 正規表達式

```go
tavern.Add(Username).RegExp("a-Z0-9")
```

## 檢查

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

```go
errs := tavern.Add(IPAddress).IP().
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

```go
err := tavern.Add(Username).Length(6, 32).Required().
	Error(errors.New("使用者帳號欄位錯誤。")).
	Check()
```

### 多形錯誤

```go
err := tavern.Add(Username).Length(6, 32).Required().
	Error(tavern.E{
		Length:   errors.New("使用者帳號的長度不對，應是 6 到 32 個字。"),
		Required: errors.New("使用者帳號為必填選項。"),
	}).
	Check()
```