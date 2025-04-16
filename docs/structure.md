# Structure

```
djtools/
├─ db/
│  ├─ library.go -> library struct, shared library functions
│  ├─ lib.go -> shared general functions
│  ├─ validate.go -> library validation
├─ [program]/
│  ├─ [program].go -> import/export framework, raw data structs, shared functions
│  ├─ import[step].go -> step-specific import functions (if necessary)
│  ├─ export[step].go -> step-specific export functions (if necessary)
├─ cli/
│  ├─ cli.go -> cli logic
├─ main.go -> cli endpoint, dev testing
```