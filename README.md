# Passport-go

## 如何生成 MVM jwt Token

**HEADER**

```json
{
  "alg": "ESMVM",
  "typ": "JWT"
}
```

**PAYLOAD**

```json
{
  "aud": "a753e0eb-3010-4c4a-a7b2-a7bda4063f62",
  "exp": 1676050274,
  "jti": "b3c99d22-959f-4f47-b029-77eb2dce3d60",
  "iss": "0x486cFFDDE71B9655D85D1c77C1ad2966a9C5ea49",
  "non": 1,
  "ver": 1
}
```

```aud```: 为应用的 appid

```exp```: 为过期时间，单位为秒

```jti```: 为 request id

```iss```: 为用户的公钥地址

```non```: 为随机数

```ver```: 为版本号

**SIGNATURE**

```javascript
sig = this.signer.signMessage(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload)
)
```

最后 ```base64UrlEncode(header) + "." + base64UrlEncode(payload) + "." + base64UrlEncode(sig)``` 即为最终 jwt token
