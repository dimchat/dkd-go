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
	InstantMessage

	_content Content
}

func NewPlainMessage(dict map[string]interface{}, head Envelope, body Content) *PlainMessage {
	if dict == nil {
		dict = head.GetMap(false)
		dict["content"] = body.GetMap(false)
	}
	msg := new(PlainMessage).Init(dict)
	if msg != nil {
		msg._env = head
		msg._content = body
	}
	return msg
}

func (msg *PlainMessage) Init(dict map[string]interface{}) *PlainMessage {
	if msg.BaseMessage.Init(dict) != nil {
		// lazy load
		msg._content = nil
	}
	return msg
}

func (msg *PlainMessage) Equal(other interface{}) bool {
	return msg.BaseMessage.Equal(other)
}

//-------- Map

func (msg *PlainMessage) Get(name string) interface{} {
	return msg.BaseMessage.Get(name)
}

func (msg *PlainMessage) Set(name string, value interface{}) {
	msg.BaseMessage.Set(name, value)
}

func (msg *PlainMessage) Keys() []string {
	return msg.BaseMessage.Keys()
}

func (msg *PlainMessage) GetMap(clone bool) map[string]interface{} {
	return msg.BaseMessage.GetMap(clone)
}

//-------- Message

func (msg *PlainMessage) Delegate() MessageDelegate {
	return msg.BaseMessage.Delegate()
}

func (msg *PlainMessage) SetDelegate(delegate MessageDelegate) {
	msg.BaseMessage.SetDelegate(delegate)
}

func (msg *PlainMessage) Envelope() Envelope {
	return msg.BaseMessage.Envelope()
}

func (msg *PlainMessage) Sender() ID {
	return msg.BaseMessage.Sender()
}

func (msg *PlainMessage) Receiver() ID {
	return msg.BaseMessage.Receiver()
}

func (msg *PlainMessage) Time() time.Time {
	msgTime := msg.Content().Time()
	if msgTime.IsZero() {
		msgTime = msg.Envelope().Time()
	}
	return msgTime
}

func (msg *PlainMessage) Group() ID {
	return msg.Content().Group()
}

func (msg *PlainMessage) Type() uint8 {
	return msg.Content().Type()
}

//-------- InstantMessage

func (msg *PlainMessage) Content() Content {
	if msg._content == nil {
		msg._content = InstantMessageGetContent(msg.GetMap(false))
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
	info := msg.GetMap(true)
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

/**
 *  General Factory
 *  ~~~~~~~~~~~~~~~
 */
type PlainMessageFactory struct {
	InstantMessageFactory
}

func (factory *PlainMessageFactory) CreateInstantMessage(head Envelope, body Content) InstantMessage {
	return NewPlainMessage(nil, head, body)
}

func (factory *PlainMessageFactory) ParseInstantMessage(msg map[string]interface{}) InstantMessage {
	return NewPlainMessage(msg, nil, nil)
}

func BuildInstantMessageFactory() InstantMessageFactory {
	factory := InstantMessageGetFactory()
	if factory == nil {
		factory = new(PlainMessageFactory)
		InstantMessageSetFactory(factory)
	}
	return factory
}

func init() {
	BuildInstantMessageFactory()
}
