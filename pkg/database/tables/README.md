# database tables

This package provides generated functions for working with tables in our
database. The best way to view the contents of this package (in my opinion) is
with `godoc`.

## Adding a new table

### Defining table and columns

Define the table and columns in `gen/tables.go`. It should be straightforward using
existing tables as examples.

Add columns to a table with the `Col` function. Columns can be modified by
using the following functions:

```
goName(name string)    // Override the generated Go field name.
jsonName(name string)  // Override the generated JSON struct tag.
validator(name string) // Validation function from our `validators` package.
pk()                   // Mark as a primary key column.
onConflictUpdate()     // Update this column on insert primary key conflict.
noUpdate()             // Don't update this column in the Update function.
noDirectUpdate()       // Don't generate a direct update function for column.
```

Add function get get rows by keys other than the primary using the table's
`GetBy` function.

**Note**: If you're adding a column with a Go type that isn't understood by the
generator, it may lead to issues. Currently primitive types, `time.Time`,
`json.RawMessage` or pointers to those should work.

### Add constructor for testing

In the `rowtypes_test.go` file, add a constructor `New{{RowType}}ForTesting`
function. This should return an insertable object. If necessary, parents should
be created and inserted.

### Generate code

Generate types and functions for the table by running `go generate` in this
directory.

### Test

If everything went well, you can test your generated functions with `go
test`.
