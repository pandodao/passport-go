# Passport-go

## MVM login message

**Format**

```
"domain" wants you to sign in with your Ethereum account:
"address"

You'll login to Pando by the signature

URI: "uri"
Version: 1
Chain ID: "chain_id"
Nonce: "nonce"
Issued At: "issued_at"
Expiration Time: "expiration_time"
Resources:
- "resources"
```

```domain```: 请求签名的 domain, 可以是域名:端口 或者 ip:端口，如 localhost:3000

```address```: 签名用户的 address, 如 0xc86d59e6De389c031e79F469Ca3A4f8B9817efC5

```chain_id```: 数字，其中 mvm 的chain id 为 73927

```nonce```: 8 位以上的 数字/字母

```issued_at```: 签名 签发时间

```expiration_time```: 签名 过期时间

```resources```: 需要访问的资源

**Example**

```
pando-apps.aspens.rocks wants you to sign in with your Ethereum account:
0xc86d59e6De389c031e79F469Ca3A4f8B9817efC5

You'll login to Pando by the signature

URI: https://pando-apps.aspens.rocks
Version: 1
Chain ID: 73927
Nonce: oCxDubPgiNZdE8z71
Issued At: 2023-04-18T11:15:20+09:00
Expiration Time: 2023-04-18T11:18:20+09:00
Resources:
- https://pando-apps.aspens.rocks
```

**Signature**

Read the [Golang test](example/main.go#L48) as an example for signing messages.