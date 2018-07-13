# Colorized and Formatted JSON

A Go program that takes any valid [JSON](http://json.org/ "json.org") as input (from standard input) and outputs (to standard output) HTML that transforms the JSON so that all tokens are consistently colored and the formatting is clean and consistent.

### Example

Input JSON fragment:
```
{"s":[2, 3], "a < b && a >= c":true}
```

Output view:
```json
{
  "s" : [2, 3],
  "a < b && a >= c" : true
}
```
