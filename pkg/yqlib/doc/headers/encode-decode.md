# Encoder / Decoder

Encode operators will take the piped in object structure and encode it as a string in the desired format. The decode operators do the opposite, they take a formatted string and decode it into the relevant object structure.

Note that you can optionally pass an indent value to the encode functions (see below).

These operators are useful to process yaml documents that have stringified embeded yaml/json/props in them.
