package primitives

import (
	"bytes"
	"fmt"
	"time"

	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives/random"
)

type FileList struct {
	FileList []File
}

func NewFileList() *FileList {
	af := new(FileList)
	af.FileList = make([]File, 0)

	return af
}

func RandomFileList(max uint32) *FileList {
	fl := NewFileList()
	l := random.RandomUInt32Between(0, max)
	fl.FileList = make([]File, l)

	for i := range fl.FileList {
		fl.FileList[i] = *(RandomFile())
	}

	return fl
}

func (af *FileList) Empty() bool {
	if len(af.FileList) == 0 {
		return true
	}

	return false
}

func (af *FileList) AddFile(filename string) error {
	f, err := NewFile(filename)
	if err != nil {
		return err
	}

	af.FileList = append(af.FileList, *f)
	return nil
}

func (fl *FileList) GetFiles() []File {
	return fl.FileList
}

func (a *FileList) IsSameAs(b *FileList) bool {
	if len(a.FileList) != len(b.FileList) {
		return false
	}

	for i := range a.FileList {
		if !a.FileList[i].IsSameAs(&(b.FileList[i])) {
			return false
		}
	}

	return true
}

func (fl *FileList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := Uint32ToBytes(uint32(len(fl.FileList)))
	buf.Write(data)

	for _, f := range fl.FileList {
		data, err := f.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	return buf.Next(buf.Len()), nil
}

func (fl *FileList) UnmarshalBinary(data []byte) error {
	_, err := fl.UnmarshalBinaryData(data)
	return err
}

func (fl *FileList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[FileName] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	u, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	newData = newData[4:]

	fl.FileList = make([]File, u)

	var i uint32 = 0
	for ; i < u; i++ {
		newData, err = fl.FileList[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}

type File struct {
	Name      string
	DocHash   Hash
	Timestamp time.Time
	DocType   uint32
	Source    string
	Location  string
}

func NewFile(filename string) (*File, error) {
	f := new(File)
	err := f.SetFileName(filename)
	if err != nil {
		return nil, err
	}

	f.Timestamp = time.Now()

	return f, nil
}

func RandomFile() *File {
	f := new(File)
	str := random.RandStringOfSize(f.MaxLength())
	f.SetFileName(str)
	f.DocHash = *RandomHash()
	f.Timestamp = time.Now()
	f.DocType = constants.DOC_TXT

	str = random.RandStringOfSize(f.MaxLength())
	f.Source = str
	str = random.RandStringOfSize(f.MaxLength())
	f.Location = str

	return f

}

/*func (f *File) GetFileName() string {
	return f.FileName
}*/
const layout = "Jan 2, 2006"

func (f *File) GetFullPath() string {
	return f.Name
}

func (f *File) GetTimeStampFormatted() string {
	return f.Timestamp.Format(layout)
}

func (f *File) SetFileName(filename string) error {
	if len(filename) > f.MaxLength() {
		return fmt.Errorf("Name given is too long, length must be under %d, given length is %d", constants.MAX_FILENAME_LEN, len(filename))
	}

	f.Name = filename
	return nil
}

func (f *File) String() string {
	return fmt.Sprintf("%s ", f.Name)
}

func (d *File) MaxLength() int {
	return constants.MAX_FILENAME_LEN
}

func (a *File) IsSameAs(b *File) bool {
	if a.Name != b.Name {
		return false
	}

	if !a.DocHash.IsSameAs(&a.DocHash) {
		return false
	}

	if a.Timestamp.Unix() != b.Timestamp.Unix() {
		return false
	}

	if a.DocType != b.DocType {
		return false
	}

	if a.Source != b.Source {
		return false
	}

	if a.Location != b.Location {
		return false
	}

	return true
}

func (f *File) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := MarshalStringToBytes(f.Name, f.MaxLength())
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = f.DocHash.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = f.Timestamp.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data = Uint32ToBytes(f.DocType)
	buf.Write(data)

	data, err = MarshalStringToBytes(f.Source, f.MaxLength())
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = MarshalStringToBytes(f.Location, f.MaxLength())
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (f *File) UnmarshalBinary(data []byte) error {
	_, err := f.UnmarshalBinaryData(data)
	return err
}

func (f *File) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	str, newData, err := UnmarshalStringFromBytesData(newData, f.MaxLength())
	if err != nil {
		return data, err
	}

	err = f.SetFileName(str)
	if err != nil {
		return data, err
	}

	newData, err = f.DocHash.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	err = f.Timestamp.UnmarshalBinary(newData[:15])
	if err != nil {
		return data, err
	}
	newData = newData[15:]

	val, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	f.DocType = val
	newData = newData[4:]

	str, newData, err = UnmarshalStringFromBytesData(newData, f.MaxLength())
	if err != nil {
		return data, err
	}
	f.Source = str

	str, newData, err = UnmarshalStringFromBytesData(newData, f.MaxLength())
	if err != nil {
		return data, err
	}
	f.Location = str

	return
}
