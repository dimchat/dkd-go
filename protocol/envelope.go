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
)

// Envelope defines a message envelope, encapsulating routing and timing information.
//
// It typically contains details such as the message's sender, receiver, and timestamp.
//
//	data format: {
//		"sender"   : "moki@xxx",
//		"receiver" : "hulk@yyy",
//		"time"     : 123
//	}
type Envelope interface {
	Mapper

	// Sender returns the message sender
	Sender() ID

	// Receiver returns the message receiver
	Receiver() ID

	// Time returns the message time
	Time() Time

	// Group returns the group ID
	//
	// When a group message was split/trimmed to a single message
	// the 'receiver' will be changed to a member ID, and
	// the group ID will be saved as 'group'.
	Group() ID
	SetGroup(group ID)

	// Type returns the message type
	//
	// Because the message content will be encrypted, so
	// the intermediate nodes(station) cannot recognize what kind of it.
	// we pick out the content type and set it in envelope
	// to let the station do its job.
	Type() MessageType
	SetType(msgType MessageType)
}

// EnvelopeFactory creates and parses envelopes
type EnvelopeFactory interface {

	// CreateEnvelope creates a new envelope
	//
	// Parameters:
	//   - from: sender ID
	//   - to: receiver ID
	//   - when: message time
	// Returns: Envelope
	CreateEnvelope(from, to ID, when Time) Envelope

	// ParseEnvelope parses a map object to envelope
	//
	// Parameters:
	//   - env: envelope info
	// Returns: Envelope
	ParseEnvelope(env StringKeyMap) Envelope
}

//
//  Factory methods
//

func CreateEnvelope(from, to ID, when Time) Envelope {
	helper := GetEnvelopeHelper()
	return helper.CreateEnvelope(from, to, when)
}

func ParseEnvelope(env interface{}) Envelope {
	helper := GetEnvelopeHelper()
	return helper.ParseEnvelope(env)
}

func GetEnvelopeFactory() EnvelopeFactory {
	helper := GetEnvelopeHelper()
	return helper.GetEnvelopeFactory()
}

func SetEnvelopeFactory(factory EnvelopeFactory) {
	helper := GetEnvelopeHelper()
	helper.SetEnvelopeFactory(factory)
}
