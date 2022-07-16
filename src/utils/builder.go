package utils

import (
	"strings"
	"fmt"
)

type StringBuilder struct {
	builder		*strings.Builder
}

func NewStringBuilder() *StringBuilder {
	var b strings.Builder
	return &StringBuilder{&b}
}

func (b *StringBuilder) Append(s string) {
	sb := b.builder
	fmt.Fprintf(sb, s)
}

func (b *StringBuilder) ToString() string {
	return b.builder.String()
}
