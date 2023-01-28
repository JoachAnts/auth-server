# API

## GET /me

### Example Request

```
curl http://localhost:8080/me
        -H "Authorization: 1"
```

### Example Response

```javascript
{
    "id": 1,
    "name" "Joe Bloggs",
    "role": "user"
}
```

## GET /card

### Example Request

```
curl http://localhost:8080/card
        -H "Authorization: 1"
```

### Example Response

```javascript
{
    "card": {
        "maskedNumber": "**** **** **** 4444",
        "limit": 10000,
        "balance": 1241,
        "exp": "12/23"
    }
}
```

# Execution

To start the auth server, run the following command:

    docker compose up
