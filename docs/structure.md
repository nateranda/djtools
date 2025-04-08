# Structure

## `main` package
This package is the entrypoint for the CLI.

### `main.go`
This file handles CLI logic and envokes other packages.

## `db` package
This package handles all database and file storage operations. Each file in this package handles operations with one program, either `djtools` itself or other programs.

To eliminate naming conflicts while keeping all database operations in the same package, every function is prefixed with a shortened name of the program it interfaces with: `dt` for djtools, `en` for Engine, `rb` for Rekordbox, `tr` for Traktor, `dj` for Algoriddim Djay, and `mi` for Mixxx. Program-agnostic functions/methods in `db.go` are not prefixed.

In each file, library data is imported in four steps:
1. **Extraction:** all necessary library data is extracted from databases and analysis files using basic SQL queries and stored in their program-native format and datatype. The data is stored in structs suited to each database or file format.
2. **Validation:** the library data, now in various structs, is validated to make sure it is in a recognized format and not corrupted.
3. **Conversion:** the library data, now validated, is processed and converted into datatypes accepted by the `Library` struct.
4. **Injection:** the library data, now in an accepted datatype, is loaded into the `Library` struct.

In each file, library data is exported in four steps:
1. **Extraction:** all necessary library data is extracted from the `Library` struct into various structs suited to the chosen database and/or file format.
2. **Validation:** the library data, now in various structs, is validated to make sure it is in a recognized format and not corrupted.
3. **Conversion:** the library data, now validated, is processed and converted into datatypes assepted by the chosen database and/or file format.
4. **Injection:** the library data, now in an accepted datatype, is loaded into the chosen database and/or file format using basic SQL queries.

### `db.go`
This file contains shared logic for the `db` package and serves as the main endpoint for database operations. It also contains the `Library` struct used to store library data.

### `djtools.go`
This file handles all database operations involving the djtools database.

### `engine.go`
This file handles all database operations involving Engine databases.