# Gin Templ Clerk

## Dependencies
- Gin
- Templ
- Clerk account and sdk
- [Swag](https://github.com/swaggo/swag)
- [chromedp](https://github.com/chromedp/chromedp)

*Optional*:
- [Air](https://github.com/cosmtrek/air) (for hot-reloading Go projects)


## Files examples


_Note: .env file is not the best idea but good for quick dev - Remember to always save your secrets in a safe place._ 


#### clerk_public_key.pem
Key in .pem format

```
-----BEGIN PUBLIC KEY-----
MIJBIjANBgkqhkiG9w...
-----END PUBLIC KEY-----
```

#### .env
```
CLERK_API_KEY=sk_test_
JWT_PUBLIC_KEY_PATH=./clerk_public_key.pem
DB_USER=
DB_PASSWORD=
DB_HOST=
DB_PORT=
DB_NAME=
```