# sentinel-proxy

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)

Easy proxy for redis sentinel.

Main purpose of the proxy is easy work with redis sentinel without changing the client code.

Especially actual for clouds providers without managed redis e.g. Yandex.Cloud and when working with its own redis sentinel cluster (vps)

Key features:
```
[x] Initial bootstrap with ping all redis servers
[x] Correct handling of cluster anomaly
[x] Selecting preferred redis using the quorum result
[x] Easy Graylog integration
[x] Actual binary releases
```

Work principles:
```
client code -> local sentinel proxy (*persistent or not, it doesn't matter)
local sentinel proxy -> remote redis sentinels (get actual redis)
local sentinel proxy -> actual redis (proxy connection in both direction)
```
