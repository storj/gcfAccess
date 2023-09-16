# Copyright (C) 2023 Storj Labs, Inc.
# See LICENSE for copying information.

import requests
import os

# Define the GCF endpoint URL
domain = os.environ.get("GCF_DOMAIN")
gcf_url = "https://" + domain + ".cloudfunctions.net/OverrideEncryption"

# Define the request payload
payload = {
    "access_grant": os.environ.get("ACCESS"),
    "passphrase": "supersecret",
        "path": {
            "bucket": "storj",
            "prefix": "/new/"
    }
}

# Send a POST request to the GCF endpoint
response = requests.post(gcf_url, json=payload)

# Check the response status code
if response.status_code == 200:
    print("Access encryption overriden successfully")
    print(response.text)
else:
    print("Error:", response.status_code)
    print(response.text)