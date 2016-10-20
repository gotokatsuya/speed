# speed

## Usage

### Enable

A: Add speed option

```bash
$ go run main.go -speed
```

B: Call EnableLogger() in main.init()

```go
package main

import "github.com/gotokatsuya/speed"

func init() {
    speed.EnableLogger()
}
```


### Trace

```go
package f

import "github.com/gotokatsuya/speed"

func (f *F) Heavy() {
    speedLogger := speed.NewLogger("Trace Heavy Function").Begin()
    defer speedLogger.End()

    ...

}
```


### Watch

```bash
$ cat $TMPDIR/speed-20161020.log 
```

```
description: Trace Heavy Function	
file: github.com/gotokatsuya/speed/heavy.go
begin_at: 2016-10-20 15:21:32.871501713 +0900 JST
end_at: 2016-10-20 15:21:33.874674989 +0900 JST
caller: github.com/gotokatsuya/speed.Heavy.func1 (17L)
seconds: 1.003173
milliseconds: 1003.173276
microseconds: 1003173.276000
```

