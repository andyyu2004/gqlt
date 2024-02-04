module.exports = grammar({
    name: "gqlt",

    extras: ($) => [/[\s\uFEFF\u0009\u0020\u000A\u000D]/, $.comment],

    rules: {
        // graphql executable grammar below
        source_file: ($) => $.document,
        document: ($) => repeat1($.definition),
        definition: ($) => choice($.executable_definition),
        executable_definition: ($) =>
            choice($.operation_definition, $.fragment_definition),
        default_value: ($) => seq("=", $.value),
        operation_definition: ($) =>
            choice(
                $.selection_set,
                seq(
                    $.operation_type,
                    optional($.name),
                    optional($.variable_definitions),
                    optional($.directives),
                    $.selection_set,
                ),
            ),
        operation_type: (_) => choice("query", "mutation", "subscription"),
        variable_definitions: ($) => seq("(", repeat1($.variable_definition), ")"),
        variable_definition: ($) =>
            seq(
                $.variable,
                ":",
                $.type,
                optional($.default_value),
                optional($.directives),
                optional($.comma),
            ),
        selection_set: ($) => seq("{", repeat1($.selection), "}"),
        selection: ($) => choice($.field, $.inline_fragment, $.fragment_spread),
        field: ($) =>
            seq(
                optional($.alias),
                $.name,
                optional($.arguments),
                optional($.directive),
                optional($.selection_set),
            ),
        alias: ($) => seq($.name, ":"),
        arguments: ($) => seq("(", repeat1($.argument), ")"),
        argument: ($) => seq($.name, ":", $.value),
        value: ($) =>
            choice(
                $.variable,
                $.string_value,
                $.int_value,
                $.float_value,
                $.boolean_value,
                $.null_value,
                $.enum_value,
                $.list_value,
                $.object_value,
            ),
        variable: ($) => seq("$", $.name),
        string_value: (_) =>
            choice(
                seq('"""', /([^"]|\n|""?[^"])*/, '"""'),
                seq('"', /[^"\\\n]*/, '"'),
            ),
        int_value: (_) => /-?(0|[1-9][0-9]*)/,
        float_value: (_) =>
            token(
                seq(
                    /-?(0|[1-9][0-9]*)/,
                    choice(
                        /\.[0-9]+/,
                        /(e|E)(\+|-)?[0-9]+/,
                        seq(/\.[0-9]+/, /(e|E)(\+|-)?[0-9]+/),
                    ),
                ),
            ),
        boolean_value: (_) => choice("true", "false"),
        null_value: (_) => "null",
        enum_value: ($) => $.name,
        list_value: ($) => seq("[", repeat($.value), "]"),
        object_value: ($) => seq("{", repeat($.object_field), "}"),
        object_field: ($) => seq($.name, ":", $.value, optional($.comma)),
        fragment_spread: ($) => seq("...", $.fragment_name, optional($.directives)),
        fragment_definition: ($) =>
            seq(
                "fragment",
                $.fragment_name,
                $.type_condition,
                optional($.directives),
                $.selection_set,
            ),
        fragment_name: ($) => $.name,
        inline_fragment: ($) =>
            seq(
                "...",
                optional($.type_condition),
                optional($.directives),
                $.selection_set,
            ),
        type_condition: ($) => seq("on", $.named_type),
        directives: ($) => repeat1($.directive),
        directive: ($) => seq("@", $.name, optional($.arguments)),
        type: ($) => choice($.named_type, $.list_type, $.non_null_type),
        named_type: ($) => $.name,
        list_type: ($) => seq("[", $.type, "]"),
        non_null_type: ($) => choice(seq($.named_type, "!"), seq($.list_type, "!")),
        name: (_) => /[_A-Za-z][_0-9A-Za-z]*/,
        comment: (_) => token(seq("#", /.*/)),
        comma: (_) => ",",
        description: ($) => $.string_value,
    },
});
