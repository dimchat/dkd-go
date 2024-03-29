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
type PlainMessage struct {
	BaseMessage

	_content Content
}

func NewInstantMessage(dict map[string]interface{}, head Envelope, body Content) InstantMessage {
	if ValueIsNil(dict) {
		dict = head.Map()
		dict["content"] = body.Map()
	}
	msg := new(PlainMessage)
	if msg.BaseMessage.Init(dict) != nil {
		msg._env = head
		msg._content = body
	}
	return msg
}

func (msg *PlainMessage) Init(dict map[string]interface{}) InstantMessage {
	if msg.BaseMessage.Init(dict) != nil {
		// lazy load
		msg._content = nil
	}
	return msg
}

//-------- IMessage

func (msg *PlainMessage) Time() Time {
	msgTime := msg.Content().Time()
	if TimeIsNil(msgTime) {
		return msg.Envelope().Time()
	} else {
		return msgTime
	}
}

func (msg *PlainMessage) Group() ID {
	return msg.Content().Group()
}

func (msg *PlainMessage) Type() ContentType {
	return msg.Content().Type()
}

//-------- IInstantMessage

func (msg *PlainMessage) Content() Content {
	if msg._content == nil {
		msg._content = InstantMessageGetContent(msg.Map())
	}
	return msg._content
}

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
 * @return SecureMessage object
 */
func (msg *PlainMessage) Encrypt(password SymmetricKey, members []ID) SecureMessage {
	// 0. check attachment for File/Image/Audio/Video message content
	//    (do it in 'core' module)

	delegate := msg.Delegate()
	content := msg.Content()

	// 1. encrypt 'message.content' to 'message.data'
	data := delegate.SerializeContent(content, password, msg)
	data = delegate.EncryptContent(data, password, msg)
	base64 := delegate.EncodeData(data, msg)
	info := msg.CopyMap(false)
	delete(info, "content")
	info["data"] = base64

	// 2. encrypt symmetric key(password) to 'message.key' or 'message.keys'
	// 2.1. serialize symmetric key
	key := delegate.SerializeKey(password, msg)
	if key == nil {
		// A) broadcast message has no key
		// B) reused key
		return SecureMessageParse(info)
	}
	// 2.2. encrypt symmetric key(s)
	if ValueIsNil(members) {
		// personal message
		key = delegate.EncryptKey(key, msg.Receiver(), msg)
		if key == nil {
			// public key for encryption not found
			// TODO: suspend this message for waiting receiver's meta
			return nil
		}
		// 2.3. encode encrypted key data
		base64 = delegate.EncodeKey(key, msg)
		// 2.4. insert as 'key'
		info["key"] = base64
	} else {
		// group message
		keys := make(map[string]string, len(members))
		count := 0
		for _, member := range members {
			data = delegate.EncryptKey(key, member, msg)
			if data == nil {
				// public key for encryption not found
				// TODO: suspend this message for waiting receiver's meta
				continue
			}
			// 2.3. encode encrypted key data
			base64 = delegate.EncodeKey(data, msg)
			// 2.4. insert to 'message.keys' with member ID
			keys[member.String()] = base64
			count++
		}
		if count > 0 {
			info["keys"] = keys
		}
	}

	// 3. pack message
	return SecureMessageParse(info)
}
