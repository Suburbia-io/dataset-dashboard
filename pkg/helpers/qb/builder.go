package qb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Builder struct {
	err          error
	mask         []bool
	str          []string
	placeholders []int
	args         []interface{}
}

var placeholderRE = regexp.MustCompile(`\$[1-9]\d*`)

func (b *Builder) Write(s string, args ...interface{}) {
	// if an error is already present do nothing
	if b.err != nil {
		return
	}

	// parse the query and abort if error
	builder, err := parse(s, args)
	if err != nil {
		b.err = err
		return
	}

	// join this builder with the parsed builder
	*b = Join(" ", *b, builder)
}

func (b Builder) Build() (s string, args []interface{}, err error) {
	// if an error is present, return it
	if b.err != nil {
		return "", nil, err
	}

	stringBuilder := strings.Builder{}

	i, j := 0, 0
	for _, yes := range b.mask {
		if yes {
			stringBuilder.WriteString("$" + strconv.Itoa(b.placeholders[i]))
			i += 1
		} else {
			stringBuilder.WriteString(b.str[j])
			j += 1
		}
	}

	return stringBuilder.String(), b.args, nil
}

func (b Builder) MustBuild() (s string, args []interface{}) {
	c, d, e := b.Build()
	if e != nil {
		panic(e)
	}
	return c, d
}

func (b *Builder) addPartial(s string) {
	b.mask = append(b.mask, false)
	b.str = append(b.str, s)
}

func (b *Builder) addPlaceholder(n int) {
	b.mask = append(b.mask, true)
	b.placeholders = append(b.placeholders, n)
}

func parse(s string, args []interface{}) (b Builder, err error) {
	phs := placeholderRE.FindAllStringIndex(s, -1)

	if len(phs) == 0 {
		b.addPartial(s)
		return b, nil
	}

	placeholderSet := map[int]bool{}

	start := 0
	n := 0
	for _, ph := range phs {
		// unless this placeholder was found at
		// the start of the string there must
		// be args partial prefix
		if start != ph[0] {
			b.addPartial(s[start:ph[0]])
		}

		// get the placeholder number as int
		n, err = strconv.Atoi(s[ph[0]+1 : ph[1]])
		if err != nil {
			return b, err
		}

		if n == 0 || n-1 >= len(args) {
			return b, fmt.Errorf("$%d not found", n)
		}

		placeholderSet[n] = true
		b.addPlaceholder(n)

		start = ph[1]
	}

	if len(placeholderSet) != len(args) {
		return b, fmt.Errorf("expected %d args saw %d", len(placeholderSet), len(args))
	}
	b.args = args

	if start != len(s) {
		b.addPartial(s[start:])
	}

	return b, nil
}

// -----------------------------------------------------------------------------
// Join
// -----------------------------------------------------------------------------
func Join(sep string, bs ...Builder) Builder {
	joined := Builder{}

	if len(bs) > 1 {
		for i := range bs[:len(bs)-1] {
			bs[i].addPartial(sep)
		}
	}

	var mLen, sLen, nLen, aLen int
	for _, b := range bs {
		mLen += len(b.mask)
		sLen += len(b.str)
		nLen += len(b.placeholders)
		aLen += len(b.args)
	}

	joined.mask = make([]bool, mLen)
	joined.str = make([]string, sLen)
	joined.placeholders = make([]int, nLen)
	joined.args = make([]interface{}, aLen)
	var i, j, k, l int
	for _, p := range bs {
		i += copy(joined.mask[i:], p.mask)
		j += copy(joined.str[j:], p.str)

		_k := copy(joined.placeholders[k:], p.placeholders)
		for x := range joined.str[k : k+_k] {
			joined.placeholders[k+x] = joined.placeholders[k+x] + l
		}
		k += _k

		l += copy(joined.args[l:], p.args)
	}

	return joined
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------
func New(s string, args ...interface{}) Builder {
	b := Builder{}
	b.Write(s, args...)
	return b
}

func WhereAnd(qbs ...Builder) Builder {
	if len(qbs) == 0 {
		return New("WHERE (TRUE)")
	}

	return Join(" ", New("WHERE"), And(qbs...))
}

func And(qbs ...Builder) Builder {
	return Paren(Join(" AND ", qbs...))
}

func Or(qbs ...Builder) Builder {
	return Paren(Join(" OR ", qbs...))
}

func Paren(b Builder) Builder {
	return Join("", New("("), b, New(")"))
}
