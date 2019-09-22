---
menuTitle: FileStore
title: "FileStore - Local Storage for Files"
---


```go
type Path string

type File inteface {
  Path() Path
  Size() int
  Seek(offset int) error
  Read([]byte) (n int, err error)
  Write([]byte) (n int, err error)
  Close() error
}

type FileStore interface {
  Open(Path) (File, error)
  Create(Path) (File, error)
}
```

TODO:

- explain why this abstraction is needed
- explain OS filesystem as basic impl
- explain that users can replace w/ other systems
- give examples:
  - networked filesystems
  - raw disk sectors - like haystack
  - databases
