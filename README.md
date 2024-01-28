# external-auth


Login Flow -
1. External Traffic (login Request)-> Envoy -> External Auth  -> Key Cloak -> External Auth (will cache the token) -> Return to the user via Envoy 

All other protected endpoint ->
2. External Traffic ( Protected endpoint)-> Envoy -> External Auth (if token is not in the cache then fwd to the keyClack) if invalid token request is more then 10 then deny the token request if not found in yhe cache 

when auth service return the response to envoy, it would add the following param - 

1. User id 
2. Group Id 
3. Tanent id 
4. 



infra -> 

1. curl Request to the k8s keycloak

curl -X POST \
  'http://192.168.86.211:32088/realms/nshub/protocol/openid-connect/token' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'grant_type=password&client_id=auth-svc&client_secret=3nUB5EUG3fenYFa9xqnF376PLFZHWxFV&username=aaron&password=aaron' -v


curl -X POST \
  'http://10.99.24.16:80/realms/nshub/protocol/openid-connect/token' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'grant_type=password&client_id=auth-svc&client_secret=3nUB5EUG3fenYFa9xqnF376PLFZHWxFV&username=aaron&password=aaron' -v




curl --location --request GET 'http://192.168.86.211:8090/logout' \
--header 'Content-Type: application/json' \
--data '{"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJkMjE4ZjQxNS0wYTFmLTQ5ZmUtYjI1NS0wM2FkMzFiMGIyY2YifQ.eyJleHAiOjE3MDY0Njc2MzgsImlhdCI6MTcwNjQ2NTgzOCwianRpIjoiMDY3NjdiZDctMDFhMy00NGJiLWFiNDYtMzc2ZDBmNzRjZTM0IiwiaXNzIjoiaHR0cDovLzE5Mi4xNjguODYuMjExOjMyMDg4L3JlYWxtcy9uc2h1YiIsImF1ZCI6Imh0dHA6Ly8xOTIuMTY4Ljg2LjIxMTozMjA4OC9yZWFsbXMvbnNodWIiLCJzdWIiOiI4ZmRiZTU0OC1jNjc2LTQyYTktYWM0Ni1lMzk5NDcyZmIzNTkiLCJ0eXAiOiJSZWZyZXNoIiwiYXpwIjoiYXV0aC1zdmMiLCJzZXNzaW9uX3N0YXRlIjoiMDE0YTY2NGMtMDk5Ny00ZjBmLThhNzctNTUwYmExYWExMTczIiwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInNpZCI6IjAxNGE2NjRjLTA5OTctNGYwZi04YTc3LTU1MGJhMWFhMTE3MyJ9.U1i01FTyNNOtqgTestD1QIOMiKix9Wb1R5nythsk-ys"}'

