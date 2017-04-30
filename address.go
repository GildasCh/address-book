package main

import yaml "gopkg.in/yaml.v2"

type Addresses map[string]([]*Address)

type Address struct {
	Name    string
	Town    string
	Address string
	Comment string
	Type    string
}

func ParseAddresses(input []byte) (as Addresses, err error) {
	var raw []*Address
	err = yaml.Unmarshal(input, &raw)
	if err != nil {
		return
	}

	as = make(map[string][]*Address)
	for _, r := range raw {
		as[r.Town] = append(as[r.Town], r)
	}

	return
}
