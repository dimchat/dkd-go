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
	"time"
)

type Message interface {
	Map

	Envelope() Envelope

	Sender() ID
	Receiver() ID
	Time() time.Time

	Group() ID
	Type() ContentType
}

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
type InstantMessage interface {
	Message

	Content() Content

	/**
	 *  Encrypt message, replace 'content' field with encrypted 'data'
	 *
	 * @param password - symmetric key
	 * @return SecureMessage object
	 */
	Encrypt(password SymmetricKey, members []ID) SecureMessage
}

/**
 *  Secure Message
 *  ~~~~~~~~~~~~~~
 *  Instant Message encrypted by a symmetric key
 *
 *  data format: {
 *      //-- envelope
 *      sender   : "moki@xxx",
 *      receiver : "hulk@yyy",
 *      time     : 123,
 *      //-- content data and key/keys
 *      data     : "...",  // base64_encode(symmetric)
 *      key      : "...",  // base64_encode(asymmetric)
 *      keys     : {
 *          "ID1": "key1", // base64_encode(asymmetric)
 *      }
 *  }
 */
type SecureMessage interface {
	Message

	EncryptedData() []byte
	EncryptedKey() []byte
	EncryptedKeys() map[string]string

	/**
	 *  Decrypt message, replace encrypted 'data' with 'content' field
	 *
	 * @return InstantMessage object
	 */
	Decrypt() InstantMessage

	/**
	 *  Sign message.data, add 'signature' field
	 *
	 * @return ReliableMessage object
	 */
	Sign() ReliableMessage

	/**
	 *  Split the group message to single person messages
	 *
	 *  @param members - group members
	 *  @return secure/reliable message(s)
	 */
	Split(members []ID) []SecureMessage

	/**
	 *  Trim the group message for a member
	 *
	 * @param member - group member ID/string
	 * @return SecureMessage
	 */
	Trim(member ID) SecureMessage
}

/**
 *  Reliable Message signed by an asymmetric key
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *  This class is used to sign the SecureMessage
 *  It contains a 'signature' field which signed with sender's private key
 *
 *  data format: {
 *      //-- envelope
 *      sender   : "moki@xxx",
 *      receiver : "hulk@yyy",
 *      time     : 123,
 *      //-- content data and key/keys
 *      data     : "...",  // base64_encode(symmetric)
 *      key      : "...",  // base64_encode(asymmetric)
 *      keys     : {
 *          "ID1": "key1", // base64_encode(asymmetric)
 *      },
 *      //-- signature
 *      signature: "..."   // base64_encode()
 *  }
 */
type ReliableMessage interface {
	SecureMessage

	Signature() []byte

	/**
	 *  Verify 'data' and 'signature' field with sender's public key
	 *
	 * @return SecureMessage object
	 */
	Verify() SecureMessage
}
