# Structure

### Tree
```
djtools/
├─ lib/
│  ├─ lib.go -> shared functions (if needed)
├─ db/
│  ├─ library.go -> library struct, shared library functions
│  ├─ lib.go -> shared general functions
│  ├─ validate.go -> library validation
├─ [program]/
│  ├─ [program].go -> import/export framework, raw data structs, shared functions
│  ├─ import[step].go -> import functions (step-specific if necessary)
│  ├─ export[step].go -> export functions (step-specific if necessary)
├─ cli/
│  ├─ cli.go -> cli logic
├─ main.go -> cli endpoint, dev testing
```