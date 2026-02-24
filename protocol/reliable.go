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
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

// ReliableMessage represents a SecureMessage signed by an asymmetric key, ensuring authenticity and integrity.
//
// It enhances security by including a digital signature, generated using the sender's private key.
// This signature allows recipients to verify the message's origin and confirm that its content
// has not been tampered with since it was signed.
//
//	data format: {
//		//-- envelope
//		"sender"   : "moki@xxx",
//		"receiver" : "hulk@yyy",
//		"time"     : 123,
//
//		//-- content data and keys
//		"data"     : "...",    // base64_encode( symmetric_encrypt(content))
//		"keys"     : {
//			"{ID}"   : "...",  // base64_encode(asymmetric_encrypt(pwd))
//			"digest" : "..."   // hash(pwd.data)
//		},
//		//-- signature
//		"signature": "..."     // base64_encode(asymmetric_sign(data))
//	}
type ReliableMessage interface {
	SecureMessage

	Signature() TransportableData
}

/**
 *  Message Factory
 */

type ReliableMessageFactory interface {

	// ParseReliableMessage parses a map object to reliable message
	//
	// Parameters:
	//   - msg: message info
	// Returns: ReliableMessage
	ParseReliableMessage(msg StringKeyMap) ReliableMessage
}

//
//  Factory method
//

func ParseReliableMessage(msg interface{}) ReliableMessage {
	helper := GetReliableMessageHelper()
	return helper.ParseReliableMessage(msg)
}

func GetReliableMessageFactory() ReliableMessageFactory {
	helper := GetReliableMessageHelper()
	return helper.GetReliableMessageFactory()
}

func SetReliableMessageFactory(factory ReliableMessageFactory) {
	helper := GetReliableMessageHelper()
	helper.SetReliableMessageFactory(factory)
}

//
//  Conveniences
//

func ReliableMessageConvert(array interface{}) []ReliableMessage {
	values := FetchList(array)
	messages := make([]ReliableMessage, 0, len(values))
	var msg ReliableMessage
	for _, item := range values {
		msg = ParseReliableMessage(item)
		if msg == nil {
			continue
		}
		messages = append(messages, msg)
	}
	return messages
}

func ReliableMessageRevert(messages []ReliableMessage) []StringKeyMap {
	array := make([]StringKeyMap, len(messages))
	for idx, msg := range messages {
		array[idx] = msg.Map()
	}
	return array
}
