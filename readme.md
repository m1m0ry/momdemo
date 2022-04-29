# HOW TO START

```
# first
## 启动nsq
### https://nsq.io/overview/quick_start.html

# in one shell
nsqlookupd
# in another shell
nsqd --lookupd-tcp-address=127.0.0.1:4160
# in another shell
nsqadmin --lookupd-http-address=127.0.0.1:4161
```

```go
# then
go run random/rand.go
go run mean/mean.go
go run variance/variance.go
go run maxAndmin/maxmin.go
go run chart/chart.go

#cd chart
## in \chart
#go run chart.go
```