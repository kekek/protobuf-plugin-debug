package messager

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Package struct {
	Data []byte
	Meta reflect.Type
}

type Person struct {
}

func (p *Person) Encode() (*Package, error) {

	prType := reflect.TypeOf(p)

	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return &Package{
		Data: data,
		Meta: prType,
	}, nil
}

func (p *Person) Decode(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *Person) Title() string {
	return "person"
}

type Animal struct{}

func (a *Animal) Encode() (*Package, error) {
	prType := reflect.TypeOf(a)

	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return &Package{
		Data: data,
		Meta: prType,
	}, nil
}

func (a *Animal) Decode(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *Animal) Title() string {
	return "animal"
}

func main() {
	list := []*Package{}

	a := Animal{}
	dataA, err := a.Encode()
	if err != nil {
		panic(err)
	}

	p := Person{}

	dataP, err := p.Encode()
	if err != nil {
		panic(err)
	}

	list = append(list, dataA, dataP)

	for k, v := range list {

		fmt.Println(k, v.Meta.PkgPath(), v.Data)
		// 补全代码，反解析出data

		// Use reflection to create a new instance of the correct type
		instance := reflect.New(v.Meta.Elem()).Interface().(Messager)

		err := instance.Decode(v.Data)
		if err != nil {
			panic(err)
		}
	}
}
