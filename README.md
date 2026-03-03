# external-auth


Login Flow -
1. External Traffic (login Request)-> Envoy -> External Auth  -> Key Cloak -> External Auth (will cache the token) -> Return to the user via Envoy 

All other protected endpoint ->
2. External Traffic ( Protected endpoint)-> Envoy -> External Auth (if token is not in the cache then fwd to the keyCloak) if invalid token request is more then 10 then deny the token request if not found in yhe cache 

when auth service return the response to envoy, it would add the following params - 

1. User id 
2. Group Id 
3. Tenent id 



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



curl -X POST "http://localhost:8089/login" -d '{"username": ""}'



