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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
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
type Message interface {
	Map

	Delegate() MessageDelegate
	SetDelegate(delegate MessageDelegate)

	Envelope() Envelope

	// properties within envelope
	Sender() ID
	Receiver() ID
	Time() Time

	Group() ID
	Type() ContentType
}

func MessageGetEnvelope(msg map[string]interface{}) Envelope {
	return EnvelopeParse(msg)
}

/**
 *  Message Delegate
 *  ~~~~~~~~~~~~~~~~
 */
type MessageDelegate interface {
	InstantMessageDelegate
	//SecureMessageDelegate
	ReliableMessageDelegate
}

/**
 *  Instant Message Delegate
 *  ~~~~~~~~~~~~~~~~~~~~~~~~
 */
type InstantMessageDelegate interface {

	//
	//  Encrypt Content
	//

	/**
	 *  1. Serialize 'message.content' to data (JsON / ProtoBuf / ...)
	 *
	 * @param iMsg - instant message object
	 * @param content - message.content
	 * @param password - symmetric key
	 * @return serialized content data
	 */
	SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte

	/**
	 *  2. Encrypt content data to 'message.data' with symmetric key
	 *
	 * @param iMsg - instant message object
	 * @param data - serialized data of message.content
	 * @param password - symmetric key
	 * @return encrypted message content data
	 */
	EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte

	/**
	 *  3. Encode 'message.data' to String (Base64)
	 *
	 * @param iMsg - instant message object
	 * @param data - encrypted content data
	 * @return String object
	 */
	EncodeData(data []byte, iMsg InstantMessage) string

	//
	//  Encrypt Key
	//

	/**
	 *  4. Serialize message key to data (JsON / ProtoBuf / ...)
	 *
	 * @param iMsg - instant message object
	 * @param password - symmetric key
	 * @return serialized key data
	 */
	SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte

	/**
	 *  5. Encrypt key data to 'message.key' with receiver's public key
	 *
	 * @param iMsg - instant message object
	 * @param data - serialized data of symmetric key
	 * @param receiver - receiver ID string
	 * @return encrypted symmetric key data
	 */
	EncryptKey(data []byte, receiver ID, iMsg InstantMessage) []byte

	/**
	 *  6. Encode 'message.key' to String (Base64)
	 *
	 * @param iMsg - instant message object
	 * @param data - encrypted symmetric key data
	 * @return String object
	 */
	EncodeKey(data []byte, iMsg InstantMessage) string
}

/**
 *  Secure Message Delegate
 *  ~~~~~~~~~~~~~~~~~~~~~~~
 */
type SecureMessageDelegate interface {

	//
	//  Decrypt Key
	//

	/**
	 *  1. Decode 'message.key' to encrypted symmetric key data
	 *
	 * @param key - base64 string object
	 * @param sMsg - secure message object
	 * @return encrypted symmetric key data
	 */
	DecodeKey(key interface{}, sMsg SecureMessage) []byte

	/**
	 *  2. Decrypt 'message.key' with receiver's private key
	 *
	 *  @param key - encrypted symmetric key data
	 *  @param sender - sender/member ID string
	 *  @param receiver - receiver/group ID string
	 *  @param sMsg - secure message object
	 *  @return serialized symmetric key
	 */
	DecryptKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) []byte

	/**
	 *  3. Deserialize message key from data (JsON / ProtoBuf / ...)
	 *
	 * @param key - serialized key data
	 * @param sender - sender/member ID string
	 * @param receiver - receiver/group ID string
	 * @param sMsg - secure message object
	 * @return symmetric key
	 */
	DeserializeKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) SymmetricKey

	//
	//  Decrypt Content
	//

	/**
	 *  4. Decode 'message.data' to encrypted content data
	 *
	 * @param data - base64 string object
	 * @param sMsg - secure message object
	 * @return encrypted content data
	 */
	DecodeData(data interface{}, sMsg SecureMessage) []byte

	/**
	 *  5. Decrypt 'message.data' with symmetric key
	 *
	 *  @param data - encrypt content data
	 *  @param password - symmetric key
	 *  @param sMsg - secure message object
	 *  @return serialized message content
	 */
	DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte

	/**
	 *  6. Deserialize message content from data (JsON / ProtoBuf / ...)
	 *
	 * @param data - serialized content data
	 * @param password - symmetric key
	 * @param sMsg - secure message object
	 * @return message content
	 */
	DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content

	//
	//  Signature
	//

	/**
	 *  1. Sign 'message.data' with sender's private key
	 *
	 *  @param data - encrypted message data
	 *  @param sender - sender ID string
	 *  @param sMsg - secure message object
	 *  @return signature of encrypted message data
	 */
	SignData(data []byte, sender ID, sMsg SecureMessage) []byte

	/**
	 *  2. Encode 'message.signature' to String (Base64)
	 *
	 * @param signature - signature of message.data
	 * @param sMsg - secure message object
	 * @return String object
	 */
	EncodeSignature(signature []byte, sMsg SecureMessage) string
}

/**
 *  Reliable Message Delegate
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~
 */
type ReliableMessageDelegate interface {
	SecureMessageDelegate

	/**
	 *  1. Decode 'message.signature' from String (Base64)
	 *
	 * @param signature - base64 string object
	 * @param rMsg - reliable message
	 * @return signature data
	 */
	DecodeSignature(signature interface{}, rMsg ReliableMessage) []byte

	/**
	 *  2. Verify the message data and signature with sender's public key
	 *
	 *  @param data - message content(encrypted) data
	 *  @param signature - signature for message content(encrypted) data
	 *  @param sender - sender ID/string
	 *  @param rMsg - reliable message object
	 *  @return YES on signature matched
	 */
	VerifyDataSignature(data []byte, signature []byte, sender ID, rMsg ReliableMessage) bool
}
