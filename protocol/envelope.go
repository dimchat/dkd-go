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
package protocol

import (
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
type Envelope interface {
	Map
	IEnvelope
}
type IEnvelope interface {

	/**
	 * Get message sender
	 */
	Sender() ID

	/**
	 * Get message receiver
	 */
	Receiver() ID

	/**
	 * Get message time
	 */
	Time() time.Time

	/*
	 *  Group ID
	 *  ~~~~~~~~
	 *  when a group message was split/trimmed to a single message
	 *  the 'receiver' will be changed to a member ID, and
	 *  the group ID will be saved as 'group'.
	 */
	Group() ID
	SetGroup(group ID)

	/*
	 *  Message Type
	 *  ~~~~~~~~~~~~
	 *  because the message content will be encrypted, so
	 *  the intermediate nodes(station) cannot recognize what kind of it.
	 *  we pick out the content type and set it in envelope
	 *  to let the station do its job.
	 */
	Type() uint8
	SetType(msgType uint8)
}

func EnvelopeGetSender(env map[string]interface{}) ID {
	return IDParse(env["sender"])
}

func EnvelopeGetReceiver(env map[string]interface{}) ID {
	return IDParse(env["receiver"])
}

func EnvelopeGetTime(env map[string]interface{}) time.Time {
	timestamp := env["time"]
	if timestamp == nil {
		return time.Time{}
	}
	return time.Unix(timestamp.(int64), 0)
}

func EnvelopeGetGroup(env map[string]interface{}) ID {
	return IDParse(env["group"])
}

func EnvelopeSetGroup(env map[string]interface{}, group ID) {
	if group == nil {
		delete(env, "group")
	} else {
		env["group"] = group.String()
	}
}

func EnvelopeGetType(env map[string]interface{}) uint8 {
	msgType := env["type"]
	if msgType == nil {
		return 0
	}
	return msgType.(uint8)
}

func EnvelopeSetType(env map[string]interface{}, msgType uint8) {
	if msgType == 0 {
		delete(env, "type")
	} else {
		env["type"] = msgType
	}
}

/**
 *  Envelope Factory
 *  ~~~~~~~~~~~~~~~~
 */
type EnvelopeFactory interface {

	/**
	 *  Create envelope
	 *
	 * @param from - sender ID
	 * @param to   - receiver ID
	 * @param when - message time
	 * @return Envelope
	 */
	CreateEnvelope(from ID, to ID, when time.Time) Envelope

	/**
	 *  Parse map object to envelope
	 *
	 * @param env - envelope info
	 * @return Envelope
	 */
	ParseEnvelope(env map[string]interface{}) Envelope
}

var envelopeFactory EnvelopeFactory = nil

func EnvelopeSetFactory(factory EnvelopeFactory) {
	envelopeFactory = factory
}

func EnvelopeGetFactory() EnvelopeFactory {
	return envelopeFactory
}

//
//  Factory methods
//
func EnvelopeCreate(from ID, to ID, when time.Time) Envelope {
	factory := EnvelopeGetFactory()
	return factory.CreateEnvelope(from, to, when)
}

func EnvelopeParse(env interface{}) Envelope {
	if env == nil {
		return nil
	}
	value, ok := env.(Envelope)
	if ok {
		return value
	}
	// get envelope info
	var info map[string]interface{}
	wrapper, ok := env.(Map)
	if ok {
		info = wrapper.GetMap(false)
	} else {
		info, ok = env.(map[string]interface{})
		if !ok {
			panic(env)
			return nil
		}
	}
	// create by envelope factory
	factory := EnvelopeGetFactory()
	return factory.ParseEnvelope(info)
}
