# teeny-weeny-url-api
An API to shorten URLs as well as fetch the actual URLs from the shortened ones

Shortening is based on MD5. The algorithm takes the first 8 bytes of the generated hash.

Contains 3 routes:
  1. /create: Creates a short_url from the URL passed in the body (json of the body: {"URL": "<url_that_needs_to_be_shortened>"}  
  2. /db: Fetches list of stored URLs
  3. /{s_url}: redirects browser to actual URl whose shortened URL-code was based as {s_url}
