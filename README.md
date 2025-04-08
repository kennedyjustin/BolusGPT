# BolusGPT

TODO:
- Get OpenAPI spec - document exact use cases (onboarding, dosing, confirm dosing, getting nutrition, etc.)
- exercise use string enums to convert
- Document here in readme
- GPT needs nutrition database file

```
DEXCOM_USERNAME="<username>" DEXCOM_PASSWORD="<password>" BEARER_TOKEN="<token>" go run .
```

```
curl -X POST -H "Authorization: Bearer <token>" localhost:8080/dose -d '{"total_grams_of_carbs": 20}'
```