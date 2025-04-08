# BolusGPT

TODO:
- Get OpenAPI spec
  - Make sure every description is accurate
  - Update with URL
  - Errors, etc.
- find hosting service
- Document here in readme
- GPT needs nutrition database file

```
DEXCOM_USERNAME="<username>" DEXCOM_PASSWORD="<password>" BEARER_TOKEN="<token>" go run .
```

```
curl -X POST -H "Authorization: Bearer <token>" localhost:8080/dose -d '{"total_grams_of_carbs": 20}'
```