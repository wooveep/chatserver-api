package uuid

import "fmt"

var KEYTable = makeTable(0x80e2)

// Table is a 256-word table representing the polynomial for efficient processing.
type Table struct {
	entries  [256]uint16
	reversed bool
	noXOR    bool
}

func MakeTable(poly uint16) *Table {
	return makeTable(poly)
}

// MakeBitsReversedTable returns the Table constructed from the specified polynomial.
func MakeBitsReversedTable(poly uint16) *Table {
	return makeBitsReversedTable(poly)
}

// MakeTableNoXOR returns the Table constructed from the specified polynomial.
// Updates happen without XOR in and XOR out.
func MakeTableNoXOR(poly uint16) *Table {
	tab := makeTable(poly)
	tab.noXOR = true
	return tab
}

// makeTable returns the Table constructed from the specified polynomial.
func makeBitsReversedTable(poly uint16) *Table {
	t := &Table{
		reversed: true,
	}
	width := uint16(16)
	for i := uint16(0); i < 256; i++ {
		crc := i << (width - 8)
		for j := 0; j < 8; j++ {
			if crc&(1<<(width-1)) != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		t.entries[i] = crc
	}
	return t
}

func makeTable(poly uint16) *Table {
	t := &Table{
		reversed: false,
	}
	for i := 0; i < 256; i++ {
		crc := uint16(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
		t.entries[i] = crc
	}
	return t
}

func updateBitsReversed(crc uint16, tab *Table, p []byte) uint16 {
	for _, v := range p {
		crc = tab.entries[byte(crc>>8)^v] ^ (crc << 8)
	}
	return crc
}

func update(crc uint16, tab *Table, p []byte) uint16 {
	crc = ^crc

	for _, v := range p {
		crc = tab.entries[byte(crc)^v] ^ (crc >> 8)
	}

	return ^crc
}

func updateNoXOR(crc uint16, tab *Table, p []byte) uint16 {
	for _, v := range p {
		crc = tab.entries[byte(crc)^v] ^ (crc >> 8)
	}

	return crc
}

func Update(crc uint16, tab *Table, p []byte) uint16 {
	if tab.reversed {
		return updateBitsReversed(crc, tab, p)
	} else if tab.noXOR {
		return updateNoXOR(crc, tab, p)
	} else {
		return update(crc, tab, p)
	}
}

func ChecksumKey(data []byte) string {

	return fmt.Sprintf("%04X", Update(0, KEYTable, data))

}
