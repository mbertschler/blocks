// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"testing"
)

func TestTypedContent(t *testing.T) {
	data := []Block{
		Text(`<b> "foo%" O'Reilly &bar;`),
		CSS(`a[href =~ "//example.com"]#foo`),
		HTML(`Hello, <b>World</b> &amp;tc!`),
		HTMLAttr(` dir="ltr"`),
		JS(`c && alert("Hello, World!");`),
		JSStr(`Hello, World & O'Reilly\x21`),
		URL(`greeting=H%69&addressee=(World)`),
	}

	// For each content sensitive escaper, see how it does on
	// each of the typed strings above.
	tests := []struct {
		// A template containing a single {{.}}.
		input func(in Block) Block
		want  []string
	}{
		{
			//`<style>{{.}} { color: blue }</style>`,
			func(in Block) Block {
				return Style(nil, in, Text(" { color: blue }"))
			},
			[]string{
				`ZgotmplZ`,
				// Allowed but not escaped.
				`a[href =~ "//example.com"]#foo`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
			},
		},
		{
			//`<div style="{{.}}">`,
			func(in Block) Block {
				return Div(Attr{{"style", in}})
			},
			[]string{
				`ZgotmplZ`,
				// Allowed and HTML escaped.
				`a[href =~ &#34;//example.com&#34;]#foo`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
			},
		},
		{
			//`{{.}}`,
			func(in Block) Block {
				return in
			},
			[]string{
				`&lt;b&gt; &#34;foo%&#34; O&#39;Reilly &amp;bar;`,
				`a[href =~ &#34;//example.com&#34;]#foo`,
				// Not escaped.
				`Hello, <b>World</b> &amp;tc!`,
				` dir=&#34;ltr&#34;`,
				`c &amp;&amp; alert(&#34;Hello, World!&#34;);`,
				`Hello, World &amp; O&#39;Reilly\x21`,
				`greeting=H%69&amp;addressee=(World)`,
			},
		},
		// attribute strings are not supported yet
		// {
		// 	`<a{{.}}>`,
		// 	[]string{
		// 		`ZgotmplZ`,
		// 		`ZgotmplZ`,
		// 		`ZgotmplZ`,
		// 		// Allowed and HTML escaped.
		// 		` dir="ltr"`,
		// 		`ZgotmplZ`,
		// 		`ZgotmplZ`,
		// 		`ZgotmplZ`,
		// 	},
		// },
		{
			//`<a title={{.}}>`,
			func(in Block) Block {
				return A(Attr{{"title", in}})
			},
			[]string{
				`&lt;b&gt;&#32;&#34;foo%&#34;&#32;O&#39;Reilly&#32;&amp;bar;`,
				`a[href&#32;&#61;~&#32;&#34;//example.com&#34;]#foo`,
				// Tags stripped, spaces escaped, entity not re-escaped.
				`Hello,&#32;World&#32;&amp;tc!`,
				`&#32;dir&#61;&#34;ltr&#34;`,
				`c&#32;&amp;&amp;&#32;alert(&#34;Hello,&#32;World!&#34;);`,
				`Hello,&#32;World&#32;&amp;&#32;O&#39;Reilly\x21`,
				`greeting&#61;H%69&amp;addressee&#61;(World)`,
			},
		},
		{
			//`<a title='{{.}}'>`,
			func(in Block) Block {
				return A(Attr{{"title", in}})
			},
			[]string{
				`&lt;b&gt; &#34;foo%&#34; O&#39;Reilly &amp;bar;`,
				`a[href =~ &#34;//example.com&#34;]#foo`,
				// Tags stripped, entity not re-escaped.
				`Hello, World &amp;tc!`,
				` dir=&#34;ltr&#34;`,
				`c &amp;&amp; alert(&#34;Hello, World!&#34;);`,
				`Hello, World &amp; O&#39;Reilly\x21`,
				`greeting=H%69&amp;addressee=(World)`,
			},
		},
		{
			//`<textarea>{{.}}</textarea>`,
			func(in Block) Block {
				return Textarea(nil, in)
			},
			[]string{
				`&lt;b&gt; &#34;foo%&#34; O&#39;Reilly &amp;bar;`,
				`a[href =~ &#34;//example.com&#34;]#foo`,
				// Angle brackets escaped to prevent injection of close tags, entity not re-escaped.
				`Hello, &lt;b&gt;World&lt;/b&gt; &amp;tc!`,
				` dir=&#34;ltr&#34;`,
				`c &amp;&amp; alert(&#34;Hello, World!&#34;);`,
				`Hello, World &amp; O&#39;Reilly\x21`,
				`greeting=H%69&amp;addressee=(World)`,
			},
		},
		{
			//`<script>alert({{.}})</script>`,
			func(in Block) Block {
				//return Script(nil, JS("alert(", in, ")"))
				return Script(nil, JS("alert("), in, JS(")"))
			},
			[]string{
				`"\u003cb\u003e \"foo%\" O'Reilly \u0026bar;"`,
				`"a[href =~ \"//example.com\"]#foo"`,
				`"Hello, \u003cb\u003eWorld\u003c/b\u003e \u0026amp;tc!"`,
				`" dir=\"ltr\""`,
				// Not escaped.
				`c && alert("Hello, World!");`,
				// Escape sequence not over-escaped.
				`"Hello, World & O'Reilly\x21"`,
				`"greeting=H%69\u0026addressee=(World)"`,
			},
		},
		{
			//`<button onclick="alert({{.}})">`,
			func(in Block) Block {
				return Button(Attr{{"onclick", Blocks{JS("alert("), in, JS(")")}}})
			},
			[]string{
				`&#34;\u003cb\u003e \&#34;foo%\&#34; O&#39;Reilly \u0026bar;&#34;`,
				`&#34;a[href =~ \&#34;//example.com\&#34;]#foo&#34;`,
				`&#34;Hello, \u003cb\u003eWorld\u003c/b\u003e \u0026amp;tc!&#34;`,
				`&#34; dir=\&#34;ltr\&#34;&#34;`,
				// Not JS escaped but HTML escaped.
				`c &amp;&amp; alert(&#34;Hello, World!&#34;);`,
				// Escape sequence not over-escaped.
				`&#34;Hello, World &amp; O&#39;Reilly\x21&#34;`,
				`&#34;greeting=H%69\u0026addressee=(World)&#34;`,
			},
		},
		{
			//`<script>alert("{{.}}")</script>`,
			func(in Block) Block {
				return Script(nil, JS("alert("), in, JS(")"))
			},
			[]string{
				`\x3cb\x3e \x22foo%\x22 O\x27Reilly \x26bar;`,
				`a[href =~ \x22\/\/example.com\x22]#foo`,
				`Hello, \x3cb\x3eWorld\x3c\/b\x3e \x26amp;tc!`,
				` dir=\x22ltr\x22`,
				`c \x26\x26 alert(\x22Hello, World!\x22);`,
				// Escape sequence not over-escaped.
				`Hello, World \x26 O\x27Reilly\x21`,
				`greeting=H%69\x26addressee=(World)`,
			},
		},
		{
			//`<button onclick='alert("{{.}}")'>`,
			func(in Block) Block {
				return Button(Attr{{"onclick", Blocks{JS("alert("), in, JS(")")}}})
			},
			[]string{
				`\x3cb\x3e \x22foo%\x22 O\x27Reilly \x26bar;`,
				`a[href =~ \x22\/\/example.com\x22]#foo`,
				`Hello, \x3cb\x3eWorld\x3c\/b\x3e \x26amp;tc!`,
				` dir=\x22ltr\x22`,
				`c \x26\x26 alert(\x22Hello, World!\x22);`,
				// Escape sequence not over-escaped.
				`Hello, World \x26 O\x27Reilly\x21`,
				`greeting=H%69\x26addressee=(World)`,
			},
		},
		{
			//`<a href="?q={{.}}">`,
			func(in Block) Block {
				return A(Attr{{"href", Blocks{URL("?q="), in}}})
			},
			[]string{
				`%3cb%3e%20%22foo%25%22%20O%27Reilly%20%26bar%3b`,
				`a%5bhref%20%3d~%20%22%2f%2fexample.com%22%5d%23foo`,
				`Hello%2c%20%3cb%3eWorld%3c%2fb%3e%20%26amp%3btc%21`,
				`%20dir%3d%22ltr%22`,
				`c%20%26%26%20alert%28%22Hello%2c%20World%21%22%29%3b`,
				`Hello%2c%20World%20%26%20O%27Reilly%5cx21`,
				// Quotes and parens are escaped but %69 is not over-escaped. HTML escaping is done.
				`greeting=H%69&amp;addressee=%28World%29`,
			},
		},
		{
			//`<style>body { background: url('?img={{.}}') }</style>`,
			func(in Block) Block {
				return Style(nil, CSS("body { background: url('?img="), in, CSS("') }"))
			},
			[]string{
				`%3cb%3e%20%22foo%25%22%20O%27Reilly%20%26bar%3b`,
				`a%5bhref%20%3d~%20%22%2f%2fexample.com%22%5d%23foo`,
				`Hello%2c%20%3cb%3eWorld%3c%2fb%3e%20%26amp%3btc%21`,
				`%20dir%3d%22ltr%22`,
				`c%20%26%26%20alert%28%22Hello%2c%20World%21%22%29%3b`,
				`Hello%2c%20World%20%26%20O%27Reilly%5cx21`,
				// Quotes and parens are escaped but %69 is not over-escaped. HTML escaping is not done.
				`greeting=H%69&addressee=%28World%29`,
			},
		},
	}

	for _, test := range tests {
		// tmpl := template.Must(template.New("x").Parse(test.input))
		// pre := strings.Index(test.input, "{{.}}")
		// post := len(test.input) - (pre + 5)
		pre, post := 0, 0
		var b bytes.Buffer
		fail := false
		for i, x := range data {
			b.Reset()
			out, err := RenderString(test.input(x))
			if err != nil {
				t.Error(err)
			}
			// if err := tmpl.Execute(&b, x); err != nil {
			// 	t.Errorf("%q with %v: %s", test.input, x, err)
			// 	continue
			// }
			//out := b.String()
			if want, got := test.want[i], out[pre:b.Len()-post]; want != got {
				t.Log("out:", out)
				t.Errorf("with %v:\nwant\n\t%q,\ngot\n\t%q\n", x, want, got)
				fail = true
			}
		}
		if fail {
			break
		}
	}
}

func TestTypedContentTemplate(t *testing.T) {
	data := []interface{}{
		`<b> "foo%" O'Reilly &bar;`,
		template.CSS(`a[href =~ "//example.com"]#foo`),
		template.HTML(`Hello, <b>World</b> &amp;tc!`),
		template.HTMLAttr(` dir="ltr"`),
		template.JS(`c && alert("Hello, World!");`),
		template.JSStr(`Hello, World & O'Reilly\x21`),
		template.URL(`greeting=H%69&addressee=(World)`),
	}

	// For each content sensitive escaper, see how it does on
	// each of the typed strings above.
	tests := []struct {
		// A template containing a single {{.}}.
		input string
		want  []string
	}{
		{
			`<style>{{.}} { color: blue }</style>`,
			[]string{
				`ZgotmplZ`,
				// Allowed but not escaped.
				`a[href =~ "//example.com"]#foo`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
			},
		},
		{
			`<div style="{{.}}">`,
			[]string{
				`ZgotmplZ`,
				// Allowed and HTML escaped.
				`a[href =~ &#34;//example.com&#34;]#foo`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
			},
		},
		{
			`{{.}}`,
			[]string{
				`&lt;b&gt; &#34;foo%&#34; O&#39;Reilly &amp;bar;`,
				`a[href =~ &#34;//example.com&#34;]#foo`,
				// Not escaped.
				`Hello, <b>World</b> &amp;tc!`,
				` dir=&#34;ltr&#34;`,
				`c &amp;&amp; alert(&#34;Hello, World!&#34;);`,
				`Hello, World &amp; O&#39;Reilly\x21`,
				`greeting=H%69&amp;addressee=(World)`,
			},
		},
		{
			`<a{{.}}>`,
			[]string{
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				// Allowed and HTML escaped.
				` dir="ltr"`,
				`ZgotmplZ`,
				`ZgotmplZ`,
				`ZgotmplZ`,
			},
		},
		{
			`<a title={{.}}>`,
			[]string{
				`&lt;b&gt;&#32;&#34;foo%&#34;&#32;O&#39;Reilly&#32;&amp;bar;`,
				`a[href&#32;&#61;~&#32;&#34;//example.com&#34;]#foo`,
				// Tags stripped, spaces escaped, entity not re-escaped.
				`Hello,&#32;World&#32;&amp;tc!`,
				`&#32;dir&#61;&#34;ltr&#34;`,
				`c&#32;&amp;&amp;&#32;alert(&#34;Hello,&#32;World!&#34;);`,
				`Hello,&#32;World&#32;&amp;&#32;O&#39;Reilly\x21`,
				`greeting&#61;H%69&amp;addressee&#61;(World)`,
			},
		},
		{
			`<a title='{{.}}'>`,
			[]string{
				`&lt;b&gt; &#34;foo%&#34; O&#39;Reilly &amp;bar;`,
				`a[href =~ &#34;//example.com&#34;]#foo`,
				// Tags stripped, entity not re-escaped.
				`Hello, World &amp;tc!`,
				` dir=&#34;ltr&#34;`,
				`c &amp;&amp; alert(&#34;Hello, World!&#34;);`,
				`Hello, World &amp; O&#39;Reilly\x21`,
				`greeting=H%69&amp;addressee=(World)`,
			},
		},
		{
			`<textarea>{{.}}</textarea>`,
			[]string{
				`&lt;b&gt; &#34;foo%&#34; O&#39;Reilly &amp;bar;`,
				`a[href =~ &#34;//example.com&#34;]#foo`,
				// Angle brackets escaped to prevent injection of close tags, entity not re-escaped.
				`Hello, &lt;b&gt;World&lt;/b&gt; &amp;tc!`,
				` dir=&#34;ltr&#34;`,
				`c &amp;&amp; alert(&#34;Hello, World!&#34;);`,
				`Hello, World &amp; O&#39;Reilly\x21`,
				`greeting=H%69&amp;addressee=(World)`,
			},
		},
		{
			`<script>alert({{.}})</script>`,
			[]string{
				`"\u003cb\u003e \"foo%\" O'Reilly \u0026bar;"`,
				`"a[href =~ \"//example.com\"]#foo"`,
				`"Hello, \u003cb\u003eWorld\u003c/b\u003e \u0026amp;tc!"`,
				`" dir=\"ltr\""`,
				// Not escaped.
				`c && alert("Hello, World!");`,
				// Escape sequence not over-escaped.
				`"Hello, World & O'Reilly\x21"`,
				`"greeting=H%69\u0026addressee=(World)"`,
			},
		},
		{
			`<button onclick="alert({{.}})">`,
			[]string{
				`&#34;\u003cb\u003e \&#34;foo%\&#34; O&#39;Reilly \u0026bar;&#34;`,
				`&#34;a[href =~ \&#34;//example.com\&#34;]#foo&#34;`,
				`&#34;Hello, \u003cb\u003eWorld\u003c/b\u003e \u0026amp;tc!&#34;`,
				`&#34; dir=\&#34;ltr\&#34;&#34;`,
				// Not JS escaped but HTML escaped.
				`c &amp;&amp; alert(&#34;Hello, World!&#34;);`,
				// Escape sequence not over-escaped.
				`&#34;Hello, World &amp; O&#39;Reilly\x21&#34;`,
				`&#34;greeting=H%69\u0026addressee=(World)&#34;`,
			},
		},
		{
			`<script>alert("{{.}}")</script>`,
			[]string{
				`\x3cb\x3e \x22foo%\x22 O\x27Reilly \x26bar;`,
				`a[href =~ \x22\/\/example.com\x22]#foo`,
				`Hello, \x3cb\x3eWorld\x3c\/b\x3e \x26amp;tc!`,
				` dir=\x22ltr\x22`,
				`c \x26\x26 alert(\x22Hello, World!\x22);`,
				// Escape sequence not over-escaped.
				`Hello, World \x26 O\x27Reilly\x21`,
				`greeting=H%69\x26addressee=(World)`,
			},
		},
		{
			`<button onclick='alert("{{.}}")'>`,
			[]string{
				`\x3cb\x3e \x22foo%\x22 O\x27Reilly \x26bar;`,
				`a[href =~ \x22\/\/example.com\x22]#foo`,
				`Hello, \x3cb\x3eWorld\x3c\/b\x3e \x26amp;tc!`,
				` dir=\x22ltr\x22`,
				`c \x26\x26 alert(\x22Hello, World!\x22);`,
				// Escape sequence not over-escaped.
				`Hello, World \x26 O\x27Reilly\x21`,
				`greeting=H%69\x26addressee=(World)`,
			},
		},
		{
			`<a href="?q={{.}}">`,
			[]string{
				`%3cb%3e%20%22foo%25%22%20O%27Reilly%20%26bar%3b`,
				`a%5bhref%20%3d~%20%22%2f%2fexample.com%22%5d%23foo`,
				`Hello%2c%20%3cb%3eWorld%3c%2fb%3e%20%26amp%3btc%21`,
				`%20dir%3d%22ltr%22`,
				`c%20%26%26%20alert%28%22Hello%2c%20World%21%22%29%3b`,
				`Hello%2c%20World%20%26%20O%27Reilly%5cx21`,
				// Quotes and parens are escaped but %69 is not over-escaped. HTML escaping is done.
				`greeting=H%69&amp;addressee=%28World%29`,
			},
		},
		{
			`<style>body { background: url('?img={{.}}') }</style>`,
			[]string{
				`%3cb%3e%20%22foo%25%22%20O%27Reilly%20%26bar%3b`,
				`a%5bhref%20%3d~%20%22%2f%2fexample.com%22%5d%23foo`,
				`Hello%2c%20%3cb%3eWorld%3c%2fb%3e%20%26amp%3btc%21`,
				`%20dir%3d%22ltr%22`,
				`c%20%26%26%20alert%28%22Hello%2c%20World%21%22%29%3b`,
				`Hello%2c%20World%20%26%20O%27Reilly%5cx21`,
				// Quotes and parens are escaped but %69 is not over-escaped. HTML escaping is not done.
				`greeting=H%69&addressee=%28World%29`,
			},
		},
	}

	for _, test := range tests {
		tmpl := template.Must(template.New("x").Parse(test.input))
		pre := strings.Index(test.input, "{{.}}")
		post := len(test.input) - (pre + 5)
		var b bytes.Buffer
		for i, x := range data {
			b.Reset()
			if err := tmpl.Execute(&b, x); err != nil {
				t.Errorf("%q with %v: %s", test.input, x, err)
				continue
			}
			if want, got := test.want[i], b.String()[pre:b.Len()-post]; want != got {
				t.Errorf("%q with %v:\nwant\n\t%q,\ngot\n\t%q\n", test.input, x, want, got)
				continue
			}
		}
	}
}

// Test that we print using the String method. Was issue 3073.
type stringer struct {
	v int
}

func (s *stringer) String() string {
	return fmt.Sprintf("string=%d", s.v)
}

type errorer struct {
	v int
}

func (s *errorer) Error() string {
	return fmt.Sprintf("error=%d", s.v)
}

func TestStringer(t *testing.T) {
	s := &stringer{3}
	b := new(bytes.Buffer)
	tmpl := template.Must(template.New("x").Parse("{{.}}"))
	if err := tmpl.Execute(b, s); err != nil {
		t.Fatal(err)
	}
	var expect = "string=3"
	if b.String() != expect {
		t.Errorf("expected %q got %q", expect, b.String())
	}
	e := &errorer{7}
	b.Reset()
	if err := tmpl.Execute(b, e); err != nil {
		t.Fatal(err)
	}
	expect = "error=7"
	if b.String() != expect {
		t.Errorf("expected %q got %q", expect, b.String())
	}
}

// https://golang.org/issue/5982
func TestEscapingNilNonemptyInterfaces(t *testing.T) {
	tmpl := template.Must(template.New("x").Parse("{{.E}}"))

	got := new(bytes.Buffer)
	testData := struct{ E error }{} // any non-empty interface here will do; error is just ready at hand
	tmpl.Execute(got, testData)

	// Use this data instead of just hard-coding "&lt;nil&gt;" to avoid
	// dependencies on the html escaper and the behavior of fmt w.r.t. nil.
	want := new(bytes.Buffer)
	data := struct{ E string }{E: fmt.Sprint(nil)}
	tmpl.Execute(want, data)

	if !bytes.Equal(want.Bytes(), got.Bytes()) {
		t.Errorf("expected %q got %q", string(want.Bytes()), string(got.Bytes()))
	}
}
