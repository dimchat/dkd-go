/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
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
package protocol

import (
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
type ReliableMessage interface {
	IReliableMessage
	SecureMessage
}
type IReliableMessage interface {

	Signature() []byte

	/**
	 *  Sender's Meta
	 *  ~~~~~~~~~~~~~
	 *  Extends for the first message package of 'Handshake' protocol.
	 *
	 * @param meta - Meta
	 */
	Meta() Meta
	SetMeta(meta Meta)

	/**
	 *  Sender's Visa
	 *  ~~~~~~~~~~~~~
	 *  Extends for the first message package of 'Handshake' protocol.
	 *
	 * @param doc - Visa
	 */
	Visa() Visa
	SetVisa(visa Visa)

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
	Verify() SecureMessage
}

func ReliableMessageGetMeta(msg map[string]interface{}) Meta {
	return MetaParse(msg["meta"])
}

func ReliableMessageSetMeta(msg map[string]interface{}, meta Meta) {
	if ValueIsNil(meta) {
		delete(msg, "meta")
	} else {
		msg["meta"] = meta.GetMap(false)
	}
}

func ReliableMessageGetVisa(msg map[string]interface{}) Visa {
	doc := DocumentParse(msg["visa"])
	visa, ok := doc.(Visa)
	if ok {
		return visa
	} else {
		return nil
	}
}

func ReliableMessageSetVisa(msg map[string]interface{}, visa Visa) {
	if ValueIsNil(visa) {
		delete(msg, "visa")
	} else {
		msg["visa"] = visa.GetMap(false)
	}
}

/**
 *  Message Factory
 *  ~~~~~~~~~~~~~~~
 */
type ReliableMessageFactory interface {
	IReliableMessageFactory
}
type IReliableMessageFactory interface {

	/**
	 *  Parse map object to message
	 *
	 * @param msg - message info
	 * @return ReliableMessage
	 */
	ParseReliableMessage(msg map[string]interface{}) ReliableMessage
}

var reliableFactory ReliableMessageFactory = nil

func ReliableMessageSetFactory(factory ReliableMessageFactory) {
	reliableFactory = factory
}

func ReliableMessageGetFactory() ReliableMessageFactory {
	return reliableFactory
}

//
//  Factory method
//
func ReliableMessageParse(msg interface{}) ReliableMessage {
	if ValueIsNil(msg) {
		return nil
	}
	value, ok := msg.(ReliableMessage)
	if ok {
		return value
	}
	// get message info
	var info map[string]interface{}
	wrapper, ok := msg.(Map)
	if ok {
		info = wrapper.GetMap(false)
	} else {
		info, ok = msg.(map[string]interface{})
		if !ok {
			panic(msg)
			return nil
		}
	}
	// create by message factory
	factory := ReliableMessageGetFactory()
	if factory == nil {
		panic("reliable message factory not found")
	}
	return factory.ParseReliableMessage(info)
}
