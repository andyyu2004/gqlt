#include "tree_sitter/parser.h"

#if defined(__GNUC__) || defined(__clang__)
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wmissing-field-initializers"
#endif

#define LANGUAGE_VERSION 14
#define STATE_COUNT 137
#define LARGE_STATE_COUNT 2
#define SYMBOL_COUNT 71
#define ALIAS_COUNT 0
#define TOKEN_COUNT 30
#define EXTERNAL_TOKEN_COUNT 0
#define FIELD_COUNT 0
#define MAX_ALIAS_SEQUENCE_LENGTH 6
#define PRODUCTION_ID_COUNT 1

enum ts_symbol_identifiers {
  anon_sym_EQ = 1,
  anon_sym_query = 2,
  anon_sym_mutation = 3,
  anon_sym_subscription = 4,
  anon_sym_LPAREN = 5,
  anon_sym_RPAREN = 6,
  anon_sym_COLON = 7,
  anon_sym_LBRACE = 8,
  anon_sym_RBRACE = 9,
  anon_sym_DOLLAR = 10,
  anon_sym_DQUOTE_DQUOTE_DQUOTE = 11,
  aux_sym_string_value_token1 = 12,
  anon_sym_DQUOTE = 13,
  aux_sym_string_value_token2 = 14,
  sym_int_value = 15,
  sym_float_value = 16,
  anon_sym_true = 17,
  anon_sym_false = 18,
  sym_null_value = 19,
  anon_sym_LBRACK = 20,
  anon_sym_RBRACK = 21,
  anon_sym_DOT_DOT_DOT = 22,
  anon_sym_fragment = 23,
  anon_sym_on = 24,
  anon_sym_AT = 25,
  anon_sym_BANG = 26,
  sym_name = 27,
  sym_comment = 28,
  sym_comma = 29,
  sym_source_file = 30,
  sym_document = 31,
  sym_definition = 32,
  sym_executable_definition = 33,
  sym_default_value = 34,
  sym_operation_definition = 35,
  sym_operation_type = 36,
  sym_variable_definitions = 37,
  sym_variable_definition = 38,
  sym_selection_set = 39,
  sym_selection = 40,
  sym_field = 41,
  sym_alias = 42,
  sym_arguments = 43,
  sym_argument = 44,
  sym_value = 45,
  sym_variable = 46,
  sym_string_value = 47,
  sym_boolean_value = 48,
  sym_enum_value = 49,
  sym_list_value = 50,
  sym_object_value = 51,
  sym_object_field = 52,
  sym_fragment_spread = 53,
  sym_fragment_definition = 54,
  sym_fragment_name = 55,
  sym_inline_fragment = 56,
  sym_type_condition = 57,
  sym_directives = 58,
  sym_directive = 59,
  sym_type = 60,
  sym_named_type = 61,
  sym_list_type = 62,
  sym_non_null_type = 63,
  aux_sym_document_repeat1 = 64,
  aux_sym_variable_definitions_repeat1 = 65,
  aux_sym_selection_set_repeat1 = 66,
  aux_sym_arguments_repeat1 = 67,
  aux_sym_list_value_repeat1 = 68,
  aux_sym_object_value_repeat1 = 69,
  aux_sym_directives_repeat1 = 70,
};

static const char * const ts_symbol_names[] = {
  [ts_builtin_sym_end] = "end",
  [anon_sym_EQ] = "=",
  [anon_sym_query] = "query",
  [anon_sym_mutation] = "mutation",
  [anon_sym_subscription] = "subscription",
  [anon_sym_LPAREN] = "(",
  [anon_sym_RPAREN] = ")",
  [anon_sym_COLON] = ":",
  [anon_sym_LBRACE] = "{",
  [anon_sym_RBRACE] = "}",
  [anon_sym_DOLLAR] = "$",
  [anon_sym_DQUOTE_DQUOTE_DQUOTE] = "\"\"\"",
  [aux_sym_string_value_token1] = "string_value_token1",
  [anon_sym_DQUOTE] = "\"",
  [aux_sym_string_value_token2] = "string_value_token2",
  [sym_int_value] = "int_value",
  [sym_float_value] = "float_value",
  [anon_sym_true] = "true",
  [anon_sym_false] = "false",
  [sym_null_value] = "null_value",
  [anon_sym_LBRACK] = "[",
  [anon_sym_RBRACK] = "]",
  [anon_sym_DOT_DOT_DOT] = "...",
  [anon_sym_fragment] = "fragment",
  [anon_sym_on] = "on",
  [anon_sym_AT] = "@",
  [anon_sym_BANG] = "!",
  [sym_name] = "name",
  [sym_comment] = "comment",
  [sym_comma] = "comma",
  [sym_source_file] = "source_file",
  [sym_document] = "document",
  [sym_definition] = "definition",
  [sym_executable_definition] = "executable_definition",
  [sym_default_value] = "default_value",
  [sym_operation_definition] = "operation_definition",
  [sym_operation_type] = "operation_type",
  [sym_variable_definitions] = "variable_definitions",
  [sym_variable_definition] = "variable_definition",
  [sym_selection_set] = "selection_set",
  [sym_selection] = "selection",
  [sym_field] = "field",
  [sym_alias] = "alias",
  [sym_arguments] = "arguments",
  [sym_argument] = "argument",
  [sym_value] = "value",
  [sym_variable] = "variable",
  [sym_string_value] = "string_value",
  [sym_boolean_value] = "boolean_value",
  [sym_enum_value] = "enum_value",
  [sym_list_value] = "list_value",
  [sym_object_value] = "object_value",
  [sym_object_field] = "object_field",
  [sym_fragment_spread] = "fragment_spread",
  [sym_fragment_definition] = "fragment_definition",
  [sym_fragment_name] = "fragment_name",
  [sym_inline_fragment] = "inline_fragment",
  [sym_type_condition] = "type_condition",
  [sym_directives] = "directives",
  [sym_directive] = "directive",
  [sym_type] = "type",
  [sym_named_type] = "named_type",
  [sym_list_type] = "list_type",
  [sym_non_null_type] = "non_null_type",
  [aux_sym_document_repeat1] = "document_repeat1",
  [aux_sym_variable_definitions_repeat1] = "variable_definitions_repeat1",
  [aux_sym_selection_set_repeat1] = "selection_set_repeat1",
  [aux_sym_arguments_repeat1] = "arguments_repeat1",
  [aux_sym_list_value_repeat1] = "list_value_repeat1",
  [aux_sym_object_value_repeat1] = "object_value_repeat1",
  [aux_sym_directives_repeat1] = "directives_repeat1",
};

static const TSSymbol ts_symbol_map[] = {
  [ts_builtin_sym_end] = ts_builtin_sym_end,
  [anon_sym_EQ] = anon_sym_EQ,
  [anon_sym_query] = anon_sym_query,
  [anon_sym_mutation] = anon_sym_mutation,
  [anon_sym_subscription] = anon_sym_subscription,
  [anon_sym_LPAREN] = anon_sym_LPAREN,
  [anon_sym_RPAREN] = anon_sym_RPAREN,
  [anon_sym_COLON] = anon_sym_COLON,
  [anon_sym_LBRACE] = anon_sym_LBRACE,
  [anon_sym_RBRACE] = anon_sym_RBRACE,
  [anon_sym_DOLLAR] = anon_sym_DOLLAR,
  [anon_sym_DQUOTE_DQUOTE_DQUOTE] = anon_sym_DQUOTE_DQUOTE_DQUOTE,
  [aux_sym_string_value_token1] = aux_sym_string_value_token1,
  [anon_sym_DQUOTE] = anon_sym_DQUOTE,
  [aux_sym_string_value_token2] = aux_sym_string_value_token2,
  [sym_int_value] = sym_int_value,
  [sym_float_value] = sym_float_value,
  [anon_sym_true] = anon_sym_true,
  [anon_sym_false] = anon_sym_false,
  [sym_null_value] = sym_null_value,
  [anon_sym_LBRACK] = anon_sym_LBRACK,
  [anon_sym_RBRACK] = anon_sym_RBRACK,
  [anon_sym_DOT_DOT_DOT] = anon_sym_DOT_DOT_DOT,
  [anon_sym_fragment] = anon_sym_fragment,
  [anon_sym_on] = anon_sym_on,
  [anon_sym_AT] = anon_sym_AT,
  [anon_sym_BANG] = anon_sym_BANG,
  [sym_name] = sym_name,
  [sym_comment] = sym_comment,
  [sym_comma] = sym_comma,
  [sym_source_file] = sym_source_file,
  [sym_document] = sym_document,
  [sym_definition] = sym_definition,
  [sym_executable_definition] = sym_executable_definition,
  [sym_default_value] = sym_default_value,
  [sym_operation_definition] = sym_operation_definition,
  [sym_operation_type] = sym_operation_type,
  [sym_variable_definitions] = sym_variable_definitions,
  [sym_variable_definition] = sym_variable_definition,
  [sym_selection_set] = sym_selection_set,
  [sym_selection] = sym_selection,
  [sym_field] = sym_field,
  [sym_alias] = sym_alias,
  [sym_arguments] = sym_arguments,
  [sym_argument] = sym_argument,
  [sym_value] = sym_value,
  [sym_variable] = sym_variable,
  [sym_string_value] = sym_string_value,
  [sym_boolean_value] = sym_boolean_value,
  [sym_enum_value] = sym_enum_value,
  [sym_list_value] = sym_list_value,
  [sym_object_value] = sym_object_value,
  [sym_object_field] = sym_object_field,
  [sym_fragment_spread] = sym_fragment_spread,
  [sym_fragment_definition] = sym_fragment_definition,
  [sym_fragment_name] = sym_fragment_name,
  [sym_inline_fragment] = sym_inline_fragment,
  [sym_type_condition] = sym_type_condition,
  [sym_directives] = sym_directives,
  [sym_directive] = sym_directive,
  [sym_type] = sym_type,
  [sym_named_type] = sym_named_type,
  [sym_list_type] = sym_list_type,
  [sym_non_null_type] = sym_non_null_type,
  [aux_sym_document_repeat1] = aux_sym_document_repeat1,
  [aux_sym_variable_definitions_repeat1] = aux_sym_variable_definitions_repeat1,
  [aux_sym_selection_set_repeat1] = aux_sym_selection_set_repeat1,
  [aux_sym_arguments_repeat1] = aux_sym_arguments_repeat1,
  [aux_sym_list_value_repeat1] = aux_sym_list_value_repeat1,
  [aux_sym_object_value_repeat1] = aux_sym_object_value_repeat1,
  [aux_sym_directives_repeat1] = aux_sym_directives_repeat1,
};

static const TSSymbolMetadata ts_symbol_metadata[] = {
  [ts_builtin_sym_end] = {
    .visible = false,
    .named = true,
  },
  [anon_sym_EQ] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_query] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_mutation] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_subscription] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_LPAREN] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_RPAREN] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_COLON] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_LBRACE] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_RBRACE] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_DOLLAR] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_DQUOTE_DQUOTE_DQUOTE] = {
    .visible = true,
    .named = false,
  },
  [aux_sym_string_value_token1] = {
    .visible = false,
    .named = false,
  },
  [anon_sym_DQUOTE] = {
    .visible = true,
    .named = false,
  },
  [aux_sym_string_value_token2] = {
    .visible = false,
    .named = false,
  },
  [sym_int_value] = {
    .visible = true,
    .named = true,
  },
  [sym_float_value] = {
    .visible = true,
    .named = true,
  },
  [anon_sym_true] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_false] = {
    .visible = true,
    .named = false,
  },
  [sym_null_value] = {
    .visible = true,
    .named = true,
  },
  [anon_sym_LBRACK] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_RBRACK] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_DOT_DOT_DOT] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_fragment] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_on] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_AT] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_BANG] = {
    .visible = true,
    .named = false,
  },
  [sym_name] = {
    .visible = true,
    .named = true,
  },
  [sym_comment] = {
    .visible = true,
    .named = true,
  },
  [sym_comma] = {
    .visible = true,
    .named = true,
  },
  [sym_source_file] = {
    .visible = true,
    .named = true,
  },
  [sym_document] = {
    .visible = true,
    .named = true,
  },
  [sym_definition] = {
    .visible = true,
    .named = true,
  },
  [sym_executable_definition] = {
    .visible = true,
    .named = true,
  },
  [sym_default_value] = {
    .visible = true,
    .named = true,
  },
  [sym_operation_definition] = {
    .visible = true,
    .named = true,
  },
  [sym_operation_type] = {
    .visible = true,
    .named = true,
  },
  [sym_variable_definitions] = {
    .visible = true,
    .named = true,
  },
  [sym_variable_definition] = {
    .visible = true,
    .named = true,
  },
  [sym_selection_set] = {
    .visible = true,
    .named = true,
  },
  [sym_selection] = {
    .visible = true,
    .named = true,
  },
  [sym_field] = {
    .visible = true,
    .named = true,
  },
  [sym_alias] = {
    .visible = true,
    .named = true,
  },
  [sym_arguments] = {
    .visible = true,
    .named = true,
  },
  [sym_argument] = {
    .visible = true,
    .named = true,
  },
  [sym_value] = {
    .visible = true,
    .named = true,
  },
  [sym_variable] = {
    .visible = true,
    .named = true,
  },
  [sym_string_value] = {
    .visible = true,
    .named = true,
  },
  [sym_boolean_value] = {
    .visible = true,
    .named = true,
  },
  [sym_enum_value] = {
    .visible = true,
    .named = true,
  },
  [sym_list_value] = {
    .visible = true,
    .named = true,
  },
  [sym_object_value] = {
    .visible = true,
    .named = true,
  },
  [sym_object_field] = {
    .visible = true,
    .named = true,
  },
  [sym_fragment_spread] = {
    .visible = true,
    .named = true,
  },
  [sym_fragment_definition] = {
    .visible = true,
    .named = true,
  },
  [sym_fragment_name] = {
    .visible = true,
    .named = true,
  },
  [sym_inline_fragment] = {
    .visible = true,
    .named = true,
  },
  [sym_type_condition] = {
    .visible = true,
    .named = true,
  },
  [sym_directives] = {
    .visible = true,
    .named = true,
  },
  [sym_directive] = {
    .visible = true,
    .named = true,
  },
  [sym_type] = {
    .visible = true,
    .named = true,
  },
  [sym_named_type] = {
    .visible = true,
    .named = true,
  },
  [sym_list_type] = {
    .visible = true,
    .named = true,
  },
  [sym_non_null_type] = {
    .visible = true,
    .named = true,
  },
  [aux_sym_document_repeat1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_variable_definitions_repeat1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_selection_set_repeat1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_arguments_repeat1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_list_value_repeat1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_object_value_repeat1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_directives_repeat1] = {
    .visible = false,
    .named = false,
  },
};

static const TSSymbol ts_alias_sequences[PRODUCTION_ID_COUNT][MAX_ALIAS_SEQUENCE_LENGTH] = {
  [0] = {0},
};

static const uint16_t ts_non_terminal_alias_map[] = {
  0,
};

static const TSStateId ts_primary_state_ids[STATE_COUNT] = {
  [0] = 0,
  [1] = 1,
  [2] = 2,
  [3] = 3,
  [4] = 4,
  [5] = 2,
  [6] = 4,
  [7] = 7,
  [8] = 8,
  [9] = 9,
  [10] = 10,
  [11] = 11,
  [12] = 12,
  [13] = 13,
  [14] = 14,
  [15] = 15,
  [16] = 16,
  [17] = 17,
  [18] = 18,
  [19] = 19,
  [20] = 20,
  [21] = 21,
  [22] = 22,
  [23] = 23,
  [24] = 24,
  [25] = 25,
  [26] = 26,
  [27] = 27,
  [28] = 27,
  [29] = 29,
  [30] = 30,
  [31] = 31,
  [32] = 32,
  [33] = 33,
  [34] = 34,
  [35] = 33,
  [36] = 36,
  [37] = 37,
  [38] = 38,
  [39] = 39,
  [40] = 40,
  [41] = 15,
  [42] = 42,
  [43] = 43,
  [44] = 44,
  [45] = 45,
  [46] = 46,
  [47] = 47,
  [48] = 48,
  [49] = 49,
  [50] = 13,
  [51] = 14,
  [52] = 19,
  [53] = 16,
  [54] = 54,
  [55] = 20,
  [56] = 12,
  [57] = 57,
  [58] = 58,
  [59] = 59,
  [60] = 60,
  [61] = 61,
  [62] = 18,
  [63] = 63,
  [64] = 64,
  [65] = 65,
  [66] = 66,
  [67] = 17,
  [68] = 68,
  [69] = 69,
  [70] = 70,
  [71] = 71,
  [72] = 72,
  [73] = 73,
  [74] = 74,
  [75] = 75,
  [76] = 76,
  [77] = 77,
  [78] = 78,
  [79] = 79,
  [80] = 79,
  [81] = 81,
  [82] = 82,
  [83] = 83,
  [84] = 84,
  [85] = 78,
  [86] = 86,
  [87] = 87,
  [88] = 88,
  [89] = 89,
  [90] = 90,
  [91] = 91,
  [92] = 60,
  [93] = 93,
  [94] = 94,
  [95] = 95,
  [96] = 96,
  [97] = 97,
  [98] = 98,
  [99] = 99,
  [100] = 100,
  [101] = 101,
  [102] = 102,
  [103] = 103,
  [104] = 104,
  [105] = 105,
  [106] = 106,
  [107] = 107,
  [108] = 108,
  [109] = 109,
  [110] = 110,
  [111] = 111,
  [112] = 112,
  [113] = 113,
  [114] = 114,
  [115] = 115,
  [116] = 116,
  [117] = 117,
  [118] = 118,
  [119] = 119,
  [120] = 120,
  [121] = 121,
  [122] = 122,
  [123] = 123,
  [124] = 124,
  [125] = 125,
  [126] = 126,
  [127] = 118,
  [128] = 83,
  [129] = 129,
  [130] = 130,
  [131] = 124,
  [132] = 122,
  [133] = 133,
  [134] = 134,
  [135] = 130,
  [136] = 120,
};

static bool ts_lex(TSLexer *lexer, TSStateId state) {
  START_LEXER();
  eof = lexer->eof(lexer);
  switch (state) {
    case 0:
      if (eof) ADVANCE(52);
      if (lookahead == '!') ADVANCE(89);
      if (lookahead == '"') ADVANCE(68);
      if (lookahead == '#') ADVANCE(104);
      if (lookahead == '$') ADVANCE(62);
      if (lookahead == '(') ADVANCE(57);
      if (lookahead == ')') ADVANCE(58);
      if (lookahead == ',') ADVANCE(105);
      if (lookahead == '-') ADVANCE(8);
      if (lookahead == '.') ADVANCE(7);
      if (lookahead == '0') ADVANCE(72);
      if (lookahead == ':') ADVANCE(59);
      if (lookahead == '=') ADVANCE(53);
      if (lookahead == '@') ADVANCE(88);
      if (lookahead == '[') ADVANCE(82);
      if (lookahead == ']') ADVANCE(83);
      if (lookahead == 'f') ADVANCE(9);
      if (lookahead == 'm') ADVANCE(42);
      if (lookahead == 'n') ADVANCE(45);
      if (lookahead == 'o') ADVANCE(26);
      if (lookahead == 'q') ADVANCE(43);
      if (lookahead == 's') ADVANCE(44);
      if (lookahead == 't') ADVANCE(35);
      if (lookahead == '{') ADVANCE(60);
      if (lookahead == '}') ADVANCE(61);
      if (('\t' <= lookahead && lookahead <= '\r') ||
          lookahead == ' ' ||
          lookahead == 65279) SKIP(0)
      if (('1' <= lookahead && lookahead <= '9')) ADVANCE(73);
      END_STATE();
    case 1:
      if (lookahead == '"') ADVANCE(68);
      if (lookahead == '#') ADVANCE(104);
      if (lookahead == '$') ADVANCE(62);
      if (lookahead == '-') ADVANCE(8);
      if (lookahead == '0') ADVANCE(72);
      if (lookahead == '[') ADVANCE(82);
      if (lookahead == ']') ADVANCE(83);
      if (lookahead == 'f') ADVANCE(90);
      if (lookahead == 'n') ADVANCE(100);
      if (lookahead == 't') ADVANCE(97);
      if (lookahead == '{') ADVANCE(60);
      if (('\t' <= lookahead && lookahead <= '\r') ||
          lookahead == ' ' ||
          lookahead == 65279) SKIP(1)
      if (('1' <= lookahead && lookahead <= '9')) ADVANCE(73);
      if (('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 2:
      if (lookahead == '"') ADVANCE(63);
      END_STATE();
    case 3:
      if (lookahead == '"') ADVANCE(67);
      if (lookahead == '#') ADVANCE(104);
      if (lookahead == '@') ADVANCE(88);
      if (lookahead == 'o') ADVANCE(96);
      if (lookahead == '{') ADVANCE(60);
      if (('\t' <= lookahead && lookahead <= '\r') ||
          lookahead == ' ' ||
          lookahead == 65279) SKIP(3)
      if (('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 4:
      if (lookahead == '"') ADVANCE(51);
      if (lookahead != 0) ADVANCE(66);
      END_STATE();
    case 5:
      if (lookahead == '#') ADVANCE(104);
      if (lookahead == '$') ADVANCE(62);
      if (lookahead == '(') ADVANCE(57);
      if (lookahead == ')') ADVANCE(58);
      if (lookahead == ',') ADVANCE(105);
      if (lookahead == '.') ADVANCE(7);
      if (lookahead == ':') ADVANCE(59);
      if (lookahead == '@') ADVANCE(88);
      if (lookahead == '[') ADVANCE(82);
      if (lookahead == '{') ADVANCE(60);
      if (lookahead == '}') ADVANCE(61);
      if (('\t' <= lookahead && lookahead <= '\r') ||
          lookahead == ' ' ||
          lookahead == 65279) SKIP(5)
      if (('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 6:
      if (lookahead == '.') ADVANCE(84);
      END_STATE();
    case 7:
      if (lookahead == '.') ADVANCE(6);
      END_STATE();
    case 8:
      if (lookahead == '0') ADVANCE(72);
      if (('1' <= lookahead && lookahead <= '9')) ADVANCE(73);
      END_STATE();
    case 9:
      if (lookahead == 'a') ADVANCE(22);
      if (lookahead == 'r') ADVANCE(10);
      END_STATE();
    case 10:
      if (lookahead == 'a') ADVANCE(18);
      END_STATE();
    case 11:
      if (lookahead == 'a') ADVANCE(38);
      END_STATE();
    case 12:
      if (lookahead == 'b') ADVANCE(36);
      END_STATE();
    case 13:
      if (lookahead == 'c') ADVANCE(34);
      END_STATE();
    case 14:
      if (lookahead == 'e') ADVANCE(76);
      END_STATE();
    case 15:
      if (lookahead == 'e') ADVANCE(78);
      END_STATE();
    case 16:
      if (lookahead == 'e') ADVANCE(29);
      END_STATE();
    case 17:
      if (lookahead == 'e') ADVANCE(33);
      END_STATE();
    case 18:
      if (lookahead == 'g') ADVANCE(25);
      END_STATE();
    case 19:
      if (lookahead == 'i') ADVANCE(30);
      END_STATE();
    case 20:
      if (lookahead == 'i') ADVANCE(32);
      END_STATE();
    case 21:
      if (lookahead == 'i') ADVANCE(31);
      END_STATE();
    case 22:
      if (lookahead == 'l') ADVANCE(37);
      END_STATE();
    case 23:
      if (lookahead == 'l') ADVANCE(80);
      END_STATE();
    case 24:
      if (lookahead == 'l') ADVANCE(23);
      END_STATE();
    case 25:
      if (lookahead == 'm') ADVANCE(16);
      END_STATE();
    case 26:
      if (lookahead == 'n') ADVANCE(86);
      END_STATE();
    case 27:
      if (lookahead == 'n') ADVANCE(55);
      END_STATE();
    case 28:
      if (lookahead == 'n') ADVANCE(56);
      END_STATE();
    case 29:
      if (lookahead == 'n') ADVANCE(39);
      END_STATE();
    case 30:
      if (lookahead == 'o') ADVANCE(27);
      END_STATE();
    case 31:
      if (lookahead == 'o') ADVANCE(28);
      END_STATE();
    case 32:
      if (lookahead == 'p') ADVANCE(41);
      END_STATE();
    case 33:
      if (lookahead == 'r') ADVANCE(47);
      END_STATE();
    case 34:
      if (lookahead == 'r') ADVANCE(20);
      END_STATE();
    case 35:
      if (lookahead == 'r') ADVANCE(46);
      END_STATE();
    case 36:
      if (lookahead == 's') ADVANCE(13);
      END_STATE();
    case 37:
      if (lookahead == 's') ADVANCE(15);
      END_STATE();
    case 38:
      if (lookahead == 't') ADVANCE(19);
      END_STATE();
    case 39:
      if (lookahead == 't') ADVANCE(85);
      END_STATE();
    case 40:
      if (lookahead == 't') ADVANCE(11);
      END_STATE();
    case 41:
      if (lookahead == 't') ADVANCE(21);
      END_STATE();
    case 42:
      if (lookahead == 'u') ADVANCE(40);
      END_STATE();
    case 43:
      if (lookahead == 'u') ADVANCE(17);
      END_STATE();
    case 44:
      if (lookahead == 'u') ADVANCE(12);
      END_STATE();
    case 45:
      if (lookahead == 'u') ADVANCE(24);
      END_STATE();
    case 46:
      if (lookahead == 'u') ADVANCE(14);
      END_STATE();
    case 47:
      if (lookahead == 'y') ADVANCE(54);
      END_STATE();
    case 48:
      if (lookahead == '+' ||
          lookahead == '-') ADVANCE(50);
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(75);
      END_STATE();
    case 49:
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(74);
      END_STATE();
    case 50:
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(75);
      END_STATE();
    case 51:
      if (lookahead != 0 &&
          lookahead != '"') ADVANCE(66);
      END_STATE();
    case 52:
      ACCEPT_TOKEN(ts_builtin_sym_end);
      END_STATE();
    case 53:
      ACCEPT_TOKEN(anon_sym_EQ);
      END_STATE();
    case 54:
      ACCEPT_TOKEN(anon_sym_query);
      END_STATE();
    case 55:
      ACCEPT_TOKEN(anon_sym_mutation);
      END_STATE();
    case 56:
      ACCEPT_TOKEN(anon_sym_subscription);
      END_STATE();
    case 57:
      ACCEPT_TOKEN(anon_sym_LPAREN);
      END_STATE();
    case 58:
      ACCEPT_TOKEN(anon_sym_RPAREN);
      END_STATE();
    case 59:
      ACCEPT_TOKEN(anon_sym_COLON);
      END_STATE();
    case 60:
      ACCEPT_TOKEN(anon_sym_LBRACE);
      END_STATE();
    case 61:
      ACCEPT_TOKEN(anon_sym_RBRACE);
      END_STATE();
    case 62:
      ACCEPT_TOKEN(anon_sym_DOLLAR);
      END_STATE();
    case 63:
      ACCEPT_TOKEN(anon_sym_DQUOTE_DQUOTE_DQUOTE);
      END_STATE();
    case 64:
      ACCEPT_TOKEN(aux_sym_string_value_token1);
      if (lookahead == '\n') ADVANCE(66);
      if (lookahead == '"') ADVANCE(103);
      if (lookahead != 0) ADVANCE(64);
      END_STATE();
    case 65:
      ACCEPT_TOKEN(aux_sym_string_value_token1);
      if (('\t' <= lookahead && lookahead <= '\r') ||
          lookahead == ' ' ||
          lookahead == 65279) ADVANCE(65);
      if (lookahead == '"') ADVANCE(4);
      if (lookahead == '#') ADVANCE(64);
      if (lookahead != 0) ADVANCE(66);
      END_STATE();
    case 66:
      ACCEPT_TOKEN(aux_sym_string_value_token1);
      if (lookahead != 0 &&
          lookahead != '"') ADVANCE(66);
      if (lookahead == '"') ADVANCE(4);
      END_STATE();
    case 67:
      ACCEPT_TOKEN(anon_sym_DQUOTE);
      END_STATE();
    case 68:
      ACCEPT_TOKEN(anon_sym_DQUOTE);
      if (lookahead == '"') ADVANCE(2);
      END_STATE();
    case 69:
      ACCEPT_TOKEN(aux_sym_string_value_token2);
      if (lookahead == '#') ADVANCE(70);
      if (lookahead == '\t' ||
          (11 <= lookahead && lookahead <= '\r') ||
          lookahead == ' ' ||
          lookahead == 65279) ADVANCE(69);
      if (lookahead != 0 &&
          lookahead != '\n' &&
          lookahead != '"' &&
          lookahead != '\\') ADVANCE(71);
      END_STATE();
    case 70:
      ACCEPT_TOKEN(aux_sym_string_value_token2);
      if (lookahead == '"' ||
          lookahead == '\\') ADVANCE(104);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(70);
      END_STATE();
    case 71:
      ACCEPT_TOKEN(aux_sym_string_value_token2);
      if (lookahead != 0 &&
          lookahead != '\n' &&
          lookahead != '"' &&
          lookahead != '\\') ADVANCE(71);
      END_STATE();
    case 72:
      ACCEPT_TOKEN(sym_int_value);
      if (lookahead == '.') ADVANCE(49);
      if (lookahead == 'E' ||
          lookahead == 'e') ADVANCE(48);
      END_STATE();
    case 73:
      ACCEPT_TOKEN(sym_int_value);
      if (lookahead == '.') ADVANCE(49);
      if (lookahead == 'E' ||
          lookahead == 'e') ADVANCE(48);
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(73);
      END_STATE();
    case 74:
      ACCEPT_TOKEN(sym_float_value);
      if (lookahead == 'E' ||
          lookahead == 'e') ADVANCE(48);
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(74);
      END_STATE();
    case 75:
      ACCEPT_TOKEN(sym_float_value);
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(75);
      END_STATE();
    case 76:
      ACCEPT_TOKEN(anon_sym_true);
      END_STATE();
    case 77:
      ACCEPT_TOKEN(anon_sym_true);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 78:
      ACCEPT_TOKEN(anon_sym_false);
      END_STATE();
    case 79:
      ACCEPT_TOKEN(anon_sym_false);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 80:
      ACCEPT_TOKEN(sym_null_value);
      END_STATE();
    case 81:
      ACCEPT_TOKEN(sym_null_value);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 82:
      ACCEPT_TOKEN(anon_sym_LBRACK);
      END_STATE();
    case 83:
      ACCEPT_TOKEN(anon_sym_RBRACK);
      END_STATE();
    case 84:
      ACCEPT_TOKEN(anon_sym_DOT_DOT_DOT);
      END_STATE();
    case 85:
      ACCEPT_TOKEN(anon_sym_fragment);
      END_STATE();
    case 86:
      ACCEPT_TOKEN(anon_sym_on);
      END_STATE();
    case 87:
      ACCEPT_TOKEN(anon_sym_on);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 88:
      ACCEPT_TOKEN(anon_sym_AT);
      END_STATE();
    case 89:
      ACCEPT_TOKEN(anon_sym_BANG);
      END_STATE();
    case 90:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'a') ADVANCE(93);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('b' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 91:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'e') ADVANCE(77);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 92:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'e') ADVANCE(79);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 93:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'l') ADVANCE(98);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 94:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'l') ADVANCE(81);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 95:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'l') ADVANCE(94);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 96:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'n') ADVANCE(87);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 97:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'r') ADVANCE(99);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 98:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 's') ADVANCE(92);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 99:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'u') ADVANCE(91);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 100:
      ACCEPT_TOKEN(sym_name);
      if (lookahead == 'u') ADVANCE(95);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 101:
      ACCEPT_TOKEN(sym_name);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(101);
      END_STATE();
    case 102:
      ACCEPT_TOKEN(sym_comment);
      if (lookahead == '\n') ADVANCE(66);
      if (lookahead == '"') ADVANCE(104);
      if (lookahead != 0) ADVANCE(64);
      END_STATE();
    case 103:
      ACCEPT_TOKEN(sym_comment);
      if (lookahead == '\n') ADVANCE(66);
      if (lookahead == '"') ADVANCE(102);
      if (lookahead != 0) ADVANCE(64);
      END_STATE();
    case 104:
      ACCEPT_TOKEN(sym_comment);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(104);
      END_STATE();
    case 105:
      ACCEPT_TOKEN(sym_comma);
      END_STATE();
    default:
      return false;
  }
}

static const TSLexMode ts_lex_modes[STATE_COUNT] = {
  [0] = {.lex_state = 0},
  [1] = {.lex_state = 0},
  [2] = {.lex_state = 1},
  [3] = {.lex_state = 1},
  [4] = {.lex_state = 1},
  [5] = {.lex_state = 1},
  [6] = {.lex_state = 1},
  [7] = {.lex_state = 1},
  [8] = {.lex_state = 1},
  [9] = {.lex_state = 1},
  [10] = {.lex_state = 0},
  [11] = {.lex_state = 0},
  [12] = {.lex_state = 1},
  [13] = {.lex_state = 1},
  [14] = {.lex_state = 1},
  [15] = {.lex_state = 1},
  [16] = {.lex_state = 1},
  [17] = {.lex_state = 1},
  [18] = {.lex_state = 1},
  [19] = {.lex_state = 1},
  [20] = {.lex_state = 1},
  [21] = {.lex_state = 3},
  [22] = {.lex_state = 5},
  [23] = {.lex_state = 5},
  [24] = {.lex_state = 5},
  [25] = {.lex_state = 5},
  [26] = {.lex_state = 5},
  [27] = {.lex_state = 5},
  [28] = {.lex_state = 5},
  [29] = {.lex_state = 0},
  [30] = {.lex_state = 5},
  [31] = {.lex_state = 5},
  [32] = {.lex_state = 0},
  [33] = {.lex_state = 5},
  [34] = {.lex_state = 5},
  [35] = {.lex_state = 5},
  [36] = {.lex_state = 0},
  [37] = {.lex_state = 5},
  [38] = {.lex_state = 0},
  [39] = {.lex_state = 5},
  [40] = {.lex_state = 0},
  [41] = {.lex_state = 5},
  [42] = {.lex_state = 5},
  [43] = {.lex_state = 5},
  [44] = {.lex_state = 0},
  [45] = {.lex_state = 5},
  [46] = {.lex_state = 0},
  [47] = {.lex_state = 0},
  [48] = {.lex_state = 0},
  [49] = {.lex_state = 0},
  [50] = {.lex_state = 5},
  [51] = {.lex_state = 5},
  [52] = {.lex_state = 5},
  [53] = {.lex_state = 5},
  [54] = {.lex_state = 0},
  [55] = {.lex_state = 5},
  [56] = {.lex_state = 5},
  [57] = {.lex_state = 0},
  [58] = {.lex_state = 0},
  [59] = {.lex_state = 0},
  [60] = {.lex_state = 0},
  [61] = {.lex_state = 0},
  [62] = {.lex_state = 5},
  [63] = {.lex_state = 0},
  [64] = {.lex_state = 0},
  [65] = {.lex_state = 0},
  [66] = {.lex_state = 0},
  [67] = {.lex_state = 5},
  [68] = {.lex_state = 5},
  [69] = {.lex_state = 0},
  [70] = {.lex_state = 0},
  [71] = {.lex_state = 5},
  [72] = {.lex_state = 0},
  [73] = {.lex_state = 5},
  [74] = {.lex_state = 5},
  [75] = {.lex_state = 0},
  [76] = {.lex_state = 5},
  [77] = {.lex_state = 0},
  [78] = {.lex_state = 5},
  [79] = {.lex_state = 5},
  [80] = {.lex_state = 5},
  [81] = {.lex_state = 0},
  [82] = {.lex_state = 5},
  [83] = {.lex_state = 5},
  [84] = {.lex_state = 5},
  [85] = {.lex_state = 5},
  [86] = {.lex_state = 5},
  [87] = {.lex_state = 5},
  [88] = {.lex_state = 5},
  [89] = {.lex_state = 5},
  [90] = {.lex_state = 5},
  [91] = {.lex_state = 5},
  [92] = {.lex_state = 5},
  [93] = {.lex_state = 0},
  [94] = {.lex_state = 5},
  [95] = {.lex_state = 5},
  [96] = {.lex_state = 5},
  [97] = {.lex_state = 5},
  [98] = {.lex_state = 0},
  [99] = {.lex_state = 5},
  [100] = {.lex_state = 5},
  [101] = {.lex_state = 0},
  [102] = {.lex_state = 0},
  [103] = {.lex_state = 0},
  [104] = {.lex_state = 0},
  [105] = {.lex_state = 0},
  [106] = {.lex_state = 0},
  [107] = {.lex_state = 0},
  [108] = {.lex_state = 0},
  [109] = {.lex_state = 0},
  [110] = {.lex_state = 5},
  [111] = {.lex_state = 0},
  [112] = {.lex_state = 5},
  [113] = {.lex_state = 0},
  [114] = {.lex_state = 5},
  [115] = {.lex_state = 0},
  [116] = {.lex_state = 5},
  [117] = {.lex_state = 0},
  [118] = {.lex_state = 5},
  [119] = {.lex_state = 5},
  [120] = {.lex_state = 69},
  [121] = {.lex_state = 0},
  [122] = {.lex_state = 3},
  [123] = {.lex_state = 0},
  [124] = {.lex_state = 0},
  [125] = {.lex_state = 5},
  [126] = {.lex_state = 0},
  [127] = {.lex_state = 5},
  [128] = {.lex_state = 0},
  [129] = {.lex_state = 5},
  [130] = {.lex_state = 65},
  [131] = {.lex_state = 0},
  [132] = {.lex_state = 3},
  [133] = {.lex_state = 0},
  [134] = {.lex_state = 0},
  [135] = {.lex_state = 65},
  [136] = {.lex_state = 69},
};

static const uint16_t ts_parse_table[LARGE_STATE_COUNT][SYMBOL_COUNT] = {
  [0] = {
    [ts_builtin_sym_end] = ACTIONS(1),
    [anon_sym_EQ] = ACTIONS(1),
    [anon_sym_query] = ACTIONS(1),
    [anon_sym_mutation] = ACTIONS(1),
    [anon_sym_subscription] = ACTIONS(1),
    [anon_sym_LPAREN] = ACTIONS(1),
    [anon_sym_RPAREN] = ACTIONS(1),
    [anon_sym_COLON] = ACTIONS(1),
    [anon_sym_LBRACE] = ACTIONS(1),
    [anon_sym_RBRACE] = ACTIONS(1),
    [anon_sym_DOLLAR] = ACTIONS(1),
    [anon_sym_DQUOTE_DQUOTE_DQUOTE] = ACTIONS(1),
    [anon_sym_DQUOTE] = ACTIONS(1),
    [sym_int_value] = ACTIONS(1),
    [sym_float_value] = ACTIONS(1),
    [anon_sym_true] = ACTIONS(1),
    [anon_sym_false] = ACTIONS(1),
    [sym_null_value] = ACTIONS(1),
    [anon_sym_LBRACK] = ACTIONS(1),
    [anon_sym_RBRACK] = ACTIONS(1),
    [anon_sym_DOT_DOT_DOT] = ACTIONS(1),
    [anon_sym_fragment] = ACTIONS(1),
    [anon_sym_on] = ACTIONS(1),
    [anon_sym_AT] = ACTIONS(1),
    [anon_sym_BANG] = ACTIONS(1),
    [sym_comment] = ACTIONS(3),
    [sym_comma] = ACTIONS(1),
  },
  [1] = {
    [sym_source_file] = STATE(134),
    [sym_document] = STATE(133),
    [sym_definition] = STATE(11),
    [sym_executable_definition] = STATE(49),
    [sym_operation_definition] = STATE(54),
    [sym_operation_type] = STATE(30),
    [sym_selection_set] = STATE(46),
    [sym_fragment_definition] = STATE(54),
    [aux_sym_document_repeat1] = STATE(11),
    [anon_sym_query] = ACTIONS(5),
    [anon_sym_mutation] = ACTIONS(5),
    [anon_sym_subscription] = ACTIONS(5),
    [anon_sym_LBRACE] = ACTIONS(7),
    [anon_sym_fragment] = ACTIONS(9),
    [sym_comment] = ACTIONS(3),
  },
};

static const uint16_t ts_small_parse_table[] = {
  [0] = 13,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(11), 1,
      anon_sym_LBRACE,
    ACTIONS(13), 1,
      anon_sym_DOLLAR,
    ACTIONS(15), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
    ACTIONS(17), 1,
      anon_sym_DQUOTE,
    ACTIONS(21), 1,
      sym_float_value,
    ACTIONS(25), 1,
      anon_sym_LBRACK,
    ACTIONS(27), 1,
      anon_sym_RBRACK,
    ACTIONS(29), 1,
      sym_name,
    ACTIONS(19), 2,
      sym_int_value,
      sym_null_value,
    ACTIONS(23), 2,
      anon_sym_true,
      anon_sym_false,
    STATE(4), 2,
      sym_value,
      aux_sym_list_value_repeat1,
    STATE(12), 6,
      sym_variable,
      sym_string_value,
      sym_boolean_value,
      sym_enum_value,
      sym_list_value,
      sym_object_value,
  [48] = 13,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(31), 1,
      anon_sym_LBRACE,
    ACTIONS(34), 1,
      anon_sym_DOLLAR,
    ACTIONS(37), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
    ACTIONS(40), 1,
      anon_sym_DQUOTE,
    ACTIONS(46), 1,
      sym_float_value,
    ACTIONS(52), 1,
      anon_sym_LBRACK,
    ACTIONS(55), 1,
      anon_sym_RBRACK,
    ACTIONS(57), 1,
      sym_name,
    ACTIONS(43), 2,
      sym_int_value,
      sym_null_value,
    ACTIONS(49), 2,
      anon_sym_true,
      anon_sym_false,
    STATE(3), 2,
      sym_value,
      aux_sym_list_value_repeat1,
    STATE(12), 6,
      sym_variable,
      sym_string_value,
      sym_boolean_value,
      sym_enum_value,
      sym_list_value,
      sym_object_value,
  [96] = 13,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(11), 1,
      anon_sym_LBRACE,
    ACTIONS(13), 1,
      anon_sym_DOLLAR,
    ACTIONS(15), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
    ACTIONS(17), 1,
      anon_sym_DQUOTE,
    ACTIONS(21), 1,
      sym_float_value,
    ACTIONS(25), 1,
      anon_sym_LBRACK,
    ACTIONS(29), 1,
      sym_name,
    ACTIONS(60), 1,
      anon_sym_RBRACK,
    ACTIONS(19), 2,
      sym_int_value,
      sym_null_value,
    ACTIONS(23), 2,
      anon_sym_true,
      anon_sym_false,
    STATE(3), 2,
      sym_value,
      aux_sym_list_value_repeat1,
    STATE(12), 6,
      sym_variable,
      sym_string_value,
      sym_boolean_value,
      sym_enum_value,
      sym_list_value,
      sym_object_value,
  [144] = 13,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(11), 1,
      anon_sym_LBRACE,
    ACTIONS(13), 1,
      anon_sym_DOLLAR,
    ACTIONS(15), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
    ACTIONS(17), 1,
      anon_sym_DQUOTE,
    ACTIONS(21), 1,
      sym_float_value,
    ACTIONS(25), 1,
      anon_sym_LBRACK,
    ACTIONS(29), 1,
      sym_name,
    ACTIONS(62), 1,
      anon_sym_RBRACK,
    ACTIONS(19), 2,
      sym_int_value,
      sym_null_value,
    ACTIONS(23), 2,
      anon_sym_true,
      anon_sym_false,
    STATE(6), 2,
      sym_value,
      aux_sym_list_value_repeat1,
    STATE(12), 6,
      sym_variable,
      sym_string_value,
      sym_boolean_value,
      sym_enum_value,
      sym_list_value,
      sym_object_value,
  [192] = 13,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(11), 1,
      anon_sym_LBRACE,
    ACTIONS(13), 1,
      anon_sym_DOLLAR,
    ACTIONS(15), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
    ACTIONS(17), 1,
      anon_sym_DQUOTE,
    ACTIONS(21), 1,
      sym_float_value,
    ACTIONS(25), 1,
      anon_sym_LBRACK,
    ACTIONS(29), 1,
      sym_name,
    ACTIONS(64), 1,
      anon_sym_RBRACK,
    ACTIONS(19), 2,
      sym_int_value,
      sym_null_value,
    ACTIONS(23), 2,
      anon_sym_true,
      anon_sym_false,
    STATE(3), 2,
      sym_value,
      aux_sym_list_value_repeat1,
    STATE(12), 6,
      sym_variable,
      sym_string_value,
      sym_boolean_value,
      sym_enum_value,
      sym_list_value,
      sym_object_value,
  [240] = 12,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(66), 1,
      anon_sym_LBRACE,
    ACTIONS(68), 1,
      anon_sym_DOLLAR,
    ACTIONS(70), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
    ACTIONS(72), 1,
      anon_sym_DQUOTE,
    ACTIONS(76), 1,
      sym_float_value,
    ACTIONS(80), 1,
      anon_sym_LBRACK,
    ACTIONS(82), 1,
      sym_name,
    STATE(95), 1,
      sym_value,
    ACTIONS(74), 2,
      sym_int_value,
      sym_null_value,
    ACTIONS(78), 2,
      anon_sym_true,
      anon_sym_false,
    STATE(56), 6,
      sym_variable,
      sym_string_value,
      sym_boolean_value,
      sym_enum_value,
      sym_list_value,
      sym_object_value,
  [284] = 12,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(66), 1,
      anon_sym_LBRACE,
    ACTIONS(68), 1,
      anon_sym_DOLLAR,
    ACTIONS(70), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
    ACTIONS(72), 1,
      anon_sym_DQUOTE,
    ACTIONS(76), 1,
      sym_float_value,
    ACTIONS(80), 1,
      anon_sym_LBRACK,
    ACTIONS(82), 1,
      sym_name,
    STATE(81), 1,
      sym_value,
    ACTIONS(74), 2,
      sym_int_value,
      sym_null_value,
    ACTIONS(78), 2,
      anon_sym_true,
      anon_sym_false,
    STATE(56), 6,
      sym_variable,
      sym_string_value,
      sym_boolean_value,
      sym_enum_value,
      sym_list_value,
      sym_object_value,
  [328] = 12,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(66), 1,
      anon_sym_LBRACE,
    ACTIONS(68), 1,
      anon_sym_DOLLAR,
    ACTIONS(70), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
    ACTIONS(72), 1,
      anon_sym_DQUOTE,
    ACTIONS(76), 1,
      sym_float_value,
    ACTIONS(80), 1,
      anon_sym_LBRACK,
    ACTIONS(82), 1,
      sym_name,
    STATE(116), 1,
      sym_value,
    ACTIONS(74), 2,
      sym_int_value,
      sym_null_value,
    ACTIONS(78), 2,
      anon_sym_true,
      anon_sym_false,
    STATE(56), 6,
      sym_variable,
      sym_string_value,
      sym_boolean_value,
      sym_enum_value,
      sym_list_value,
      sym_object_value,
  [372] = 10,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(84), 1,
      ts_builtin_sym_end,
    ACTIONS(89), 1,
      anon_sym_LBRACE,
    ACTIONS(92), 1,
      anon_sym_fragment,
    STATE(30), 1,
      sym_operation_type,
    STATE(46), 1,
      sym_selection_set,
    STATE(49), 1,
      sym_executable_definition,
    STATE(10), 2,
      sym_definition,
      aux_sym_document_repeat1,
    STATE(54), 2,
      sym_operation_definition,
      sym_fragment_definition,
    ACTIONS(86), 3,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
  [407] = 10,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    ACTIONS(9), 1,
      anon_sym_fragment,
    ACTIONS(95), 1,
      ts_builtin_sym_end,
    STATE(30), 1,
      sym_operation_type,
    STATE(46), 1,
      sym_selection_set,
    STATE(49), 1,
      sym_executable_definition,
    STATE(10), 2,
      sym_definition,
      aux_sym_document_repeat1,
    STATE(54), 2,
      sym_operation_definition,
      sym_fragment_definition,
    ACTIONS(5), 3,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
  [442] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(97), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(99), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [462] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(101), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(103), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [482] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(105), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(107), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [502] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(109), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(111), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [522] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(113), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(115), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [542] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(117), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(119), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [562] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(121), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(123), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [582] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(125), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(127), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [602] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(129), 6,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
      sym_float_value,
      anon_sym_LBRACK,
      anon_sym_RBRACK,
    ACTIONS(131), 6,
      anon_sym_DQUOTE,
      sym_int_value,
      anon_sym_true,
      anon_sym_false,
      sym_null_value,
      sym_name,
  [622] = 10,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    ACTIONS(135), 1,
      anon_sym_on,
    ACTIONS(137), 1,
      anon_sym_AT,
    ACTIONS(139), 1,
      sym_name,
    STATE(43), 1,
      sym_fragment_name,
    STATE(63), 1,
      sym_type_condition,
    STATE(90), 1,
      sym_selection_set,
    STATE(102), 1,
      sym_directives,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [654] = 9,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    ACTIONS(141), 1,
      anon_sym_LPAREN,
    ACTIONS(143), 1,
      anon_sym_COLON,
    STATE(39), 1,
      sym_arguments,
    STATE(73), 1,
      sym_directive,
    STATE(87), 1,
      sym_selection_set,
    ACTIONS(145), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [684] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(141), 1,
      anon_sym_LPAREN,
    STATE(37), 1,
      sym_arguments,
    ACTIONS(147), 8,
      anon_sym_RPAREN,
      anon_sym_LBRACE,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_DOT_DOT_DOT,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [704] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(151), 1,
      anon_sym_AT,
    STATE(24), 2,
      sym_directive,
      aux_sym_directives_repeat1,
    ACTIONS(149), 7,
      anon_sym_RPAREN,
      anon_sym_LBRACE,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_DOT_DOT_DOT,
      sym_name,
      sym_comma,
  [724] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(137), 1,
      anon_sym_AT,
    STATE(24), 2,
      sym_directive,
      aux_sym_directives_repeat1,
    ACTIONS(154), 7,
      anon_sym_RPAREN,
      anon_sym_LBRACE,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_DOT_DOT_DOT,
      sym_name,
      sym_comma,
  [744] = 8,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    ACTIONS(141), 1,
      anon_sym_LPAREN,
    STATE(42), 1,
      sym_arguments,
    STATE(74), 1,
      sym_directive,
    STATE(99), 1,
      sym_selection_set,
    ACTIONS(156), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [771] = 7,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(158), 1,
      anon_sym_RBRACE,
    ACTIONS(160), 1,
      anon_sym_DOT_DOT_DOT,
    ACTIONS(162), 1,
      sym_name,
    STATE(129), 1,
      sym_alias,
    STATE(31), 2,
      sym_selection,
      aux_sym_selection_set_repeat1,
    STATE(89), 3,
      sym_field,
      sym_fragment_spread,
      sym_inline_fragment,
  [796] = 7,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(160), 1,
      anon_sym_DOT_DOT_DOT,
    ACTIONS(162), 1,
      sym_name,
    ACTIONS(164), 1,
      anon_sym_RBRACE,
    STATE(129), 1,
      sym_alias,
    STATE(31), 2,
      sym_selection,
      aux_sym_selection_set_repeat1,
    STATE(89), 3,
      sym_field,
      sym_fragment_spread,
      sym_inline_fragment,
  [821] = 8,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(137), 1,
      anon_sym_AT,
    ACTIONS(166), 1,
      anon_sym_EQ,
    ACTIONS(170), 1,
      sym_comma,
    STATE(44), 1,
      sym_default_value,
    STATE(93), 1,
      sym_directives,
    ACTIONS(168), 2,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [848] = 9,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    ACTIONS(172), 1,
      anon_sym_LPAREN,
    ACTIONS(174), 1,
      sym_name,
    STATE(65), 1,
      sym_selection_set,
    STATE(69), 1,
      sym_variable_definitions,
    STATE(108), 1,
      sym_directives,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [877] = 7,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(176), 1,
      anon_sym_RBRACE,
    ACTIONS(178), 1,
      anon_sym_DOT_DOT_DOT,
    ACTIONS(181), 1,
      sym_name,
    STATE(129), 1,
      sym_alias,
    STATE(31), 2,
      sym_selection,
      aux_sym_selection_set_repeat1,
    STATE(89), 3,
      sym_field,
      sym_fragment_spread,
      sym_inline_fragment,
  [902] = 8,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    ACTIONS(172), 1,
      anon_sym_LPAREN,
    STATE(47), 1,
      sym_variable_definitions,
    STATE(64), 1,
      sym_selection_set,
    STATE(115), 1,
      sym_directives,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [928] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(160), 1,
      anon_sym_DOT_DOT_DOT,
    ACTIONS(162), 1,
      sym_name,
    STATE(129), 1,
      sym_alias,
    STATE(28), 2,
      sym_selection,
      aux_sym_selection_set_repeat1,
    STATE(89), 3,
      sym_field,
      sym_fragment_spread,
      sym_inline_fragment,
  [950] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(184), 8,
      anon_sym_RPAREN,
      anon_sym_LBRACE,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_DOT_DOT_DOT,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [964] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(160), 1,
      anon_sym_DOT_DOT_DOT,
    ACTIONS(162), 1,
      sym_name,
    STATE(129), 1,
      sym_alias,
    STATE(27), 2,
      sym_selection,
      aux_sym_selection_set_repeat1,
    STATE(89), 3,
      sym_field,
      sym_fragment_spread,
      sym_inline_fragment,
  [986] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(186), 8,
      anon_sym_EQ,
      anon_sym_RPAREN,
      anon_sym_LBRACE,
      anon_sym_DOLLAR,
      anon_sym_RBRACK,
      anon_sym_AT,
      anon_sym_BANG,
      sym_comma,
  [1000] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(188), 8,
      anon_sym_RPAREN,
      anon_sym_LBRACE,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_DOT_DOT_DOT,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1014] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(192), 1,
      anon_sym_BANG,
    ACTIONS(190), 6,
      anon_sym_EQ,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
      anon_sym_RBRACK,
      anon_sym_AT,
      sym_comma,
  [1029] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    STATE(74), 1,
      sym_directive,
    STATE(99), 1,
      sym_selection_set,
    ACTIONS(156), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1050] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(194), 7,
      anon_sym_EQ,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
      anon_sym_RBRACK,
      anon_sym_AT,
      anon_sym_BANG,
      sym_comma,
  [1063] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(109), 7,
      anon_sym_RPAREN,
      anon_sym_COLON,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1076] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    STATE(71), 1,
      sym_directive,
    STATE(94), 1,
      sym_selection_set,
    ACTIONS(196), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1097] = 5,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(137), 1,
      anon_sym_AT,
    STATE(97), 1,
      sym_directives,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
    ACTIONS(198), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1116] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(137), 1,
      anon_sym_AT,
    ACTIONS(202), 1,
      sym_comma,
    STATE(98), 1,
      sym_directives,
    ACTIONS(200), 2,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [1137] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(204), 1,
      anon_sym_LBRACK,
    ACTIONS(206), 1,
      sym_name,
    STATE(29), 1,
      sym_type,
    STATE(61), 1,
      sym_non_null_type,
    STATE(38), 2,
      sym_named_type,
      sym_list_type,
  [1157] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(208), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1169] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    STATE(66), 1,
      sym_selection_set,
    STATE(104), 1,
      sym_directives,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [1189] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(210), 6,
      anon_sym_EQ,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
      anon_sym_RBRACK,
      anon_sym_AT,
      sym_comma,
  [1201] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(212), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1213] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(101), 6,
      anon_sym_RPAREN,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1225] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(105), 6,
      anon_sym_RPAREN,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1237] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(125), 6,
      anon_sym_RPAREN,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1249] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(113), 6,
      anon_sym_RPAREN,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1261] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(214), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1273] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(129), 6,
      anon_sym_RPAREN,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1285] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(97), 6,
      anon_sym_RPAREN,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1297] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(216), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1309] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    STATE(59), 1,
      sym_selection_set,
    STATE(106), 1,
      sym_directives,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [1329] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(218), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1341] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(220), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1353] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(190), 6,
      anon_sym_EQ,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
      anon_sym_RBRACK,
      anon_sym_AT,
      sym_comma,
  [1365] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(121), 6,
      anon_sym_RPAREN,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1377] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    STATE(88), 1,
      sym_selection_set,
    STATE(109), 1,
      sym_directives,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [1397] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(222), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1409] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(224), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1421] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(226), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1433] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(117), 6,
      anon_sym_RPAREN,
      anon_sym_RBRACE,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_name,
      sym_comma,
  [1445] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(204), 1,
      anon_sym_LBRACK,
    ACTIONS(206), 1,
      sym_name,
    STATE(61), 1,
      sym_non_null_type,
    STATE(117), 1,
      sym_type,
    STATE(38), 2,
      sym_named_type,
      sym_list_type,
  [1465] = 6,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    ACTIONS(137), 1,
      anon_sym_AT,
    STATE(64), 1,
      sym_selection_set,
    STATE(115), 1,
      sym_directives,
    STATE(25), 2,
      sym_directive,
      aux_sym_directives_repeat1,
  [1485] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(228), 6,
      ts_builtin_sym_end,
      anon_sym_query,
      anon_sym_mutation,
      anon_sym_subscription,
      anon_sym_LBRACE,
      anon_sym_fragment,
  [1497] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    STATE(96), 1,
      sym_selection_set,
    ACTIONS(230), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1512] = 5,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(232), 1,
      anon_sym_RPAREN,
    ACTIONS(234), 1,
      anon_sym_DOLLAR,
    STATE(121), 1,
      sym_variable,
    STATE(72), 2,
      sym_variable_definition,
      aux_sym_variable_definitions_repeat1,
  [1529] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    STATE(99), 1,
      sym_selection_set,
    ACTIONS(156), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1544] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    STATE(94), 1,
      sym_selection_set,
    ACTIONS(196), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1559] = 5,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(68), 1,
      anon_sym_DOLLAR,
    ACTIONS(237), 1,
      anon_sym_RPAREN,
    STATE(121), 1,
      sym_variable,
    STATE(72), 2,
      sym_variable_definition,
      aux_sym_variable_definitions_repeat1,
  [1576] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(239), 1,
      anon_sym_RPAREN,
    ACTIONS(241), 1,
      sym_name,
    STATE(76), 2,
      sym_argument,
      aux_sym_arguments_repeat1,
  [1590] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(68), 1,
      anon_sym_DOLLAR,
    STATE(121), 1,
      sym_variable,
    STATE(75), 2,
      sym_variable_definition,
      aux_sym_variable_definitions_repeat1,
  [1604] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(244), 1,
      anon_sym_RBRACE,
    ACTIONS(246), 1,
      sym_name,
    STATE(84), 2,
      sym_object_field,
      aux_sym_object_value_repeat1,
  [1618] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(246), 1,
      sym_name,
    ACTIONS(248), 1,
      anon_sym_RBRACE,
    STATE(85), 2,
      sym_object_field,
      aux_sym_object_value_repeat1,
  [1632] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(246), 1,
      sym_name,
    ACTIONS(250), 1,
      anon_sym_RBRACE,
    STATE(78), 2,
      sym_object_field,
      aux_sym_object_value_repeat1,
  [1646] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(252), 4,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
      anon_sym_AT,
      sym_comma,
  [1656] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(254), 1,
      anon_sym_RPAREN,
    ACTIONS(256), 1,
      sym_name,
    STATE(76), 2,
      sym_argument,
      aux_sym_arguments_repeat1,
  [1670] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(258), 4,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      anon_sym_AT,
      sym_name,
  [1680] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(260), 1,
      anon_sym_RBRACE,
    ACTIONS(262), 1,
      sym_name,
    STATE(84), 2,
      sym_object_field,
      aux_sym_object_value_repeat1,
  [1694] = 4,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(246), 1,
      sym_name,
    ACTIONS(265), 1,
      anon_sym_RBRACE,
    STATE(84), 2,
      sym_object_field,
      aux_sym_object_value_repeat1,
  [1708] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(267), 4,
      anon_sym_LPAREN,
      anon_sym_LBRACE,
      anon_sym_AT,
      sym_name,
  [1718] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(156), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1727] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(269), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1736] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(271), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1745] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(273), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1754] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(256), 1,
      sym_name,
    STATE(82), 2,
      sym_argument,
      aux_sym_arguments_repeat1,
  [1765] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(220), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1774] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(202), 1,
      sym_comma,
    ACTIONS(200), 2,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
  [1785] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(230), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1794] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(277), 1,
      sym_comma,
    ACTIONS(275), 2,
      anon_sym_RBRACE,
      sym_name,
  [1805] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(279), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1814] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(281), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1823] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(285), 1,
      sym_comma,
    ACTIONS(283), 2,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
  [1834] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(196), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1843] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(287), 3,
      anon_sym_RBRACE,
      anon_sym_DOT_DOT_DOT,
      sym_name,
  [1852] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(289), 2,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
  [1860] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    STATE(88), 1,
      sym_selection_set,
  [1870] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(291), 1,
      anon_sym_on,
    STATE(58), 1,
      sym_type_condition,
  [1880] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    STATE(57), 1,
      sym_selection_set,
  [1890] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(293), 2,
      anon_sym_LBRACE,
      anon_sym_AT,
  [1898] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    STATE(70), 1,
      sym_selection_set,
  [1908] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(283), 2,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
  [1916] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    STATE(64), 1,
      sym_selection_set,
  [1926] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(133), 1,
      anon_sym_LBRACE,
    STATE(100), 1,
      sym_selection_set,
  [1936] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(206), 1,
      sym_name,
    STATE(111), 1,
      sym_named_type,
  [1946] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(295), 2,
      anon_sym_LBRACE,
      anon_sym_AT,
  [1954] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(297), 1,
      sym_name,
    STATE(103), 1,
      sym_fragment_name,
  [1964] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(200), 2,
      anon_sym_RPAREN,
      anon_sym_DOLLAR,
  [1972] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(299), 2,
      anon_sym_RBRACE,
      sym_name,
  [1980] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(7), 1,
      anon_sym_LBRACE,
    STATE(66), 1,
      sym_selection_set,
  [1990] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(301), 2,
      anon_sym_RPAREN,
      sym_name,
  [1998] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(303), 1,
      anon_sym_RBRACK,
  [2005] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(305), 1,
      sym_name,
  [2012] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(307), 1,
      sym_name,
  [2019] = 2,
    ACTIONS(309), 1,
      aux_sym_string_value_token2,
    ACTIONS(311), 1,
      sym_comment,
  [2026] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(313), 1,
      anon_sym_COLON,
  [2033] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(315), 1,
      anon_sym_DQUOTE,
  [2040] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(317), 1,
      anon_sym_COLON,
  [2047] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(315), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
  [2054] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(319), 1,
      sym_name,
  [2061] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(321), 1,
      anon_sym_COLON,
  [2068] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(323), 1,
      sym_name,
  [2075] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(258), 1,
      anon_sym_on,
  [2082] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(325), 1,
      sym_name,
  [2089] = 2,
    ACTIONS(311), 1,
      sym_comment,
    ACTIONS(327), 1,
      aux_sym_string_value_token1,
  [2096] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(329), 1,
      anon_sym_DQUOTE_DQUOTE_DQUOTE,
  [2103] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(329), 1,
      anon_sym_DQUOTE,
  [2110] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(331), 1,
      ts_builtin_sym_end,
  [2117] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(333), 1,
      ts_builtin_sym_end,
  [2124] = 2,
    ACTIONS(311), 1,
      sym_comment,
    ACTIONS(335), 1,
      aux_sym_string_value_token1,
  [2131] = 2,
    ACTIONS(311), 1,
      sym_comment,
    ACTIONS(337), 1,
      aux_sym_string_value_token2,
};

static const uint32_t ts_small_parse_table_map[] = {
  [SMALL_STATE(2)] = 0,
  [SMALL_STATE(3)] = 48,
  [SMALL_STATE(4)] = 96,
  [SMALL_STATE(5)] = 144,
  [SMALL_STATE(6)] = 192,
  [SMALL_STATE(7)] = 240,
  [SMALL_STATE(8)] = 284,
  [SMALL_STATE(9)] = 328,
  [SMALL_STATE(10)] = 372,
  [SMALL_STATE(11)] = 407,
  [SMALL_STATE(12)] = 442,
  [SMALL_STATE(13)] = 462,
  [SMALL_STATE(14)] = 482,
  [SMALL_STATE(15)] = 502,
  [SMALL_STATE(16)] = 522,
  [SMALL_STATE(17)] = 542,
  [SMALL_STATE(18)] = 562,
  [SMALL_STATE(19)] = 582,
  [SMALL_STATE(20)] = 602,
  [SMALL_STATE(21)] = 622,
  [SMALL_STATE(22)] = 654,
  [SMALL_STATE(23)] = 684,
  [SMALL_STATE(24)] = 704,
  [SMALL_STATE(25)] = 724,
  [SMALL_STATE(26)] = 744,
  [SMALL_STATE(27)] = 771,
  [SMALL_STATE(28)] = 796,
  [SMALL_STATE(29)] = 821,
  [SMALL_STATE(30)] = 848,
  [SMALL_STATE(31)] = 877,
  [SMALL_STATE(32)] = 902,
  [SMALL_STATE(33)] = 928,
  [SMALL_STATE(34)] = 950,
  [SMALL_STATE(35)] = 964,
  [SMALL_STATE(36)] = 986,
  [SMALL_STATE(37)] = 1000,
  [SMALL_STATE(38)] = 1014,
  [SMALL_STATE(39)] = 1029,
  [SMALL_STATE(40)] = 1050,
  [SMALL_STATE(41)] = 1063,
  [SMALL_STATE(42)] = 1076,
  [SMALL_STATE(43)] = 1097,
  [SMALL_STATE(44)] = 1116,
  [SMALL_STATE(45)] = 1137,
  [SMALL_STATE(46)] = 1157,
  [SMALL_STATE(47)] = 1169,
  [SMALL_STATE(48)] = 1189,
  [SMALL_STATE(49)] = 1201,
  [SMALL_STATE(50)] = 1213,
  [SMALL_STATE(51)] = 1225,
  [SMALL_STATE(52)] = 1237,
  [SMALL_STATE(53)] = 1249,
  [SMALL_STATE(54)] = 1261,
  [SMALL_STATE(55)] = 1273,
  [SMALL_STATE(56)] = 1285,
  [SMALL_STATE(57)] = 1297,
  [SMALL_STATE(58)] = 1309,
  [SMALL_STATE(59)] = 1329,
  [SMALL_STATE(60)] = 1341,
  [SMALL_STATE(61)] = 1353,
  [SMALL_STATE(62)] = 1365,
  [SMALL_STATE(63)] = 1377,
  [SMALL_STATE(64)] = 1397,
  [SMALL_STATE(65)] = 1409,
  [SMALL_STATE(66)] = 1421,
  [SMALL_STATE(67)] = 1433,
  [SMALL_STATE(68)] = 1445,
  [SMALL_STATE(69)] = 1465,
  [SMALL_STATE(70)] = 1485,
  [SMALL_STATE(71)] = 1497,
  [SMALL_STATE(72)] = 1512,
  [SMALL_STATE(73)] = 1529,
  [SMALL_STATE(74)] = 1544,
  [SMALL_STATE(75)] = 1559,
  [SMALL_STATE(76)] = 1576,
  [SMALL_STATE(77)] = 1590,
  [SMALL_STATE(78)] = 1604,
  [SMALL_STATE(79)] = 1618,
  [SMALL_STATE(80)] = 1632,
  [SMALL_STATE(81)] = 1646,
  [SMALL_STATE(82)] = 1656,
  [SMALL_STATE(83)] = 1670,
  [SMALL_STATE(84)] = 1680,
  [SMALL_STATE(85)] = 1694,
  [SMALL_STATE(86)] = 1708,
  [SMALL_STATE(87)] = 1718,
  [SMALL_STATE(88)] = 1727,
  [SMALL_STATE(89)] = 1736,
  [SMALL_STATE(90)] = 1745,
  [SMALL_STATE(91)] = 1754,
  [SMALL_STATE(92)] = 1765,
  [SMALL_STATE(93)] = 1774,
  [SMALL_STATE(94)] = 1785,
  [SMALL_STATE(95)] = 1794,
  [SMALL_STATE(96)] = 1805,
  [SMALL_STATE(97)] = 1814,
  [SMALL_STATE(98)] = 1823,
  [SMALL_STATE(99)] = 1834,
  [SMALL_STATE(100)] = 1843,
  [SMALL_STATE(101)] = 1852,
  [SMALL_STATE(102)] = 1860,
  [SMALL_STATE(103)] = 1870,
  [SMALL_STATE(104)] = 1880,
  [SMALL_STATE(105)] = 1890,
  [SMALL_STATE(106)] = 1898,
  [SMALL_STATE(107)] = 1908,
  [SMALL_STATE(108)] = 1916,
  [SMALL_STATE(109)] = 1926,
  [SMALL_STATE(110)] = 1936,
  [SMALL_STATE(111)] = 1946,
  [SMALL_STATE(112)] = 1954,
  [SMALL_STATE(113)] = 1964,
  [SMALL_STATE(114)] = 1972,
  [SMALL_STATE(115)] = 1980,
  [SMALL_STATE(116)] = 1990,
  [SMALL_STATE(117)] = 1998,
  [SMALL_STATE(118)] = 2005,
  [SMALL_STATE(119)] = 2012,
  [SMALL_STATE(120)] = 2019,
  [SMALL_STATE(121)] = 2026,
  [SMALL_STATE(122)] = 2033,
  [SMALL_STATE(123)] = 2040,
  [SMALL_STATE(124)] = 2047,
  [SMALL_STATE(125)] = 2054,
  [SMALL_STATE(126)] = 2061,
  [SMALL_STATE(127)] = 2068,
  [SMALL_STATE(128)] = 2075,
  [SMALL_STATE(129)] = 2082,
  [SMALL_STATE(130)] = 2089,
  [SMALL_STATE(131)] = 2096,
  [SMALL_STATE(132)] = 2103,
  [SMALL_STATE(133)] = 2110,
  [SMALL_STATE(134)] = 2117,
  [SMALL_STATE(135)] = 2124,
  [SMALL_STATE(136)] = 2131,
};

static const TSParseActionEntry ts_parse_actions[] = {
  [0] = {.entry = {.count = 0, .reusable = false}},
  [1] = {.entry = {.count = 1, .reusable = false}}, RECOVER(),
  [3] = {.entry = {.count = 1, .reusable = true}}, SHIFT_EXTRA(),
  [5] = {.entry = {.count = 1, .reusable = true}}, SHIFT(86),
  [7] = {.entry = {.count = 1, .reusable = true}}, SHIFT(35),
  [9] = {.entry = {.count = 1, .reusable = true}}, SHIFT(112),
  [11] = {.entry = {.count = 1, .reusable = true}}, SHIFT(79),
  [13] = {.entry = {.count = 1, .reusable = true}}, SHIFT(127),
  [15] = {.entry = {.count = 1, .reusable = true}}, SHIFT(135),
  [17] = {.entry = {.count = 1, .reusable = false}}, SHIFT(136),
  [19] = {.entry = {.count = 1, .reusable = false}}, SHIFT(12),
  [21] = {.entry = {.count = 1, .reusable = true}}, SHIFT(12),
  [23] = {.entry = {.count = 1, .reusable = false}}, SHIFT(20),
  [25] = {.entry = {.count = 1, .reusable = true}}, SHIFT(2),
  [27] = {.entry = {.count = 1, .reusable = true}}, SHIFT(17),
  [29] = {.entry = {.count = 1, .reusable = false}}, SHIFT(19),
  [31] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(79),
  [34] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(127),
  [37] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(135),
  [40] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(136),
  [43] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(12),
  [46] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(12),
  [49] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(20),
  [52] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(2),
  [55] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_list_value_repeat1, 2),
  [57] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_list_value_repeat1, 2), SHIFT_REPEAT(19),
  [60] = {.entry = {.count = 1, .reusable = true}}, SHIFT(13),
  [62] = {.entry = {.count = 1, .reusable = true}}, SHIFT(67),
  [64] = {.entry = {.count = 1, .reusable = true}}, SHIFT(50),
  [66] = {.entry = {.count = 1, .reusable = true}}, SHIFT(80),
  [68] = {.entry = {.count = 1, .reusable = true}}, SHIFT(118),
  [70] = {.entry = {.count = 1, .reusable = true}}, SHIFT(130),
  [72] = {.entry = {.count = 1, .reusable = false}}, SHIFT(120),
  [74] = {.entry = {.count = 1, .reusable = false}}, SHIFT(56),
  [76] = {.entry = {.count = 1, .reusable = true}}, SHIFT(56),
  [78] = {.entry = {.count = 1, .reusable = false}}, SHIFT(55),
  [80] = {.entry = {.count = 1, .reusable = true}}, SHIFT(5),
  [82] = {.entry = {.count = 1, .reusable = false}}, SHIFT(52),
  [84] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_document_repeat1, 2),
  [86] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_document_repeat1, 2), SHIFT_REPEAT(86),
  [89] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_document_repeat1, 2), SHIFT_REPEAT(35),
  [92] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_document_repeat1, 2), SHIFT_REPEAT(112),
  [95] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_document, 1),
  [97] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_value, 1),
  [99] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_value, 1),
  [101] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_list_value, 3),
  [103] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_list_value, 3),
  [105] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_string_value, 3),
  [107] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_string_value, 3),
  [109] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_variable, 2),
  [111] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_variable, 2),
  [113] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_object_value, 3),
  [115] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_object_value, 3),
  [117] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_list_value, 2),
  [119] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_list_value, 2),
  [121] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_object_value, 2),
  [123] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_object_value, 2),
  [125] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_enum_value, 1),
  [127] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_enum_value, 1),
  [129] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_boolean_value, 1),
  [131] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_boolean_value, 1),
  [133] = {.entry = {.count = 1, .reusable = true}}, SHIFT(33),
  [135] = {.entry = {.count = 1, .reusable = false}}, SHIFT(110),
  [137] = {.entry = {.count = 1, .reusable = true}}, SHIFT(125),
  [139] = {.entry = {.count = 1, .reusable = false}}, SHIFT(83),
  [141] = {.entry = {.count = 1, .reusable = true}}, SHIFT(91),
  [143] = {.entry = {.count = 1, .reusable = true}}, SHIFT(119),
  [145] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_field, 1),
  [147] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_directive, 2),
  [149] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_directives_repeat1, 2),
  [151] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_directives_repeat1, 2), SHIFT_REPEAT(125),
  [154] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_directives, 1),
  [156] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_field, 2),
  [158] = {.entry = {.count = 1, .reusable = true}}, SHIFT(60),
  [160] = {.entry = {.count = 1, .reusable = true}}, SHIFT(21),
  [162] = {.entry = {.count = 1, .reusable = true}}, SHIFT(22),
  [164] = {.entry = {.count = 1, .reusable = true}}, SHIFT(92),
  [166] = {.entry = {.count = 1, .reusable = true}}, SHIFT(8),
  [168] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_variable_definition, 3),
  [170] = {.entry = {.count = 1, .reusable = true}}, SHIFT(113),
  [172] = {.entry = {.count = 1, .reusable = true}}, SHIFT(77),
  [174] = {.entry = {.count = 1, .reusable = true}}, SHIFT(32),
  [176] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_selection_set_repeat1, 2),
  [178] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_selection_set_repeat1, 2), SHIFT_REPEAT(21),
  [181] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_selection_set_repeat1, 2), SHIFT_REPEAT(22),
  [184] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_arguments, 3),
  [186] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_named_type, 1),
  [188] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_directive, 3),
  [190] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_type, 1),
  [192] = {.entry = {.count = 1, .reusable = true}}, SHIFT(48),
  [194] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_list_type, 3),
  [196] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_field, 3),
  [198] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_fragment_spread, 2),
  [200] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_variable_definition, 4),
  [202] = {.entry = {.count = 1, .reusable = true}}, SHIFT(107),
  [204] = {.entry = {.count = 1, .reusable = true}}, SHIFT(68),
  [206] = {.entry = {.count = 1, .reusable = true}}, SHIFT(36),
  [208] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_operation_definition, 1),
  [210] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_non_null_type, 2),
  [212] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_definition, 1),
  [214] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_executable_definition, 1),
  [216] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_operation_definition, 5),
  [218] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_fragment_definition, 4),
  [220] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_selection_set, 3),
  [222] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_operation_definition, 3),
  [224] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_operation_definition, 2),
  [226] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_operation_definition, 4),
  [228] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_fragment_definition, 5),
  [230] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_field, 4),
  [232] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_variable_definitions_repeat1, 2),
  [234] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_variable_definitions_repeat1, 2), SHIFT_REPEAT(118),
  [237] = {.entry = {.count = 1, .reusable = true}}, SHIFT(105),
  [239] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_arguments_repeat1, 2),
  [241] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_arguments_repeat1, 2), SHIFT_REPEAT(126),
  [244] = {.entry = {.count = 1, .reusable = true}}, SHIFT(53),
  [246] = {.entry = {.count = 1, .reusable = true}}, SHIFT(123),
  [248] = {.entry = {.count = 1, .reusable = true}}, SHIFT(18),
  [250] = {.entry = {.count = 1, .reusable = true}}, SHIFT(62),
  [252] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_default_value, 2),
  [254] = {.entry = {.count = 1, .reusable = true}}, SHIFT(34),
  [256] = {.entry = {.count = 1, .reusable = true}}, SHIFT(126),
  [258] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_fragment_name, 1),
  [260] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_object_value_repeat1, 2),
  [262] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_object_value_repeat1, 2), SHIFT_REPEAT(123),
  [265] = {.entry = {.count = 1, .reusable = true}}, SHIFT(16),
  [267] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_operation_type, 1),
  [269] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_inline_fragment, 3),
  [271] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_selection, 1),
  [273] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_inline_fragment, 2),
  [275] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_object_field, 3),
  [277] = {.entry = {.count = 1, .reusable = true}}, SHIFT(114),
  [279] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_field, 5),
  [281] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_fragment_spread, 3),
  [283] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_variable_definition, 5),
  [285] = {.entry = {.count = 1, .reusable = true}}, SHIFT(101),
  [287] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_inline_fragment, 4),
  [289] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_variable_definition, 6),
  [291] = {.entry = {.count = 1, .reusable = true}}, SHIFT(110),
  [293] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_variable_definitions, 3),
  [295] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_type_condition, 2),
  [297] = {.entry = {.count = 1, .reusable = true}}, SHIFT(128),
  [299] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_object_field, 4),
  [301] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_argument, 3),
  [303] = {.entry = {.count = 1, .reusable = true}}, SHIFT(40),
  [305] = {.entry = {.count = 1, .reusable = true}}, SHIFT(41),
  [307] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_alias, 2),
  [309] = {.entry = {.count = 1, .reusable = false}}, SHIFT(122),
  [311] = {.entry = {.count = 1, .reusable = false}}, SHIFT_EXTRA(),
  [313] = {.entry = {.count = 1, .reusable = true}}, SHIFT(45),
  [315] = {.entry = {.count = 1, .reusable = true}}, SHIFT(51),
  [317] = {.entry = {.count = 1, .reusable = true}}, SHIFT(7),
  [319] = {.entry = {.count = 1, .reusable = true}}, SHIFT(23),
  [321] = {.entry = {.count = 1, .reusable = true}}, SHIFT(9),
  [323] = {.entry = {.count = 1, .reusable = true}}, SHIFT(15),
  [325] = {.entry = {.count = 1, .reusable = true}}, SHIFT(26),
  [327] = {.entry = {.count = 1, .reusable = true}}, SHIFT(124),
  [329] = {.entry = {.count = 1, .reusable = true}}, SHIFT(14),
  [331] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_source_file, 1),
  [333] = {.entry = {.count = 1, .reusable = true}},  ACCEPT_INPUT(),
  [335] = {.entry = {.count = 1, .reusable = true}}, SHIFT(131),
  [337] = {.entry = {.count = 1, .reusable = false}}, SHIFT(132),
};

#ifdef __cplusplus
extern "C" {
#endif
#ifdef _WIN32
#define extern __declspec(dllexport)
#endif

extern const TSLanguage *tree_sitter_gqlt(void) {
  static const TSLanguage language = {
    .version = LANGUAGE_VERSION,
    .symbol_count = SYMBOL_COUNT,
    .alias_count = ALIAS_COUNT,
    .token_count = TOKEN_COUNT,
    .external_token_count = EXTERNAL_TOKEN_COUNT,
    .state_count = STATE_COUNT,
    .large_state_count = LARGE_STATE_COUNT,
    .production_id_count = PRODUCTION_ID_COUNT,
    .field_count = FIELD_COUNT,
    .max_alias_sequence_length = MAX_ALIAS_SEQUENCE_LENGTH,
    .parse_table = &ts_parse_table[0][0],
    .small_parse_table = ts_small_parse_table,
    .small_parse_table_map = ts_small_parse_table_map,
    .parse_actions = ts_parse_actions,
    .symbol_names = ts_symbol_names,
    .symbol_metadata = ts_symbol_metadata,
    .public_symbol_map = ts_symbol_map,
    .alias_map = ts_non_terminal_alias_map,
    .alias_sequences = &ts_alias_sequences[0][0],
    .lex_modes = ts_lex_modes,
    .lex_fn = ts_lex,
    .primary_state_ids = ts_primary_state_ids,
  };
  return &language;
}
#ifdef __cplusplus
}
#endif
