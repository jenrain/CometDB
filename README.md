# CometDB


## 兼容 Redis 协议
| RedisObject | command     |
|-------------|-------------|
| String      | Set、Get、StrLen |
| Hash        | HSet、HGet、HDel、HLen |
| Set         | SAdd、SIsMember、SRem、SCard |
| List        | LPush、RPush、LPop、RPop、LLen |
| ZSet        | ZAdd、ZScore |
| Key         | Ping、Del、Type |