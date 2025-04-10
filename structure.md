# Structure

### Tree
```
djtools/
├─ lib/
│  ├─ lib.go -> shared functions (if needed)
├─ db/
│  ├─ db.go -> library struct, shared database functions
│  ├─ validate.go -> library validation
├─ [program]/
│  ├─ [program].go -> import/export framework, raw data structs, shared functions
│  ├─ import.go -> import framework
│  ├─ import[step].go -> step-specific import stuff (if needed)
│  ├─ export.go -> export framework
│  ├─ export[step].go -> step-specific export stuff (if needed)
├─ cli/
│  ├─ cli.go -> cli logic
├─ main.go -> cli endpoint, dev testing
```