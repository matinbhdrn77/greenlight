## Chapter2

Any packages which live under *internal* directory can only be imported by code inside the parent of the internal directory.

`func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request)` => and idiomatic way to make dependencies available to our handlers without resorting to global variables or closures — any dependency that the healthcheckHandler needs can simply be included as a field in the application struct when we initialize it in main() .

The *httprouter* package provides a few configuration options that you can use to customize the behavior of your application further, including enabling *trailing slash redirects* and *enabling automatic URL path cleaning*.

## Chapter3 Sending JSON Response

JSON is just text. key must be string. value can be any JS object.
You would write any other text response: using `w.Write()` , `ioWriteString()` or one of the `fmt.Fprint` functions. But you have to set `Content-Type: application/json`