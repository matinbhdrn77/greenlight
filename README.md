## Chapter2

Any packages which live under *internal* directory can only be imported by code inside the parent of the internal directory.

`func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request)` => and idiomatic way to make dependencies available to our handlers without resorting to global variables or closures — any dependency that the healthcheckHandler needs can simply be included as a field in the application struct when we initialize it in main() .