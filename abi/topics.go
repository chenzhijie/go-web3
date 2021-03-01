package abi

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ParseLog(args *Type, log *types.Log) (map[string]interface{}, error) {
	var indexed, nonIndexed []*TupleElem

	for _, arg := range args.TupleElems() {
		if arg.Indexed {
			indexed = append(indexed, arg)
		} else {
			nonIndexed = append(nonIndexed, arg)
		}
	}

	indexedObjs, err := ParseTopics(&Type{kind: KindTuple, tuple: indexed}, log.Topics[1:])
	if err != nil {
		return nil, err
	}

	nonIndexedRaw, err := Decode(&Type{kind: KindTuple, tuple: nonIndexed}, log.Data)
	if err != nil {
		return nil, err
	}
	nonIndexedObjs, ok := nonIndexedRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("bad decoding")
	}

	res := map[string]interface{}{}
	for _, arg := range args.TupleElems() {
		if arg.Indexed {
			res[arg.Name] = indexedObjs[0]
			indexedObjs = indexedObjs[1:]
		} else {
			res[arg.Name] = nonIndexedObjs[arg.Name]
		}
	}

	return res, nil
}

func ParseTopics(args *Type, topics []common.Hash) ([]interface{}, error) {
	if args.kind != KindTuple {
		return nil, fmt.Errorf("expected a tuple type")
	}
	if len(args.TupleElems()) != len(topics) {
		return nil, fmt.Errorf("bad length")
	}

	elems := []interface{}{}
	for indx, arg := range args.TupleElems() {
		elem, err := ParseTopic(arg.Elem, topics[indx])
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}

	return elems, nil
}

func ParseTopic(t *Type, topic common.Hash) (interface{}, error) {
	switch t.kind {
	case KindBool:
		if bytes.Equal(topic[:], topicTrue[:]) {
			return true, nil
		} else if bytes.Equal(topic[:], topicFalse[:]) {
			return false, nil
		}
		return true, fmt.Errorf("is not a boolean")

	case KindInt, KindUInt:
		return readInteger(t, topic[:]), nil

	case KindAddress:
		return readAddr(topic[:])

	default:
		return nil, fmt.Errorf("Topic parsing for type %s not supported", t.String())
	}
}

func EncodeTopic(t *Type, val interface{}) (common.Hash, error) {
	return encodeTopic(t, reflect.ValueOf(val))
}

func encodeTopic(t *Type, val reflect.Value) (common.Hash, error) {
	switch t.kind {
	case KindBool:
		return encodeTopicBool(val)

	case KindUInt, KindInt:
		return encodeTopicNum(t, val)

	case KindAddress:
		return encodeTopicAddress(val)

	}
	return common.Hash{}, fmt.Errorf("not found")
}

var topicTrue, topicFalse common.Hash

func init() {
	topicTrue[31] = 1
}

func encodeTopicAddress(val reflect.Value) (res common.Hash, err error) {
	var b []byte
	b, err = encodeAddress(val)
	if err != nil {
		return
	}
	copy(res[:], b[:])
	return
}

func encodeTopicNum(t *Type, val reflect.Value) (res common.Hash, err error) {
	var b []byte
	b, err = encodeNum(val)
	if err != nil {
		return
	}
	copy(res[:], b[:])
	return
}

func encodeTopicBool(v reflect.Value) (res common.Hash, err error) {
	if v.Kind() != reflect.Bool {
		return common.Hash{}, encodeErr(v, "bool")
	}
	if v.Bool() {
		return topicTrue, nil
	}
	return topicFalse, nil
}

func encodeTopicErr(val reflect.Value, str string) error {
	return fmt.Errorf("cannot encode %s as %s", val.Type().String(), str)
}
