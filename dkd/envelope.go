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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
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
type MessageEnvelope struct {
	Dictionary
	IdentifierDelegateHolder

	_delegate IdentifierDelegate

	_sender ID
	_receiver ID
	_group ID

	// message type: text, image, ...
	_type ContentType

	// message time
	_time time.Time
}

func (env *MessageEnvelope) Init(dictionary map[string]interface{}) *MessageEnvelope {
	if env.Dictionary.Init(dictionary) != nil {
		// lazy load
		env._sender = nil
		env._receiver = nil
		env._group = nil
		env._type = 0
		env._time = time.Unix(0, 0)
	}
	return env
}

func (env *MessageEnvelope) InitWithSender(sender ID, receiver ID, when time.Time) *MessageEnvelope {
	if env.Dictionary.Init(nil) != nil {
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

		env.Set("sender", sender.String())
		env.Set("receiver", receiver.String())
		env.Set("time", when.Unix())
	}
	return env
}

func (env MessageEnvelope) Delegate() IdentifierDelegate {
	return env._delegate
}

func (env *MessageEnvelope) SetDelegate(delegate IdentifierDelegate) {
	env._delegate = delegate
}

func (env *MessageEnvelope) Sender() ID {
	if env._sender == nil {
		sender := env.Get("sender")
		delegate := env.Delegate()
		env._sender = delegate.GetID(sender)
	}
	return env._sender
}

func (env *MessageEnvelope) Receiver() ID {
	if env._receiver == nil {
		receiver := env.Get("receiver")
		delegate := env.Delegate()
		env._receiver = delegate.GetID(receiver)
	}
	return env._receiver
}

func (env *MessageEnvelope) Time() time.Time {
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
func (env *MessageEnvelope) Group() ID {
	if env._group == nil {
		group := env.Get("group")
		if group != nil {
			delegate := env.Delegate()
			env._group = delegate.GetID(group)
		}
	}
	return env._group
}

func (env *MessageEnvelope) SetGroup(group ID)  {
	env.Set("group", group.String())
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
func (env *MessageEnvelope) Type() ContentType {
	if env._type == 0 {
		t := env.Get("type")
		if t != nil {
			env._type = ContentType(t.(uint8))
		}
	}
	return env._type
}

func (env *MessageEnvelope) SetType(t ContentType)  {
	env.Set("type", uint8(t))
	env._type = t
}
