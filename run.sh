# Replace the following placeholders with your actual values
KEYCLOAK_URL="http://192.168.86.211:32088/auth"
REALM_NAME="nshub"
CLIENT_ID="auth-svc"
CLIENT_SECRET="3nUB5EUG3fenYFa9xqnF376PLFZHWxFV"
USERNAME="aaron"
PASSWORD="aaron"


# Get an access token using client credentials flow
TOKEN_RESPONSE=$(curl -X POST \
  "${KEYCLOAK_URL}/realms/${REALM_NAME}/protocol/openid-connect/token" \
  -d "client_id=${CLIENT_ID}" \
  -d "client_secret=${CLIENT_SECRET}" \
  -d "grant_type=client_credentials")

ACCESS_TOKEN=$(echo "${TOKEN_RESPONSE}" | jq -r .access_token)

# Get user information including groups
USER_INFO=$(curl -X GET \
  "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/users?username=${USERNAME}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

# Extract user groups from the response
USER_GROUPS=$(echo "${USER_INFO}" | jq -r '.[0].groups | join(", ")')

# Print the user groups
echo "User Groups: ${USER_GROUPS}"
