# Insiight API Reverse Engineering Project

## Project Overview
This project reverse engineers the Insiight follower app API to access insulin pump data from a Twiist pump (Sequel Med Tech). The goal is to build a Go client that can authenticate and fetch real-time glucose/pump data, eventually for use in a Garmin watch app.

## Authentication Flow
- **Platform**: AWS Cognito User Pool authentication
- **Region**: us-east-1
- **User Pool ID**: `us-east-1_fnkWvSdfv`
- **Client ID**: `65ev2vbkr2mle7uu4cqkn7ohgl`
- **Auth Method**: `USER_PASSWORD_AUTH` flow
- **Token Expiry**: 3600 seconds (1 hour)
- **User Group**: "Follower"

## API Endpoints

### Authentication
```
POST https://cognito-idp.us-east-1.amazonaws.com/us-east-1_fnkWvSdfv
Headers:
  Content-Type: application/x-amz-json-1.1
  X-Amz-Target: AWSCognitoIdentityProviderService.InitiateAuth

Body:
{
  "AuthFlow": "USER_PASSWORD_AUTH",
  "ClientId": "65ev2vbkr2mle7uu4cqkn7ohgl",
  "AuthParameters": {
    "USERNAME": "your_email",
    "PASSWORD": "your_password"
  }
}
```

### Token Refresh
```
POST https://cognito-idp.us-east-1.amazonaws.com/us-east-1_fnkWvSdfv
Headers:
  Content-Type: application/x-amz-json-1.1
  X-Amz-Target: AWSCognitoIdentityProviderService.InitiateAuth

Body:
{
  "AuthFlow": "REFRESH_TOKEN_AUTH",
  "ClientId": "65ev2vbkr2mle7uu4cqkn7ohgl",
  "AuthParameters": {
    "REFRESH_TOKEN": "your_refresh_token"
  }
}
```

### Data Retrieval
```
GET https://follower-service.mytwiistportal.com/pwd/overviews
Headers:
  Authorization: Bearer {access_token}
  Content-Type: application/json
```

## Sample Response Data
```json
[{
  "pwdId": "4a32d18e-f868-4847-9789-5f7fa33a32eb",
  "pwdNickname": "Andrew McCall",
  "status": {
    "date": "2025-08-27T11:54:42Z",
    "summary": {
      "glucoseDate": "2025-08-27T11:54:39Z",
      "glucoseUnit": "mg/dL",
      "cgmRateArrow": "→",
      "glucoseQuantity": "101.0",
      "pumpBatteryLevel": "0.86",
      "closedLoopEnabled": true,
      "netBasal_UnitsPerHour": "0.58",
      "pumpCassetteVolume_Units": "126.39"
    }
  }
}]
```

## Go Implementation Status
✅ **Complete and Working**:
- Authentication with username/password
- Token refresh mechanism
- API data fetching
- JSON parsing and display

## Usage
```bash
# Run the client
go run main.go your_email@example.com your_password

# Example output:
# Login successful! Token expires in 3600 seconds
# Found 1 PWD overview(s):
# Nickname: Andrew McCall
# Glucose: 101.0 mg/dL (→)
# Battery: 0.86
# Last update: 2025-08-27T11:59:38Z
```

## Key Data Points Available
- **Glucose**: Current reading with trend arrow (→, ↑, ↓, etc.)
- **Battery**: Pump battery level (0.0-1.0)
- **Pump Status**: Basal rates, cassette volume, closed-loop status
- **Timestamps**: Last glucose reading time

## Next Steps for Garmin Integration
1. Create Garmin Connect IQ app project
2. Implement secure token storage on device
3. Add periodic data refresh (every 1-5 minutes)
4. Design watch face display for glucose/battery
5. Handle network connectivity and error states
6. Add alerts for high/low glucose readings

## Security Notes
- Refresh tokens should be stored securely on device
- Access tokens expire every hour and must be refreshed
- Username/password should never be stored in production
- API calls should include proper error handling and timeouts
