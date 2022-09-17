# Kong API Catalog Harvester

## Description

Kong plugin to harvest your API catalog. Turn actual API traffic into valuable data. 

## Key Objectives

- Provide a way to collect & reconstruct an API catalog.
- Provide a way to feed a data catalog service to govern the data being processed by your APIs.
- Turn actual API traffic into an API specification for your service.
- Detect API changes and check conformance from your live services.
- Track sensitive data flows

## Getting Started

How is this working ?

To demonstrate the API harvesting through our plugin we have created an end-to-end testing Docker compose that put the following into action.

- From a given sample API specification (aka the Swagger Pet Store) we generate a collection with Postman 
- A Docker Newman runtime will be used to execute the collection
- From the same sample API specification, we generate an API mock with Stoplight Prism
- Kong will route the generated API request to the mock
- This plugin will harvest a specification based on the captured API traffic

In the `repository` folder you will find the harvested API specification and in the `snapshots` you will find the learning path of this plugin.
A new API snapshot being created each time an update is provided by the harvesting.

![alt text](./how-is-this-iworking.png "Description")

### References

- https://www.postman.com
- https://konghq.com
- https://learning.postman.com/docs/running-collections/using-newman-cli/command-line-integration-with-newman/
- https://stoplight.io/open-source/prism

Build the required images

``` bash
docker-compose build
```

Run the stack

``` bash
docker-compose up
```

Debug the stack

``` bash
docker-compose  -f docker-compose.yml -f docker-compose-debug.yml build
```

``` bash
docker-compose  -f docker-compose.yml -f docker-compose-debug.yml up
```

Then attach your remote Go debugger to localhost:40000 !

See below end-to-end test execution

![alt text](./compose.gif "Up")

View the harvested API specification and track change(s) with https://api-diff.io !

![alt text](./diff.gif "Compare the specifications")

## Goal

The aim of this plugin is to create a way to extract relevant data from Kong in order to be processed externally. 

## Inspirations & Integrations

See the following vendor's / solutions that are linked to this plugin goal.

https://www.useoptic.com/

Optic makes it easy to Track and Review all your API changes before they get released. Start working API-first and ship better APIs, faster.

https://getmizu.io/

A simple-yet-powerful API traffic viewer for Kubernetes to help you troubleshoot and debug your APIs.

https://www.levo.ai/

Take control of API sprawl & proactively mitigate API risk. Ship secure and resilient APIs into production.

https://www.traceable.ai/

Automatic and Continuous API discovery that provides comprehensive visibility into all APIs, sensitive data flows, and risk posture – even as your environment changes.

https://www.talend.com/products/data-catalog/

Crawl, profile, organize, link, and enrich all your data at speed

https://www.apiclarity.io/

Open source for API traffic visibility in K8s clusters

https://stoplight.io/open-source/prism

Accelerate API development with realistic mock servers, powered by OpenAPI documents.

https://roadmap.stoplight.io/c/66-learning-recording?utm_source=github&utm_medium=prism&utm_campaign=readme

https://api-diff.io

Diff two API versions in seconds and see what has changed.

https://www.akitasoftware.com

Powered by eBPF and API traffic analysis, Akita makes it possible for you to understand your API behavior, without having to change code or make custom dashboards.

