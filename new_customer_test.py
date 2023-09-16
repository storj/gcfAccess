# Copyright (C) 2023 Storj Labs, Inc.
# See LICENSE for copying information.

import requests
import os

# Define the GCF endpoint URL
domain = os.environ.get("GCF_DOMAIN")
gcf_url = "https://" + domain + ".cloudfunctions.net/NewS3Customer"

# Define the request payload
payload = {
    "access_grant": os.environ.get("ACCESS"),
    "path": {
        "bucket": "storj",
        "prefix": "/"
    },
    "permission": {
        "allowDownload": True,
        "allowUpload": True,
        "allowList": True,
        "allowDelete": True,
    },
    "passphrase": "supersecret",
    "public": False
}


# Send a POST request to the GCF endpoint
response = requests.post(gcf_url, json=payload)

# Check the response status code
if response.status_code == 200:
    print("New Gateway-MT customer credentials created successfully")
    print(response.text)
else:
    print("Error:", response.status_code)
    print(response.text)