# Replace the following placeholders with your actual values
# KEYCLOAK_URL="http://your-keycloak-server/auth"
# REALM_NAME="your-realm"
# CLIENT_ID="your-client-id"
# CLIENT_SECRET="your-client-secret"
# USER_SESSION_ID="user-session-id"  # Obtain the user session ID from the session


# Replace the following placeholders with your actual values
KEYCLOAK_URL="http://192.168.86.211:32088/auth"
REALM_NAME="nshub"
CLIENT_ID="auth-svc"
CLIENT_SECRET="3nUB5EUG3fenYFa9xqnF376PLFZHWxFV"
USERNAME="aaron"
PASSWORD="aaron"
USER_SESSION_ID="eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwdXZPdzZoMzNQYVNETFQzSVprZnRGdXpGRzBrOTBHU3hWZzhfYmt4c2t3In0.eyJleHAiOjE3MDY0MjYzOTgsImlhdCI6MTcwNjQyNjA5OCwianRpIjoiMzdmNTM1YWMtYzQ2Zi00N2NiLTg4Y2MtODk5N2E5ZGZmNzZlIiwiaXNzIjoiaHR0cDovLzE5Mi4xNjguODYuMjExOjMyMDg4L3JlYWxtcy9uc2h1YiIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJhMTJmZmZlOC1iZTIxLTQ5NWYtOGZhZS1hYWYxYTc5MmJhMGIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhdXRoLXN2YyIsInNlc3Npb25fc3RhdGUiOiJmYjI1MmEyMy0wYjBiLTRiMjktOTkwMi1mMzg4ZjhhNDM1MDIiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIi8qIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLW5zaHViIiwib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoiZW1haWwgcHJvZmlsZSIsInNpZCI6ImZiMjUyYTIzLTBiMGItNGIyOS05OTAyLWYzODhmOGE0MzUwMiIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYWFyb24ifQ.nDKQnIQGIzJbr7fejvaod-_zGuIAsbG5TgrehfuvhsOp21Jv5R7tkPrcCoxDhilHcfhSgDnf_W8ZGlL-kmap_59jL1_7Evxwb9O_WjfqV83dR1nqRrrWiWjmb8HIqfnzNbpV_5W1irvcaD7XWcKBl8MiT3lTWdBc3tvtDLGgkTPiqKqoMhi5cM5ntFrfBrrE1ph2zoSec7NTpsAhmQC-vaL5xg0LHoeVr1WeCLrA7boQmzgge9l0UA_ASmgoVL4RH3t3Fe4MoGcJ-1d2uvwbuvj3gPMoxcO8L5aGIuEPDpwR7hp12nzesjDsIN9s45oDDVqK_dLlym8De9ksTBJM2w"


# Perform the logout using the Keycloak endpoint
curl -X POST \
  "${KEYCLOAK_URL}/realms/${REALM_NAME}/protocol/openid-connect/logout" \
  -d "client_id=${CLIENT_ID}" \
  -d "client_secret=${CLIENT_SECRET}" \
  -d "refresh_token=${USER_SESSION_ID}"


# echo "Logging In... \n"

# # Get an access token using client credentials flow
# TOKEN_RESPONSE=$(curl -X POST \
#   "${KEYCLOAK_URL}/realms/${REALM_NAME}/protocol/openid-connect/token" \
#   -d "client_id=${CLIENT_ID}" \
#   -d "client_secret=${CLIENT_SECRET}" \
#   -d "grant_type=client_credentials")