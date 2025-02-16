package tests

import (
	"strings"
	"testing"

	"github.com/VirusTotal/gyp"
	"github.com/VirusTotal/gyp/ast"
	"github.com/stretchr/testify/assert"
)

var testRules = `
rule BASIC_BOOL {
  condition:
    true
}

rule BASIC_BOOL2 {
  condition:
    false
}

rule OPERATORS1 {
  condition:
    (1 + 1 - 1) * 2 == 4 % 2
}

rule OPERATORS2 {
  condition:
    -1 != 0
}

rule HEX_STRING1 {
  strings:
    $h1 = { 01 23 45 67 89 AB }
    $h2 = { CD EF 01 23 45 67 }
  condition:
    any of ($h*)
}

rule HEX_STRING2 {
  strings:
    $h1 = { 01 23 ( 45 67 | 89 AB | CD ) ?? ?A ?B }
    $h2 = { CD EF 01 [10-20] 23 45 [-] 67 }
    $h3 = { CD EF 01 [10-20] 23 45 [30-] 67 }
    $h4 = { CD ?? 01 [5] 23 }
    $h5 = { 01 23 ( 45 [30-35] 67 | 89 [40] AB [50-60] CD ) ?? ?A ?B }
  condition:
    any of ($h*)
}

rule HEX_STRING3 {
  strings:
    $a = { AA BB ~CC DD }
    $b = { AA BB ~?C DD }
  condition:
    any of them
}

rule REGEX1 {
  strings:
    $r1 = /first regex/
  condition:
    $r1
}

rule REGEX2 {
  strings:
    $r1 = /regex with mod i/i
    $r2 = /regex with mod s/s
  condition:
    $r1 or $r2
}

rule STRING1 {
  strings:
    $s1 = "ABCDEFG"
  condition:
    $s1
}

rule STRING2 {
  strings:
    $s1 = "ABCDEFG"
    $s2 = "HIJKLMN"
  condition:
    $s1 or $s2
}

rule STRING_MODIFIERS {
  strings:
    $s1 = "foo" ascii wide nocase fullword private base64 base64wide xor
    $s2 = "bar" xor(1)
    $s3 = "baz" xor(2-4)
    $s4 = "qux" base64("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
    $s5 = "qux" base64wide("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
  condition:
    all of them
}

rule STRING_ESCAPED_CHARS {
  strings:
    $s1 = "C:\\Foo\"\\Bar\r\n"
    $s2 = "\""
    $s3 = "\\"
    $s4 = "\t"
    $s5 = "\x00"
    $s6 = "\x7f"
  condition:
    all of them
}

rule TAG : tag1 {
  condition:
    true
}

rule TAG_STRING : tag2 {
  strings:
    $s1 = "ABCDEFG"
  condition:
    $s1
}

rule TAGS : tag1 tag2 tag3 {
  condition:
    true
}

global rule GLOBAL {
  condition:
    true
}

private rule PRIVATE {
  condition:
    true
}

rule META {
  meta:
    meta_str = "abcñ\n\t\x01☺"
    meta_int = 42
    meta_neg = -42
    meta_true = true
    meta_false = false
  condition:
    true
}

rule XOR {
  strings:
    $xor1 = "xor!" xor
    $xor2 = "xor?" nocase xor
    $no_xor1 = "no xor :(" wide
    $no_xor2 = "no xor >:(" ascii nocase
  condition:
    any of them
}

rule OCCURRENCES {
  strings:
    $a = "str1"
    $b = "str2"
    $c = "str3"
  condition:
    #a == 20 and #b < 5 and #c >= 30
}

rule FOR_IN1 {
  strings:
    $a = "str1"
    $b = "str2"
    $c = "str3"
  condition:
    for any i in (5, 10, 15) : (@a[i] % 6 == @c[i * 2])
}

rule FOR_IN2 {
  strings:
    $a = "str1"
    $b = "str2"
    $c = "str3"
  condition:
    for any k,v in some_dict : (k == "foo" and v == "bar")
}

rule FOR_OF {
  meta:
    description = "for..of rule"
  strings:
    $a = "str"
    $b = /regex/
    $c = { 00 11 22 }
  condition:
    for all of ($a, $b, $c) : ($ at entrypoint)
}

rule INTEGER_FUNCTION {
  condition:
    uint8(500) == 3470 and uint16(uint32(100)) == 275
}

rule DEFINED1 {
  condition:
    defined some_value
}

rule DEFINED2 {
  condition:
    not defined some_value
}

rule DEFINED3 {
  condition:
    defined not some_value
}

rule MATCHES {
  condition:
    some_string matches /[a-z0-9]*/i
}

rule CONTAINS {
  condition:
    some_string contains "this string"
}

rule ICONTAINS {
  condition:
    some_string icontains "this string"
}

rule IEQUALS {
  condition:
    some_string iequals "this string"
}

rule STARTSWITH {
  condition:
    some_string startswith "this string"
}

rule ISTARTSWITH {
  condition:
    some_string istartswith "this string"
}

rule ENDSWITH {
  condition:
    some_string endswith "this string"
}

rule IENDSWITH {
  condition:
    some_string iendswith "this string"
}

rule NOT {
  condition:
    not that_var and this_var < 500
}

rule PRECEDENCE_NO_PARENS {
  condition:
    "foo" | "bar" >> 5
}

rule PRECEDENCE_PARENS {
  condition:
    ("foo" | "bar") >> 5
}

rule RANGE {
  strings:
    $a = "str1"
    $b = "str2"
  condition:
    $a in (0..100) and $b in (100..filesize)
}

rule SET_OF_STRINGS {
  strings:
    $foo1 = "foo1"
    $foo2 = "foo2"
    $foo3 = "foo3"
    $foo4 = "foo4"
  condition:
    2 of ($foo1, $foo2, $foo4*)
}

rule AND_OR_PRECEDENCE_NO_PARENS {
  strings:
    $foo1 = "foo1"
    $foo2 = /foo2/
    $foo3 = { AA BB CC }
    $foo4 = "I AM A STRING! ;)"
  condition:
    $foo1 or $foo2 or $foo3 and $foo4
}

rule AND_OR_PRECEDENCE_PARENS {
  strings:
    $foo1 = "foo1"
    $foo2 = /foo2/
    $foo3 = { AA BB CC }
    $foo4 = "I AM A STRING! ;)"
  condition:
    ($foo1 or $foo2 or $foo3) and $foo4
}

rule STRING_LENGTH {
  strings:
    $foo1 = /foo(1)+/
  condition:
    for all i in (5, 10, 15) : (!foo1[i] >= 20)
}

rule MODULE {
  condition:
    foo.bar(1, 10, 3.141592, 4.000000, "foo") != /(test){1}/
}

rule KEYWORD_NONE {
  strings:
    $s = "AXSERS"
  condition:
    none of them
}

rule PERCENT_OF {
  strings:
    $s0 = "AXSERS"
    $s1 = "WXSMTS"
  condition:
    50% of them
}

rule PERCENT_OF_IDENTIFIER {
  strings:
    $s0 = "AXSERS"
    $s1 = "WXSMTS"
  condition:
    x% of them
}

rule STRING_COUNT_IN_RANGE {
  strings:
    $s0 = "AXSERS"
  condition:
    #s0 in (0..100) == 2
}

rule STRING_SET_IN_RANGE {
  strings:
    $s0 = "AXSERS"
    $s1 = "WXSMTS"
  condition:
    any of ($s*) in (0..100)
}

rule RULE_SET {
  condition:
    any of (KEYWORD_NONE, STRING_*)
}

rule RULE_SET_PERCENT {
  condition:
    50% of (PERCENT*)
}

rule RULE_SET_IDENTIFIER_ISSUE_1631 {
  condition:
    for all a in (0..3) : (a * 0 == 0)
}

rule TEXT_STRING_SETS {
  condition:
    for any s in ("a", "b") : (s == "b")
}

rule RULE_STRING_SET_AT_EXPRESSION {
  strings:
    $s0 = "AXSERS"
    $s1 = "WXSMTS"
  condition:
    any of them at 0
}
`

func TestRulesetParsing(t *testing.T) {
	ruleset, err := gyp.ParseString(testRules)
	assert.NoError(t, err)

	var b strings.Builder
	serializer := gyp.NewSerializer(&b)
	serializer.SetIndent("  ")
	err = serializer.Serialize(ruleset.AsProto())
	assert.NoError(t, err)

	output := b.String()
	assert.Equal(t, testRules, output)
}

func TestParsing(t *testing.T) {
	// Parse rule and build AST.
	ruleset, err := gyp.ParseString(testRules)
	assert.NoError(t, err)
	// Write rules from AST back to text.
	var b strings.Builder
	err = ruleset.WriteSource(&b)
	assert.NoError(t, err)
	// Make sure they are equal to the original sources.
	output := b.String()
	assert.Equal(t, testRules, output)
}

func TestProtoSerialization(t *testing.T) {
	// Parse rule and build AST.
	ruleset, err := gyp.ParseString(testRules)
	assert.NoError(t, err)
	// Convert AST to proto.
	pbRuleset := ruleset.AsProto()
	assert.NotNil(t, pbRuleset)
	// Convert the proto back to AST.
	ruleset = ast.RuleSetFromProto(pbRuleset)
	assert.NotNil(t, ruleset)
	// Recover rules sources from AST.
	var b strings.Builder
	err = ruleset.WriteSource(&b)
	assert.NoError(t, err)
	// Make sure they are equal to the original sources.
	assert.Equal(t, testRules, b.String())
}

func TestDuplicateRuleIdentifiers(t *testing.T) {
	_, err := gyp.ParseString(`
	rule DUPLICATE_RULE_IDENTIFIER { condition: true }
	rule DUPLICATE_RULE_IDENTIFIER { condition: true }`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 3: duplicate rule "DUPLICATE_RULE_IDENTIFIER"`, err.Error())
	}
}

func TestDuplicateStringIdentifiers(t *testing.T) {
	_, err := gyp.ParseString(`
	rule DUPLICATE_STRING_IDENTIFIERS {
		strings:
			$s0 = "AXSERS"
			$s0 = "WXSMTS"
		condition:
			all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 5: rule "DUPLICATE_STRING_IDENTIFIERS": duplicate string identifier "s0"`, err.Error())
	}
}

// Make sure anonymous strings are not caught in duplicate checks.
func TestDuplicateStringAnonymous(t *testing.T) {
	_, err := gyp.ParseString(`
	rule DUPLICATE_STRING_IDENTIFIERS {
		strings:
			$ = "AXSERS"
			$ = "WXSMTS"
		condition:
			all of them
	}`)
	assert.NoError(t, err)
}

func TestDuplicateStringModifiers(t *testing.T) {
	_, err := gyp.ParseString(`
	rule DUPLICATE_STRING_MODIFIERS {
		strings:
			$s0 = "AXSERS" xor xor
		condition:
			all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 5: duplicate modifier", err.Error())
	}
}

func TestUndefinedStringIdentifier(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_IDENTIFIER {
		strings:
			$a = "AXSERS"
		condition:
			$s
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 7: undefined string identifier: $s`, err.Error())
	}
}

func TestUndefinedStringKeywordNoStrings(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_KEYWORD_NO_STRINGS {
		condition:
			1 of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: undefined string identifier: them`, err.Error())
	}
}

func TestUndefinedStringIdentifierWildcard(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_IDENTIFIER_WITH_WILDCARD {
		condition:
			1 of ($s*)
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: undefined string identifier: $s*`, err.Error())
	}
}

// Make sure wildcard expansion still works.
func TestUndefinedStringIdentifierWildcardExpansion(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_IDENTIFIER_WITH_WILDCARD_2 {
		strings:
			$s0 = "AXSERS"
		condition:
			1 of ($s*) and $x
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 7: undefined string identifier: $x`, err.Error())
	}
}

func TestUndefinedStringCount(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_COUNT {
		condition:
			#s > 0
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: undefined string identifier: #s`, err.Error())
	}
}

// Make sure anonymous string counts are allowed.
func TestUndefinedStringCountAnonymous(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_COUNT_ANONYMOUS {
		strings:
			$s = "AXSERS"
		condition:
			for any of them: (# > 10)
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringOffset(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_OFFSET {
		condition:
			@s[0] > 0
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: undefined string identifier: @s`, err.Error())
	}
}

// Make sure anonymous string offsets are allowed.
func TestUndefinedStringOffsetAnonymous(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_OFFSET_ANONYMOUS {
		strings:
			$s = "AXSERS"
		condition:
			for any of them: (@ > 10)
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringOffsetAnonymousExpression(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_OFFSET_ANONYMOUS_EXPRESSION {
		strings:
			$s = "AXSERS"
		condition:
			for any of them: (@[1] > 10)
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringLength(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_LENGTH {
		condition:
			!s == 40
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: undefined string identifier: !s`, err.Error())
	}
}

// Make sure anonymous string lengths are allowed.
func TestUndefinedStringLengthAnonymous(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_LENGTH_ANONYMOUS {
		strings:
			$s = "AXSERS"
		condition:
			for any of them: (! > 5)
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringLengthAnonymousExpression(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_LENGTH_ANONYMOUS_EXPRESSION {
		strings:
			$s = "AXSERS"
		condition:
			for any of them: (![1] > 5)
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringAt(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_AT {
		condition:
			$s at 10
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 5: undefined string identifier: $s`, err.Error())
	}
}

// Make sure anonymous string at are allowed.
func TestUndefinedStringAtAnonymous(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_AT_ANONYMOUS {
		strings:
			$s = "AXSERS"
		condition:
			for any of them: ($ at 5)
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringInRange(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_IN_RANGE {
		condition:
			$s in (0..10)
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: undefined string identifier: $s`, err.Error())
	}
}

// Make sure anonymous string in range is allowed.
func TestUndefinedStringInRangeAnonymous(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_IN_RANGE_ANONYMOUS {
		strings:
		  $s = "AXSERS"
		condition:
		    for any of them: ($ in (0..10))
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringInRangeAnonymousExpression(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_IN_RANGE_ANONYMOUS_EXPRESSION {
		strings:
			$s = "AXSERS"
		condition:
			for any of them: ($ in (0..10))
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringCountInRange(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_COUNT_IN_RANGE {
		condition:
			#s in (0..10) == 2
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: undefined string identifier: #s`, err.Error())
	}
}

// Make sure anonymous string count in range is allowed.
func TestUndefinedStringCountInRangeAnonymous(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_COUNT_IN_RANGE_ANONYMOUS {
		strings:
			$s = "AXSERS"
		condition:
			for any of them: (# in (0..10) == 2)
	}`)
	assert.NoError(t, err)
}

func TestUndefinedStringEnumerationAnonymous(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_STRING_ENUMERATION_ANONYMOUS {
		condition:
			any of ($)
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: undefined string identifier: $`, err.Error())
	}
}

func TestUndefinedRuleIdentifier(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_RULE_IDENTIFIER {
		condition:
			any of (NON_EXISTENT_RULE)
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 4: undefined rule identifier: NON_EXISTENT_RULE", err.Error())
	}
}

func TestUndefinedRuleIdentifierWildcard(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNDEFINED_RULE_IDENTIFIER_WILDCARD {
		condition:
			any of (NON_EXISTENT_RULE*)
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 4: undefined rule identifier: NON_EXISTENT_RULE*", err.Error())
	}
}

// Special case where someone does:
//
// rule a { condition: true }
// rule b { condition: any of (a*) }
// rule a2 { condition: true }
//
// This is an error because they have defined a rule (a2) that matches a
// previously defined rule with wildcards.
func TestUndefinedRuleIdentifierWildcardRuleAfter(t *testing.T) {
	_, err := gyp.ParseString(`
	rule a { condition: true }
	rule b { condition: any of (a*) }
	rule a2 { condition : true }
	`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: rule identifier "a2" matches previously used wildcard rule set`, err.Error())
	}
}

func TestBase64AlphabetLength(t *testing.T) {
	_, err := gyp.ParseString(`
	rule BASE64 {
		strings:
			$foo = "foo" base64("baz")
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 4: length of base64 alphabet must be 64", err.Error())
	}
}

func TestNegativeForExpression(t *testing.T) {
	_, err := gyp.ParseString(`
    rule NEGATIVE_FOR_EXPRESSION {
      condition:
        for -1 i in (1) : (i)
    }`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 4: invalid value in condition: -1", err.Error())
	}
}

func TestStringForExpression(t *testing.T) {
	_, err := gyp.ParseString(`
    rule STRING_FOR_EXPRESSION {
      condition:
        for "AXS" i in (1) : (i)
    }`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: invalid value in condition: "AXS"`, err.Error())
	}
}

func TestRegexpForExpression(t *testing.T) {
	_, err := gyp.ParseString(`
    rule REGEXP_FOR_EXPRESSION {
      condition:
        for /AXS/ i in (1) : (i)
    }`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: invalid value in condition: /AXS/`, err.Error())
	}
}

func TestFloatingForExpression(t *testing.T) {
	_, err := gyp.ParseString(`
    rule FLOAT_FOR_EXPRESSION {
      condition:
        for 3.14159 i in (1) : (i)
    }`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: invalid value in condition: 3.141590`, err.Error())
	}
}

func TestNegativeLowerRange(t *testing.T) {
	_, err := gyp.ParseString(`
    rule NEGATIVE_LOWER_RANGE {
	  strings:
	    $a = "AXSERS"
      condition:
        $a in (-1..10)
    }`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 6: lower bound can not be negative", err.Error())
	}
}

func TestNegativeUpperRange(t *testing.T) {
	_, err := gyp.ParseString(`
    rule NEGATIVE_UPPER_RANGE {
	  strings:
	    $a = "AXSERS"
      condition:
        $a in (1..-1)
    }`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 6: upper bound can not be negative", err.Error())
	}
}
func TestInvalidRange(t *testing.T) {
	_, err := gyp.ParseString(`
    rule INVALID_RANGE {
	  strings:
	    $a = "AXSERS"
      condition:
        $a in (10..1)
    }`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 6: lower bound must be less than upper bound", err.Error())
	}
}

// { ~?? } is an error
func TestHexNotWildcard(t *testing.T) {
	_, err := gyp.ParseString(`
  rule TEST {
    strings:
      $ = { AA ~?? CC }
    condition:
      all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 4: invalid not operator (~) in hex string", err.Error())
	}
}

// { AA BB ~ } is an error
func TestHexNotAlone(t *testing.T) {
	_, err := gyp.ParseString(`
  rule TEST {
    strings:
      $ = { AA BB ~ }
    condition:
      all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 4: invalid not operator (~) in hex string", err.Error())
	}
}

func TestUnevenNumberOfDigits(t *testing.T) {
	_, err := gyp.ParseString(`
	rule UNEVEN_HEX_STRING {
		strings:
			$s1 = {012 010203}
		condition:
			all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 4: uneven number of digits in hex string", err.Error())
	}
	_, err = gyp.ParseString(`
	rule UNEVEN_HEX_STRING {
		strings:
			$s1 = {12233}
		condition:
			all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, "line 4: uneven number of digits in hex string", err.Error())
	}
}

func TestBuiltInFuncCall(t *testing.T) {
	ruleset, err := gyp.ParseString(`
	rule RULE_BUILTIN_FUNC_CALL {
		condition:
		  uint32(400)
	}`)
	assert.NoError(t, err)
	assert.Equal(t, true, ruleset.Rules[0].Condition.(*ast.FunctionCall).Builtin)
}

func TestNonBuiltInFuncCall(t *testing.T) {
	ruleset, err := gyp.ParseString(`
	rule RULE_BUILTIN_FUNC_CALL {
		condition:
		  foo(400)
	}`)
	assert.NoError(t, err)
	assert.Equal(t, false, ruleset.Rules[0].Condition.(*ast.FunctionCall).Builtin)
}

func TestInvalidCharacters(t *testing.T) {
	// Non-UTF8 characters are NOT accepted in text strings.
	_, err := gyp.ParseString(`
	rule TEST {
		strings:
			$s1 = "foo` + "\xe8" + `bar"
		condition:
			all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: invalid UTF-8 character "\xe8"`, err.Error())
	}

	// Non-UTF8 characters are NOT accepted in regexps.
	_, err = gyp.ParseString(`
	rule TEST {
		strings:
			$s1 = /foo` + "\xC3" + `bar/
		condition:
			all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: invalid UTF-8 character "\xc3"`, err.Error())
	}

	// Non-ascii characters are NOT accepted in argument to base64 modifier.
	_, err = gyp.ParseString(`
	rule TEST {
		strings:
			$s1 = "foo" base64("é")
		condition:
			all of them
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 4: invalid ASCII character "\xc3"`, err.Error())
	}

	// UTF8 characters are accepted in literal strings and regexps.
	_, err = gyp.ParseString(`
	rule TEST {
        strings:
            $ = "★bebé★"
            $ = /★★★★★/
		condition:
			"abc" != "ñbc"
	}`)
	assert.NoError(t, err)

	// Non-ascii characters are NOT accepted in imports.
	_, err = gyp.ParseString(`
	import "ñoño"
	rule TEST {
		condition:
			false
	}`)
	if assert.Error(t, err) {
		assert.Equal(t, `line 2: invalid ASCII character "\xc3"`, err.Error())
	}

	// Non-ascii characters are accepted in meta strings.
	_, err = gyp.ParseString(`
	rule TEST {
        meta:
			test = "foo` + "\xe8" + `bar"
		condition:
			false
	}`)
	assert.NoError(t, err)

	// Non-ascii and non-UTF8 characters are accepted in comments.
	_, err = gyp.ParseString(`
	rule TEST {
        // This is a comment with a non-ascii é character: ` + "\xe8" + `
		condition:
			false
	}`)
	assert.NoError(t, err)
}

// Test case for https://github.com/VirusTotal/gyp/issues/50
func TestIssue50(t *testing.T) {
	_, err := gyp.ParseString(`
	rule TEST {
   		strings:
       		$a = "foo"
   		condition:
       		for all i in (0..10 : ($a at i))
	}`)
	assert.Error(t, err)
}

// Test case for https://github.com/VirusTotal/gyp/issues/57
func TestIssue57(t *testing.T) {
	_, err := gyp.ParseString(`
  rule TEST {
    strings:
      $ = {
            01 [5]    /* comment */
            02        /* comment */
      }
    condition:
      all of them
	}`)
	assert.NoError(t, err)
}
