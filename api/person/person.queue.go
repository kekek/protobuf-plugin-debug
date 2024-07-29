// source: api/person/person.proto

package person

// 测试文件 person
import (
	"encoding/json"
	"reflect"

	msg "protobuf-plugin-debug/internal/messager"
)

type Person struct {
}

func (a *Person) Encode() (*msg.Package, error) {
	prType := reflect.TypeOf(a)

	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return &msg.Package{
		Data: data,
		Meta: prType,
	}, nil
}

func (a *Person) Decode(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *Person) Title() string {
	return "Person"
}

type Animal struct {
}

func (a *Animal) Encode() (*msg.Package, error) {
	prType := reflect.TypeOf(a)

	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return &msg.Package{
		Data: data,
		Meta: prType,
	}, nil
}

func (a *Animal) Decode(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *Animal) Title() string {
	return "Animal"
}
