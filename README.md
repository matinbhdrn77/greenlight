## Chapter2

Any packages which live under *internal* directory can only be imported by code inside the parent of the internal directory.

`func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request)` => and idiomatic way to make dependencies available to our handlers without resorting to global variables or closures — any dependency that the healthcheckHandler needs can simply be included as a field in the application struct when we initialize it in main() .

The *httprouter* package provides a few configuration options that you can use to customize the behavior of your application further, including enabling *trailing slash redirects* and *enabling automatic URL path cleaning*.

## Chapter3 Sending JSON Response

JSON is just text. key must be string. value can be any JS object.
You would write any other text response: using `w.Write()` , `ioWriteString()` or one of the `fmt.Fprint` functions. But you have to set `Content-Type: application/json`

Go’s **encoding/json** package provides two options for encoding things toJSON. You can either call the *json.Marshal()* function, or you can declare and use a *json.Encoder* type.

`func Marshal(v interface{}) ([]byte, error)`
`err := json.NewEncoder(w).Encode(data)` => write to writer in a single step means no time for setting header

All the fields in our struct are exported (i.e. start with a capital letter), which is necessary for them to be visible to Go’s `encoding/json` package. Any fields which aren’t exported won’t be included when encoding a struct to JSON.

If a struct field doesn’t have an explicit value set, then the JSON-encoding of the zero value for the field will appear in the output.

Customize the JSON by annotating the fields with *struct tags*. must common use is to change the key names. ``json:"title"``

Control the visibility of individual struct fields in the JSON by using the omitempty and - *struct tag directives*.

In contrast the omitempty directive hides a field in the JSON output if and only if the struct field value is empty, where empty is defined as being:
- Equal to false , 0 , or ""
- An empty array , slice or map
- A nil pointer or a nil interface value

*Struct tag directive string*: You can use this on  individual struct fields to force the data to be represented as a string in the JSON output. ``json:"runtime,omitempty,string"`` but work only on uint*, int*, float*, bool.

When Go is encoding a particular type to JSON, it looks to see if the type has a `MarshalJSON()` method implemented on it. If it has, then Go will call this method to determine how to encode it.
```go
    type Marshaler interface {
        MarshalJSON() ([]byte, error)
    }
```
If the type doesn’t have a `MarshalJSON()` method, then Go will fall back to trying to encode it to JSON based on its own internal set of rules. So, if we want to customize how something is encoded, all we need to do is implement a `MarshalJSON()` method on it which returns a custom JSON representation of itself in a `[]byte` slice.

The rule about pointers vs. values for receivers is that value methods can be invoked on pointers and values, but pointer methods can only be invoked on pointers.

## Chapter4 Parsing JSON Requests

Using json.Decoder is generally the best choice. It’s more efficient than json.Unmarshal() , requires less code, and offers some helpful settings that you can use to tweak its behavior.


`err := json.NewDecoder(r.Body).Decode(&input)`
1. When decoding a JSON object into a struct, the key/value pairs in the JSON are mapped to the struct fields based on the struct tag names. If there is no matching struct tag, Go will attempt to decode the value into a field that matches the key name (exact matchesare preferred, but it will fall back to a case-insensitive match). Any JSON key/value pairs which cannot be successfully mapped to the struct fields will be silently ignored.
2. There is no need to close `r.Body` after it has been read. This will be done automatically by Go’s `http.Server` , so you don’t have too.

If we omit a particular key/value pair in our JSON request body. => it save thst field as a zero value, how can you tell the difference between a client not providing a key/value pair, and providing a key/value pair but deliberately setting it to its zero value?

**Error**
Two classes of error that your application might encounter:
1. Expected errors: Occur during normal operation. for example those caused by a database query timeout, a networkre source being unavailable, or bad user input. These errors don’t necessarily mean there is a problem with your program itself — in fact they’re often caused by things outside the control of your program practice to return these kinds of errors and handle them gracefully.
2. Unexpected errors: which should not happen during normal operation, and if they do it is probably the result of a developer mistake or a logical error in your codebase. These errors are truly exceptional, and using panic in these circumstances is more widely accepted. In fact, the Go standard library frequently does this when you trying to access an out-of-bounds index in a slice, or trying to close an already-closed channel.

`json.InvalidUnmarshalError` at runtime it’s because we as the developers have passed an unsupported value to Decode(). This is firmly an unexpected error which we shouldn’t see under normal operation, and is something that should be picked up in development and tests long before deployment.

Go’s json.Decoder provides a `DisallowUnknownFields()` setting for handling unwanted fields in request.

`json.Decoder` is designed to support streams of JSON data. When we call `Decode()` on our request body, it actually reads the first JSON value only from the body and decodes it. If we made a second call to `Decode()` , it would read and decode the second value and so on. To ensure that there are no additional JSON values (or any other content) in the request body, we will need to call `Decode()` a second time in our readJSON() helper and check that it returns an `io.EOF` (end of file) error.

if there is no ensure that there are no additional JSON values (or any other content) in the request body, we will need to call `Decode()` a second time in our `readJSON()` helper and check that
it returns an `io.EOF` (end of file) error.

Go is decoding some JSON, it will check to see if the destination type satisfies the json.Unmarshaler interface. If it does satisfy the interface, then Go will call it’s `UnmarshalJSON()` method to determine how to decode the provided JSON into the target type.

## Chapter5 Database Setup and Configuration

we’ll use the `sql.Open()` function to establish a new `sql.DB` connection pool, then — because connections to the database are established lazily as and when needed for the first time — we will also need to use the `db.PingContext()` method to actually create a connection and verify that everything is set up correctly.

`export GREENLIGHT_DB_DSN='postgres://reenlight:pa55word@localhost/greenlight'`
```go
flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
```
`psql $GREENLIGHT_DB_DSN`

**How does the sql.DB connection pool work?**
`sql.DB` pool contains two types of connections:
A connection is marked as **in-use** when you are using it to perform a database task, such as executing a SQL statement or querying rows, and when the task is complete the connection is then marked as **idle**.

When you instruct Go to perform a database task, it will first check if any idle connections are available in the pool. If one is available, then Go will reuse this existing connection and mark it as in-use for the duration of the task. If there are no idle connections in the pool when you need one, then Go will create a new additional connection.

When Go reuses an idle connection from the pool, any problems with the connection are handled gracefully. Bad connections will automatically be re-tried twice before giving up, at which point Go will remove the bad connection from the pool and create a new one to carry out the task.

**Configuring the pool**
The connection pool has four methods that we can use to configure its behavior:
1. `SetMaxOpenConns()` method : set an upper limit on the number of ‘open’ connections (in-use + idle connections) in the pool. By default is unlimited. The higher MaxOpenConns limit, the more database queries can be performed concurrently and the lower the risk is that the connection pool itself will be a bottleneck in your application. If the MaxOpenConns limit is reached, and all connections are in-use, then any further database tasks will be forced to wait until a connection becomes free and marked as idle so it’s important to always set a timeout on database tasks using a context.Context object. You should tweak this value for your hardware depending on the results of benchmarking and load-testing.
2. `SetMaxIdleConns()` method : upper limit on the number of idle connections in the pool. Deafault is 2. Keeping an idle connection takes up memory. *you only want to keep a connection idle if you’re likely to be using it again soon.*
3. `SetConnMaxLifetime()` method : limit the maximum length of time that a connection can be reused for. By default, there’s no maximum lifetime and connections will be reused forever.
- This doesn’t guarantee that a connection will exist in the pool for a whole hour; it’s possible that a connection will become unusable for some reason and be automatically closed before then.
- A connection can still be in use more than one hour after being created — it just cannot start to be reused after that time.
4. `SetConnMaxIdleTime()` method : limit the maximum length of time that a connection can be idle for before it is marked as expired. By default there’s no limit.

## Chapter6 SQL Migrations
For every change that you want to make to your database schema (like creating a table, adding a column, or removing an unused index) you create a pair of migration files. One file is the ‘up’ migration which contains the SQL statements necessary to implement the change, and the other is a ‘down’ migration which contains the SQL statements to reverse (or roll-back) the change.

`$ migrate create -seq -ext=.sql -dir=./migrations create_movies_table` -seq => use sequential number 0001, 0002, ...

Executing migration : `$ migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up`
`\dt` => listing tables `\d movies` => list movies table

`$ migrate -path=./migrations -database=$EXAMPLE_DSN version` => to see which migration version your database is currently on 

`$ migrate -path=./migrations -database=$EXAMPLE_DSN goto 1` => migrate up or down to a specific version

`$ migrate -path=./migrations -database =$EXAMPLE_DSN down ` => rolling bacl all migrations

Patter for models: 

inside internal/data/models.go:
```go
    type Models struct {
        Movies MovieModel
        User UserModel
    }

    func NewModels(db *sql.DB) Models {
        return Models{
            Movies: MovieModel{DB: db},
        }
    }
```
inside internal/data/movies.go:
```go
    type MovieModel struct {
        DB *sql.DB
    }

    func (m MovieModel) Insert(movie *Movie) error {
        return nil
    }
```
inside main.go:
```go
type application struct {
	...
	models data.Models
}
```
execute actions on our movies table will be very clear and readable from the perspective of our API handlers. :
`app.models.Movies.Insert(...)`

## Chapter8 Advanced CRUP Operations

Change the fields in our input struct to be pointers. Then to see if a client has provided a particular key/value pair in the JSON, we can simply check whether the corresponding field in the input struct equals nil or not.

**SQL Query Timeouts**
Go also provides context-aware variants of these Exec(), `QueryRow()` methods: `ExecContext()` and `QueryRowContext()`. These variants accept a context.Context instance as the first
parameter which you can leverage to terminate running database queries.

```go
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		...
	)
```

It’s possible that the timeout deadline will be hit before the PostgreSQL query even starts.

## Chapter8 Filtering, Sorting, and Pagination

`r.URL.Query()` returns a url.Values type, which is a map holding the query string data. Using the `Get()` method return thev alue for a specific key as a string type, or the empty string "" if no matching key exists.

The hardest part of building a dynamic filtering feature like this is the SQL query to retrieve the data — we need it to work with no filters, filters on both title and genres , or a filter on only one of them.

```SQL
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (LOWER(title) = LOWER($1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY id
```
This SQL query is designed so that each of the filters behaves like it is ‘optional’. `(LOWER(title) = LOWER($1) OR $1 = '')` will evaluate as `true` if the placeholder parameter $1 is a case-insensitive match for the movie title or the placeholder parameter equals ''.

The `(genres @> $2 OR $2 = '{}')` condition works in the same way. The @> symbol is the ‘contains’ operator for PostgreSQL arrays, and this condition will return true if all values in the placeholder parameter `$2` are contained in the database genres field or the placeholder parameter contains an empty array.
https://www.postgresql.org/docs/9.6/functions-array.html

**Partial Serching**:
```SQL
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
```
1.  `to_tsvector('simple', title)` function takes a movie title and splits it into lexemes. We specify the simple configuration, which means that the lexemes are just lowercase versions of the words in the title†. "The Breakfast" => "the", "breakfast"
2. `plainto_tsquery('simple', $1)` function takes a search value and turns it into a formatted query term that PostgreSQL full-text search can understand. It normalizes thesearch value (again using the simple configuration), strips any special characters, and inserts the and operator & between the words. "The Club" => query term 'the' & 'club' .
3. The @@ operator is the matches operator. In our statement we are using it to check whether the generated query term matches the lexemes. To continue the example, the query term 'the' & 'club' will match rows which contain both lexemes 'the' and 'club'.

**Adding indexes**
To keep our SQL query performing quickly as the dataset grows, it’s sensible to use indexes to help avoid full table scans and avoid generating the lexemes for the title field every time the query is run.

**GIN** indexes are “inverted indexes” which are appropriate for data values that contain multiple component values, such as arrays. An inverted index contains a separate entry for each component value, and can efficiently handle queries that test for the presence of specific component values.

### Paginating Lists

The LIMIT clause allows you to set the maximum number of records that a SQL query should return, and OFFSET allows you to ‘skip’ a specific number of rows before starting to return records from the query.

```go
LIMIT = page_size
OFFSET = (page - 1) * page_size
```

## Chapter9 Structured Logging and Error Handling

We want to write log entries in this format :
`{"level":"INFO","time":"2020-12-16T10:53:35Z","message":"starting server","properties":{"addr":":4000","env":"development"}}`

Our `Logger` type is a fairly thin wrapper around an io.Writer . We have some helper methods like `PrintInfo()` and `PrintError()` which accept some data for the log entry, encode this data to JSON, and then write it to the io.Writer. You can also use zerolog package as a third-party package for logging.

## Chapter10 Panic Recovery

Panics in our API handlers will be recovered automatically by Go’s http.Server => Unwind the stack for the affected goroutine (calling deferred functions along the way), close the underlying HTTP connection, and log an error message and stack trace. Create a middleware to send 500 server error if panic happen.

## Chapter11 Rate Limiting

rate limiting to prevent clients from making too many requests too quickly, and putting excessive strain on your server.

Create middleware to check how many requests have been received in the last ‘N’ seconds and — if there have been too many — then it should send the client a 429 Too Many Requests response. We’ll position this middleware before our main application handlers, so that it carries out this check before we do any expensive processing like decoding a JSON request body or querying our database.

Create middleware to rate-limit requests to your API endpoints, first by making a single rate global limiter, then extending it to support per-client limiting based on IP address.

Make rate limiter behavior configurable at runtime, including disabling the rate limiter altogether for testing purposes.

**Globa Rate Limiting**
This will consider all the requests that our API receives (rather than having separate rate limiters for every individual client).

`x/time/rate` provides a tried-and-tested implementation of a token bucket rate limiter.

*How token-bucket rate limiters work?*

A Limiter controls how frequently events are allowed to happen. It implements a “token bucket” of size b , initially full and refilled at rate r tokens per second.

1. We will have a bucket that starts with b tokens in it.
2. Each time we receive a HTTP request, we will remove one token from the bucket.
3. Every 1/r seconds, a token is added back to the bucket — up to a maximum of b total tokens.
4. If we receive a HTTP request and the bucket is empty, then we should return a 429 Too Many Requests response.

Our application would allow a maximum ‘burst’ of b HTTP requests in quick succession, but over time it would allow an average of r requests per second.

```go
// Note that the Limit type is an 'alias' for float64.
func NewLimiter(r Limit, b int) *Limiter
```

```go
func (app *application) exampleMiddleware(next http.Handler) http.Handler {
// Any code here will run only once, when we wrap something with the middleware.
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Any code here will run for every request that the middleware handles.
    next.ServeHTTP(w, r)
    })
}
```

Make a rateLimit() middleware method which creates a new rate limiter as part of the ‘initialization’ code, and then uses this rate limiter for every request that it subsequently handlers.

```go
func (app *application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```
`Allow()` method on the rate limiter exactly one token will be consumed from the bucket. If there are no tokens left in the bucket, then `Allow()` will return false.

Code behind the Allow() method is protected by a mutex and is safe for concurrent use.