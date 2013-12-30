package main

import ("testing"
        "bytes"
)

func TestMakeBinary(t *testing.T){
  tests := []struct {
    value uint32
    expected []byte
  }{
    {3, []byte{0,0,0,3}},
    {256, []byte{0,0,1,0}},
    {123456, []byte{0,1,226,64}},
    {1239923456, []byte{73,231,187,0}},
  }
  
  for _, pair := range tests {
    v := makeBinary(pair.value)
    if !bytes.Equal(v,pair.expected) {
      t.Error("Expected",pair.expected ,
      "For", pair.value,
      "Got", v)
    }
  }
}

func TestMakeBinary64(t *testing.T){
  tests := []struct {
    value uint64
    expected []byte
  }{
    {3, []byte{0,0,0,0,0,0,0,3}},
    {256, []byte{0,0,0,0,0,0,1,0}},
    {123456, []byte{0,0,0,0,0,1,226,64}},
    {1239923456, []byte{0,0,0,0,73,231,187,0}},
    {221239923456, []byte{0,0,0,51,130,237,83,0}},
  }
  
  for _, pair := range tests {
    v := makeBinary(pair.value)
    if !bytes.Equal(v,pair.expected) {
      t.Error("Expected",pair.expected ,
      "For", pair.value,
      "Got", v)
    }
  }
}

func TestDecodeFromAddress(t *testing.T){
  tests := []struct {
    address string
    values SimpleSend
  }{
    {"17vrMab8gQx72eCEaUxJzL4fg5VwEUumJQ", SimpleSend{currency_id: 2, sequence: 76, transaction_type: 0, amount: 50}},
    {"15NoSD4F1ULYHPfSiV1dp1kr9n2bBffGGd", SimpleSend{currency_id: 2, sequence: 48, transaction_type: 0, amount: 12382}},
    {"15NoSD4F1ULYHGW3TK6khj6NEZsPAmHf41", SimpleSend{currency_id: 1, sequence: 48, transaction_type: 0, amount: 102382}},
  }

  for _, pair := range tests{
    v := DecodeFromAddress(pair.address)
    if v != pair.values {
      t.Error("For", pair.address,
              "Expected", pair.values,
              "Got", v,
      )

    }
  }
}

func TestEncodeToAddress(t *testing.T){
  tests := []struct {
    address string
    values SimpleSend
  }{
    {"17vrMab8gQx72eCEaUxJzL4fg5VwDuND4T", SimpleSend{currency_id: 2, sequence: 76, transaction_type: 0, amount: 50}},
    {"15NoSD4F1ULYHPfSiV1dp1kr9n2b9Npxf1", SimpleSend{currency_id: 2, sequence: 48, transaction_type: 0, amount: 12382}},
    {"15NoSD4F1ULYHGW3TK6khj6NEZsP9ariEK", SimpleSend{currency_id: 1, sequence: 48, transaction_type: 0, amount: 102382}},
  }

  for _, pair := range tests{
    v := pair.values.EncodeToAddress()
    if v != pair.address {
      t.Error("For", pair.address,
              "Expected", pair.address,
              "Got", v,
      )

    }
  }
}
