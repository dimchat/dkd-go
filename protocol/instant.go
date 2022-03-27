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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Instant Message
 *  ~~~~~~~~~~~~~~~
 *
 *  data format: {
 *      //-- envelope
 *      sender   : "moki@xxx",
 *      receiver : "hulk@yyy",
 *      time     : 123,
 *      //-- content
 *      content  : {...}
 *  }
 */
type InstantMessage interface {
	IInstantMessage
	Message
}
type IInstantMessage interface {

	Content() Content

	/*
	 *  Encrypt the Instant Message to Secure Message
	 *
	 *    +----------+      +----------+
	 *    | sender   |      | sender   |
	 *    | receiver |      | receiver |
	 *    | time     |  ->  | time     |
	 *    |          |      |          |
	 *    | content  |      | data     |  1. data = encrypt(content, PW)
	 *    +----------+      | key/keys |  2. key  = encrypt(PW, receiver.PK)
	 *                      +----------+
	 */

	/**
	 *  Encrypt message, replace 'content' field with encrypted 'data'
	 *
	 * @param password - symmetric key
	 * @param members  - group members; nil for personal message
	 * @return SecureMessage object
	 */
	Encrypt(password SymmetricKey, members []ID) SecureMessage
}

func InstantMessageGetContent(msg map[string]interface{}) Content {
	return ContentParse(msg["content"])
}

/**
 *  Message Factory
 *  ~~~~~~~~~~~~~~~
 */
type InstantMessageFactory interface {
	IInstantMessageFactory
}
type IInstantMessageFactory interface {

	/**
	 *  Create instant message with envelope & content
	 *
	 * @param head - message envelope
	 * @param body - message content
	 * @return InstantMessage
	 */
	CreateInstantMessage(head Envelope, body Content) InstantMessage

	/**
	 *  Parse map object to message
	 *
	 * @param msg - message info
	 * @return InstantMessage
	 */
	ParseInstantMessage(msg map[string]interface{}) InstantMessage
}

var instantFactory InstantMessageFactory = nil

func InstantMessageSetFactory(factory InstantMessageFactory) {
	instantFactory = factory
}

func InstantMessageGetFactory() InstantMessageFactory {
	return instantFactory
}

//
//  Factory methods
//
func InstantMessageCreate(head Envelope, body Content) InstantMessage {
	factory := InstantMessageGetFactory()
	if factory == nil {
		panic("instant message factory not found")
	}
	return factory.CreateInstantMessage(head, body)
}

func InstantMessageParse(msg interface{}) InstantMessage {
	if ValueIsNil(msg) {
		return nil
	}
	value, ok := msg.(InstantMessage)
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
	factory := InstantMessageGetFactory()
	if factory == nil {
		panic("instant message factory not found")
	}
	return factory.ParseInstantMessage(info)
}
