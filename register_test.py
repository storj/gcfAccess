# Copyright (C) 2023 Storj Labs, Inc.
# See LICENSE for copying information.

import requests
import os

# Define the GCF endpoint URL
domain = os.environ.get("GCF_DOMAIN")
gcf_url = "https://" + domain + ".cloudfunctions.net/RegisterAccess"

# Define the request payload
payload = {
    "access_grant": os.environ.get("ACCESS"),
    "public": True,
}

# Send a POST request to the GCF endpoint
response = requests.post(gcf_url, json=payload)

# Check the response status code
if response.status_code == 200:
    print("Access registered successfully")
    print(response.text)
else:
    print("Error:", response.status_code)
    print(response.text)