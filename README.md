# access
Google Cloud Function definitions for manipulating Access Grants

Real world Storj workflows often use S3 for all of their upload and 
download needs. However, there is still the burden of customer 
credential management.  Various language-specific and community 
maintained bindings are up to this task.  These bindings can be 
difficult to keep up to date.

This tool exposes a simple option for HTTP + JSON based Storj 
Access Grant management that doesn't require the complexities of 
uplink-c.