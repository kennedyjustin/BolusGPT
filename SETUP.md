# Setup Guide

```
DEXCOM_USERNAME="<username>" DEXCOM_PASSWORD="<password>" BEARER_TOKEN="<token>" go run .
```

```
curl -X POST -H "Authorization: Bearer <token>" http://localhost:80/dose -d '{"total_grams_of_carbs": 20}'
```