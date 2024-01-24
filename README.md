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

1. 

