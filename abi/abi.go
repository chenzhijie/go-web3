package abi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"strings"
	"sync"

	"golang.org/x/crypto/sha3"
)

type ABI struct {
	Constructor *Method
	Methods     map[string]*Method
	Events      map[string]*Event
}

func NewABI(s string) (*ABI, error) {
	return NewABIFromReader(bytes.NewReader([]byte(s)))
}

func MustNewABI(s string) *ABI {
	a, err := NewABI(s)
	if err != nil {
		panic(err)
	}
	return a
}

func NewABIFromReader(r io.Reader) (*ABI, error) {
	var abi *ABI
	dec := json.NewDecoder(r)
	if err := dec.Decode(&abi); err != nil {
		return nil, err
	}
	return abi, nil
}

func (a *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type            string
		Name            string
		Constant        bool
		Anonymous       bool
		StateMutability string
		Inputs          arguments
		Outputs         arguments
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	a.Methods = make(map[string]*Method)
	a.Events = make(map[string]*Event)

	for _, field := range fields {
		switch field.Type {
		case "constructor":
			if a.Constructor != nil {
				return fmt.Errorf("multiple constructor declaration")
			}
			a.Constructor = &Method{
				Inputs: field.Inputs.Type(),
			}

		case "function", "":
			c := field.Constant
			if field.StateMutability == "view" || field.StateMutability == "pure" {
				c = true
			}

			a.Methods[field.Name] = &Method{
				Name:    field.Name,
				Const:   c,
				Inputs:  field.Inputs.Type(),
				Outputs: field.Outputs.Type(),
			}

		case "event":
			a.Events[field.Name] = &Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    field.Inputs.Type(),
			}

		case "fallback":

		default:
			return fmt.Errorf("unknown field type '%s'", field.Type)
		}
	}
	return nil
}

type Method struct {
	Name    string
	Const   bool
	Inputs  *Type
	Outputs *Type
}

func (m *Method) Sig() string {
	return buildSignature(m.Name, m.Inputs)
}

func (m *Method) ID() []byte {
	k := acquireKeccak()
	k.Write([]byte(m.Sig()))
	dst := k.Sum(nil)[:4]
	releaseKeccak(k)
	return dst
}

func (m *Method) EncodeABI(args ...interface{}) ([]byte, error) {
	if len(args) == 0 {
		return m.ID(), nil
	}
	data, err := Encode(args, m.Inputs)
	if err != nil {
		return nil, err
	}
	return append(m.ID(), data...), nil
}

func buildSignature(name string, typ *Type) string {
	types := make([]string, len(typ.tuple))
	for i, input := range typ.tuple {
		types[i] = input.Elem.raw
	}
	return fmt.Sprintf("%v(%v)", name, strings.Join(types, ","))
}

type argument struct {
	Name    string
	Type    *Type
	Indexed bool
}

type arguments []*argument

func (a *arguments) Type() *Type {
	inputs := []*TupleElem{}
	for _, i := range *a {
		inputs = append(inputs, &TupleElem{
			Name:    i.Name,
			Elem:    i.Type,
			Indexed: i.Indexed,
		})
	}

	tt := &Type{
		kind:  KindTuple,
		raw:   "tuple",
		tuple: inputs,
	}
	return tt
}

func (a *argument) UnmarshalJSON(data []byte) error {
	var arg *ArgumentStr
	if err := json.Unmarshal(data, &arg); err != nil {
		return fmt.Errorf("argument json err: %v", err)
	}

	t, err := NewTypeFromArgument(arg)
	if err != nil {
		return err
	}

	a.Type = t
	a.Name = arg.Name
	a.Indexed = arg.Indexed
	return nil
}

type ArgumentStr struct {
	Name       string
	Type       string
	Indexed    bool
	Components []*ArgumentStr
}

var keccakPool = sync.Pool{
	New: func() interface{} {
		return sha3.NewLegacyKeccak256()
	},
}

func acquireKeccak() hash.Hash {
	return keccakPool.Get().(hash.Hash)
}

func releaseKeccak(k hash.Hash) {
	k.Reset()
	keccakPool.Put(k)
}
