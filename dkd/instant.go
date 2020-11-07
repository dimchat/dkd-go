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
	"time"
)

type PlainMessage struct {
	BaseMessage
	InstantMessage

	_content Content
}

func CreateInstantMessage(dictionary map[string]interface{}) InstantMessage {
	return new(PlainMessage).Init(dictionary)
}

func (msg *PlainMessage) Init(dictionary map[string]interface{}) *PlainMessage {
	if msg.BaseMessage.Init(dictionary) != nil {
		// lazy load
		msg._content = nil
	}
	return msg
}

func (msg *PlainMessage) InitWithEnvelope(env Envelope, body Content) *PlainMessage {
	if msg.BaseMessage.InitWithEnvelope(env) != nil {
		msg._content = body
		msg.Set("content", body.GetMap(false))
	}
	return msg
}

func (msg PlainMessage) Delegate() InstantMessageDelegate {
	delegate := msg.BaseMessage.Delegate()
	return delegate.(InstantMessageDelegate)
}

func (msg *PlainMessage) Content() Content {
	if msg._content == nil {
		body := msg.Get("content")
		handler := msg.Delegate()
		msg._content = handler.GetContent(body)
	}
	return msg._content
}

func (msg *PlainMessage) Time() time.Time {
	t := msg.Content().Time()
	if t.IsZero() {
		t = msg.Envelope().Time()
	}
	return t
}

func (msg *PlainMessage) Group() ID {
	return msg.Content().Group()
}

func (msg *PlainMessage) Type() ContentType {
	return msg.Content().Type()
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
	info := msg.GetMap(true)
	delete(info, "content")
	info["data"] = base64

	// 2. encrypt symmetric key(password) to 'message.key' or 'message.keys'
	// 2.1. serialize symmetric key
	key := delegate.SerializeKey(password, msg)
	if key == nil {
		// A) broadcast message has no key
		// B) reused key
		return CreateSecureMessage(info)
	}
	// 2.2. encrypt symmetric key(s)
	if members == nil {
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
		keys := make(map[string]string)
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
	return CreateSecureMessage(info)
}
