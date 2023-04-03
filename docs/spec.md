# TIPO Specifications

## Definitions

### Data structures

- Row is an in-memory representation of a JSONLine, key are ordered
- Node is a generic term to design an element in a row, it can be an Object, an Array or a Value
- Object is like a Row, but not necessary the root object of a JSONLine (it can be deeply nested)
- Array is an indexed list of Node
- Value is a single primitive value (string, int, float)

### Processing structures

- Frame is an indexed buffer that can contains Rows, works like a FIFO list, read Rows from an io.Reader
- Tuple is an ordered set of values of any types (can be a primitive type like int or string, or can be an array or an object)
- Cache is a buffer of tuples of any type (all tuples share the same type in a single cache), tuples are pulled randomly from the cache, tuples are added to the cache with a Selector
- CacheList is a set of Caches
- Selector is an object that defines how to extract a Tuple from a Row, and how to update a Row via a Provider
- Provider is an object that can provide a stream of Tuples
- RowCollector is an object that collects Rows and sends them to an io.Writer
- RowReader is an object that read Rows from an io.Reader

## General Algorithm

```text
Parameters :
- N is the frame size, must be >0

Variables :
- collector is a RowCollector plugged to stdout
- reader is a RowReader plugged to stdin
- frame is a Frame of size N
- caches is a CacheList configured by the user

Main algorithm :
- build caches from the configuration provided by user
- while frame can FillUp from reader
  - nextRow is pulled from the frame (PullFirst)
  - for each cache in caches
    - get the selector assigned to the cache (GetSelector)
    - use the selector to Update the nextRow with a Provider that uses cache.PullRandom as source
  - use the collector to Collect the nextRow
```

## Methods

### **Frame** structure

- _**FillUp**(reader RowReader, callback func(Row)) bool_ : fills the frame with Rows read from reader, calls the callback function for each row added to the frame, and finally returns false if there is no more values in the Frame and it can't read from the reader.
- _**PullFirst**() Row_ : get the first row and removes it from the frame.

### **Cache** structure

- _**Post**(row Row)_ : post the row to feed the cache with new tuples.
- _**PullRandom**() Tuple_ : removes a random tuple from the cache and returns it. This method can serve as a tuple Provider.
- _**GetSelector**() Selector_ : returns the selector used by the cache.

### **CacheList** structure

- _**Add**(cache Cache)_ : add a cache to the cache list.
- _**All**() []Cache_ : returns all caches in the cache list.
- _**Post**(row Row)_ : post the row to feed all caches with new tuples.

### **Selector** structure

- _**Select**(row Row) []Tuple_ : select all tuples from the given row that are matched by the selector.
- _**Update**(row Row, provider Provider)_ : update all tuples matched by the selector with tuples provided by the provider.

### **Provider** interface

- _**func**() Tuple_ : provides a tuple.

### **RowReader** interface

- _**HasRow**() bool_ : returns true if a Row is ready.
- _**ReadRow**() (Row, error)_ : returns the next Row.

### **RowCollector** interface

- _**Collect**(row Row)_ : collects a row.
