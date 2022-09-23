package rpce

import (
	"bytes"
	"encoding/binary"

	"github.com/C-Sto/goWMIExec/pkg/uuid"
)

//https://pubs.opengroup.org/onlinepubs/9629399/toc.pdf

type hyper uint64
type long uint32
type short uint16
type small uint8

//12.6.3.1

/* start 8-octet aligned
u_int8  rpc_vers = 5;        /* 00:01 RPC version
u_int8  rpc_vers_minor;      /* 01:01 minor version
u_int8  PTYPE;               /* 02:01 packet type
u_int8  pfc_flags;           /* 03:01 flags (see PFC_... )
byte    packed_drep[4];   /* 04:04 NDR data representation format label
u_int16 frag_length;         /* 08:02 total length of fragment
u_int16 auth_length;         /* 10:02 length of auth_value
u_int32 call_id;             /* 12:04 call identifier
*/

//CommonHead appears in all PDU types. Alignment etc in comments. (page 590)
type CommonHead struct {
	Version            byte // 00:01 =5. RPC Version
	VersionMinor       byte
	PacketType         byte
	PFCFlags           PFCFlags
	DataRepresentation [4]byte
	FragLength         uint16
	AuthLength         uint16
	CallID             uint32
}

func NewCommonHeader(ptype byte, flags PFCFlags, callID uint32) CommonHead {
	return CommonHead{
		Version:            5,
		VersionMinor:       0, //shruggy guy emoji?
		PacketType:         ptype,
		PFCFlags:           flags,
		DataRepresentation: [4]byte{0x10, 0, 0, 0}, //PDRep seems hard?
		CallID:             callID,
	}
}

//
type PFCFlags byte

/*
#define PFC_FIRST_FRAG           0x01/* First fragment
#define PFC_LAST_FRAG            0x02/* Last fragment
#define PFC_PENDING_CANCEL       0x04/* Cancel was pending at sender
#define PFC_RESERVED_1           0x08#define PFC_CONC_MPX             0x10/* supports concurrent multiplexing* of a single connection.
#define PFC_DID_NOT_EXECUTE      0x20/* only meaningful on ‘fault’ packet;* if true, guaranteed call did not* execute.
#define PFC_MAYBE                0x40/* ‘maybe’ call semantics requested
#define PFC_OBJECT_UUID          0x80/* if true, a non-nil object UUID* was specified in the handle, and* is present in the optional object* field. If false, the object field* is omitted.
*/

//
const (
	PFCFirstFrag PFCFlags = (iota + 1)
	PFCLastFrag  PFCFlags = iota << 1
	PFCCancelPending
	PFCReserved
	PFCMultiplex
	PFCDidNotExecute
	PFCMaybe
	PFCObject
)

type ContextID uint16

type PSyntaxID struct {
	UUID    uuid.UUID
	Version int32
}

func NewPSyntaxID(uuid uuid.UUID, version int32) PSyntaxID {
	return PSyntaxID{
		UUID:    uuid,
		Version: version,
	}
}

type PContextElem struct {
	ContextID     uint16 // unclear if this should be 0 indexed per context?
	NumTransItems byte
	Unkown2       byte
	AbSyntax      PSyntaxID
	TfSyntax      []PSyntaxID //supposedly this can hold more than 1?
	//Interface         [16]byte
	//InterfaceVers     uint16
	//InterfaceVerMinor uint16
	//TransferSyntax    [16]byte
	//TransferSyntaxVer uint32
}

func NewPcontextElem(id uint16, abSyntax PSyntaxID, tfSyntax []PSyntaxID) PContextElem {
	r := PContextElem{}
	r.ContextID = id
	r.NumTransItems = byte(len(tfSyntax))
	r.AbSyntax = abSyntax
	r.TfSyntax = tfSyntax
	return r
}

func (p PContextElem) Bytes() []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.LittleEndian, p.ContextID)
	binary.Write(&buff, binary.LittleEndian, p.NumTransItems)
	binary.Write(&buff, binary.LittleEndian, p.Unkown2)
	binary.Write(&buff, binary.LittleEndian, p.AbSyntax)
	for _, x := range p.TfSyntax {
		binary.Write(&buff, binary.LittleEndian, x)
	}
	return buff.Bytes()
}

type PContextList struct {
	NumContexts   byte
	Reserved      byte
	Reserved2     uint16
	PContextElems []PContextElem
}

func NewPcontextList() *PContextList {
	return &PContextList{}
}

func (p *PContextList) AddContext(c PContextElem) {
	p.NumContexts++
	p.PContextElems = append(p.PContextElems, c)
}

func (p PContextList) Bytes() []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.LittleEndian, p.NumContexts)
	binary.Write(&buff, binary.LittleEndian, p.Reserved)
	binary.Write(&buff, binary.LittleEndian, p.Reserved2)
	for _, x := range p.PContextElems {
		binary.Write(&buff, binary.LittleEndian, x.Bytes())
	}
	return buff.Bytes()
}

//SizeOf returns the total size in bytes of the Pcontextlist to be put on the wire
func (p PContextList) SizeOf() int {
	/*
	   //arguably, this is the most efficient way... but I don't trust it..
	   	accum := 4 //header

	   	for _, x := range p.PContextElems {
	   		accum += 4                           //header
	   		accum += 20                          //abstract syntax
	   		accum += (int(x.NumTransItems) * 20) //tfsyntax
	   	}

	   	return accum
	*/
	return len(p.Bytes())
}

type NegResult short

const (
	NegResultAccept NegResult = iota
	NegResultUserReject
	NegResultProviderReject
)

type NegRejectReason short

const (
	NegRejectNotSpecified NegRejectReason = iota
	NegRejectAbSyntaxNotSpecified
	NegRejectProposedTransferSyntaxNotSupported
	NegRejectLocalLimitExceeded
)

type PResult struct {
	NegResult       NegResult
	NegRejectReason NegRejectReason
	TransferSyntax  PSyntaxID //should be 0's if result not accepted
}

type PResultList struct {
	NResults  byte
	Reserved  byte
	Reserved2 uint16
	PResults  []PResult
}

type Version struct {
	Major uint8
	Minor uint8
}

type PortAny struct {
	Length   uint16
	PortSpec []byte
}
