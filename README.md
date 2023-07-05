# CometDB

## 简介

CometDB 是一个高性能 KV 存储引擎，基于 bitcask 存储模型实现，数据以 append only 的方式追加写磁盘保证顺序 IO，内存中维护一棵 B+ 树索引帮助快速定位数据。

比内存数据库例如 Redis 有更高的可靠性，比磁盘数据库例如 MySQL 有更高的性能，比内存映射数据库例如 BoltDB 无随机 IO。

## 特性

- 采用顺序 IO 写数据，读写低延时。
- 索引维护在内存中，一次磁盘 IO 就可以快速访问到数据。
- 性能稳定，每次写都直接追加到当前数据文件末尾，写入速度不受数据量大小限制。
- 支持事务操作，事务满足 ACID 特性。
- 数据 merge 机制，清理无效数据，防止数据文件无限膨胀。
- 兼容 Redis 协议，满足多样化的需求，可以无缝平替 Redis。
- 提供 HTTP 接口，方便使用者访问。

## 使用场景

- 缓存系统：高性能和低读写，单机万级qps，适合存储频繁访问的热数据，提供快速的缓存读取操作，可作为 Redis 的一个替代方案。
- 日志存储：采用 append only 的方式高效写入数据，并且持久化存储到磁盘，可以确保日志的可靠存储和后续分析。
- 大 key 场景：将 key 和对应的索引都维护在内存中，而 value 存储在磁盘，可以利用磁盘的大空间来存储 value。

## 压测

### 环境

| Parameter  | Value            |
|------------|------------------|
| Go version | 1.20             |
| Machine    | MacBook Air 2020 |
| System     | MacOS 13.3       |
| CPU        | Apple M1         |
| Memory     | 16 GB            |

### BenchMark

**PUT**

```
BenchmarkPut-8   133208   9261 ns/op   4673 B/op   10 allocs/op
```

**GET**

```
Benchmark_Get-8   4411189   278.6 ns/op   135 B/op   4 allocs/op
```

**DEL**

```
Benchmark_Del-8   7592085   154.1 ns/op   135 B/op   4 allocs/op
```

### wrk压测http接口

wrk 每次 mock 一对随机的 k-v 来请求比较麻烦，所以我把生成随机 k-v 的逻辑放在了服务端：

- PUT: 服务端每接收到一条请求，都会生成一对 0-10000 的随机值 put 进数据库。
- GET: 服务端随机生成一个 0-10000 的 key， 然后从数据库中读取。
- DELETE: 做法与 get 请求一样。

```go
# PUT请求的 handler
k := strconv.Itoa(rand.Intn(10000))
v := strconv.Itoa(rand.Intn(10000))
db.Put([]byte("key"+k), []byte("value"+v))

# GET请求的 handler
k := strconv.Itoa(rand.Intn(10000))
v, _ := db.Get([]byte("key"+k))
```

**PUT**

```
wrk -t8 -c100 -d20s -s test.lua http://127.0.0.1:8080/put
```

test.lua：
```lua
wrk.method = "PUT"
wrk.body = '{"key": "value"}'
wrk.headers["Content-Type"] = "application/json"
```

测试结果
```
Running 20s test @ http://127.0.0.1:8080/put
  8 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    18.41ms   41.44ms 419.41ms   91.21%
    Req/Sec     1.78k     1.46k    6.63k    83.27%
  270102 requests in 20.04s, 38.90MB read
Requests/sec:  13476.52
Transfer/sec:      1.94MB
```

**GET**

```
wrk -t8 -c100 -d20s -s test.lua "http://127.0.0.1:8080/get?key=key"
```

测试结果
```
Running 20s test @ http://127.0.0.1:8080/get?key=key
  8 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    22.18ms   51.75ms 484.06ms   90.90%
    Req/Sec     1.75k     1.56k    6.90k    84.48%
  261294 requests in 20.04s, 34.58MB read
Requests/sec:  13038.23
Transfer/sec:      1.73MB
```

**DEL**

```
wrk -t8 -c100 -d20s -s test.lua "http://127.0.0.1:8080/del?key=key"
```

test.lua:
```lua
wrk.method = "DELETE"
```

测试结果
```
Running 20s test @ http://127.0.0.1:8080/del?key=key
  8 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    20.49ms   45.71ms 391.10ms   90.73%
    Req/Sec     1.76k     1.59k    7.35k    85.85%
  267350 requests in 20.04s, 35.70MB read
Requests/sec:  13338.90
Transfer/sec:      1.78MB
```

## 兼容 Redis 命令
| RedisObject | command     |
|-------------|-------------|
| String      | Set、Get、StrLen |
| Hash        | HSet、HGet、HDel、HLen |
| Set         | SAdd、SIsMember、SRem、SCard |
| List        | LPush、RPush、LPop、RPop、LLen |
| ZSet        | ZAdd、ZScore |
| Key         | Ping、Del、Type |

## 参考
[bitcask-intro](https://riak.com/assets/bitcask-intro.pdf)