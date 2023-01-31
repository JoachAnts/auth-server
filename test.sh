echo -e "====Test 1====="
curl http://localhost:8080/me \
        -H "Authorization: 1"
echo -e ''

echo -e "====Test 2====="
curl http://localhost:8080/card \
        -H "Authorization: 1"
echo -e ''

echo -e "====Test 3====="
curl http://localhost:8080/card \
        -H "Content-Type: application/json" \
        -H "Authorization: 2" \
        -d '{"UserID": "1", "CompanyID": "1", "NewLimit": 30000}'
echo -e ''

echo -e "====Test 4====="
curl http://localhost:8080/card \
        -H "Authorization: 1"
echo -e ''