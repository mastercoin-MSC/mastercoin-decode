package main

import "fmt"
import "strconv"
import "strings"
import "math/big"
import "encoding/binary"

import "github.com/conformal/btcutil"
import "github.com/op/go-logging"

var log = logging.MustGetLogger("")

func main(){
  logging.SetFormatter(logging.MustStringFormatter("[%{level}] %{message}"))
  logging.SetLevel(logging.DEBUG, "")

  ss := DecodeFromAddress("15NoSD4F1ULYHGW3TK6khe1rLSS2qoysaX")
  ss.Explain()

  a := ss.SerializeToAddress()

  ss2 := DecodeFromAddress(a)
  ss2.Explain()

  ss2.SerializeToKey()
}
type Message struct {
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
  var val uint64

  amount := 4

  if v, ok := value.(uint64); ok{
    amount = 8
    val = v
  }else if v, ok := value.(uint32); ok{
    val = uint64(v)
  }else {
    panic(fmt.Sprintf("makeBinary requires a value that's either a uint32 or an uint64, got: %s ", value))
  }

  str := strconv.FormatUint(val, 10)

  number := new(big.Int)
  number.SetString(str, 10)

  template := make([]byte, amount)

  x := number.Bytes()
  z = append(template[:(amount-len(x))],x...)

  return z
}

// Converts a number to a string array
// i.e. 100,8 => [0,0,0,0,0,1,0,0]
// i.e. 66 ,4 => [0,0,6,6]
func makeStringArray(value string, length int) []string{
  z := make([]string,length)
  for i, _ := range z {
    z[i] = "0"
  }

  pointer := length-len(value)
  for _, val := range value {
    z[pointer] = fmt.Sprintf("%c",val)
    pointer++
  }
  return z
}

// Takes SerializeToKey output and builds a valid, obfuscated public key
func (ss *SimpleSend) SerializeToCompressedPublicKey(xor_target string) string{
  return "nothing"
}

// Encodes as Class B
// Encodes the data to a format that will be used as Obfuscate source

func (ss *SimpleSend) SerializeToKey() string{
  log.Info("Encoding data to KEY")

  raw := make([]string, 62)
  for i, _ := range raw {
    raw[i] = "0"

    // This is the 'fake' sequence number, which we don't really need for Class B
    if i == 1{
      raw[i] = "1"
    }
  }

  transaction_type := makeStringArray(strconv.FormatUint(uint64(ss.transaction_type), 16), 8)
  log.Debug("Transaction type: ",transaction_type)

  currency_id :=  makeStringArray(strconv.FormatUint(uint64(ss.currency_id), 16), 8)
  log.Debug("Currency ID: ", currency_id)

  amount := makeStringArray(strconv.FormatUint(ss.amount, 16), 16)
  log.Debug("Amount: ",amount)

  // Start of the data
  pointer := 2

  // TODO: Perhaps make this a bit DRYer if there is no other way of doing it
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

  rawString := strings.Join(raw,"")

  log.Debug("Raw string: ", rawString)

  return rawString
}

// Encodes as Class A
func (ss *SimpleSend) SerializeToAddress() string{
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
