# Changelog

## v1.2.2

- migration to new unified steadybit actionIds and targetTypes

## v1.2.1

- update dependencies

## v1.2.0

 - Print build information on extension startup.

## v1.1.0

- Support creation of a TLS server through the environment variables `STEADYBIT_EXTENSION_TLS_SERVER_CERT` and `STEADYBIT_EXTENSION_TLS_SERVER_KEY`. Both environment variables must refer to files containing the certificate and key in PEM format.
- Support mutual TLS through the environment variable `STEADYBIT_EXTENSION_TLS_CLIENT_CAS`. The environment must refer to a comma-separated list of files containing allowed clients' CA certificates in PEM format.

## v1.0.0

 - Initial release