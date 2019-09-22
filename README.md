# go-api-example

This example application provides a set of http endpoints for CRUD operations against a `message` resource (`src`). Included is tooling for running migrations (`tasks/migrate.go`), library functions for common database logic (`lib/db.go`), and an integration test suite (`tests`) to validate the API's functionality.

## Architecture

### Infrastructure

Two primary components exist for this application's infrastructure:
* a Go http server (`net/http`)
* a `postgres` database instance

This is a fairly typical set up for a generalist CRUD api. 

Using Go as the implementation language for the web server and business logic gives us performance characteristics which are comparable to Java or C# (ie, fast enough for most things) and access to a (subjectively) enjoyable ecosystem of tooling.

Using postgres for the data store gives us a generalist set of functionality with pleasant guarantees around CRUD operations (e.g, transactions, though they are not necessary for this application) while being able to satisfy our performance requirements (i.e, this API doesn't need to be "webscale"). Moreover, SQL is the most popular query language, meaning nobody on a given software team should be too grumpy about having to use it.

### Software

There are two primary components of this application to discuss with regards to the software architecture: the API itself, and the test suite.

#### API

Given the API requirements are fairly simple, it was developed to use a minimal architecture with idiomatic abstractions. In other words, nothing should be too scary for the unfamiliar reader of this code base - if it is, my apologies.

At a high level and following the happy path, an incoming request is sent, and then:
* the request maps to a route
* the route maps to a handler
* the handler serializes the request to a dto (optional)
* dto validation occurs (optional)
* a repository performs a CRUD operations against the database
* the results of that query are returned as an http response

The `Server` struct (`src/api/server.go`) provides a mechanism for dependency injection of anything we might want to share across handlers, and is the top level object for the application.

#### Test suite

This repository contains an integration test suite which validates the business logic against a real database, using real http requests. Integration tests, while being the most burdensome form testing to set up (and run in CI/CD), also offer the greatest guarantee that something does what it's supposed to, and are reasonably simple to write (i.e, no setting up and maintaining mocks). Given this, I opted to use them.

For some functionality, writing mock tests for the separate components (e.g mocking `sqlx.DB`) and unit tests for logic with no external dependencies (e.g `IsPalindrome()` on the `message` struct) would be desirable as well.

## Running the api locally

If you want to get things running to send http requests to the API, follow these steps. This application has been developed to make this as painless as possible.

Dependencies:
* Docker
* a macOS laptop (I haven't tested this on Linux or Windows, so no guarantees)

Commands:
* `make build`
    * This builds the docker images defined in `docker-compose.yml` using the Dockerfiles in `docker`
* `make up`
    * This runs the images (database, api, tests, migrate). Observe stdout to confirm migrations run properly and that the test suite passes
* `curl -i localhost:3000/healthz`
    * This confirms that you can talk to the api. You should see an `ok` response from the `resource` field

## REST API documentation

For a production application, we'd probably use Swagger or OpenAPI or some form of tooling to auto generate this based on code we've written. For this example repo, I will provide a route, it's verb, a description, and the request body parameters where appropriate.

* `GET /healthz`
    * Health check route; returns ok if the server is alive
    * No request body
    * Returns `{ "resource": "ok", "error": null }`
    
* `GET /messages`
    * Returns a list of all messages
    * No request body
    * Returns `{ "resource": [{ ...message }], "error": null }`
    
* `GET /message/{id}`
    * Returns the message specified by `{id}`
    * No request body
    * Returns `{ "resource": { ...message }, "error": null }`
    
* `POST /message`
    * Creates a message
    * Body: `{ "content": string }`
    * Returns `{ "resource": { ...message }, "error": null }`
    
* `PATCH /message/{id}`
    * Updates a message specified by `{id}`
    * Body: `{ "content": string }`
    * Returns `{ "resource": { ...message }, "error": null }`
    
* `DELETE /message/{id}`
    * Deletes a message specified by `{id}`
    * No request body
    * Returns `{ "resource": null, "error": null }`
    

The shape of a `message` is the following:

```json
{
    "id": 10,
    "content": "firetruck",
    "isPalindrome": false,
    "createdAt": "2019-09-22T18:43:03.284233Z",
    "updatedAt": "2019-09-22T18:44:50.182602Z"
}
```
