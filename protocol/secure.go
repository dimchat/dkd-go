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
 *  Secure Message
 *  ~~~~~~~~~~~~~~
 *  Instant Message encrypted by a symmetric key
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
 *      }
 *  }
 */
type SecureMessage interface {
	Message

	EncryptedData() []byte
	EncryptedKey() []byte
	EncryptedKeys() map[string]string

	/*
	 *  Decrypt the Secure Message to Instant Message
	 *
	 *    +----------+      +----------+
	 *    | sender   |      | sender   |
	 *    | receiver |      | receiver |
	 *    | time     |  ->  | time     |
	 *    |          |      |          |  1. PW      = decrypt(key, receiver.SK)
	 *    | data     |      | content  |  2. content = decrypt(data, PW)
	 *    | key/keys |      +----------+
	 *    +----------+
	 */

	/**
	 *  Decrypt message, replace encrypted 'data' with 'content' field
	 *
	 * @return InstantMessage object
	 */
	Decrypt() InstantMessage

	/*
	 *  Sign the Secure Message to Reliable Message
	 *
	 *    +----------+      +----------+
	 *    | sender   |      | sender   |
	 *    | receiver |      | receiver |
	 *    | time     |  ->  | time     |
	 *    |          |      |          |
	 *    | data     |      | data     |
	 *    | key/keys |      | key/keys |
	 *    +----------+      | signature|  1. signature = sign(data, sender.SK)
	 *                      +----------+
	 */

	/**
	 *  Sign message.data, add 'signature' field
	 *
	 * @return ReliableMessage object
	 */
	Sign() ReliableMessage

	/*
	 *  Split/Trim group message
	 *
	 *  for each members, get key from 'keys' and replace 'receiver' to member ID
	 */

	/**
	 *  Split the group message to single person messages
	 *
	 *  @param members - group members
	 *  @return secure/reliable message(s)
	 */
	Split(members []ID) []SecureMessage

	/**
	 *  Trim the group message for a member
	 *
	 * @param member - group member ID/string
	 * @return SecureMessage
	 */
	Trim(member ID) SecureMessage
}

/**
 *  Message Factory
 *  ~~~~~~~~~~~~~~~
 */
type SecureMessageFactory interface {

	/**
	 *  Parse map object to message
	 *
	 * @param msg - message info
	 * @return SecureMessage
	 */
	ParseSecureMessage(msg map[string]interface{}) SecureMessage
}

var secureFactory SecureMessageFactory = nil

func SecureMessageSetFactory(factory SecureMessageFactory) {
	secureFactory = factory
}

func SecureMessageGetFactory() SecureMessageFactory {
	return secureFactory
}

//
//  Factory method
//
func SecureMessageParse(msg interface{}) SecureMessage {
	if msg == nil {
		return nil
	}
	var info map[string]interface{}
	value := ObjectValue(msg)
	switch value.(type) {
	case SecureMessage:
		return value.(SecureMessage)
	case Map:
		info = value.(Map).GetMap(false)
	case map[string]interface{}:
		info = value.(map[string]interface{})
	default:
		panic(msg)
	}
	factory := SecureMessageGetFactory()
	return factory.ParseSecureMessage(info)
}
