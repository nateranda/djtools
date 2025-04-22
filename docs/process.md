# Process
Below are the steps I follow when adding support for a new DJ software.

### Step 1: Basic features
Develop a basic feature set that serves as a minimum viable product. The program should convert tracks, including metadata, beatgrids, cues, and loops, as well as playlists, including folders and nested playlists. At this point, the program should support converting any *reasonable and accurate* library, but it doesn't have to handle deliberately corrupted library data. The program should be non-destructive, meaning it should not overwrite any existing library files unless explicitly specified.

### Step 2: Testing
Next, write integration tests for all exported functions. These tests should cover all options and test one complex beahvior at a time, like changed grids/cues/loops, playlists, nested playlists, and smartlists. This is done to make it easier to identify problems with individual parts of the program when tests do go wrong. Unit tests are optional but recommended for internal functions with edge cases that cannot be adequately tested through integration tests.

Test data should be supplied either in the test itself or through a `testdata` directory in each package folder, which may include database fixtures and library stubs. After tests are developed, every new commit should pass all tests. If a new (major) database schema is released for the DJ program, the new schema should not replace the old schema in testing, rather both versions should be tested.

### Step 3: Extended features
Next, develop any extended features as needed. Most of the time, these features should be put behind non-default options so they do not affect the existing test framework. For each feature, repeat the above steps: develop a minimum viable product and write tests to consider any edge cases.

For every feature, create a new branch and merge the branch when the feature is complete. These extended features do not have to pass tests until they are merged to the main branch.