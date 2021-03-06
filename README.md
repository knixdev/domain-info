# Domain-Info

A Simple REST API built in Gin-Gonic, allows to you check if a domain
is on the Majestic Million or is identified as a known Dynamic DNS
provider using a precompiled list of more than ~30000 Dynamic DNS providers.
This small stack utilizes Redis to speed up searches.

Steps to use it:
1. Clone this repository
2. Ensure docker-compose is installed, and run docker-compose up from the root
3. Browser to http://127.0.0.1:8080/majestic/google.com
   or http://127.0.0.1:8080/dynamicdns/google.com

**NOTE:** If you wish to update the .aof file checked into this repository
use the command line flags specified in data/main.go and update the .aof file.
