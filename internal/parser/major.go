package parser

// UnsignedInteger converts  the CBOR-encoded data into a valuid unsigned int of 64 bits
type UnsignedInteger struct {
	data  []byte // CBOR-encoded data
	value uint64 // decoded uinsigend integer value
}

// Decode decodes its data into a valid value
// func (ui *UnsignedInteger) Decode() (uint64, error) {

// 	if !valid(ui.data) {
// 		return ui.value, fmt.Errorf("unsigned integer CBOR-encoded data %v does not appears to be valid", ui.data)
// 	}

// }
