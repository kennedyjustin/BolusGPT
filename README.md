# BolusGPT

TODO:
- Get OpenAPI spec
  - Make sure every description is accurate
  - Update with URL
  - Errors, etc.
- Figure out hosting
    - Simple EC2 Instance
    - No push to deploy server, dead simple, public subnet

- Document here in readme
- GPT needs nutrition database file

```
DEXCOM_USERNAME="<username>" DEXCOM_PASSWORD="<password>" BEARER_TOKEN="<token>" go run .
```

```
curl -k -X POST -H "Authorization: Bearer <token>" https://localhost:443/dose -d '{"total_grams_of_carbs": 20}'
```

```
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
