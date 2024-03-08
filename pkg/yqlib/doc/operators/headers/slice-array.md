# Slice/Splice Array

The slice array operator takes an array as input and returns a subarray. Like the `jq` equivalent, `.[10:15]` will return an array of length 5, starting from index 10 inclusive, up to index 15 exclusive. Negative numbers count backwards from the end of the array.

You may leave out the first or second number, which will refer to the start or end of the array respectively.
