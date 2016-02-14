// See the file LICENSE for copyright and licensing information.

// Derived from FUSE's fuse_kernel.h
/*
   This file defines the kernel interface of FUSE
   Copyright (C) 2001-2007  Miklos Szeredi <miklos@szeredi.hu>


   This -- and only this -- header file may also be distributed under
   the terms of the BSD Licence as follows:

   Copyright (C) 2001-2007 Miklos Szeredi. All rights reserved.

   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:
   1. Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
   2. Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.

   THIS SOFTWARE IS PROVIDED BY AUTHOR AND CONTRIBUTORS ``AS IS'' AND
   ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
   IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
   ARE DISCLAIMED.  IN NO EVENT SHALL AUTHOR OR CONTRIBUTORS BE LIABLE
   FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
   DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
   OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
   HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
   LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
   OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
   SUCH DAMAGE.
*/

package fuse

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// Version is the FUSE version implemented by the package.
const Version = "7.8"

const (
	kernelVersion      = 7
	kernelMinorVersion = 8
	rootID             = 1
)

struct kstatfs {
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Bsize   uint32
	Namelen uint32
	Frsize  uint32
	Padding uint32
	Spare   [6]uint32
}

struct fileLock {
	Start uint64
	End   uint64
	Type  uint32
	Pid   uint32
}

// The SetattrValid are bit flags describing which fields in the SetattrRequest
// are included in the change.
type SetattrValid uint32

const (
	SetattrMode   SetattrValid = 1 << 0
	SetattrUid    SetattrValid = 1 << 1
	SetattrGid    SetattrValid = 1 << 2
	SetattrSize   SetattrValid = 1 << 3
	SetattrAtime  SetattrValid = 1 << 4
	SetattrMtime  SetattrValid = 1 << 5
	SetattrHandle SetattrValid = 1 << 6

	// Linux only(?)
	SetattrAtimeNow  SetattrValid = 1 << 7
	SetattrMtimeNow  SetattrValid = 1 << 8
	SetattrLockOwner SetattrValid = 1 << 9 // http://www.mail-archive.com/git-commits-head@vger.kernel.org/msg27852.html

	// OS X only
	SetattrCrtime   SetattrValid = 1 << 28
	SetattrChgtime  SetattrValid = 1 << 29
	SetattrBkuptime SetattrValid = 1 << 30
	SetattrFlags    SetattrValid = 1 << 31
)

func (fl SetattrValid) Mode() bool      { return fl&SetattrMode != 0 }
func (fl SetattrValid) Uid() bool       { return fl&SetattrUid != 0 }
func (fl SetattrValid) Gid() bool       { return fl&SetattrGid != 0 }
func (fl SetattrValid) Size() bool      { return fl&SetattrSize != 0 }
func (fl SetattrValid) Atime() bool     { return fl&SetattrAtime != 0 }
func (fl SetattrValid) Mtime() bool     { return fl&SetattrMtime != 0 }
func (fl SetattrValid) Handle() bool    { return fl&SetattrHandle != 0 }
func (fl SetattrValid) AtimeNow() bool  { return fl&SetattrAtimeNow != 0 }
func (fl SetattrValid) MtimeNow() bool  { return fl&SetattrMtimeNow != 0 }
func (fl SetattrValid) LockOwner() bool { return fl&SetattrLockOwner != 0 }
func (fl SetattrValid) Crtime() bool    { return fl&SetattrCrtime != 0 }
func (fl SetattrValid) Chgtime() bool   { return fl&SetattrChgtime != 0 }
func (fl SetattrValid) Bkuptime() bool  { return fl&SetattrBkuptime != 0 }
func (fl SetattrValid) Flags() bool     { return fl&SetattrFlags != 0 }

func (fl SetattrValid) String() string {
	return flagString(uint32(fl), setattrValidNames)
}

var setattrValidNames = []flagName{
	{uint32(SetattrMode), "SetattrMode"},
	{uint32(SetattrUid), "SetattrUid"},
	{uint32(SetattrGid), "SetattrGid"},
	{uint32(SetattrSize), "SetattrSize"},
	{uint32(SetattrAtime), "SetattrAtime"},
	{uint32(SetattrMtime), "SetattrMtime"},
	{uint32(SetattrHandle), "SetattrHandle"},
	{uint32(SetattrAtimeNow), "SetattrAtimeNow"},
	{uint32(SetattrMtimeNow), "SetattrMtimeNow"},
	{uint32(SetattrLockOwner), "SetattrLockOwner"},
	{uint32(SetattrCrtime), "SetattrCrtime"},
	{uint32(SetattrChgtime), "SetattrChgtime"},
	{uint32(SetattrBkuptime), "SetattrBkuptime"},
	{uint32(SetattrFlags), "SetattrFlags"},
}

// OpenFlags are the O_FOO flags passed to open/create/etc calls. For
// example, os.O_WRONLY | os.O_APPEND.
type OpenFlags uint32

func (fl OpenFlags) String() string {
	// O_RDONLY, O_RWONLY, O_RDWR are not flags
	s := accModeName(uint32(fl) & syscall.O_ACCMODE)
	flags := uint32(fl) &^ syscall.O_ACCMODE
	if flags != 0 {
		s = s + "+" + flagString(flags, openFlagNames)
	}
	return s
}

func accModeName(flags uint32) string {
	switch flags {
	case uint32(os.O_RDONLY):
		return "O_RDONLY"
	case uint32(os.O_WRONLY):
		return "O_WRONLY"
	case uint32(os.O_RDWR):
		return "O_RDWR"
	default:
		return ""
	}
}

var openFlagNames = []flagName{
	{uint32(os.O_CREATE), "O_CREATE"},
	{uint32(os.O_EXCL), "O_EXCL"},
	{uint32(os.O_TRUNC), "O_TRUNC"},
	{uint32(os.O_APPEND), "O_APPEND"},
	{uint32(os.O_SYNC), "O_SYNC"},
}

// The OpenResponseFlags are returned in the OpenResponse.
type OpenResponseFlags uint32

const (
	OpenDirectIO    OpenResponseFlags = 1 << 0 // bypass page cache for this open file
	OpenKeepCache   OpenResponseFlags = 1 << 1 // don't invalidate the data cache on open
	OpenNonSeekable OpenResponseFlags = 1 << 2 // (Linux?)

	OpenPurgeAttr OpenResponseFlags = 1 << 30 // OS X
	OpenPurgeUBC  OpenResponseFlags = 1 << 31 // OS X
)

func (fl OpenResponseFlags) String() string {
	return flagString(uint32(fl), openResponseFlagNames)
}

var openResponseFlagNames = []flagName{
	{uint32(OpenDirectIO), "OpenDirectIO"},
	{uint32(OpenKeepCache), "OpenKeepCache"},
	{uint32(OpenPurgeAttr), "OpenPurgeAttr"},
	{uint32(OpenPurgeUBC), "OpenPurgeUBC"},
}

// The InitFlags are used in the Init exchange.
type InitFlags uint32

const (
	InitAsyncRead  InitFlags = 1 << 0
	InitPosixLocks InitFlags = 1 << 1

	InitCaseSensitive InitFlags = 1 << 29 // OS X only
	InitVolRename     InitFlags = 1 << 30 // OS X only
	InitXtimes        InitFlags = 1 << 31 // OS X only
)

struct flagName {
	bit  uint32
	name string
}

var initFlagNames = []flagName{
	{uint32(InitAsyncRead), "InitAsyncRead"},
	{uint32(InitPosixLocks), "InitPosixLocks"},
	{uint32(InitCaseSensitive), "InitCaseSensitive"},
	{uint32(InitVolRename), "InitVolRename"},
	{uint32(InitXtimes), "InitXtimes"},
}

func (fl InitFlags) String() string {
	return flagString(uint32(fl), initFlagNames)
}

func flagString(f uint32, names []flagName) string {
	var s string

	if f == 0 {
		return "0"
	}

	for _, n := range names {
		if f&n.bit != 0 {
			s += "+" + n.name
			f &^= n.bit
		}
	}
	if f != 0 {
		s += fmt.Sprintf("%+#x", f)
	}
	return s[1:]
}

// The ReleaseFlags are used in the Release exchange.
type ReleaseFlags uint32

const (
	ReleaseFlush ReleaseFlags = 1 << 0
)

func (fl ReleaseFlags) String() string {
	return flagString(uint32(fl), releaseFlagNames)
}

var releaseFlagNames = []flagName{
	{uint32(ReleaseFlush), "ReleaseFlush"},
}

// Opcodes
const (
	opLookup      = 1
	opForget      = 2 // no reply
	opGetattr     = 3
	opSetattr     = 4
	opReadlink    = 5
	opSymlink     = 6
	opMknod       = 8
	opMkdir       = 9
	opUnlink      = 10
	opRmdir       = 11
	opRename      = 12
	opLink        = 13
	opOpen        = 14
	opRead        = 15
	opWrite       = 16
	opStatfs      = 17
	opRelease     = 18
	opFsync       = 20
	opSetxattr    = 21
	opGetxattr    = 22
	opListxattr   = 23
	opRemovexattr = 24
	opFlush       = 25
	opInit        = 26
	opOpendir     = 27
	opReaddir     = 28
	opReleasedir  = 29
	opFsyncdir    = 30
	opGetlk       = 31
	opSetlk       = 32
	opSetlkw      = 33
	opAccess      = 34
	opCreate      = 35
	opInterrupt   = 36
	opBmap        = 37
	opDestroy     = 38
	opIoctl       = 39 // Linux?
	opPoll        = 40 // Linux?

	// OS X
	opSetvolname = 61
	opGetxtimes  = 62
	opExchange   = 63
)

// The read buffer is required to be at least 8k but may be much larger
const minReadBuffer = 8192

struct entryOut {
	outHeader
	Nodeid         uint64 // Inode ID
	Generation     uint64 // Inode generation
	EntryValid     uint64 // Cache timeout for the name
	AttrValid      uint64 // Cache timeout for the attributes
	EntryValidNsec uint32
	AttrValidNsec  uint32
	Attr           attr
}

struct forgetIn {
	Nlookup uint64
}

struct attrOut {
	outHeader
	AttrValid     uint64 // Cache timeout for the attributes
	AttrValidNsec uint32
	Dummy         uint32
	Attr          attr
}

// OS X
struct getxtimesOut {
	outHeader
	Bkuptime     uint64
	Crtime       uint64
	BkuptimeNsec uint32
	CrtimeNsec   uint32
}

struct mknodIn {
	Mode uint32
	Rdev uint32
	// "filename\x00" follows.
}

struct mkdirIn {
	Mode    uint32
	Padding uint32
	// filename follows
}

struct renameIn {
	Newdir uint64
	// "oldname\x00newname\x00" follows
}

// OS X
struct exchangeIn {
	Olddir  uint64
	Newdir  uint64
	Options uint64
}

struct linkIn {
	Oldnodeid uint64
}

struct setattrInCommon {
	Valid     uint32
	Padding   uint32
	Fh        uint64
	Size      uint64
	LockOwner uint64 // unused on OS X?
	Atime     uint64
	Mtime     uint64
	Unused2   uint64
	AtimeNsec uint32
	MtimeNsec uint32
	Unused3   uint32
	Mode      uint32
	Unused4   uint32
	Uid       uint32
	Gid       uint32
	Unused5   uint32
}

struct openIn {
	Flags  uint32
	Unused uint32
}

struct openOut {
	outHeader
	Fh        uint64
	OpenFlags uint32
	Padding   uint32
}

struct createIn {
	Flags uint32
	Mode  uint32
}

struct createOut {
	outHeader

	Nodeid         uint64 // Inode ID
	Generation     uint64 // Inode generation
	EntryValid     uint64 // Cache timeout for the name
	AttrValid      uint64 // Cache timeout for the attributes
	EntryValidNsec uint32
	AttrValidNsec  uint32
	Attr           attr

	Fh        uint64
	OpenFlags uint32
	Padding   uint32
}

struct releaseIn {
	Fh           uint64
	Flags        uint32
	ReleaseFlags uint32
	LockOwner    uint32
}

struct flushIn {
	Fh         uint64
	FlushFlags uint32
	Padding    uint32
	LockOwner  uint64
}

struct readIn {
	Fh      uint64
	Offset  uint64
	Size    uint32
	Padding uint32
}

struct writeIn {
	Fh         uint64
	Offset     uint64
	Size       uint32
	WriteFlags uint32
}

struct writeOut {
	outHeader
	Size    uint32
	Padding uint32
}

// The WriteFlags are returned in the WriteResponse.
type WriteFlags uint32

func (fl WriteFlags) String() string {
	return flagString(uint32(fl), writeFlagNames)
}

var writeFlagNames = []flagName{}

const compatStatfsSize = 48

struct statfsOut {
	outHeader
	St kstatfs
}

struct fsyncIn {
	Fh         uint64
	FsyncFlags uint32
	Padding    uint32
}

struct setxattrInCommon {
	Size  uint32
	Flags uint32
}

func (setxattrInCommon) position() uint32 {
	return 0
}

struct getxattrInCommon {
	Size    uint32
	Padding uint32
}

func (getxattrInCommon) position() uint32 {
	return 0
}

struct getxattrOut {
	outHeader
	Size    uint32
	Padding uint32
}

struct lkIn {
	Fh    uint64
	Owner uint64
	Lk    fileLock
}

struct lkOut {
	outHeader
	Lk fileLock
}

struct accessIn {
	Mask    uint32
	Padding uint32
}

struct initIn {
	Major        uint32
	Minor        uint32
	MaxReadahead uint32
	Flags        uint32
}

const initInSize = int(unsafe.Sizeof(initIn{}))

struct initOut {
	outHeader
	Major        uint32
	Minor        uint32
	MaxReadahead uint32
	Flags        uint32
	Unused       uint32
	MaxWrite     uint32
}

struct interruptIn {
	Unique uint64
}

struct bmapIn {
	Block     uint64
	BlockSize uint32
	Padding   uint32
}

struct bmapOut {
	outHeader
	Block uint64
}

struct inHeader {
	Len     uint32
	Opcode  uint32
	Unique  uint64
	Nodeid  uint64
	Uid     uint32
	Gid     uint32
	Pid     uint32
	Padding uint32
}

const inHeaderSize = int(unsafe.Sizeof(inHeader{}))

struct outHeader {
	Len    uint32
	Error  int32
	Unique uint64
}

struct dirent {
	Ino     uint64
	Off     uint64
	Namelen uint32
	Type    uint32
	Name    [0]byte
}

const direntSize = 8 + 8 + 4 + 4
