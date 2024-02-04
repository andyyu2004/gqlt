(named_type
  (name) @type)

(directive) @type

; Properties
;-----------

(field
  (name) @property)

(field
  (alias
    (name) @property))

(object_value
  (object_field
    (name) @property))

(enum_value
  (name) @property)

; Constants
;----------

(string_value) @string

(int_value) @number

(float_value) @float

(boolean_value) @boolean

(comment) @comment

; Keywords
;----------

[
  "query"
  "mutation"
  "subscription"
  "fragment"
  "on"
] @keyword

; Punctuation
;------------

[
 "("
 ")"
 "["
 "]"
 "{"
 "}"
] @punctuation.bracket

"=" @operator

":" @punctuation.delimiter
"..." @punctuation.special
"!" @punctuation.special
