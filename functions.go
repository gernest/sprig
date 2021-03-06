/*
Sprig: Template functions for Go.

This package contains a number of utility functions for working with data
inside of Go `html/template` and `text/template` files.

To add these functions, use the `template.Funcs()` method:

	t := templates.New("foo").Funcs(sprig.FuncMap())

Note that you should add the function map before you parse any template files.

Date Functions

	- date: Format a date, where a date is an integer type or a time.Time type, and
	  format is a time.Format formatting string.
	- date_modify: Given a date, modify it with a duration: `date_modify "-1.5h" now`. If the duration doesn't
	parse, it returns the time unaltered. See `time.ParseDuration` for info on duration strings.
	- now: Current time.Time, for feeding into date-related functions.

String Functions

	- trim: strings.TrimSpace
	- upper: strings.ToUpper
	- lower: strings.ToLower
	- title: strings.Title
	- repeat: strings.Repeat, but with the arguments switched: `repeat count str`. (This simplifies common pipelines)

String Slice Functions:

	- join: strings.Join, but as `join SEP SLICE`

Conversions:

	- atoi: Convert a string to an integer. 0 if the integer could not be parsed.

Math Functions:

	- add1: Increment an integer by 1
	- add: Sum two integers
	- sub: Subtract the second integer from the first
	- div: Divide the first integer by the second
	- mod: Module of first integer divided by second
	- mul: Multiply two integers
	- biggest: Return the biggest of two integers

REMOVED (implemented in Go 1.2)

	- gt: Greater than (integer)
	- lt: Less than (integer)
	- gte: Greater than or equal to (integer)
	- lte: Less than or equal to (integer)

*/
package sprig

import (
	"html/template"
	ttemplate "text/template"
	"time"
	"strings"
	"strconv"
)

// Produce the function map.
//
// Use this to pass the functions into the template engine:
//
// 	tpl := template.New("foo").Funcs(sprig.FuncMap))
//
func FuncMap() template.FuncMap {
	return template.FuncMap(genericMap)
}

// TextFuncMap returns a 'text/template'.FuncMap
func TxtFuncMap() ttemplate.FuncMap {
	return ttemplate.FuncMap(genericMap)
}

// HtmlFuncMap returns an 'html/template'.Funcmap
func HtmlFuncMap() template.FuncMap {
	return template.FuncMap(genericMap)
}

var  genericMap = map[string]interface{} {
	"hello": func () string { return "Hello!" },

	// Date functions
	"date": date,
	"date_in_zone": dateInZone,
	"date_modify": dateModify,
	"now": func () time.Time { return time.Now() },

	// Strings
	"trim": strings.TrimSpace,
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"title": strings.Title,
	// Switch or so that "foo" | repeat 5
	"repeat": func (count int, str string) string { return strings.Repeat(str, count) },

	// Wrap Atoi to stop errors.
	"atoi": func (a string) int { i, _ := strconv.Atoi(a); return i },

	//"gt": func(a, b int) bool {return a > b},
	//"gte": func(a, b int) bool {return a >= b},
	//"lt": func(a, b int) bool {return a < b},
	//"lte": func(a, b int) bool {return a <= b},


	// VERY basic arithmetic.
	"add1": func (i int) int {return i + 1},
	"add": func (a, b int) int { return a + b },
	"sub": func (a, b int) int { return a - b },
	"div": func (a, b int) int { return a / b },
	"mod": func (a, b int) int { return a % b },
	"mul": func (a, b int) int { return a * b },
	"biggest": biggest,

	// string slices. Note that we reverse the order b/c that's better
	// for template processing.
	"join": func(sep string, ss []string) string {return strings.Join(ss, sep)},
}

// Given a format and a date, format the date string.
//
// Date can be a `time.Time` or an `int, int32, int64`.
// In the later case, it is treated as seconds since UNIX
// epoch.
func date(fmt string, date interface{}) string {
	return dateInZone(fmt, date, "Local")
}

func dateInZone(fmt string, date interface{}, zone string) string {
	var t time.Time
	switch date := date.(type) {
	default:
		t = time.Now()
	case time.Time:
		t = date
	case int64:
		t = time.Unix(date, 0)
	case int:
		t = time.Unix(int64(date), 0)
	case int32:
		t = time.Unix(int64(date), 0)
	}

	loc, err := time.LoadLocation(zone)
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}

	return t.In(loc).Format(fmt)
}

func dateModify(fmt string, date time.Time) time.Time {
	d, err := time.ParseDuration(fmt)
	if err != nil {
		return date
	}
	return date.Add(d)
}

func biggest(a, b int) int {
	if a > b {
		return a
	}
	return b
}
