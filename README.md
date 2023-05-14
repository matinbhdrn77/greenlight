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

*Struct tag directive string*: You can use this on  individual struct fields to force the data to be represented as a string in the JSON output. ``json:"runtime,omitempty,string"`` but work only on uint*, int*, float*, bool