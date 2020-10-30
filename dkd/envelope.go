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
	"dkd-go/types"
	"time"
)

/**
 *  Envelope for message
 *  ~~~~~~~~~~~~~~~~~~~~
 *  This class is used to create a message envelope
 *  which contains 'sender', 'receiver' and 'time'
 *
 *  data format: {
 *      sender   : "moki@xxx",
 *      receiver : "hulk@yyy",
 *      time     : 123
 *  }
 */
type Envelope struct {
	types.Dictionary

	_delegate *MessageDelegate

	_sender interface{}
	_receiver interface{}
	_group interface{}

	// message type: text, image, ...
	_type protocol.ContentType

	// message time
	_time time.Time
}

func (env *Envelope)LoadEnvelope(dictionary *map[string]interface{}) *Envelope {
	if env.LoadDictionary(dictionary) != nil {
		// lazy load
		env._sender = nil
		env._receiver = nil
		env._group = nil
		env._type = 0
		env._time = time.Unix(0, 0)
	}
	return env
}

func (env *Envelope) InitEnvelope(sender interface{}, receiver interface{}, when time.Time) *Envelope {
	if env.InitDictionary() != nil {
		// message time
		if when.IsZero() {
			when = time.Now()
		}
		// initialize
		env._sender = sender
		env._receiver = receiver
		env._group = nil
		env._type = 0
		env._time = when

		env.Set("sender", sender)
		env.Set("receiver", receiver)
		env.Set("time", when.Unix())
	}
	return env
}

func (env *Envelope) GetDelegate() *MessageDelegate {
	return env._delegate
}

func (env *Envelope) SetDelegate(delegate *MessageDelegate) {
	env._delegate = delegate
}

func (env *Envelope) GetSender() interface{} {
	if env._sender == nil {
		sender := env.Get("sender")
		handler := *env.GetDelegate()
		env._sender = handler.GetID(sender)
	}
	return env._sender
}

func (env *Envelope) GetReceiver() interface{} {
	if env._receiver == nil {
		receiver := env.Get("receiver")
		handler := *env.GetDelegate()
		env._receiver = handler.GetID(receiver)
	}
	return env._receiver
}

func (env *Envelope) GetTime() time.Time {
	if env._time.IsZero() {
		timestamp := env.Get("time")
		env._time = time.Unix(timestamp.(int64), 0)
	}
	return env._time
}

/*
 *  Group ID
 *  ~~~~~~~~
 *  when a group message was split/trimmed to a single message
 *  the 'receiver' will be changed to a member ID, and
 *  the group ID will be saved as 'group'.
 */
func (env *Envelope) GetGroup() interface{} {
	if env._group == nil {
		group := env.Get("group")
		if group != nil {
			handler := *env.GetDelegate()
			env._group = handler.GetID(group)
		}
	}
	return env._group
}

func (env *Envelope) SetGroup(group interface{})  {
	env.Set("group", group.(string))
	env._group = group
}

/*
 *  Message Type
 *  ~~~~~~~~~~~~
 *  because the message content will be encrypted, so
 *  the intermediate nodes(station) cannot recognize what kind of it.
 *  we pick out the content type and set it in envelope
 *  to let the station do its job.
 */
func (env *Envelope) GetType() protocol.ContentType {
	if env._type == 0 {
		t := env.Get("type")
		if t != nil {
			env._type = protocol.ContentType(t.(uint8))
		}
	}
	return env._type
}

func (env *Envelope) SetType(t protocol.ContentType)  {
	env.Set("type", uint8(t))
	env._type = t
}
