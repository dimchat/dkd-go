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
	"dkd-go/protocol"
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
type InstantMessage struct {
	Message

	_content *Content
}

func CreateInstantMessage(dictionary *map[string]interface{}) *InstantMessage {
	return new(InstantMessage).LoadInstantMessage(dictionary)
}

func (msg *InstantMessage)LoadInstantMessage(dictionary *map[string]interface{}) *InstantMessage {
	if msg.LoadMessage(dictionary) != nil {
		// lazy load
		msg._content = nil
	}
	return msg
}

func (msg *InstantMessage) InitInstantMessage(env *Envelope, body *Content) *InstantMessage {
	dict := env.GetMap()
	if msg.LoadMessage(&dict) != nil {
		msg._content = body
	}
	return msg
}

func (msg *InstantMessage) GetDelegate() *InstantMessageDelegate {
	delegate := msg.GetEnvelope().GetDelegate()
	handler := (*delegate).(InstantMessageDelegate)
	return &handler
}

func (msg *InstantMessage) SetDelegate(delegate *InstantMessageDelegate) {
	handler := (MessageDelegate)(*delegate)
	msg.GetContent().SetDelegate(&handler)
	msg.GetEnvelope().SetDelegate(&handler)
}

func (msg *InstantMessage) GetContent() *Content {
	if msg._content == nil {
		body := msg.Get("content")
		handler := *msg.GetDelegate()
		msg._content = handler.GetContent(body)
	}
	return msg._content
}

func (msg *InstantMessage) GetTime() time.Time {
	t := msg.GetContent().GetTime()
	if t.IsZero() {
		t = msg.GetEnvelope().GetTime()
	}
	return t
}

func (msg *InstantMessage) GetGroup() interface{} {
	return msg.GetContent().GetGroup()
}

func (msg *InstantMessage) GetType() protocol.ContentType {
	return msg.GetContent().GetType()
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
func (msg *InstantMessage) Encrypt(password interface{}, members []interface{}) *SecureMessage {
	// 0. check attachment for File/Image/Audio/Video message content
	//    (do it in 'core' module)

	handler := *msg.GetDelegate()
	content := *msg.GetContent()

	// 1. encrypt 'message.content' to 'message.data'
	data := handler.SerializeContent(content, password, msg)
	data = handler.EncryptContent(data, password, msg)
	base64 := handler.EncodeData(data, msg)
	info := msg.CopyMap()
	delete(info, "content")
	info["data"] = base64

	// 2. encrypt symmetric key(password) to 'message.key' or 'message.keys'
	// 2.1. serialize symmetric key
	key := handler.SerializeKey(password, msg)
	if key == nil {
		// A) broadcast message has no key
		// B) reused key
		return CreateSecureMessage(&info)
	}
	// 2.2. encrypt symmetric key(s)
	if members == nil {
		// personal message
		key = handler.EncryptKey(key, msg.GetReceiver(), msg)
		if key == nil {
			// public key for encryption not found
			// TODO: suspend this message for waiting receiver's meta
			return nil
		}
		// 2.3. encode encrypted key data
		base64 = handler.EncodeKey(key, msg)
		// 2.4. insert as 'key'
		info["key"] = base64
	} else {
		// group message
		keys := make(map[interface{}]string)
		count := 0
		for _, member := range members {
			data = handler.EncryptKey(key, member, msg)
			if data == nil {
				// public key for encryption not found
				// TODO: suspend this message for waiting receiver's meta
				continue
			}
			// 2.3. encode encrypted key data
			base64 = handler.EncodeKey(data, msg)
			// 2.4. insert to 'message.keys' with member ID
			keys[member] = base64
			count++
		}
		if count > 0 {
			info["keys"] = keys
		}
	}

	// 3. pack message
	return CreateSecureMessage(&info)
}
