#!/bin/sh

# linting apply
## none-auth
curl -w '\n' -XPOST http://127.0.0.1:8080/linting/apply/linting-111
## missing-parameter
curl -w '\n' -XPOST -H "Cookie:user-id=20001" http://127.0.0.1:8080/linting/apply/linting-111
## none-administrator
curl -w '\n' -XPOST -H "Cookie:user-id=20002" http://127.0.0.1:8080/linting/apply/linting-111
## success
curl -w '\n' -XPOST -H "Cookie:user-id=20001" -d 'projId=proj-222' http://127.0.0.1:8080/linting/apply/linting-111

# linting create
## non-administrator
curl  -w '\n' -XPOST -H "Cookie:user-id=20002" http://127.0.0.1:8080/linting/create/new --form file=@yaml/openapi.yaml
## success-new
curl  -w '\n' -XPOST -H "Cookie:user-id=20001" http://127.0.0.1:8080/linting/create/new --form file=@yaml/openapi.yaml
## success-update, (replace correct linting-id with x)
curl  -w '\n' -XPOST -H "Cookie:user-id=20001" http://127.0.0.1:8080/linting/create/x --form file=@yaml/openapi.yaml


# linting view
## none-auth
curl  -w '\n' -XGET  http://127.0.0.1:8080/linting/view/fe298f67-be31-4973-9e85-b28bf3c63273
## no-admin preview
curl  -w '\n' -XGET -H "Cookie:user-id=20001" http://127.0.0.1:8080/linting/view/fe298f67-be31-4973-9e85-b28bf3c63273
## admin-preview
curl  -w '\n' -XGET -H "Cookie:user-id=20002" http://127.0.0.1:8080/linting/view/fe298f67-be31-4973-9e85-b28bf3c63273