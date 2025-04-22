# Structure

```
djtools/
├─ lib/
│  ├─ library.go -> library struct, shared library functions
│  ├─ lib.go -> shared general functions
│  ├─ test.go -> shared test functions
│  ├─ validate.go -> library validation (TBI)
├─ [program]/
│  ├─ [program].go -> import/export framework, raw data structs, shared functions
│  ├─ import[step].go -> step-specific import functions (if necessary)
│  ├─ export[step].go -> step-specific export functions (if necessary)
│  ├─ [program]_test.go -> test package
├─ cli/ (TBI)
│  ├─ cli.go -> cli logic
├─ main.go -> cli endpoint
```