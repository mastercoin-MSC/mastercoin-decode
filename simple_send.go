package main

import "fmt"
import "strconv"
import "math/big"
import "encoding/binary"

import "github.com/conformal/btcutil"
import "github.com/op/go-logging"

var log = logging.MustGetLogger("")

func main(){
  logging.SetFormatter(logging.MustStringFormatter("[%{level}] %{message}"))
  logging.SetLevel(logging.INFO, "")

  ss := DecodeFromAddress("15NoSD4F1ULYHGW3TK6khe1rLSS2qoysaX")
  ss.Explain()

  a := ss.EncodeToAddress()

  ss2 := DecodeFromAddress(a)
  ss2.Explain()
}

type SimpleSend struct {
  currency_id uint32
  sequence byte
  transaction_type uint32
  amount uint64
}

func (ss *SimpleSend) Explain(){
  fmt.Printf("This is a Simple Send transaction for currency id %d and amount %d\n", ss.currency_id, ss.amount)
}

// Converts a number to a binary byte slice
// i.e. 4 => [0,0,0,4]
// 256 => [0,0,1,0]
// or in the case of 64
// 4 = > [0,0,0,0,0,0,0,4]
func makeBinary(value interface{}) []byte{
  var z []byte

  if val, ok := value.(uint32); ok{
    str := strconv.FormatUint(uint64(val), 10)

    number := new(big.Int)
    number.SetString(str, 10)

    template := []byte{0,0,0,0}

    x := number.Bytes()
    z = append(template[:(4-len(x))],x...)
  } else if val, ok := value.(uint64); ok{
    str := strconv.FormatUint(val, 10)

    number := new(big.Int)
    number.SetString(str, 10)

    template := []byte{0,0,0,0,0,0,0,0}

    x := number.Bytes()
    z = append(template[:(8-len(x))],x...)
  } else {
    panic(fmt.Sprintf("makeBinary requires a value that's either a uint32 or an uint64, got: %s ", value))
  }

  return z
}

// Encodes Class A - Simple Sends
func (ss *SimpleSend) EncodeToAddress() string{
  log.Info("Encoding data to address")

  raw := make([]byte, 25)
  var sequence byte = ss.sequence
  raw[1] = sequence

  transaction_type := makeBinary(ss.transaction_type)
  currency_id := makeBinary(ss.currency_id)
  amount := makeBinary(ss.amount)

  //TODO: Can we optimise this?
  pointer := 2
  for _, value := range transaction_type {
    raw[pointer] = value
    pointer++
  }
  for _, value := range currency_id {
    raw[pointer] = value
    pointer++
  }
  for _, value := range amount {
    raw[pointer] = value
    pointer++
  }
  //////////////////////////////

  rawData := btcutil.Base58Encode(raw)
  log.Debug("Raw information: ", raw)
  log.Debug("Encoded to address", rawData)
  return rawData
}


// Decodes Class A - Simple Sends
func DecodeFromAddress(address string) SimpleSend{
  log.Info("Decoding address '%s'.\n", address)

  rawData := btcutil.Base58Decode(address)

  log.Debug("Base58 decoded data: %v \n", rawData)

  sequence := rawData[1]
  log.Debug("Sequence %v", sequence)

  transaction_type := binary.BigEndian.Uint32(rawData[2:6])
  log.Debug("Transaction type: %v",transaction_type)

  currency_id := binary.BigEndian.Uint32(rawData[6:10])
  log.Debug("Currency id: %v ", currency_id)

  amount := binary.BigEndian.Uint64(rawData[10:18])
  log.Debug("Amount: %v", amount)

  ss := SimpleSend{amount: amount, currency_id: currency_id, transaction_type: transaction_type, sequence: sequence}
  return ss
}
