
# TODOs for the .id spec/ipld dsl

- Bugs
  - [x] type def not generated if one of the dependency types is not found in same file
    - sectorset test.
    - update: this seems to be fixed?
  - [ ] cannot have `struct{ cid CID }` because there's an implicit `cid` field.
  - [ ] cannot have `interface{}` parameters or returns
- Parsing rules
  - [ ] looks like a function invocation can be split across lines. doing so removes the commas (at least in the fmt output)
    - we should keep commas, not doing so is error prone IMO.
      - make it OK to have (or not have) a trailing comma in param definition
      - Go maps do something like this. if the definition is inline, no last comma allowed. if the definition is multiline, last comma is required
    - (if possible) spacing should not be used (or used the least possible) to define rules. in this case this means keeping the comma.
- Formatting code
  - [x] fmt file
  - [x] fmt ./... (recursive in a directory)
  - [ ] fmt rules use tabs for scope indentation, spaces only to space out things evenl
  - [ ] fmt rules align fields across multiple lines differently
- Documentation
  - [x] codeGen Usage, -h output
  - [ ] codeGen tool Readme
- Tests
  - [ ] run test suite before committing codeGen changes
