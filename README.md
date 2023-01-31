# API

## `GET /me`

### Example Request

```
curl http://localhost:8080/me \
        -H "Authorization: 1"
```

### Example Response
```javascript
{
    "id": 1,
    "name" "Joe Bloggs",
    "roles": [
        {
            "companyID": "1",
            "role": "user"
        },
        {
            "companyID": "2",
            "role": "admin"
        }
    ]
}
```

## `GET /card`

### Example Request

```
curl http://localhost:8080/card \
        -H "Authorization: 1"
```

### Example Response

```javascript
{
    "cards": [
        {
            "companyID": "1",
            "card": {
                "maskedNumber": "**** **** **** 4444",
                "limit": 10000,
                "balance": 1241,
                "exp": "12/23"
            }
        },
        {
            "companyID": "2",
            "card": {
                "maskedNumber": "**** **** **** 1111",
                "limit": 500,
                "balance": 76,
                "exp": "10/23"
            }
        }
    ]
}
```

## `POST /card`

### Example Request

```
curl http://localhost:8080/card \
        -H "Content-Type: application/json" \
        -H "Authorization: 2" \
        -d '{"UserID": "1", "CompanyID": "1", "NewLimit": 20000}'
```

### Example Response

```javascript
{
    "maskedNumber": "**** **** **** 4444",
    "limit": 20000,
    "balance": 11241,
    "exp": "12/23"
}
```

# Execution

To start the auth server, run the following command:

    docker compose up

# TODO

- [ ] Run tests in Docker compose
- [ ] Use card repository for card handler
- [ ] Implement limit change API
- [ ] Think about different currencies

# Out of Scope

1. Secure authorization tokens (e.g. JWT)
1. DB layer