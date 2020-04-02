package sftp

import (
	"encoding"
)

// all incoming packets
type requestPacket interface {
	encoding.BinaryUnmarshaler
	id() uint32
}

type responsePacket interface {
	encoding.BinaryMarshaler
	id() uint32
}

// interfaces to group types
type hasPath interface {
	requestPacket
	getPath() string
}

type hasHandle interface {
	requestPacket
	getHandle() string
}

type notReadOnly interface {
	notReadOnly()
}

//// define types by adding methods
// hasPath
func (p sshFxpLstatPacket) getPath() string    { return p.Path }
func (p sshFxpStatPacket) getPath() string     { return p.Path }
func (p sshFxpRealpathPacket) getPath() string { return p.Path }
func (p sshFxpOpendirPacket) getPath() string  { return p.Path }
func (p sshFxpOpenPacket) getPath() string     { return p.Path }

// getHandle
func (p sshFxpFstatPacket) getHandle() string   { return p.Handle }
func (p sshFxpReadPacket) getHandle() string    { return p.Handle }
func (p sshFxpReaddirPacket) getHandle() string { return p.Handle }
func (p sshFxpClosePacket) getHandle() string   { return p.Handle }

// some packets with ID are missing id()
func (p sshUnknownPacket) id() uint32   { return p.ID }
func (p sshFxpDataPacket) id() uint32   { return p.ID }
func (p sshFxpStatusPacket) id() uint32 { return p.ID }
func (p sshFxpStatResponse) id() uint32 { return p.ID }
func (p sshFxpNamePacket) id() uint32   { return p.ID }
func (p sshFxpHandlePacket) id() uint32 { return p.ID }
func (p sshFxVersionPacket) id() uint32 { return 0 }

// take raw incoming packet data and build packet objects
func makePacket(p rxPacket) (requestPacket, error) {
	var pkt requestPacket
	switch p.pktType {
	case ssh_FXP_INIT:
		pkt = &sshFxInitPacket{}
	case ssh_FXP_LSTAT:
		pkt = &sshFxpLstatPacket{}
	case ssh_FXP_OPEN:
		pkt = &sshFxpOpenPacket{}
	case ssh_FXP_CLOSE:
		pkt = &sshFxpClosePacket{}
	case ssh_FXP_READ:
		pkt = &sshFxpReadPacket{}
	case ssh_FXP_FSTAT:
		pkt = &sshFxpFstatPacket{}
	case ssh_FXP_OPENDIR:
		pkt = &sshFxpOpendirPacket{}
	case ssh_FXP_READDIR:
		pkt = &sshFxpReaddirPacket{}
	case ssh_FXP_REALPATH:
		pkt = &sshFxpRealpathPacket{}
	case ssh_FXP_STAT:
		pkt = &sshFxpStatPacket{}
	default:
		pkt = &sshUnknownPacket{}
	}
	if err := pkt.UnmarshalBinary(p.pktBytes); err != nil {
		// Return partially unpacked packet to allow callers to return
		// error messages appropriately with necessary id() method.
		return pkt, err
	}
	return pkt, nil
}
