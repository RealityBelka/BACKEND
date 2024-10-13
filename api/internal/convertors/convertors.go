package convertors

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

func FromStringToIntSlice(src string) ([]int, error) {
	splitted := strings.Split(src, " ")

	numbers := make([]int, len(splitted))
	for i, digit := range splitted {
		num, err := strconv.Atoi(digit)
		if err != nil {
			return []int{}, err
		}

		numbers[i] = num
	}

	return numbers, nil
}

func FromBinNumbersPairToBinary(bin []byte, numbers []int) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, int32(len(numbers)))
	if err != nil {
		return []byte{}, err
	}

	for _, num := range numbers {
		err = binary.Write(buf, binary.BigEndian, int32(num))
		if err != nil {
			return []byte{}, err
		}
	}

	err = binary.Write(buf, binary.BigEndian, int32(len(bin)))
	if err != nil {
		return []byte{}, err
	}

	_, err = buf.Write(bin)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
