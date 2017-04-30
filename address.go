package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Addresses map[string]*Address
type AddrByTown map[string][]*Address

type Address struct {
	Name    string
	Town    string
	Address string
	Comment string
	Type    string
}

var addresses Addresses

func Add(a *Address) error {
	if addresses == nil {
		addresses = make(map[string]*Address)
	}

	key := Key(a)
	if _, ok := addresses[key]; ok {
		return errors.New("This Town/Name combinasion exists already.")
	}

	addresses[key] = a
	return nil
}

func Key(a *Address) string {
	return a.Town + "/" + a.Name
}

func ParseAddresses(input []byte) (as []*Address, err error) {
	err = yaml.Unmarshal(input, &as)

	return
}

func AddFromDB() {

}

func AddFromYaml(input string) error {
	af, err := os.Open(input)
	if err != nil {
		return err
	}
	ab, err := ioutil.ReadAll(af)
	if err != nil {
		return err
	}

	as, err := ParseAddresses(ab)
	if err != nil {
		return err
	}

	for _, a := range as {
		err := Add(a)
		if err != nil {
			fmt.Println(a, "could not be added:", err)
		}
	}

	return nil
}

func ByTown() AddrByTown {
	ret := make(AddrByTown)

	for _, a := range addresses {
		ret[a.Town] = append(ret[a.Town], a)
	}

	return ret
}
