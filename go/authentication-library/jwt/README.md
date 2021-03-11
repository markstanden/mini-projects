# jwt
## JWT creation library to be shared across my apps, with sane defaults

A library for basic jwt creation and parsing written in go

I have gone back to basics reading the JWT spec and getting inspiration from jwt.io and Auth0.com

The idea was to create a standard JWT token to be used across multiple projects, using security best practices and sane defaults.

also, at least initially the idea is to have a minimal number of exportable functions, to allow the workings to invisibly change as required.
