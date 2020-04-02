// Package sftp implements the SSH File Transfer Protocol as described in
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-02
package sftp

import (
	"fmt"
)

const (
	ssh_FXP_INIT     = 1
	ssh_FXP_VERSION  = 2
	ssh_FXP_OPEN     = 3
	ssh_FXP_CLOSE    = 4
	ssh_FXP_READ     = 5
	ssh_FXP_LSTAT    = 7
	ssh_FXP_FSTAT    = 8
	ssh_FXP_OPENDIR  = 11
	ssh_FXP_READDIR  = 12
	ssh_FXP_REALPATH = 16
	ssh_FXP_STAT     = 17
	ssh_FXP_STATUS   = 101
	ssh_FXP_HANDLE   = 102
	ssh_FXP_DATA     = 103
	ssh_FXP_NAME     = 104
	ssh_FXP_ATTRS    = 105
)

const (
	ssh_FX_OK                = 0
	ssh_FX_EOF               = 1
	ssh_FX_NO_SUCH_FILE      = 2
	ssh_FX_PERMISSION_DENIED = 3
	ssh_FX_FAILURE           = 4
	ssh_FX_BAD_MESSAGE       = 5
	ssh_FX_NO_CONNECTION     = 6
	ssh_FX_CONNECTION_LOST   = 7
	ssh_FX_OP_UNSUPPORTED    = 8
)

const (
	ssh_FXF_READ   = 0x00000001
	ssh_FXF_WRITE  = 0x00000002
	ssh_FXF_APPEND = 0x00000004
	ssh_FXF_CREAT  = 0x00000008
	ssh_FXF_TRUNC  = 0x00000010
	ssh_FXF_EXCL   = 0x00000020
)

type fxp uint8

func (f fxp) String() string {
	switch f {
	case ssh_FXP_INIT:
		return "SSH_FXP_INIT"
	case ssh_FXP_VERSION:
		return "SSH_FXP_VERSION"
	case ssh_FXP_OPEN:
		return "SSH_FXP_OPEN"
	case ssh_FXP_CLOSE:
		return "SSH_FXP_CLOSE"
	case ssh_FXP_READ:
		return "SSH_FXP_READ"
	case ssh_FXP_LSTAT:
		return "SSH_FXP_LSTAT"
	case ssh_FXP_FSTAT:
		return "SSH_FXP_FSTAT"
	case ssh_FXP_OPENDIR:
		return "SSH_FXP_OPENDIR"
	case ssh_FXP_READDIR:
		return "SSH_FXP_READDIR"
	case ssh_FXP_REALPATH:
		return "SSH_FXP_REALPATH"
	case ssh_FXP_STAT:
		return "SSH_FXP_STAT"
	case ssh_FXP_STATUS:
		return "SSH_FXP_STATUS"
	case ssh_FXP_HANDLE:
		return "SSH_FXP_HANDLE"
	case ssh_FXP_DATA:
		return "SSH_FXP_DATA"
	case ssh_FXP_NAME:
		return "SSH_FXP_NAME"
	case ssh_FXP_ATTRS:
		return "SSH_FXP_ATTRS"
	default:
		return "unknown"
	}
}

type fx uint8

func (f fx) String() string {
	switch f {
	case ssh_FX_OK:
		return "SSH_FX_OK"
	case ssh_FX_EOF:
		return "SSH_FX_EOF"
	case ssh_FX_NO_SUCH_FILE:
		return "SSH_FX_NO_SUCH_FILE"
	case ssh_FX_PERMISSION_DENIED:
		return "SSH_FX_PERMISSION_DENIED"
	case ssh_FX_FAILURE:
		return "SSH_FX_FAILURE"
	case ssh_FX_BAD_MESSAGE:
		return "SSH_FX_BAD_MESSAGE"
	case ssh_FX_NO_CONNECTION:
		return "SSH_FX_NO_CONNECTION"
	case ssh_FX_CONNECTION_LOST:
		return "SSH_FX_CONNECTION_LOST"
	case ssh_FX_OP_UNSUPPORTED:
		return "SSH_FX_OP_UNSUPPORTED"
	default:
		return "unknown"
	}
}

// A StatusError is returned when an SFTP operation fails, and provides
// additional information about the failure.
type StatusError struct {
	Code      uint32
	msg, lang string
}

func (s *StatusError) Error() string { return fmt.Sprintf("sftp: %q (%v)", s.msg, fx(s.Code)) }

func debug(fmt string, args ...interface{}) {}
