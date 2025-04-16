# Process
This is the basic process to follow for implementing application conversion support. Follow this same structure for both import and export functionality.

### Step 1: Basic features
Develop a basic feature set that serves as a minimum viable product. The program should convert tracks, including metadata, beatgrids, cues, and loops, as well as playlists, including folders and nested playlists. At this point, the program should support converting any *reasonable and accurate* library, but it doesn't have to handle deliberately wrong or corrupted data. The program should be non-destructive, meaning it should not overwrite any existing library files unless explicitly specified.

### Step 2: Edge case analysis
Next, analyze all edge cases and assumptions for evey function used. Write them down in the program's respective `edge-cases` entry and note if the assumption is handled effectively by the program. This should include corrupted or missing data, but also wrong datatypes and incorrect string formatting.

### Step 3: Testing
Next, write unit tests for every function. They should consider the edge cases found in the last step. All large-scale tests (`Import()` or `Export()` and their derivative functions `Import[step]()` and `Export[step]()`) should be conducted with default options, but each option should be tested where the option applies. This is done to cut down existing test rewrites when adding new features.

Update the feature until it passes all unit tests. Test data should be stored either in the test itself or in a `test` directory in each package folder which may include database mocks and library gobs. After tests are developed, every new commit should pass all tests.

If a new database schema is released for the DJ program, the new schema should not replace the old schema in testing, rather both versions should be tested.

### Step 4: Extended features
Next, develop any extended features such as: mp3 offset readjustment, selective data extraction, and data reorganization. Most of the time, these features should be put behind non-default options so they do not affect the existing test framework. However, these features should still be tested in the functions they are applied in. For each feature, repeat the above steps: develop a minimum viable product, analyze edge cases, and write tests to consider these edge cases.

For every feature, create a new branch and merge the branch when the feature is complete. These extended features do not have to pass tests until they are merged to the main branch.