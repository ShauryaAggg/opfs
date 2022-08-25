package utils

import (
	"errors"

	"github.com/ShauryaAg/opfs/types"
	"github.com/amzn/ion-go/ion"
)

func ConvertEhrToBinary(data types.Ehr) ([]byte, error) {
	bytes, err := ion.MarshalBinary(data)
	if err != nil {
		return nil, errors.New("error marshalling data to binary")
	}

	return bytes, nil
}

func ConvertBinaryToEhr(data []byte) (types.Ehr, error) {
	var ehr types.Ehr
	err := ion.Unmarshal(data, &ehr)
	if err != nil {
		return types.Ehr{}, errors.New("error unmarshalling data to binary")
	}

	return ehr, nil
}
