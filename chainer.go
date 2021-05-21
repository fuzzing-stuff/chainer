package chainer

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"strconv"
	"strings"
)

type (
	//Chain struct which contains all information about test case
	Chain struct {
		ID       string   // new id of test case
		FilePath string   // File for test case or for permutation in case of IsMutate flag
		Content  []string // Content of test case
		IsMutate bool     // IsMutate flag
	}
)

func NewChain() *Chain {
	return &Chain{}
}

func LoadChains(reader io.Reader) ([]Chain, error) {
	chains := make([]Chain, 0)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	oneChain := Chain{}
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			oneChain.Content = append(oneChain.Content, scanner.Text())
		} else {
			chains = append(chains, oneChain)
			oneChain = Chain{}
		}
	}
	chains = append(chains, oneChain)
	return chains, nil
}

func (m *Chain) Marshal() ([]byte, error) {
	buffer := bytes.Buffer{}
	for _, part := range m.Content {
		splitted := strings.Split(part, ":")
		if len(splitted) < 2 {
			return []byte{}, errors.New("Error in data format" + m.ID)
		}
		splittedData := strings.Join(splitted[1:], "")
		switch splitted[0] {
		case "b", "+b":
			shrinked := strings.ReplaceAll(splittedData, " ", "")
			data, err := hex.DecodeString(shrinked)
			if err != nil {
				return []byte{}, err
			}
			buffer.Write(data)
		case "bs", "+bs":
			shrinked := strings.ReplaceAll(splittedData, " ", "")
			data, err := hex.DecodeString(shrinked)
			if err != nil {
				return []byte{}, err
			}
			uint32Buffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(uint32Buffer, uint32(len(data)))
			buffer.Write(uint32Buffer)
			buffer.Write(data)
		case "d", "+d":
			uint32Buffer := make([]byte, 4)
			value, err := strconv.ParseUint(splittedData, 10, 32)
			if err != nil {
				return []byte{}, err
			}
			binary.LittleEndian.PutUint32(uint32Buffer, uint32(value))
			buffer.Write(uint32Buffer)
		case "s", "+s":
			data := []byte(splittedData)
			buffer.Write(data)
		case "ss", "+ss":
			data := []byte(splittedData)
			uint32Buffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(uint32Buffer, uint32(len(data)))
			buffer.Write(uint32Buffer)
			buffer.Write(data)
		case "g":
			size, err := strconv.ParseUint(splittedData, 10, 32)
			if err != nil {
				return []byte{}, err
			}
			data := make([]byte, size)
			rand.Read(data)
			buffer.Write(data)
		case "gs":
			size, err := strconv.ParseUint(splittedData, 10, 32)
			if err != nil {
				return []byte{}, err
			}
			data := make([]byte, size)
			uint32Buffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(uint32Buffer, uint32(size))
			buffer.Write(uint32Buffer)
			rand.Read(data)
			buffer.Write(data)
		case "+gs":
			maxLength, err := strconv.ParseUint(splittedData, 10, 64)
			if err != nil {
				return []byte{}, err
			}
			bigSize, err := rand.Int(rand.Reader, big.NewInt(int64(maxLength)))
			if err != nil {
				return []byte{}, err
			}
			size := uint(bigSize.Uint64())
			data := make([]byte, size)
			uint32Buffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(uint32Buffer, uint32(size))
			buffer.Write(uint32Buffer)
			rand.Read(data)
			buffer.Write(data)
		}
	}
	bufferWithSize := bytes.Buffer{}
	uint32Buffer := make([]byte, 4)

	binary.LittleEndian.PutUint32(uint32Buffer, uint32(len(m.ID)))
	bufferWithSize.Write(uint32Buffer)
	bufferWithSize.Write([]byte(m.ID))
	bufferWithSize.Write(buffer.Bytes())

	return bufferWithSize.Bytes(), nil
}

func (m *Chain) Unmarshal([]byte) error {
	return errors.New("Not implemented")
}

func (ch *Chain) Permutate(permutator func([]byte) ([]byte, error)) error {
	if !ch.IsMutate {
		return nil
	}
	mutated := make([]int, 0)
	for i, s := range ch.Content {
		if strings.HasPrefix(s, "+") {
			mutated = append(mutated, i)
		}
	}
	if len(mutated) == 0 {
		return nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(mutated))))
	if err != nil {
		return err
	}
	selectedField := mutated[n.Uint64()]

	splitted := strings.Split(ch.Content[selectedField], ":")
	if len(splitted) < 2 {
		return errors.New("Error in data format" + ch.ID)
	}
	splittedData := strings.Join(splitted[1:], "")
	switch splitted[0] {
	case "+b", "+bs":
		shrinked := strings.ReplaceAll(splittedData, " ", "")
		data, err := hex.DecodeString(shrinked)
		if err != nil {
			return err
		}
		permutatedData, err := permutator(data)
		if err != nil {
			return err
		}

		ch.Content[selectedField] = splitted[0] + ":" + hex.EncodeToString(permutatedData)
	case "+d", "+s", "+ss":
		data := []byte(splittedData)
		permutatedData, err := permutator(data)
		if err != nil {
			return err
		}
		ch.Content[selectedField] = splitted[0] + ":" + string(permutatedData)
	}
	return nil
}
