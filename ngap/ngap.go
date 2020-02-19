packag ngap

import (
	"../encoding/per"
	"fmt"
)

const (
	reject = iota
	ignore
	notify
)

const (
	initiatingMessage = iota
	sucessfulOutcome
	unsuccessfulOutcome
)

// Elementary Procedures constants
const (
	procCodeInitialUEMessage = 15
	procCodeNGSetup          = 21
)

const (
	idDefaultPagingDRX = 21
	idGlobalRANNodeID  = 27
    idSupportedTAList  = 102
)

const (
	globalGNB = iota
	globalNgGNB
	globalN3IWF
)

/*
NGAP-PDU ::= CHOICE {
    initiatingMessage           InitiatingMessage,
    successfulOutcome           SuccessfulOutcome,
    unsuccessfulOutcome         UnsuccessfulOutcome,
    ...
}
ProcedureCode ::= INTEGER (0..255)
Criticality   ::= ENUMERATED { reject, ignore, notify }

InitiatingMessage ::= SEQUENCE {
    procedureCode   NGAP-ELEMENTARY-PROCEDURE.&procedureCode        ({NGAP-ELEMENTARY-PROCEDURES}),
    criticality     NGAP-ELEMENTARY-PROCEDURE.&criticality          ({NGAP-ELEMENTARY-PROCEDURES}{@procedureCode}),
    value           NGAP-ELEMENTARY-PROCEDURE.&InitiatingMessage    ({NGAP-ELEMENTARY-PROCEDURES}{@procedureCode})
}
*/
func encNgapPdu(pduType, procCode, criticality int) (pdu []uint8) {
	pdu, _, _ = per.EncChoice(pduType, 0, 2, true)
	v, _, _ := per.EncInteger(procCode, 0, 255, false)
	pdu = append(pdu, v...)
	v, _, _ = per.EncEnumerated(criticality, 0, 2, false)
	pdu = append(pdu, v...)

	return
}

/*
NGSetupRequest ::= SEQUENCE {
    protocolIEs     ProtocolIE-Container        { {NGSetupRequestIEs} },
    ...
}
ProtocolIE-Container {NGAP-PROTOCOL-IES : IEsSetParam} ::= 
    SEQUENCE (SIZE (0..maxProtocolIEs)) OF
    ProtocolIE-Field {{IEsSetParam}}

maxProtocolIEs                          INTEGER ::= 65535
*/
func encProtocolIEContainer(num int) (container []uint8) {
	const maxProtocolIEs = 65535
	container, _, _ = per.EncSequence(true, 0, 0)
	v, _, _ := per.EncSequenceOf(num, 0, maxProtocolIEs, false)
	container = append(container, v...)

	return
}

// 9.2.6.1 NG SETUP REQUEST
/*
NGSetupRequestIEs NGAP-PROTOCOL-IES ::= {
    { ID id-GlobalRANNodeID         CRITICALITY reject  TYPE GlobalRANNodeID            PRESENCE mandatory  }|
    { ID id-RANNodeName             CRITICALITY ignore  TYPE RANNodeName                PRESENCE optional}|
    { ID id-SupportedTAList         CRITICALITY reject  TYPE SupportedTAList            PRESENCE mandatory  }|
    { ID id-DefaultPagingDRX        CRITICALITY ignore  TYPE PagingDRX                  PRESENCE mandatory  }|
    { ID id-UERetentionInformation  CRITICALITY ignore  TYPE UERetentionInformation     PRESENCE optional   },
    ...
}
*/
func MakeNGSetupRequest() {
	pdu := encNgapPdu(initiatingMessage, procCodeNGSetup, reject)
	fmt.Printf("result: pdu = %02x\n", pdu)

	v := encProtocolIEContainer(3)
	fmt.Printf("result: ie container = %02x\n", v)

	v, _ = encGlobalRANNodeID()
	fmt.Printf("result: global RAN Node ID = %02x\n", v)

	v, _ = encSupportedTAList()
	fmt.Printf("result: Supported TA List = %02x\n", v)
}

/*
BroadcastPLMNList ::= SEQUENCE (SIZE(1..maxnoofBPLMNs)) OF BroadcastPLMNItem
    maxnoofBPLMNs                       INTEGER ::= 12
 */
func encBroadcastPLMNList() (v []uint8) {
	const maxnoofBPLMNs = 12
	v, _, _ = per.EncSequenceOf(1, 1, maxnoofBPLMNs, false)
	v = append(v, encBroadcastPLMNItem()...)
	return
}

/*
BroadcastPLMNItem ::= SEQUENCE {
    pLMNIdentity            PLMNIdentity,
    tAISliceSupportList     SliceSupportList,
    iE-Extensions           ProtocolExtensionContainer { {BroadcastPLMNItem-ExtIEs} } OPTIONAL,
    ...
}
 */
func encBroadcastPLMNItem() (v []uint8) {
	v, _, _ = per.EncSequence(true, 1, 0)
	v = append(v, encPLMNIdentity(123, 45)...)
	v = append(v, encSliceSupportList()...)
	return
}

func encProtocolIE(id, criticality int) (v []uint8, err error) {

	v1, _, _ := per.EncInteger(id, 0, 65535, false)
	v2, _, _ := per.EncEnumerated(criticality, 0, 2, false)
	v = append(v1, v2...)

	return
}

// 9.3.1.5 Global RAN Node ID
/*
  It returns only GNB-ID for now.
   GlobalRANNodeID ::= CHOICE {
       globalGNB-ID        GlobalGNB-ID,
       globalNgENB-ID      GlobalNgENB-ID,
       globalN3IWF-ID      GlobalN3IWF-ID,
       choice-Extensions   ProtocolIE-SingleContainer { {GlobalRANNodeID-ExtIEs} }
   }
 */
func encGlobalRANNodeID() (v []uint8, err error) {

	v, err = encProtocolIE(idGlobalRANNodeID, reject)

	// NG-ENB and N3IWF are not implemented yet...
	pv, plen, _ := per.EncChoice(globalGNB, 0, 2, false)
	pv2, plen2, v2 := encGlobalGNBId()
	pv, plen = per.MergeBitField(pv, plen, pv2, plen2)
	pv = append(pv, v2...)

	v3, _, _ := per.EncLengthDeterminant(len(pv), 0)
	v = append(v, v3...)
	v = append(v, pv...)

	return
}


// 9.3.1.6 Global gNB ID
/*
   GlobalGNB-ID ::= SEQUENCE {
       pLMNIdentity        PLMNIdentity,
       gNB-ID              GNB-ID,
       iE-Extensions       ProtocolExtensionContainer { {GlobalGNB-ID-ExtIEs} } OPTIONAL,
       ...
   }
 */
func encGlobalGNBId() (pv []uint8, plen int, v []uint8) {
	//temp value: MCC = 123, MNC = 45
	pv, plen, _ = per.EncSequence(true, 1, 0)
	v = append(v, encPLMNIdentity(123, 45)...)

	pv2, _ := encGNBId()
	v = append(v, pv2...)
	return
}

/*
   GNB-ID ::= CHOICE {
       gNB-ID                  BIT STRING (SIZE(22..32)),
       choice-Extensions       ProtocolIE-SingleContainer { {GNB-ID-ExtIEs} }
   }
 */
func encGNBId() (pv []uint8, plen int) {
	//GNB-ID = 1
	pv, plen, _ = per.EncChoice(0, 0, 1, false)
	pv2, plen2, _ := per.EncBitString([]uint8{0x00, 0x00, 0x01},
		22, 22, 32, false)
	pv, plen = per.MergeBitField(pv, plen, pv2, plen2)
	return
}

// 9.3.1.90 PagingDRX
/*
func encPagingDRX(drx string) (val []uint8) {
	n := 0
	switch drx {
	case "v32":
		n = 0
	case "v64":
		n = 1
	case "v128":
		n = 2
	case "v256":
		n = 3
	default:
		fmt.Printf("encPagingDRX: no such DRX value(%s)", drx)
		return
	}
	val = per.EncEnumerated(n)
	return
}
*/

// 9.3.3.5 PLMN Identity
/*
PLMNIdentity ::= OCTET STRING (SIZE(3)) 
 */
func encPLMNIdentity(mcc, mnc int) (v []uint8) {

	v = make([]uint8, 3, 3)
	v[0] = uint8(mcc % 1000 / 100)
	v[0] |= uint8(mcc%100/10) << 4

	v[1] = uint8(mcc % 10)
	v[1] |= 0xf0 // filler digit

	v[2] = uint8(mnc % 100 / 10)
	v[2] |= uint8(mnc%10) << 4

	_, _, v, _ = per.EncOctetString(v, 3, 3, false)

	return
}

/*
SliceSupportList ::= SEQUENCE (SIZE(1..maxnoofSliceItems)) OF SliceSupportItem
    maxnoofSliceItems                   INTEGER ::= 1024
 */
func encSliceSupportList() (v []uint8) {
	v, _, _ = per.EncSequenceOf(1, 1, 1024, false)
	v = append(v, encSliceSupportItem()...)
	return
}

/*
SliceSupportItem ::= SEQUENCE {
    s-NSSAI             S-NSSAI,
    iE-Extensions       ProtocolExtensionContainer { {SliceSupportItem-ExtIEs} }    OPTIONAL,
    ...
}
 */
func encSliceSupportItem() (v []uint8) {
	/*
	ex.1
	    .    .   .          .    .   .    .   .    .   .
	00 000 00000001	00 010 00000010 00000000 00000011 00001000
	0000 0000 0000 1xxx
	                000 1000 0000 10xx xxxx 00000000 00000011 11101000
	0x00 0x08 0x80 0x80 0x00 0x03 0xe1

	ex.2
	    .    .   .    .        .        .        .
	00 010 00000001 xxx 00000000 00000000 01111011
	0001 0000 0000 1xxx 00000000 00000000 11101000
	0x10 0x08 0x80 0x00 0x00 0x7b
	*/
	pv, plen, _ := per.EncSequence(true, 1, 0)

	pv2, plen2, v := encSNSSAI([]uint8{1}, []uint8{0, 0, 123})
	pv, plen = per.MergeBitField(pv, plen, pv2, plen2)
	v = append(pv, v...)
	return
}

// 9.3.1.24 S-NSSAI
/*
S-NSSAI ::= SEQUENCE {
    sST           SST,
    sD            SD                                                  OPTIONAL,
    iE-Extensions ProtocolExtensionContainer { { S-NSSAI-ExtIEs} }    OPTIONAL,
    ...
}

SST ::= OCTET STRING (SIZE(1))
SD ::= OCTET STRING (SIZE(3))
*/
func encSNSSAI(sst, sd []uint8) (pv []uint8, plen int, v []uint8) {
	pv, plen, _ = per.EncSequence(true, 2, 0x02)

	pv2, plen2, _, _ := per.EncOctetString(sst, 1, 1, false)

	pv, plen = per.MergeBitField(pv, plen, pv2, plen2)

	_, _, v, _ = per.EncOctetString(sd, 3, 3, false)
	return
}

// Supported TA List
/*
SupportedTAList ::= SEQUENCE (SIZE(1..maxnoofTACs)) OF SupportedTAItem
 */
func encSupportedTAList() (v []uint8, err error) {

	v, err = encProtocolIE(idSupportedTAList, reject)

	// maxnoofTACs INTEGER ::= 256
	const maxnoofTACs = 256
	pv, _, _ := per.EncSequenceOf(1, 1, maxnoofTACs, false)
	v = append(v, pv...)

	v = append(v, encSupportedTAItem()...)

	return
}

// Supported TA Item
/*
SupportedTAItem ::= SEQUENCE {
    tAC                     TAC,
    broadcastPLMNList       BroadcastPLMNList,
    iE-Extensions           ProtocolExtensionContainer { {SupportedTAItem-ExtIEs} } OPTIONAL,
    ...
}
 */
func encSupportedTAItem() (v []uint8) {
	//pv, plen, _ := per.EncSequence(true, 1, 0)

	//TAC
	tac := []uint8{0x00, 0x01, 0x02}
	v = encTAC(tac)

	//BroadcasePLMNList
	v = append(v, encBroadcastPLMNList()...)
	return
}

// 9.3.3.10 TAC
/*
TAC ::= OCTET STRING (SIZE(3))
 */
func encTAC(tac []uint8) (v []uint8) {
	const tacSize = 3
	_, _, v, _ = per.EncOctetString(tac, tacSize, tacSize, false)
	return
}

