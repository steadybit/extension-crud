
# :warning: This repo is archived. :warning:
### Please have a look at our [extension-scaffold](https://github.com/steadybit/extension-scaffold). 




# Steadybit extension-crud

A [Steadybit](https://www.steadybit.com/) discovery and action implementation acting as a CRUD backend. Useful to experiment freely with discovery and action terminology.

## Capabilities

 - Create, update, delete new entities to the in-memory database via actions.
 - List all added created entities as targets.

## Configuration

This extension is configurable via environment variables:

 - `STEADYBIT_EXTENSION_PORT` (defaults to `8091`)
 - `STEADYBIT_EXTENSION_INSTANCE_NAME` (defaults to `Animal Shelter`)
 - `STEADYBIT_EXTENSION_TARGET_TYPE` (defaults to `dog`)
 - `STEADYBIT_EXTENSION_TARGET_TYPE_LABEL` (defaults to `Dog`)

## Deployment

We recommend that you deploy the extension with our [official Helm chart](https://github.com/steadybit/helm-charts/tree/main/charts/steadybit-extension-crud).

## Agent Configuration

The Steadybit AWS agent needs to be configured to interact with the AWS extension by adding the following environment variables:

```shell
# Make sure to adapt the URLs and indices in the environment variables names as necessary for your setup

STEADYBIT_AGENT_ACTIONS_EXTENSIONS_0_URL=http://steadybit-extension-crud.steadybit-extension.svc.cluster.local:8080
STEADYBIT_AGENT_DISCOVERIES_EXTENSIONS_0_URL=http://steadybit-extension-crud.steadybit-extension.svc.cluster.local:8080
```

When leveraging our official Helm charts, you can set the configuration through additional environment variables on the agent:

```
--set agent.env[0].name=STEADYBIT_AGENT_ACTIONS_EXTENSIONS_0_URL \
--set agent.env[0].value="http://steadybit-extension-crud.steadybit-extension.svc.cluster.local:8080" \
--set agent.env[1].name=STEADYBIT_AGENT_DISCOVERIES_EXTENSIONS_0_URL \
--set agent.env[1].value="http://steadybit-extension-crud.steadybit-extension.svc.cluster.local:8080"
```
