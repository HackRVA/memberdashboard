// datastore defines an interface for persisting data.
// packages within the datastore dir implement this interface.
// This allows us to have an actual database and an in_memory database.
// The in_memory database allows us to easily run the server with minimal dependencies
// -- which is nice for development.
package datastore
