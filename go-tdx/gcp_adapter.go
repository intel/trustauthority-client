//go:build !test

/*
 *   Copyright (c) 2023 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */
package tdx

import (
	"crypto/sha512"
	"encoding/json"
	"syscall"
	"unsafe"

	"github.com/intel/trustauthority-client/go-connector"
	"github.com/pkg/errors"
)

// GCPAdapter manages TDX Quote collection from GCP TDX platform
type gcpAdapter struct {
	uData       []byte
	EvLogParser EventLogParser
}

// NewEvidenceAdapter returns a new GCP Adapter instance
func NewEvidenceAdapter(udata []byte, evLogParser EventLogParser) (connector.EvidenceAdapter, error) {
	return &gcpAdapter{
		uData:       udata,
		EvLogParser: evLogParser,
	}, nil
}

type TdxReportRequest struct {
	ReportData [TdxReportDataLen]byte
	TdReport   [TdxReportLen]byte
}

type TdxQuoteHeader struct {
	Version uint64
	Status  uint64
	InLen   uint32
	OutLen  uint32
	Data    [ReqBufSize]byte
}

type TdxQuoteRequest struct {
	Buf *TdxQuoteHeader
	Len uint64
}

// TdxQuoteReqABI is Linux's tdx-guest ABI for quote response
type TdxQuoteReqABI struct {
	Buffer unsafe.Pointer
	Length uint64
}

func IOC(dir, t, nr, size uintptr) uintptr {
	return (dir << IocDirshift) |
		(t << IocTypeShift) |
		(nr << IocNrShift) |
		(size << IocSizeShift)
}

func IOR(t, nr, size uintptr) uintptr {
	return IOC(IocRead, t, nr, size)
}

func IOWR(t, nr, size uintptr) uintptr {
	return IOC(IocWrite|IocRead, t, nr, size)
}

func TdxCmdGetReportIO() uintptr {
	return IOWR('T', 1, unsafe.Sizeof(TdxReportRequest{}))
}

func TdxCmdGetQuoteIO() uintptr {
	return IOWR('T', 2, unsafe.Sizeof(TdxQuoteReqABI{}))
}

// CollectEvidence is used to get TDX quote using GCP Quote Generation service
func (adapter *gcpAdapter) CollectEvidence(nonce []byte) (*connector.Evidence, error) {

	hash := sha512.New()
	_, err := hash.Write(nonce)
	if err != nil {
		return nil, err
	}
	_, err = hash.Write(adapter.uData)
	if err != nil {
		return nil, err
	}
	reportData := hash.Sum(nil)

	var tdrequest TdxReportRequest
	copy(tdrequest.ReportData[:], []byte(reportData))

	fd, err := syscall.Open(TdxAttestDevPath, syscall.O_RDWR|syscall.O_SYNC, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(fd)

	cmd := TdxCmdGetReportIO()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), cmd, uintptr(unsafe.Pointer(&tdrequest)))
	if errno != 0 {
		return nil, syscall.Errno(errno)
	}

	report := make([]byte, TdxReportLen)
	copy(report, tdrequest.TdReport[:])

	tdxQuoteHdr := &TdxQuoteHeader{
		Status:  0,
		Version: 1,
		InLen:   TdxReportLen,
		OutLen:  0,
	}
	copy(tdxQuoteHdr.Data[:], report)

	tdxrequest := &TdxQuoteRequest{
		Buf: tdxQuoteHdr,
		Len: ReqBufSize,
	}

	tdxrequestabi := TdxQuoteReqABI{
		Buffer: unsafe.Pointer(unsafe.Pointer(tdxrequest.Buf)),
		Length: tdxrequest.Len,
	}

	cmd = TdxCmdGetQuoteIO()
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), cmd, uintptr(unsafe.Pointer(&tdxrequestabi)))
	if errno != 0 {
		return nil, syscall.Errno(errno)
	}

	quote := make([]byte, tdxQuoteHdr.OutLen)
	copy(quote, tdxQuoteHdr.Data[:])

	var eventLog []byte
	if adapter.EvLogParser != nil {
		rtmrEventLogs, err := adapter.EvLogParser.GetEventLogs()
		if err != nil {
			return nil, errors.Wrap(err, "There was an error while collecting RTMR Event Log Data")
		}

		eventLog, err = json.Marshal(rtmrEventLogs)
		if err != nil {
			return nil, errors.Wrap(err, "Error while marshalling RTMR Event Log Data")
		}
	}

	return &connector.Evidence{
		Type:     1,
		Evidence: quote,
		UserData: adapter.uData,
		EventLog: eventLog,
	}, nil
}
