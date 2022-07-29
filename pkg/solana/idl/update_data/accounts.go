// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package update_data

import (
	"fmt"
	ag_binary "github.com/gagliardetto/binary"
)

type MyAccount struct {
	Data uint64
}

var MyAccountDiscriminator = [8]byte{246, 28, 6, 87, 251, 45, 50, 42}

func (obj MyAccount) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(MyAccountDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Data` param:
	err = encoder.Encode(obj.Data)
	if err != nil {
		return err
	}
	return nil
}

func (obj *MyAccount) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(MyAccountDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[246 28 6 87 251 45 50 42]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Data`:
	err = decoder.Decode(&obj.Data)
	if err != nil {
		return err
	}
	return nil
}
