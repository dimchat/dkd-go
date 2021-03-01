/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
 *
 *                                Written in 2020 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package dkd

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Reliable Message signed by an asymmetric key
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *  This class is used to sign the SecureMessage
 *  It contains a 'signature' field which signed with sender's private key
 *
 *  data format: {
 *      //-- envelope
 *      sender   : "moki@xxx",
 *      receiver : "hulk@yyy",
 *      time     : 123,
 *      //-- content data and key/keys
 *      data     : "...",  // base64_encode(symmetric)
 *      key      : "...",  // base64_encode(asymmetric)
 *      keys     : {
 *          "ID1": "key1", // base64_encode(asymmetric)
 *      },
 *      //-- signature
 *      signature: "..."   // base64_encode()
 *  }
 */
type RelayMessage struct {
	EncryptedMessage
	IReliableMessage

	_signature []byte

	_meta Meta
	_visa Visa
}

func NewRelayMessage(dict map[string]interface{}) *RelayMessage {
	msg := new(RelayMessage).Init(dict)
	ObjectRetain(msg)
	return msg
}

func (msg *RelayMessage) Init(dict map[string]interface{}) *RelayMessage {
	if msg.EncryptedMessage.Init(dict) != nil {
		// lazy load
		msg._signature = nil

		msg.setMeta(nil)
		msg.setVisa(nil)
	}
	return msg
}

func (msg *RelayMessage) Release() int {
	cnt := msg.EncryptedMessage.Release()
	if cnt == 0 {
		// this object is going to be destroyed,
		// release children
		msg.setMeta(nil)
		msg.setVisa(nil)
	}
	return cnt
}

func (msg *RelayMessage) setMeta(meta Meta)  {
	if meta != msg._meta {
		ObjectRetain(meta)
		ObjectRelease(msg._meta)
		msg._meta = meta
	}
}

func (msg *RelayMessage) setVisa(visa Visa)  {
	if visa != msg._visa {
		ObjectRetain(visa)
		ObjectRelease(msg._visa)
		msg._visa = visa
	}
}

//-------- IReliableMessage

func (msg *RelayMessage) Signature() []byte {
	if msg._signature == nil {
		base64 := msg.Get("signature")
		delegate := msg.Delegate()
		msg._signature = delegate.DecodeSignature(base64.(string), msg)
	}
	return msg._signature
}

func (msg *RelayMessage) Meta() Meta {
	if msg._meta == nil {
		msg.setMeta(ReliableMessageGetMeta(msg.GetMap(false)))
	}
	return msg._meta
}

func (msg *RelayMessage) SetMeta(meta Meta) {
	ReliableMessageSetMeta(msg.GetMap(false), meta)
	msg.setMeta(meta)
}

func (msg *RelayMessage) Visa() Visa {
	if msg._visa == nil {
		msg.setVisa(ReliableMessageGetVisa(msg.GetMap(false)))
	}
	return msg._visa
}

func (msg *RelayMessage) SetVisa(visa Visa) {
	ReliableMessageSetVisa(msg.GetMap(false), visa)
	msg.setVisa(visa)
}

/*
 *  Verify the Reliable Message to Secure Message
 *
 *    +----------+      +----------+
 *    | sender   |      | sender   |
 *    | receiver |      | receiver |
 *    | time     |  ->  | time     |
 *    |          |      |          |
 *    | data     |      | data     |  1. verify(data, signature, sender.PK)
 *    | key/keys |      | key/keys |
 *    | signature|      +----------+
 *    +----------+
 */

/**
 *  Verify 'data' and 'signature' field with sender's public key
 *
 * @return SecureMessage object
 */
func (msg *RelayMessage) Verify() SecureMessage {
	data := msg.EncryptedData()
	if data == nil {
		panic("failed to decode content data")
	}
	signature := msg.Signature()
	if signature == nil {
		panic("failed to decode message signature")
	}
	sender := msg.Sender()
	// 1. verify data signature with sender's public key
	delegate := msg.Delegate()
	if delegate.VerifyDataSignature(data, signature, sender, msg) {
		// 2. pack message
		info := msg.GetMap(true)
		delete(info, "signature")
		return SecureMessageParse(info)
	} else {
		//panic("message signature not match")
		return nil
	}
}

/**
 *  General Factory
 *  ~~~~~~~~~~~~~~~
 */
type RelayMessageFactory struct {
	ReliableMessageFactory
}

func (factory *RelayMessageFactory) ParseSecureMessage(msg map[string]interface{}) ReliableMessage {
	rMsg := NewRelayMessage(msg)
	ObjectAutorelease(rMsg)
	return rMsg
}

func BuildReliableMessageFactory() ReliableMessageFactory {
	factory := ReliableMessageGetFactory()
	if factory == nil {
		factory = new(RelayMessageFactory)
		ReliableMessageSetFactory(factory)
	}
	return factory
}
