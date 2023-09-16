# Experimental!!

Google Cloud Function definitions for manipulating Access Grants

Real world Storj workflows often use S3 for all of their upload and 
download needs. However, S3 does not address other user needs, such 
as the ability to manage Access Grants.

This tool exposes a simple option for HTTP + JSON based Storj 
Access Grant management.  Note that using a 3rd party service such as 
GCF breaks some of the end-to-end ideals of the Storj network.  

It is recommended that you use a language specific library based on 
uplink-c whenever possible instead.
