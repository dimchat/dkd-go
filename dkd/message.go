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
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/types"
	"time"
	"unsafe"
)

/*
 *  Message Transforming
 *  ~~~~~~~~~~~~~~~~~~~~
 *
 *     Instant Message <-> Secure Message <-> Reliable Message
 *     +-------------+     +------------+     +--------------+
 *     |  sender     |     |  sender    |     |  sender      |
 *     |  receiver   |     |  receiver  |     |  receiver    |
 *     |  time       |     |  time      |     |  time        |
 *     |             |     |            |     |              |
 *     |  content    |     |  data      |     |  data        |
 *     +-------------+     |  key/keys  |     |  key/keys    |
 *                         +------------+     |  signature   |
 *                                            +--------------+
 *     Algorithm:
 *         data      = password.encrypt(content)
 *         key       = receiver.public_key.encrypt(password)
 *         signature = sender.private_key.sign(data)
 */

/**
 *  Message with Envelope
 *  ~~~~~~~~~~~~~~~~~~~~~
 *  Base classes for messages
 *  This class is used to create a message
 *  with the envelope fields, such as 'sender', 'receiver', and 'time'
 *
 *  data format: {
 *      //-- envelope
 *      sender   : "moki@xxx",
 *      receiver : "hulk@yyy",
 *      time     : 123,
 *      //-- body
 *      ...
 *  }
 */
type Message struct {
	Dictionary

	_env *Envelope
}

func CreateMessage(dictionary map[string]interface{}) *Message {
	if _, exists := dictionary["content"]; exists {
		// this should be an instant message
		msg := CreateInstantMessage(dictionary)
		return (*Message)(unsafe.Pointer(msg))
	} else if _, exists := dictionary["signature"]; exists {
		// this should be a reliable message
		msg := CreateReliableMessage(dictionary)
		return (*Message)(unsafe.Pointer(msg))
	} else if _, exists := dictionary["data"]; exists {
		// this should be a secure message
		msg := CreateSecureMessage(dictionary)
		return (*Message)(unsafe.Pointer(msg))
	}
	//panic("message error")
	return new(Message).Init(dictionary)
}

func (msg *Message)Init(dictionary map[string]interface{}) *Message {
	if msg.Dictionary.Init(dictionary) != nil {
		// lazy load
		msg._env = nil
	}
	return msg
}

func (msg *Message) InitWithEnvelope(env *Envelope) *Message {
	dict := env.GetMap(false)
	if msg.Dictionary.Init(dict) != nil {
		msg._env = env
	}
	return msg
}

func (msg *Message) GetDelegate() *MessageDelegate {
	return msg.GetEnvelope().GetDelegate()
}

func (msg *Message) SetDelegate(delegate *MessageDelegate) {
	msg.GetEnvelope().SetDelegate(delegate)
}

func (msg *Message) GetEnvelope() *Envelope {
	if msg._env == nil {
		dict := msg.GetMap(false)
		msg._env = new(Envelope).Init(dict)
	}
	return msg._env
}

func (msg *Message) GetSender() *ID {
	return msg.GetEnvelope().GetSender()
}

func (msg *Message) GetReceiver() *ID {
	return msg.GetEnvelope().GetReceiver()
}

func (msg *Message) GetTime() time.Time {
	return msg.GetEnvelope().GetTime()
}

func (msg *Message) GetGroup() *ID {
	return msg.GetEnvelope().GetGroup()
}

func (msg *Message) GetType() ContentType {
	return msg.GetEnvelope().GetType()
}
