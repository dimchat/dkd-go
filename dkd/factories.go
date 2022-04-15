/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
 *
 *                                Written in 2022 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Albert Moky
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
	"math/rand"
)

/**
 *  General Factory
 *  ~~~~~~~~~~~~~~~
 */
type MessageEnvelopeFactory struct {}

func (factory *MessageEnvelopeFactory) Init() EnvelopeFactory {
	return factory
}

//-------- IEnvelopeFactory

func (factory *MessageEnvelopeFactory) CreateEnvelope(from ID, to ID, when Time) Envelope {
	return NewEnvelope(nil, from, to, when)
}

func (factory *MessageEnvelopeFactory) ParseEnvelope(env map[string]interface{}) Envelope {
	if env["sender"] == nil {
		// env.sender should not empty
		return nil
	} else {
		return NewEnvelope(env, nil, nil, nil)
	}
}

/**
 *  General Factory
 *  ~~~~~~~~~~~~~~~
 */
type PlainMessageFactory struct {}

func (factory *PlainMessageFactory) Init() InstantMessageFactory {
	return factory
}

//-------- IInstantMessageFactory

func (factory *PlainMessageFactory) GenerateSerialNumber(_ ContentType, _ Time) uint64 {
	// because we must make sure all messages in a same chat box won't have
	// same serial numbers, so we can't use time-related numbers, therefore
	// the best choice is a totally random number, maybe.
	sn := rand.Uint32()
	if sn == 0 {
		// ZERO? do it again!
		sn = 9527 + 9394
	}
	return uint64(sn)
}

func (factory *PlainMessageFactory) CreateInstantMessage(head Envelope, body Content) InstantMessage {
	return NewInstantMessage(nil, head, body)
}

func (factory *PlainMessageFactory) ParseInstantMessage(msg map[string]interface{}) InstantMessage {
	// msg.content should not empty
	if msg["content"] == nil {
		return nil
	}
	return NewInstantMessage(msg, nil, nil)
}

/**
 *  General Factory
 *  ~~~~~~~~~~~~~~~
 */
type EncryptedMessageFactory struct {}

func (factory *EncryptedMessageFactory) Init() SecureMessageFactory {
	return factory
}

//-------- ISecureMessageFactory

func (factory *EncryptedMessageFactory) ParseSecureMessage(msg map[string]interface{}) SecureMessage {
	if _, exists := msg["signature"]; exists {
		// this should be a reliable message
		return NewReliableMessage(msg)
	}
	return NewSecureMessage(msg)
}

/**
 *  General Factory
 *  ~~~~~~~~~~~~~~~
 */
type RelayMessageFactory struct {}

func (factory *RelayMessageFactory) Init() ReliableMessageFactory {
	return factory
}

//-------- IReliableMessageFactory

func (factory *RelayMessageFactory) ParseReliableMessage(msg map[string]interface{}) ReliableMessage {
	// msg.sender should not empty
	// msg.data should not empty
	// msg.signature should not empty
	if msg["sender"] == nil || msg["data"] == nil || msg["signature"] == nil {
		return nil
	}
	return NewReliableMessage(msg)
}
