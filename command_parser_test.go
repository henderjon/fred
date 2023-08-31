package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parser(t *testing.T) {
	tests := []struct {
		input      string
		expCommand *command
		expErr     bool
	}{
		// shellAction
		{"!grep -riF \"fatty fatpants\" .", &command{action: shellAction, argument: "grep -riF \"fatty fatpants\" ."}, false},
		// getMarkAction
		{"\"b", &command{action: getMarkAction, argument: "b"}, false},
		{"\"bar", &command{action: getMarkAction, argument: "bar"}, false},
		{"\"b", &command{action: getMarkAction, argument: "b"}, false},
		// searchAction
		{"/re/p", &command{action: searchAction, subCommand: 'p', addrPattern: `re`}, false},
		{"/re/>5p", &command{action: searchAction, subCommand: 'p', addrIncr: ">", addrEnd: "5", addrPattern: `re`}, false},
		{"/re/m35", &command{action: searchAction, subCommand: 'm', addrPattern: `re`, destination: "35"}, false},
		{"//", &command{action: searchAction, addrPattern: ""}, false},
		// putMarkAction
		{"10'b", &command{addrStart: "10", action: putMarkAction, argument: "b"}, false},
		{"10'bar", &command{addrStart: "10", action: putMarkAction, argument: "bar"}, false},
		{"10,25'b", &command{addrStart: "10", addrEnd: "25", action: putMarkAction, argument: "b"}, false},
		{"10>5'b", &command{addrStart: "10", addrEnd: "5", addrIncr: ">", action: putMarkAction, argument: "b"}, false},
		{"10<5'b", &command{addrStart: "10", addrEnd: "5", addrIncr: "<", action: putMarkAction, argument: "b"}, false},
		{".<5'b", &command{addrStart: ".", addrEnd: "5", addrIncr: "<", action: putMarkAction, argument: "b"}, false},
		{">5'b", &command{addrEnd: "5", addrIncr: ">", action: putMarkAction, argument: "b"}, false},
		{">'b", &command{addrIncr: ">", action: putMarkAction, argument: "b"}, false},
		// searchRevAction
		// eqAction
		{"=", &command{action: eqAction}, false},
		// appendAction
		{"25a", &command{addrStart: "25", action: appendAction}, false},
		{"12a", &command{addrStart: "12", action: appendAction}, false},
		{"12,230a", &command{addrStart: "12", addrEnd: "230", action: appendAction}, false},
		{"   12 a  ", &command{addrStart: "12", action: appendAction}, false},
		{"   12,a  ", &command{addrStart: "12", addrEnd: "$", action: appendAction}, false},
		// breakAction
		{">5b", &command{addrEnd: "5", addrIncr: ">", action: breakAction}, false},
		{"1,5b/e./g", &command{addrStart: "1", addrEnd: "5", action: breakAction, pattern: "e.", replaceNum: "-1"}, false},
		{">5b/really/", &command{addrEnd: "5", addrIncr: ">", action: breakAction, pattern: "really"}, false},
		{">5b/really/", &command{addrEnd: "5", addrIncr: ">", action: breakAction, pattern: "really"}, false},
		{">5b/really/", &command{addrEnd: "5", addrIncr: ">", action: breakAction, pattern: "really"}, false},
		{">5b/really/g", &command{addrEnd: "5", addrIncr: ">", action: breakAction, pattern: "really", replaceNum: "-1"}, false},
		// changeAction
		{"10,25c", &command{addrStart: "10", addrEnd: "25", action: changeAction}, false},
		// deleteAction
		{"dp", &command{action: deleteAction, subCommand: 'p'}, false},
		{"10,25d", &command{addrStart: "10", addrEnd: "25", action: deleteAction}, false},
		// editAction
		{"e path/file.ext", &command{action: editAction, argument: "path/file.ext"}, false},
		{"e !grep -riF \"fatty fatpants\" .", &command{action: editAction, subCommand: '!', argument: "grep -riF \"fatty fatpants\" ."}, false},
		{"e path/file.ext", &command{action: editAction, argument: "path/file.ext"}, false},
		// reallyEditAction
		{"E path/file.ext", &command{action: reallyEditAction, argument: "path/file.ext"}, false},
		// filenameAction
		{"f path/file.ext", &command{action: filenameAction, argument: "path/file.ext"}, false},
		// globalIntSearchAction
		// globalReplaceAction
		// globalSearchAction
		// helpAction
		{"h", &command{action: helpAction}, false},
		// insertAction
		{"10,25i", &command{addrStart: "10", addrEnd: "25", action: insertAction}, false},
		{"+12i", &command{addrStart: "+12", action: insertAction}, false},
		// joinAction
		{"10,15j|mm|", &command{addrStart: "10", addrEnd: "15", pattern: "mm", action: joinAction}, false},
		// copyAction
		{"10,25k35", &command{addrStart: "10", addrEnd: "25", action: copyAction, destination: "35"}, false},
		// printLiteralAction
		// mirrorAction
		// moveAction
		{"10,25m35", &command{addrStart: "10", addrEnd: "25", action: moveAction, destination: "35"}, false},
		{"10,25g/mm/m35", &command{addrStart: "10", addrEnd: "25", addrPattern: "mm", globalPrefix: 'g', action: moveAction, destination: "35"}, false},
		{"10,25g|mm|m35", &command{addrStart: "10", addrEnd: "25", addrPattern: "mm", globalPrefix: 'g', action: moveAction, destination: "35"}, false},
		// printNumsAction
		{",n", &command{addrStart: "1", addrEnd: "$", action: printNumsAction}, false},
		{",$n", &command{addrStart: "1", addrEnd: "$", action: printNumsAction}, false},
		{"0,$n", &command{addrStart: "0", addrEnd: "$", action: printNumsAction}, false}, //  buffer should error
		{"0,.n", &command{addrStart: "0", addrEnd: ".", action: printNumsAction}, false}, //  buffer should error
		{".,$n", &command{addrStart: ".", addrEnd: "$", action: printNumsAction}, false},
		{"$n", &command{addrStart: "$", action: printNumsAction}, false}, // parsing the command does not validate the command
		{".n", &command{addrStart: ".", action: printNumsAction}, false}, // parsing the command does not validate the command
		{"10,25n", &command{addrStart: "10", addrEnd: "25", action: printNumsAction}, false},
		// printAction
		{"10,25p", &command{addrStart: "10", addrEnd: "25", action: printAction}, false},
		{"5,8g/^f[ob]ar/p", &command{addrStart: "5", addrEnd: "8", action: printAction, globalPrefix: 'g', addrPattern: `^f[ob]ar`}, false},
		// quitAction
		{"q", &command{action: quitAction}, false},
		// reallyQuitAction
		// readAction
		{"r path/file.ext", &command{action: readAction, argument: "path/file.ext"}, false},
		{"5r path/file.ext", &command{addrStart: "5", action: readAction, argument: "path/file.ext"}, false},
		{"r !grep -riF \"fatty fatpants\" .", &command{action: readAction, subCommand: '!', argument: "grep -riF \"fatty fatpants\" ."}, false},
		// simpleReplaceAction
		{"10,25g/mm/s/and/for/p", &command{addrStart: "10", addrEnd: "25", addrPattern: "mm", globalPrefix: 'g', action: simpleReplaceAction, pattern: "and", substitution: "for", subCommand: 'p'}, false},
		{"10,25g|mm|s!and!for!p", &command{addrStart: "10", addrEnd: "25", addrPattern: "mm", globalPrefix: 'g', action: simpleReplaceAction, pattern: "and", substitution: "for", subCommand: 'p'}, false},
		{"10,25g/mm/s/and/for/g", &command{addrStart: "10", addrEnd: "25", addrPattern: "mm", globalPrefix: 'g', replaceNum: "-1", action: simpleReplaceAction, pattern: "and", substitution: "for"}, false},
		{"10,25g/mm/s/and/for/3", &command{addrStart: "10", addrEnd: "25", addrPattern: "mm", globalPrefix: 'g', replaceNum: "3", action: simpleReplaceAction, pattern: "and", substitution: "for"}, false},
		{"10,25g/mm/s/and/for/", &command{addrStart: "10", addrEnd: "25", addrPattern: "mm", globalPrefix: 'g', action: simpleReplaceAction, pattern: "and", substitution: "for"}, false},
		{"10,25g/mm/s/and//", &command{addrStart: "10", addrEnd: "25", addrPattern: "mm", globalPrefix: 'g', action: simpleReplaceAction, pattern: "and"}, false},
		// regexReplaceAction
		{"2,$S/^.$//", &command{addrStart: "2", addrEnd: "$", action: regexReplaceAction, pattern: "^.$", substitution: ""}, false},
		// transliterateAction
		{"10,15t|foo|bar|", &command{addrStart: "10", addrEnd: "15", pattern: "foo", substitution: "bar", action: transliterateAction}, false},
		// globalNegIntSearchAction
		// globalNegSearchAction
		// superWriteAction
		{"W path/file.ext", &command{action: superWriteAction, argument: "path/file.ext"}, false},
		// writeAction
		{"wq", &command{action: writeAction, subCommand: 'q'}, false},
		{"w path/file.ext", &command{action: writeAction, argument: "path/file.ext"}, false},
		{"4w path/file.ext", &command{addrStart: "4", action: writeAction, argument: "path/file.ext"}, false},
		{"4,12w path/file.ext", &command{addrStart: "4", addrEnd: "12", action: writeAction, argument: "path/file.ext"}, false},
		{"w !grep -riF \"fatty fatpants\" .", &command{action: writeAction, subCommand: '!', argument: "grep -riF \"fatty fatpants\" ."}, false},
		// setPagerAction
		// globalPrefix
		{",g/^f[ob]ar/", &command{addrStart: "1", addrEnd: "$", globalPrefix: 'g', addrPattern: `^f[ob]ar`}, false},
		{"5,g/^f[ob]ar/", &command{addrStart: "5", addrEnd: "$", globalPrefix: 'g', addrPattern: `^f[ob]ar`}, false},
		{"5,8g/^f[ob]ar/", &command{addrStart: "5", addrEnd: "8", globalPrefix: 'g', addrPattern: `^f[ob]ar`}, false},
		{"g/^f[ob]ar/", &command{globalPrefix: 'g', addrPattern: `^f[ob]ar`}, false},
		{"g/b.g/", &command{globalPrefix: 'g', addrPattern: `b.g`}, false},
		{"g//", &command{globalPrefix: 'g', addrPattern: ""}, false}, //itemEmptyPattern
		// testing edge cases //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		{"12", &command{addrStart: "12"}, false},
		{"-12", &command{addrStart: "-12"}, false},
		{"+12", &command{addrStart: "+12"}, false},
		{"   12  ", &command{addrStart: "12"}, false},
		{"   12,  ", &command{addrStart: "12", addrEnd: "$"}, false},
		{"   12,13  ", &command{addrStart: "12", addrEnd: "13"}, false},
		{"   12,13,14,15  ", &command{addrStart: "12", addrEnd: "15"}, false}, // parses correctly but invalid address? ... // NOTE: should we let commands be applied to specific lines?
		{"   12,13,14,11  ", &command{addrStart: "12", addrEnd: "11"}, false}, // parsing the command does not validate the command
		{"  ,12  ", &command{addrStart: "1", addrEnd: "12"}, false},
		{"   12 o  ", nil, true},                                                     // unknown command
		{"g/^f[ob]ar/", &command{globalPrefix: 'g', addrPattern: `^f[ob]ar`}, false}, // missing address
		{"g/b.z", nil, true}, //itemMissingDelim

	} //itemEmpty

	for _, test := range tests {
		c, err := (&parser{}).run(test.input)

		if diff := cmp.Diff(c, test.expCommand, cmp.AllowUnexported(command{})); diff != "" {
			t.Errorf("given: %s; -got/+want\n%s[%s]", test.input, diff, err)
			// t.Errorf("given: %s; -got/+want\n%s", test.input, diff)
		}

		if test.expErr && (err == nil) {
			t.Errorf("given: %s; error: %s\n", test.input, err)
		}

	}
}
