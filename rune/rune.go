package rune

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/wyixin/runes-go/utils"

	"github.com/btcsuite/btcd/txscript"
)

type Transaction struct {
	Hash      string    `json:"hash"`
	Transfers Transfers `json:"transfers"`
	Issuance  *Rune     `json:"issuance,omitempty"`
}

type Transfers []*Assignment

type Assignment struct {
	ID     uint64 `json:"id"`
	Output uint64 `json:"output"`
	Amount uint64 `json:"amount"`
}

type Rune struct {
	Symbol   string `json:"symbol"`
	Decimals uint64 `json:"decimals"`
}

const RUNE_TRANSFER_PARTS_LEN = 11

// func IsRuneTransaction(txHash string) error {}

// func DecodeRuneTransaction()

// There are two common implementation versions of the Rune protocol.
// 1. op_return "R" [{Assignment}] {Issuance} // https://rodarmor.com/blog/runes
// 2. OP_RETURN OP_PUSHBYTES_1 52 OP_PUSHBYTES_11 0001ff00752b7d00000000 OP_PUSHBYTES_10 ff987806010000000012
func ExtractRuneDataFromScriptPubKeyHexStr(pbKey string) (*Transaction, error) {
	hash, err := hex.DecodeString(pbKey)
	if err != nil {
		return nil, err
	}

	tokenizer := txscript.MakeScriptTokenizer(0, hash)
	if err := tokenizer.Err(); err != nil {
		return nil, err
	}

	opQunun := newQueue()
	for tokenizer.Next() {
		utils.PrettyPrint(utils.OpcodeByCode(tokenizer.Opcode()))
		utils.PrettyPrint(hex.EncodeToString(tokenizer.Data()))
		//		utils.PrettyPrint(utils.HexToBase26(hex.EncodeToString(tokenizer.Data())))
		item := &opcode{
			Code: tokenizer.Opcode(),
			Data: tokenizer.Data(),
		}

		opQunun.Enqueue(item)
	}

	opLen := len(*opQunun)
	if opLen != 3 && opLen != 4 {
		return nil, errors.New("Parse error on rune protocol on length. It must contain two or three parts.")
	}

	if opOne, _ := opQunun.Dequeue(); opOne.Code != txscript.OP_RETURN {
		utils.PrettyPrint(utils.OpcodeByCode(opOne.Code))
		return nil, errors.New("Parse error on rune protocol. It must start with op_return and follow with an 'R'")
	}

	if opTwo, _ := opQunun.Dequeue(); !bytes.Equal(opTwo.Data, []byte{'R'}) {
		return nil, errors.New("Parse error on rune protocol. It must start with op_return and follow with an 'R'")
	}

	// Start parsing transfer part data.
	// op_data_x -> data
	// x = len(data)?
	runeData := &Transaction{
		Transfers: make(Transfers, 0),
	}
	if opThree, _ := opQunun.Dequeue(); !strings.HasPrefix(utils.OpcodeByCode(opThree.Code), "OP_DATA_") {
		return nil, errors.New("Parse error on rune protocol. It must be a datapush operation in part two.")
	} else {
		runeTrans, err := DecodeTransfers(opThree.Data)
		if err != nil {
			return nil, errors.New("Parse error on rune protocol. It must be a datapush operation in part two.")
		}

		runeData.Transfers = runeTrans
	}

	// Issuance
	if opLen == 4 {
		if opFour, _ := opQunun.Dequeue(); !strings.HasPrefix(utils.OpcodeByCode(opFour.Code), "OP_DATA_") {
			return nil, errors.New("Parse error on rune protocol. It must be a datapush operation in part three.")
		} else {
			issuance, err := DecodeIssuance(opFour.Data)

			if err != nil {
				return nil, errors.New("Parse error on rune protocol. It must be a datapush operation in part two.")
			}

			runeData.Issuance = issuance
		}
	}

	return nil, errors.New("runes not found")
}

// decode bits into varint
func decodeTransfer(bits []byte) (id, output, amount uint64, err error) {
	if len(bits) != RUNE_TRANSFER_PARTS_LEN {
		err = errors.New("Parse error on rune protocol, transfer part len error")
		return
	}

	id = uint64(bits[0])
	output = uint64(bits[1])
	encodeLittleEndian := bits[2] == 0xff

	if encodeLittleEndian {
		amount = uint64(binary.LittleEndian.Uint64(bits[3:]))
	} else {
		amount = uint64(binary.BigEndian.Uint64(bits[3:]))
	}

	return
}

func DecodeTransfers(b []byte) (Transfers, error) {
	// 00  01  ff 00752b7d00000000
	// id | output | little endian | amount

	parts := len(b) / RUNE_TRANSFER_PARTS_LEN
	if remainder := len(b) % RUNE_TRANSFER_PARTS_LEN; remainder != 0 {
		return nil, errors.New("Parse error on rune protocol. The transfer message is invalid.")
	}

	trans := make(Transfers, 0)
	for i := 0; i < parts; i++ {
		index := i * RUNE_TRANSFER_PARTS_LEN
		j := index + RUNE_TRANSFER_PARTS_LEN
		id, output, amount, err := decodeTransfer(b[index:j])
		if err != nil {
			return nil, err
		}

		assign := &Assignment{
			ID:     id,
			Output: output,
			Amount: amount,
		}

		trans = append(trans, assign)
	}

	utils.PrettyPrint(trans)
	return trans, nil
}

func reverseBits(input []byte) []byte {
	reversed := make([]byte, len(input))

	j := 0
	for i := len(input) - 1; i >= 0; i-- {
		reversed[j] = input[i]
		j++
	}

	return reversed
}

func DecodeIssuance(b []byte) (*Rune, error) {

	var arr []byte
	dec := uint64(b[len(b)-1])

	encodeLittleEndian := b[0] == 0xff
	if encodeLittleEndian {
		arr = reverseBits(b[1 : len(b)-1])
	} else {
		arr = b[1 : len(b)-1]
	}

	symbol, err := base26Encode(arr)
	if err != nil {
		return nil, err
	}

	return &Rune{
		Symbol:   symbol,
		Decimals: dec,
	}, nil
}
