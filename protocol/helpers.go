/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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

/**
 *  Content Helper
 */
type ContentHelper interface {
	SetContentFactory(msgType ContentType, factory ContentFactory)
	GetContentFactory(msgType ContentType) ContentFactory

	ParseContent(content interface{}) Content
}

var sharedContentHelper ContentHelper = nil

func SetContentHelper(helper ContentHelper) {
	sharedContentHelper = helper
}

func GetContentHelper() ContentHelper {
	return sharedContentHelper
}

/**
 *  Envelope Helper
 */
type EnvelopeHelper interface {
	SetEnvelopeFactory(factory EnvelopeFactory)
	GetEnvelopeFactory() EnvelopeFactory

	ParseEnvelope(env interface{}) Envelope

	CreateEnvelope(from, to ID, when Time) Envelope
}

var sharedEnvelopeHelper EnvelopeHelper = nil

func SetEnvelopeHelper(helper EnvelopeHelper) {
	sharedEnvelopeHelper = helper
}

func GetEnvelopeHelper() EnvelopeHelper {
	return sharedEnvelopeHelper
}

/**
 *  InstantMessage Helper
 */
type InstantMessageHelper interface {
	SetInstantMessageFactory(factory InstantMessageFactory)
	GetInstantMessageFactory() InstantMessageFactory

	ParseInstantMessage(msg interface{}) InstantMessage

	CreateInstantMessage(head Envelope, body Content) InstantMessage

	GenerateSerialNumber(msgType ContentType, now Time) SerialNumberType
}

var sharedInstantMessageHelper InstantMessageHelper = nil

func SetInstantMessageHelper(helper InstantMessageHelper) {
	sharedInstantMessageHelper = helper
}

func GetInstantMessageHelper() InstantMessageHelper {
	return sharedInstantMessageHelper
}

/**
 *  SecureMessage Helper
 */
type SecureMessageHelper interface {
	SetSecureMessageFactory(factory SecureMessageFactory)
	GetSecureMessageFactory() SecureMessageFactory

	ParseSecureMessage(msg interface{}) SecureMessage
}

var sharedSecureMessageHelper SecureMessageHelper = nil

func SetSecureMessageHelper(helper SecureMessageHelper) {
	sharedSecureMessageHelper = helper
}

func GetSecureMessageHelper() SecureMessageHelper {
	return sharedSecureMessageHelper
}

/**
 *  ReliableMessage Helper
 */
type ReliableMessageHelper interface {
	SetReliableMessageFactory(factory ReliableMessageFactory)
	GetReliableMessageFactory() ReliableMessageFactory

	ParseReliableMessage(msg interface{}) ReliableMessage
}

var sharedReliableMessageHelper ReliableMessageHelper = nil

func SetReliableMessageHelper(helper ReliableMessageHelper) {
	sharedReliableMessageHelper = helper
}

func GetReliableMessageHelper() ReliableMessageHelper {
	return sharedReliableMessageHelper
}
