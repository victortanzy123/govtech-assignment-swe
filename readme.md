# GovTech - Golang Assignment (Full Name: Tan Zuu-Yuaan@Victor Tan)

## How to run the application

1.  Clone the application with `git@github.com:victortanzy123/govtech-assignment-swe.git`

2.  Navigate to `sql-dump` folder, use the mySQL dump files `sql-teach-dump.sql`, `sql-suspend-dump.sql` & `sql-notification-dump.sql` to create the respective tables within database, create the tables without inserting any data.

3.  Once this application is cloned and mySQL database has been set up accordingly (with 3 tables), amend the Connection String inside `config.go` which is located within `config` folder to the appropriate mysql username, password and database name on line 13.

4.  To run the application, please use the following command -

        ```shell
            go run main.go
        ```

    > Note: By default the port number has been set to **8080**.

5.  All schemas/struct can be found in `model.go` inside `model` folder, whereas all API endpoint logic are located within `controller.go` inside `controller` folder.

## User Story Endpoints Description

### Student Registration

#### As a teacher, I want to register one or more students to a specified teacher.

```
    Endpoint: POST http://localhost:8080/api/register
    Headers: Content-Type: application/json
    Success response status: HTTP 204
    Body - (content-type = application/json)
```

```JSON
    {
    "teacher": "t1@gmail.com",
    "students": ["s1@gmail.com","s2@gmail.com"]
    }

```

### Fetch Common Students

#### As a teacher, I want to retrieve a list of students common to a given list of teachers (i.e. retrieve students who are registered to ALL of the given teachers).

```
    URL - *http://localhost:12345/api/entry?id=1*
    Method - GET
    Endpoint: GET http://localhost:8080/api/commonstudents
    Success response status: HTTP 200
    Request example 1: GET /api/commonstudents?teacher=t1%40gmail.com

```

### Suspend Student

#### As a teacher, I want to suspend a specified student.

```
    Endpoint: POST http://localhost:8080/api/suspend
    Headers: Content-Type: application/json
    Success response status: HTTP 204
    Body - (content-type = application/json)
```

```JSON
    {
    "student": "s1@gmail.com"
    }
```

### Retrieve For Notification

#### As a teacher, I want to retrieve a list of students who can receive a given notification.

```
    Endpoint: POST http://localhost:8080/api/retrievefornotifications
    Headers: Content-Type: application/json
    Success response status: HTTP 200
    Body - (content-type = application/json)
```

```JSON
    {
    "teacher": "t1@gmail.com",
    "notification": "hello world bye @s1@gmail.com @s2@gmail.com @s3@gmail.com"
    }
```

## Unit Test Cases (All Endpoints)

To run all the unit test cases, please do the following -

1. Ensure the database with all respective tables has been set up and are all empty.
2. Ensure connection with the database can be establish via `go run main.go`
3. `go test ./controller`

## Remarks:

1. Ensure that the database (with 3 tables - Teach, Suspend & Notification) is deliberately chosen given how all teachers and students are represented by their email, which is unique to every entity and can be used as a primary key to represent their identity which adequately serves the required user stories. In reality, a `Students` and `Teachers` would be created with an `id` as a primary key to store all of the personal relevant information.

2. The unit tests are deliberately designed in sequential order, where some of the test cases will require the actions of the previous unit test case to simulate an entire user story flow, and hence **before each run, an empty database with the required tables** from `sql-dump` is required for setting up for all test cases to pass. Additionally, given that the assignment document has specified that Govtech's end will be running your own set of test cases, initial population of table data would not be required.

3. The unit test cases also include cases that should fail i.e. bad request, attempt to insert an existing entry and should emit their respective HTTP code & customised error message. More details can be found in the comments & code from `controller.go` inside the `controller` folder.

### Lastly, should there be any further clarifications, please feel free to reach out to me via my email. Thanks and I hope everything works well for you!
