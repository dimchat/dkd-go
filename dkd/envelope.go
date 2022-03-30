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

	_sender ID
	_receiver ID
	_time time.Time
}

func NewMessageEnvelope(dict map[string]interface{}, from ID, to ID, when time.Time) *MessageEnvelope {
	if ValueIsNil(dict) {
		if when.IsZero() {
			when = time.Now()
		}
		dict = make(map[string]interface{})
		dict["sender"] = from.String()
		dict["receiver"] = to.String()
		dict["time"] = when.Unix()
	}
	env := new(MessageEnvelope)
	if env.Init(dict) != nil {
		env._sender = from
		env._receiver = to
		env._time = when
	}
	return env
}

func (env *MessageEnvelope) Init(dict map[string]interface{}) *MessageEnvelope {
	if env.Dictionary.Init(dict) != nil {
		// lazy load
		env._sender = nil
		env._receiver = nil
		env._time = time.Unix(0, 0)
	}
	return env
}

//-------- IEnvelope

func (env *MessageEnvelope) Sender() ID {
	if env._sender == nil {
		env._sender = EnvelopeGetSender(env.GetMap(false))
	}
	return env._sender
}

func (env *MessageEnvelope) Receiver() ID {
	if env._receiver == nil {
		receiver := EnvelopeGetReceiver(env.GetMap(false))
		if receiver == nil {
			env._receiver = ANYONE
		} else {
			env._receiver = receiver
		}
	}
	return env._receiver
}

func (env *MessageEnvelope) Time() time.Time {
	if env._time.IsZero() {
		env._time = EnvelopeGetTime(env.GetMap(false))
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
	return EnvelopeGetGroup(env.GetMap(false))
}

func (env *MessageEnvelope) SetGroup(group ID)  {
	EnvelopeSetGroup(env.GetMap(false), group)
}

/*
 *  Message Type
 *  ~~~~~~~~~~~~
 *  because the message content will be encrypted, so
 *  the intermediate nodes(station) cannot recognize what kind of it.
 *  we pick out the content type and set it in envelope
 *  to let the station do its job.
 */
func (env *MessageEnvelope) Type() uint8 {
	return EnvelopeGetType(env.GetMap(false))
}

func (env *MessageEnvelope) SetType(msgType uint8)  {
	EnvelopeSetType(env.GetMap(false), msgType)
}

/**
 *  General Factory
 *  ~~~~~~~~~~~~~~~
 */
type MessageEnvelopeFactory struct {}

//-------- IEnvelopeFactory

func (factory *MessageEnvelopeFactory) CreateEnvelope(from ID, to ID, when time.Time) Envelope {
	return NewMessageEnvelope(nil, from, to, when)
}

func (factory *MessageEnvelopeFactory) ParseEnvelope(env map[string]interface{}) Envelope {
	if env["sender"] == nil {
		// env.sender should not empty
		return nil
	} else {
		return NewMessageEnvelope(env, nil, nil, time.Time{})
	}
}

func BuildEnvelopeFactory() EnvelopeFactory {
	factory := EnvelopeGetFactory()
	if factory == nil {
		factory = new(MessageEnvelopeFactory)
		EnvelopeSetFactory(factory)
	}
	return factory
}
