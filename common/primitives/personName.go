package primitives

import (
	"fmt"

	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives/random"
)

type PersonName string

func NewPersonName(description string) (*PersonName, error) {
	d := new(PersonName)

	err := d.SetString(description)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *PersonName) Empty() bool {
	if d.String() == "" {
		return true
	}
	return false
}

func (d *PersonName) SetString(name string) error {
	if len(name) > d.MaxLength() {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			d.MaxLength(), len(name))
	}

	*d = PersonName(name)
	return nil
}

func (d *PersonName) String() string {
	return string(*d)
}

func (a *PersonName) IsSameAs(b *PersonName) bool {
	return a.String() == b.String()
}

func (d *PersonName) MaxLength() int {
	return constants.MAX_NAME_LENGTH
}

func (d *PersonName) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(d.String(), d.MaxLength())
}

func (d *PersonName) UnmarshalBinary(data []byte) error {
	_, err := d.UnmarshalBinaryData(data)
	return err
}

func (d *PersonName) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[LongDesc] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	str, newData, err := UnmarshalStringFromBytesData(newData, d.MaxLength())
	if err != nil {
		return data, err
	}

	err = d.SetString(str)
	if err != nil {
		return data, err
	}

	return newData, nil
}

func RandomName() *PersonName {
	l, _ := NewPersonName("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}
