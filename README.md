# Tavern

```go
tavern.Add(Username).Required()
```

```go
tavern.Add(Username).Length(0, 3)
```

```go
tavern.Add(Username).Longer(10)
tavern.Add(Username).Shorter(30)
```

```go
tavern.Add(Number).Range(0, 26)
```

```go
tavern.Add(Age).Greater(30)
tavern.Add(Age).Lesser(99)
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
tavern.Add(Gender).URL()
tavern.Add(Gender).URL("https://", "http://")
```

```go
tavern.Add(ConfirmPassword).Equal(Password)
```

```go
tavern.Add(Username).RegExp("a-Z0-9")
```

```go
tavern.Add(IPAddress).IP("v4").
      .Add(Age).Max(99)
```