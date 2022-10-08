package smb2

import (
	"Orca_Puppet/pkg/psexec/common"
	"Orca_Puppet/pkg/psexec/encoder"
	"Orca_Puppet/pkg/psexec/gss"
	"Orca_Puppet/pkg/psexec/ms"
	"Orca_Puppet/pkg/psexec/ntlm"
	"Orca_Puppet/pkg/psexec/smb"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
)

// 此文件提供smb连接方法

type Client struct {
	common.Client
}

func NewSMB2Header() smb.SMB2Header {
	return smb.SMB2Header{
		ProtocolId:    []byte(smb.ProtocolSMB2),
		StructureSize: 64,
		CreditCharge:  0,
		Status:        0,
		Command:       0,
		Credits:       0,
		Flags:         0,
		NextCommand:   0,
		MessageId:     0,
		Reserved:      0,
		TreeId:        0,
		SessionId:     0,
		Signature:     make([]byte, 16),
	}
}

// 协商版本请求初始化
func (c *Client) NewNegotiateRequest() smb.SMB2NegotiateRequestStruct {
	// 初始化
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_NEGOTIATE
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.CreditCharge = 1
	return smb.SMB2NegotiateRequestStruct{
		SMB2Header:      smb2Header,
		StructureSize:   36,
		DialectCount:    1,
		SecurityMode:    smb.SecurityModeSigningEnabled, // 必须开启签名
		Reserved:        0,
		Capabilities:    0,
		ClientGuid:      make([]byte, 16),
		ClientStartTime: 0,
		Dialects: []uint16{
			uint16(smb.SMB2_1_Dialect),
		},
	}
}

// 协商版本响应初始化
func NewNegotiateResponse() smb.SMB2NegotiateResponseStruct {
	smb2Header := NewSMB2Header()
	return smb.SMB2NegotiateResponseStruct{
		SMB2Header:           smb2Header,
		StructureSize:        0,
		SecurityMode:         0,
		DialectRevision:      0,
		Reserved:             0,
		ServerGuid:           make([]byte, 16),
		Capabilities:         0,
		MaxTransactSize:      0,
		MaxReadSize:          0,
		MaxWriteSize:         0,
		SystemTime:           0,
		ServerStartTime:      0,
		SecurityBufferOffset: 0,
		SecurityBufferLength: 0,
		Reserved2:            0,
		SecurityBlob:         &gss.NegTokenInit{},
	}
}

// 质询请求初始化
func (c *Client) NewSessionSetupRequest() (smb.SMB2SessionSetupRequestStruct, error) {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_SESSION_SETUP
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.SessionId = c.GetSessionId()

	ntlmsspneg := ntlm.NewNegotiate(c.GetOptions().Domain, c.GetOptions().Workstation)
	data, err := encoder.Marshal(ntlmsspneg)
	if err != nil {
		return smb.SMB2SessionSetupRequestStruct{}, err
	}

	if c.GetSessionId() != 0 {
		return smb.SMB2SessionSetupRequestStruct{}, errors.New("Bad session ID for session setup 1 message")
	}

	// Initial session setup request
	init, err := gss.NewNegTokenInit()
	if err != nil {
		return smb.SMB2SessionSetupRequestStruct{}, err
	}
	init.Data.MechToken = data

	return smb.SMB2SessionSetupRequestStruct{
		SMB2Header:           smb2Header,
		StructureSize:        25,
		Flags:                0x00,
		SecurityMode:         byte(smb.SecurityModeSigningEnabled),
		Capabilities:         0,
		Channel:              0,
		SecurityBufferOffset: 88,
		SecurityBufferLength: 0,
		PreviousSessionID:    0,
		SecurityBlob:         &init,
	}, nil
}

// 质询响应初始化
func NewSessionSetupResponse() (smb.SMB2SessionSetupResponseStruct, error) {
	smb2Header := NewSMB2Header()
	resp, err := gss.NewNegTokenResp()
	if err != nil {
		return smb.SMB2SessionSetupResponseStruct{}, err
	}
	ret := smb.SMB2SessionSetupResponseStruct{
		SMB2Header:   smb2Header,
		SecurityBlob: &resp,
	}
	return ret, nil
}

// 认证请求初始化
func (c *Client) NewSessionSetup2Request() (smb.SMB2SessionSetup2RequestStruct, error) {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_SESSION_SETUP
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.SessionId = c.GetSessionId()

	ntlmsspneg := ntlm.NewNegotiate(c.GetOptions().Domain, c.GetOptions().Workstation)
	data, err := encoder.Marshal(ntlmsspneg)
	if err != nil {
		return smb.SMB2SessionSetup2RequestStruct{}, err
	}

	if c.GetSessionId() == 0 {
		return smb.SMB2SessionSetup2RequestStruct{}, errors.New("Bad session ID for session setup 2 message")
	}

	// Session setup request #2
	resp, err := gss.NewNegTokenResp()
	if err != nil {
		return smb.SMB2SessionSetup2RequestStruct{}, err
	}
	resp.ResponseToken = data

	return smb.SMB2SessionSetup2RequestStruct{
		SMB2Header:           smb2Header,
		StructureSize:        25,
		Flags:                0x00,
		SecurityMode:         byte(smb.SecurityModeSigningEnabled),
		Capabilities:         0,
		Channel:              0,
		SecurityBufferOffset: 88,
		SecurityBufferLength: 0,
		PreviousSessionID:    0,
		SecurityBlob:         &resp,
	}, nil
}

func (c *Client) NegotiateProtocol() error {
	// 第一步 发送协商请求
	c.Debug("Sending Negotiate request", nil)
	negReq := c.NewNegotiateRequest()
	buf, err := c.Send(negReq)
	if err != nil {
		c.Debug("", err)
		return err
	}
	negRes := NewNegotiateResponse()
	if err = encoder.Unmarshal(buf, &negRes); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
		return err
	}
	if negRes.SMB2Header.Status != ms.STATUS_SUCCESS {
		return errors.New(fmt.Sprintf("NT Status Error: %d\n", negRes.SMB2Header.Status))
	}
	// Check SPNEGO security blob
	//spnegoOID, err := encoder.ObjectIDStrToInt(encoder.SpnegoOid)
	//if err != nil {
	//	c.Debug(err)
	//	return err
	//}
	//oid := negRes.SecurityBlob.OID
	//fmt.Println(oid)
	// 检查是否存在ntlmssp
	hasNTLMSSP := false
	ntlmsspOID, err := gss.ObjectIDStrToInt(ntlm.NTLMSSPMECHTYPEOID)
	if err != nil {
		return err
	}
	for _, mechType := range negRes.SecurityBlob.Data.MechTypes {
		if mechType.Equal(ntlmsspOID) {
			hasNTLMSSP = true
			break
		}
	}
	if !hasNTLMSSP {
		return errors.New("Server does not support NTLMSSP")
	}
	// 设置会话安全模式
	c.WithSecurityMode(negRes.SecurityMode)
	// 设置会话协议
	c.WithDialect(negRes.DialectRevision)
	// 签名开启/关闭
	mode := c.GetSecurityMode()
	if mode&smb.SecurityModeSigningEnabled > 0 {
		if mode&smb.SecurityModeSigningRequired > 0 {
			c.IsSigningRequired = true
		} else {
			c.IsSigningRequired = false
		}
	} else {
		c.IsSigningRequired = false
	}
	// 第二步 发送质询
	c.Debug("Sending SessionSetup1 request", nil)
	ssreq, err := c.NewSessionSetupRequest()
	if err != nil {
		c.Debug("", err)
		return err
	}
	ssres, err := NewSessionSetupResponse()
	if err != nil {
		c.Debug("", err)
		return err
	}
	buf, err = encoder.Marshal(ssreq)
	if err != nil {
		c.Debug("", err)
		return err
	}

	buf, err = c.Send(ssreq)
	if err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
		return err
	}

	c.Debug("Unmarshalling SessionSetup1 response", nil)
	if err = encoder.Unmarshal(buf, &ssres); err != nil {
		c.Debug("", err)
		return err
	}

	challenge := ntlm.NewChallenge()
	resp := ssres.SecurityBlob
	if err = encoder.Unmarshal(resp.ResponseToken, &challenge); err != nil {
		c.Debug("", err)
		return err
	}

	if ssres.SMB2Header.Status != ms.STATUS_MORE_PROCESSING_REQUIRED {
		status, _ := ms.StatusMap[negRes.SMB2Header.Status]
		return errors.New(fmt.Sprintf("NT Status Error: %c\n", status))
	}
	c.WithSessionId(ssres.SMB2Header.SessionId)

	c.Debug("Sending SessionSetup2 request", nil)
	// 第三步 认证
	ss2req, err := c.NewSessionSetup2Request()
	if err != nil {
		c.Debug("", err)
		return err
	}

	var auth ntlm.NTLMv2Authentication
	if c.GetOptions().Hash != "" {
		// Hash present, use it for auth
		c.Debug("Performing hash-based authentication", nil)
		auth = ntlm.NewAuthenticateHash(c.GetOptions().Domain, c.GetOptions().User, c.GetOptions().Workstation, c.GetOptions().Hash, challenge)
	} else {
		// No hash, use password
		c.Debug("Performing password-based authentication", nil)
		auth = ntlm.NewAuthenticatePass(c.GetOptions().Domain, c.GetOptions().User, c.GetOptions().Workstation, c.GetOptions().Password, challenge)
	}

	responseToken, err := encoder.Marshal(auth)
	if err != nil {
		c.Debug("", err)
		return err
	}
	resp2 := ss2req.SecurityBlob
	resp2.ResponseToken = responseToken
	ss2req.SecurityBlob = resp2
	ss2req.SMB2Header.Credits = 127
	buf, err = encoder.Marshal(ss2req)
	if err != nil {
		c.Debug("", err)
		return err
	}

	buf, err = c.Send(ss2req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Unmarshalling SessionSetup2 response", nil)
	var authResp smb.SMB2Header
	if err = encoder.Unmarshal(buf, &authResp); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
		return err
	}
	if authResp.Status != ms.STATUS_SUCCESS {
		// authResp.Status 十进制表示
		status, _ := ms.StatusMap[authResp.Status]
		return errors.New(fmt.Sprintf("NT Status Error: %c\n", status))
	}
	c.IsAuthenticated = true

	c.Debug("Completed NegotiateProtocol and SessionSetup", nil)
	return nil
}

// SMB2连接封装
func NewSession(opt common.ClientOptions, debug bool) (client *Client, err error) {
	address := fmt.Sprintf("%s:%d", opt.Host, opt.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	client = &Client{}
	client.WithOptions(&opt)
	client.WithConn(conn)
	client.WithDebug(debug)

	err = client.NegotiateProtocol()
	if err != nil {
		return
	}

	return client, nil
}

func (c *Client) Close() {
	c.Debug("Closing session", nil)
	trees := c.GetTrees()
	for k, _ := range trees {
		c.TreeDisconnect(k)
	}
	c.GetConn().Close()
	c.Debug("Session close completed", nil)
}
