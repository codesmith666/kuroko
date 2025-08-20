package object

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
)

/*
 * 整数
 */
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

/*
 * 浮動小数点
 */
type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: uint64(math.Float64bits(f.Value))}
}

/*
 * 複素数
 */
type Complex struct {
	Value complex128
}

func (c *Complex) Type() ObjectType { return COMPLEX_OBJ }
func (c *Complex) Inspect() string {
	return fmt.Sprintf("%g+%gi", real(c.Value), imag(c.Value))
}
func (c *Complex) HashKey() HashKey {
	h := fnv.New64a()
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf[0:8], math.Float64bits(real(c.Value)))
	binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(imag(c.Value)))
	h.Write(buf)
	return HashKey{Type: c.Type(), Value: h.Sum64()}
}
