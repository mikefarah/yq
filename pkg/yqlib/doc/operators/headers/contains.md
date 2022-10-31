# Contains

This returns `true` if the context contains the passed in parameter, and false otherwise. For arrays, this will return true if the passed in array is contained within the array. For strings, it will return true if the string is a substring.

{% hint style="warning" %}

_Note_ that, just like jq, when checking if an array of strings `contains` another, this will use `contains` and _not_ equals to check each string. This means an array like `contains(["cat"])` will return true for an array `["cats"]`.

See the "Array has a subset array" example below on how to check for a subset.

{% endhint %}
