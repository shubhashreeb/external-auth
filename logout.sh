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
USER_SESSION_ID="user-session-id"


# Perform the logout using the Keycloak endpoint
curl -X POST \
  "${KEYCLOAK_URL}/realms/${REALM_NAME}/protocol/openid-connect/logout" \
  -d "client_id=${CLIENT_ID}" \
  -d "client_secret=${CLIENT_SECRET}" \
  -d "refresh_token=${USER_SESSION_ID}"


# echo "Logging In... \n"

# # Get an access token using client credentials flow
# TOKEN_RESPONSE=$(curl -X POST \
#   "http://192.168.86.211:32088/auth/realms/${REALM_NAME}/protocol/openid-connect/logout" \
#   -d "client_id=auth-svc" \
#   -d "client_secret=3nUB5EUG3fenYFa9xqnF376PLFZHWxFV" \
#   -d "refresh_token="user-session-id"
#   

# curl -X POST \
#   "http://192.168.86.211:32088/auth/realms/nshub/protocol/openid-connect/logout" \
#  -d "client_id=auth-svc" \
#  -d "client_secret=3nUB5EUG3fenYFa9xqnF376PLFZHWxFV" \
#  -d "refresh_token=eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI5NjVmYjRjZS00MTYyLTRiNDItOTZlNy01NTEwZDQ0ZDI0MDYifQ.eyJleHAiOjE3MDY0MzQxODksImlhdCI6MTcwNjQzMjM4OSwianRpIjoiZWQ4NDZjNzUtNmU5Yy00YzQzLWE1MjQtZDMyOGEzN2U5NGZjIiwiaXNzIjoiaHR0cDovLzE5Mi4xNjguODYuMjExOjMyMDg4L3JlYWxtcy9uc2h1YiIsImF1ZCI6Imh0dHA6Ly8xOTIuMTY4Ljg2LjIxMTozMjA4OC9yZWFsbXMvbnNodWIiLCJzdWIiOiJhMTJmZmZlOC1iZTIxLTQ5NWYtOGZhZS1hYWYxYTc5MmJhMGIiLCJ0eXAiOiJSZWZyZXNoIiwiYXpwIjoiYXV0aC1zdmMiLCJzZXNzaW9uX3N0YXRlIjoiMDk1NmU1NGEtMmNlNi00NzRhLTg1ZjUtZTU5MzE2ZmRjM2ExIiwic2NvcGUiOiJlbWFpbCBwcm9maWxlIiwic2lkIjoiMDk1NmU1NGEtMmNlNi00NzRhLTg1ZjUtZTU5MzE2ZmRjM2ExIn0.mBbCaZ8p9l_Mw2hJZWaQ2UqS0C5CdeApMQsgUyDzwaM"