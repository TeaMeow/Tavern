# Tavern

```go
tavern.Add(Username).Required()
```

```go
tavern.Add(Username).Min(10).Max(30)
```

```go
tavern.Add(Username).Length(10, 30)
```

```go
tavern.Add(Number).Range(10, 30)
```

```go
tavern.Add(Birthday).Date("2016-01-02")
```


```go
tavern.Add(Birthday).Email()
```

```go
tavern.Add(Gender).In(1, 2, 3)
```

```go
tavern.Add(IPAddress).IP()
tavern.Add(IPAddress).IP("v4")
tavern.Add(IPAddress).IP("v6")
```

```go
tavern.Add(Website).URL()
tavern.Add(Website).URL("https://", "http://")
```

```go
tavern.Add(ConfirmPassword).Equal(Password)
```

```go
tavern.Add(Username).RegExp("a-Z0-9")
```

```
err := tavern.Add(IPAddress).IP().
	Add(Website).URL().
	Add(Username).Required().
	Check()
if err != nil {
	panic(err)
}
```

```
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

```
err := tavern.Add(Username).Length(6, 32).Required().
	Error(tavern.ErrorMessages{
		Length:   "使用者帳號的長度不對，應是 6 到 32 個字。",
		Required: "使用者帳號為必填選項。",
	}).
	Check()
```