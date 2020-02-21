package tools

import (
	"fmt"
	"github.com/google/uuid"
)



type UUIDGenerator interface {
	Generate() (uuid.UUID, error)
}

func NewUUIDGenerator() UUIDGenerator {
	return UUIDGen{}
}

type UUIDGen struct{}

func (u UUIDGen) Generate() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return [16]byte{}, err
	}
	return id, nil
}



type CodeGenerator interface {
	Generate() string
}

func NewCodeGenerator() CodeGenerator{
	return CodeGen{}
}

type CodeGen struct {}


func (c CodeGen) Generate() string {
	// TODO implement proper code generator
	code, _ := uuid.NewRandom()

	q1 := c.ParseCode(sumQuarter(code[0:4]))
	q2 := c.ParseCode(sumQuarter(code[4:8]))
	q3 := c.ParseCode(sumQuarter(code[8:12]))
	q4 := c.ParseCode(sumQuarter(code[12:16]))

	return fmt.Sprintf("%c%c%c%c", q1, q2, q3, q4)
}

func (c CodeGen) ParseCode(b byte) rune{
	chars := []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z',
	}
	mod := b % byte(len(chars))
	return chars[mod]
}

func sumQuarter(quarter []byte) byte {
	count := byte(0)
	for _, b := range quarter {
		count += b
	}
	return count
}
