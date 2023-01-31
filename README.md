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

    docker compose build && docker compose up

This will also run the tests as part of the build.

# Sanity Test

A simple bash script is included, to test the API. This can be run like so:

    ./sanity-test.sh

# Known Limitations / Future Considerations

1. May want to move to a more powerful routing library if use of path variables is needed [e.g. chi](https://github.com/go-chi/chi)
1. May need to consider how different currencies and card formats are handled
1. For the purposes of limiting scope, an actual DB was not used
1. For the purposes of simplicity, a user need only provide their user ID as an auth token. This is obviously not secure.