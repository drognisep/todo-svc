# TODO Service

This is meant to represent an existing service in a wider system with which the Kotlin UI side will integrate.
It will expose a JSON REST and gRPC interface to a shared todo list.

## Data

### Todo Item
A Todo Item represents some unit of work.

#### Attributes

| Name    | Type   | Description                                           |
|---------|--------|-------------------------------------------------------|
| ID      | uint64 | The unique identifier for a Task Item                 |
| Summary | string | A human readable string that identifies the Task Item |
| Done    | bool   | Whether this Todo Item has been completed or not      |

#### Operations
The scope of this service is to expose basic CRUD operations over Todo Items

* Create a new Todo Item
* Retrieve a list or specific Todo Item
* Update a specific Todo Item
  * Update a Todo Item's Summary attribute
  * Update a Todo Item's Done attribute
* Delete a specific Todo Item
  * No batch deletion operations will be exposed
