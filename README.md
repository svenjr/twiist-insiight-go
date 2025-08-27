# Twiist Insiight Go

## Overview

This is a simple Go script which will use your Insiight credentials to make a call to get all your followings data.

I use this as a self-following mechanism so I can pull my own Twiist data in an easy way. The script sets up the OAuth handshake and gets a token for you that lasts an hour. There is also a function to refresh the token if you wish to alter this to setup a polling.

This was an exercise for me as I wanted to see how easy it would be to get this data for a possible Garmin watch integration (view only - no action).


## Usage
To run, you must obviously have go installed. Other than that, it should all JustWork.

Call the main go script with your email and password (this is not stored anywhere) and you will get a return of the full data response back.

```bash
❯ go run main.go notreal@email.com 'notArealPassword'
Logging in...
Login successful! Token expires in 3600 seconds
Fetching PWD overviews...
Raw API Response:
==================
[
  {
    "pwdId": "4a32d18e-f868-4847-9789-5f7fa33a32eb",
    "pwdNickname": "Svenjr",
    "status": {
      "date": "2025-08-27T12:27:35Z",
      "summary": {
        "glucoseDate": "2025-08-27T12:26:38Z",
        "glucoseUnit": "mg/dL",
        "cgmRateArrow": "→",
        "isBasalActive": true,
        "loopRingColor": "green",
        "glucoseQuantity": "88.0",
        "pumpBatteryLevel": "0.86",
        "closedLoopEnabled": true,
        "pumpEventsComplete": true,
        "glucoseHighlightState": "NoHighlight",
        "netBasal_UnitsPerHour": "-0.66",
        "lastCassetteChangeDate": "2025-08-26T14:54:35Z",
        "pumpCassetteVolume_Units": "125.59",
        "maximumBasalRate_UnitsPerHour": "6.0",
        "pumpCassetteFilledVolume_Units": "200.0"
      }
    }
  }
]
```
