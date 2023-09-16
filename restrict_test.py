# Copyright (C) 2023 Storj Labs, Inc.
# See LICENSE for copying information.

import requests
import os

# Define the GCF endpoint URL
domain = os.environ.get("GCF_DOMAIN")
gcf_url = "https://" + domain + ".cloudfunctions.net/RestrictAccess"

# Define the request payload
payload = {
    "access_grant": os.environ.get("ACCESS"),
    "paths": [
        {
            "bucket": "storj",
            "prefix": "/"
        }
    ],
    "permission": {
        "allowDownload": True,
        "allowUpload": True,
        "allowList": True,
        "allowDelete": True,
        "notBefore": "2023-09-15T12:34:56Z",
        "notAfter": "2023-09-16T12:34:56Z",
        "maxObjectTTL": 3600
    }
}

# Send a POST request to the GCF endpoint
response = requests.post(gcf_url, json=payload)

# Check the response status code
if response.status_code == 200:
    print("Access grant restricted successfully")
    print(response.text)
else:
    print("Error:", response.status_code)
    print(response.text)