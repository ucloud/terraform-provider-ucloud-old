package ucloud

import "fmt"

const EnumUnknownString = "unknown"
const EnumUnknownInt = -1

type intConverter struct {
	c map[int]string
	r map[string]int
}

func newIntConverter(input map[int]string) intConverter {
	reversed := make(map[string]int)
	for k, v := range input {
		reversed[v] = k
	}
	return intConverter{
		c: input,
		r: reversed,
	}
}

func (c intConverter) mustConvert(src int) string {
	v, _ := c.convert(src)
	return v
}

func (c intConverter) mustUnconvert(dst string) int {
	v, _ := c.unconvert(dst)
	return v
}

func (c intConverter) convert(src int) (string, error) {
	if dst, ok := c.c[src]; ok {
		return dst, nil
	}
	return EnumUnknownString, fmt.Errorf("")
}

func (c intConverter) unconvert(dst string) (int, error) {
	if src, ok := c.r[dst]; ok {
		return src, nil
	}
	return EnumUnknownInt, fmt.Errorf("")
}

type boolConverter struct {
	c map[bool]string
	r map[string]bool
}

func newBoolConverter(input map[bool]string) boolConverter {
	reversed := make(map[string]bool)
	for k, v := range input {
		reversed[v] = k
	}
	return boolConverter{
		c: input,
		r: reversed,
	}
}

func (c boolConverter) mustConvert(src bool) string {
	v, _ := c.convert(src)
	return v
}

func (c boolConverter) mustUnconvert(dst string) bool {
	v, _ := c.unconvert(dst)
	return v
}

func (c boolConverter) convert(src bool) (string, error) {
	if dst, ok := c.c[src]; ok {
		return dst, nil
	}
	return EnumUnknownString, fmt.Errorf("")
}

func (c boolConverter) unconvert(dst string) (bool, error) {
	if src, ok := c.r[dst]; ok {
		return src, nil
	}
	return false, fmt.Errorf("")
}

type stringConverter struct {
	c map[string]string
	r map[string]string
}

func newStringConverter(input map[string]string) stringConverter {
	reversed := make(map[string]string)
	for k, v := range input {
		reversed[v] = k
	}
	return stringConverter{
		c: input,
		r: reversed,
	}
}

func (c stringConverter) mustConvert(src string) string {
	if dst, ok := c.c[src]; ok {
		return dst
	}
	return src
}

func (c stringConverter) mustUnconvert(dst string) string {
	if src, ok := c.r[dst]; ok {
		return src
	}
	return dst
}
