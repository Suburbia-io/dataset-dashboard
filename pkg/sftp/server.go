package sftp

// sftp server counterpart

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"

	"github.com/Suburbia-io/dashboard/pkg/database"
)

const (
	sftpProtocolVersion   = 3 // http://tools.ietf.org/html/draft-ietf-secsh-filexfer-02
	SftpServerWorkerCount = 8
)

type OpenFile struct {
	file *os.File
	path string
}

// Server is an SSH File Transfer Protocol (sftp) server.
// This is intended to provide the sftp subsystem to an ssh server daemon.
// This implementation currently supports most of sftp server protocol version 3,
// as specified at http://tools.ietf.org/html/draft-ietf-secsh-filexfer-02
type Server struct {
	*serverConn
	debugStream   io.Writer
	readOnly      bool
	pktMgr        *packetManager
	openFiles     map[string]OpenFile
	rootFolder    string
	openFilesLock sync.RWMutex
	handleCount   int
	maxTxPacket   uint32
	DBAL          database.DBAL
	sessionId     string
}

func (svr *Server) nextHandle(file *os.File, path string) string {
	svr.openFilesLock.Lock()
	defer svr.openFilesLock.Unlock()
	svr.handleCount++
	handle := strconv.Itoa(svr.handleCount)
	svr.openFiles[handle] = OpenFile{
		file: file,
		path: path,
	}
	return handle
}

func (svr *Server) closeHandle(handle string) error {
	svr.openFilesLock.Lock()
	defer svr.openFilesLock.Unlock()
	if f, ok := svr.openFiles[handle]; ok {
		delete(svr.openFiles, handle)
		return f.file.Close()
	}

	return syscall.EBADF
}

func (svr *Server) getHandle(handle string) (OpenFile, bool) {
	svr.openFilesLock.RLock()
	defer svr.openFilesLock.RUnlock()
	f, ok := svr.openFiles[handle]
	return f, ok
}

// NewServer creates a new Server instance around the provided streams, serving
// content from the root of the filesystem.  Optionally, ServerOption
// functions may be specified to further configure the Server.
//
// A subsequent call to Serve() is required to begin serving files over SFTP.

func NewServer(rwc io.ReadWriteCloser, rootFolder string, sessionId string, DBAL database.DBAL, options ...ServerOption) (*Server, error) {
	if len(rootFolder) <= 0 {
		panic("No root folder for SFTP server provided. This is insecure and reveals access to server files!")
	}

	if len(sessionId) <= 0 {
		panic("No session id provided. Needed for folder permission auth.")
	}

	svrConn := &serverConn{
		conn: conn{
			Reader:      rwc,
			WriteCloser: rwc,
		},
	}

	s := &Server{
		serverConn:  svrConn,
		debugStream: ioutil.Discard,
		pktMgr:      newPktMgr(svrConn),
		openFiles:   make(map[string]OpenFile),
		maxTxPacket: 1 << 15,
		rootFolder:  rootFolder,
		sessionId:   sessionId,
		DBAL:        DBAL,
	}

	for _, o := range options {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// A ServerOption is a function which applies configuration to a Server.
type ServerOption func(*Server) error

// WithDebug enables Server debugging output to the supplied io.Writer.
func WithDebug(w io.Writer) ServerOption {
	return func(s *Server) error {
		s.debugStream = w
		return nil
	}
}

// ReadOnly configures a Server to serve files in read-only mode.
func ReadOnly() ServerOption {
	return func(s *Server) error {
		s.readOnly = true
		return nil
	}
}

type rxPacket struct {
	pktType  fxp
	pktBytes []byte
}

// Up to N parallel servers
func (svr *Server) sftpServerWorker(pktChan chan orderedRequest) error {
	for pkt := range pktChan {
		if err := handlePacket(svr, pkt); err != nil {
			return err
		}
	}
	return nil
}

func getAllowedPath(s *Server, path string) (realPath string, dataset *database.Dataset, err error) {
	f, err := filepath.Abs(path)
	path = cleanPath(f)
	if !strings.HasPrefix(path, s.rootFolder) {
		return "", nil, syscall.EPERM
	}

	pathWithoutRoot := strings.Replace(path, s.rootFolder, "", 1)
	if strings.HasPrefix(pathWithoutRoot, "/") {
		pathWithoutRoot = strings.Replace(pathWithoutRoot, "/", "", 1)
	}

	folderName := pathWithoutRoot
	if slashesCount := strings.IndexByte(pathWithoutRoot, '/'); slashesCount >= 0 {
		// get everything up until the first forward slash (without root folder path)
		// and consider that the folder name the user is trying to access
		folderName = pathWithoutRoot[:strings.IndexByte(pathWithoutRoot, '/')]
	}

	if folderName == "" {
		// user is accessing the root folder, that's ok!
		return path, nil, nil
	}

	canAccessDataset, dataset := canAccessDataset(s, folderName)
	if !canAccessDataset {
		return "", nil, syscall.EPERM
	}

	return path, dataset, nil
}

func canAccessDataset(s *Server, datasetName string) (result bool, dataset *database.Dataset) {
	sftpSession, err := s.DBAL.SftpSessionGetByToken(s.sessionId)
	if err != nil {
		return false, nil
	}
	customer, err := s.DBAL.CustomerGet(sftpSession.User.CustomerID)
	if err != nil {
		return false, nil
	}
	if customer.ArchivedAt != nil {
		return false, nil
	}

	ds, err := s.DBAL.DatasetGetBySlug(datasetName)
	if err != nil {
		return false, nil
	}

	_, err = s.DBAL.CustomerDatasetMappingGet(customer.CustomerID, ds.DatasetID)
	if err != nil {
		return false, nil
	}

	return true, &ds
}

func handlePacket(s *Server, p orderedRequest) error {
	var rpkt responsePacket
	switch p := p.requestPacket.(type) {
	case *sshFxInitPacket:
		rpkt = sshFxVersionPacket{
			Version: sftpProtocolVersion,
		}
	case *sshFxpStatPacket:
		// Get allowed path
		requestPath, _, err := getAllowedPath(s, p.Path)
		if err != nil {
			rpkt = statusFromError(p, syscall.EPERM)
			break
		}

		// stat the requested file
		info, err := os.Stat(requestPath)
		rpkt = sshFxpStatResponse{
			ID:   p.ID,
			info: info,
		}
		if err != nil {
			rpkt = statusFromError(p, err)
		}
	case *sshFxpLstatPacket:
		// Get allowed path
		requestPath, _, err := getAllowedPath(s, p.Path)
		if err != nil {
			rpkt = statusFromError(p, syscall.EPERM)
			break
		}

		// stat the requested file
		info, err := os.Lstat(requestPath)
		rpkt = sshFxpStatResponse{
			ID:   p.ID,
			info: info,
		}
		if err != nil {
			rpkt = statusFromError(p, err)
		}
	case *sshFxpFstatPacket:
		f, ok := s.getHandle(p.Handle)
		var err error = syscall.EBADF
		var info os.FileInfo
		if ok {
			info, err = f.file.Stat()
			rpkt = sshFxpStatResponse{
				ID:   p.ID,
				info: info,
			}
		}
		if err != nil {
			rpkt = statusFromError(p, err)
		}
	case *sshFxpClosePacket:
		rpkt = statusFromError(p, s.closeHandle(p.Handle))
	case *sshFxpRealpathPacket:
		if p.Path == "." {
			p.Path, _ = filepath.Abs(s.rootFolder)
			rpkt = sshFxpNamePacket{
				ID: p.ID,
				NameAttrs: []sshFxpNameAttr{{
					Name:     p.Path,
					LongName: p.Path,
					Attrs:    emptyFileStat,
				}},
			}
			break
		}

		// Get allowed path
		requestPath, _, err := getAllowedPath(s, p.Path)
		if err != nil {
			rpkt = statusFromError(p, syscall.EPERM)
			break
		}

		rpkt = sshFxpNamePacket{
			ID: p.ID,
			NameAttrs: []sshFxpNameAttr{{
				Name:     requestPath,
				LongName: requestPath,
				Attrs:    emptyFileStat,
			}},
		}
	case *sshFxpOpenPacket:
		// Get allowed path
		requestPath, dataset, err := getAllowedPath(s, p.Path)
		if err != nil {
			rpkt = statusFromError(p, syscall.EPERM)
			break
		}

		sftpSession, err := s.DBAL.SftpSessionGetByToken(s.sessionId)

		if err != nil {
			rpkt = statusFromError(p, syscall.EPERM)
			break
		}

		pathWithoutRoot := strings.Replace(requestPath, s.rootFolder, "", 1)
		defer s.DBAL.AuditTrailBySftpUserInsertAsync(sftpSession, tables.Datasets.Table(), dataset.DatasetID, "SftpFileDownload", pathWithoutRoot)

		file, err := os.OpenFile(requestPath, os.O_RDONLY, 0644)
		if err != nil {
			rpkt = statusFromError(p, err)
			break
		}

		handle := s.nextHandle(file, requestPath)
		rpkt = sshFxpHandlePacket{ID: p.id(), Handle: handle}
	case *sshFxpOpendirPacket:
		// Get allowed path
		requestPath, _, err := getAllowedPath(s, p.Path)
		if err != nil {
			rpkt = statusFromError(p, syscall.EPERM)
			break
		}

		if stat, err := os.Stat(requestPath); err != nil {
			rpkt = statusFromError(p, err)
		} else if !stat.IsDir() {
			rpkt = statusFromError(p, &os.PathError{
				Path: requestPath, Err: syscall.ENOTDIR})
		} else {
			f, err := os.OpenFile(requestPath, os.O_RDONLY, 0644)
			if err != nil {
				rpkt = statusFromError(p, err)
				break
			}
			handle := s.nextHandle(f, requestPath)
			rpkt = sshFxpHandlePacket{ID: p.id(), Handle: handle}
		}
	case *sshFxpReaddirPacket:
		f, ok := s.getHandle(p.Handle)
		if !ok {
			rpkt = statusFromError(p, syscall.EBADF)
			break
		}

		dirname := f.file.Name()
		dirents, err := f.file.Readdir(128)
		if err != nil {
			rpkt = statusFromError(p, err)
			break
		}

		ret := sshFxpNamePacket{ID: p.ID}

		for _, dirent := range dirents {
			if f.path == s.rootFolder {
				_, _, err := getAllowedPath(s, f.path+dirent.Name())
				if err == nil {
					ret.NameAttrs = append(ret.NameAttrs, sshFxpNameAttr{
						Name:     dirent.Name(),
						LongName: runLs(dirname, dirent),
						Attrs:    []interface{}{dirent},
					})
				}
			} else {
				ret.NameAttrs = append(ret.NameAttrs, sshFxpNameAttr{
					Name:     dirent.Name(),
					LongName: runLs(dirname, dirent),
					Attrs:    []interface{}{dirent},
				})
			}
		}
		rpkt = ret

	case *sshFxpReadPacket:
		var err error = syscall.EBADF
		f, ok := s.getHandle(p.Handle)
		if ok {
			err = nil
			data := make([]byte, clamp(p.Len, s.maxTxPacket))
			n, _err := f.file.ReadAt(data, int64(p.Offset))
			if _err != nil && (_err != io.EOF || n == 0) {
				err = _err
			}
			rpkt = sshFxpDataPacket{
				ID:     p.ID,
				Length: uint32(n),
				Data:   data[:n],
			}
		}
		if err != nil {
			rpkt = statusFromError(p, err)
		}
	default:
		rpkt = sshFxpStatusPacket{
			ID: p.id(),
			StatusError: StatusError{
				Code: ssh_FX_OP_UNSUPPORTED,
			},
		}
	}
	s.pktMgr.readyPacket(s.pktMgr.newOrderedResponse(rpkt, p.orderId()))
	return nil
}

// Serve serves SFTP connections until the streams stop or the SFTP subsystem
// is stopped.
func (svr *Server) Serve() error {
	var wg sync.WaitGroup
	runWorker := func(ch chan orderedRequest) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := svr.sftpServerWorker(ch); err != nil {
				svr.conn.Close() // shuts down recvPacket
			}
		}()
	}
	pktChan := svr.pktMgr.workerChan(runWorker)

	var err error
	var pkt requestPacket
	var pktType uint8
	var pktBytes []byte
	for {
		pktType, pktBytes, err = svr.recvPacket()
		if err != nil {
			break
		}

		pkt, err = makePacket(rxPacket{fxp(pktType), pktBytes})
		if err != nil {
			debug("makePacket err: %v", err)
			svr.conn.Close() // shuts down recvPacket
			break
		}

		pktChan <- svr.pktMgr.newOrderedRequest(pkt)
	}

	close(pktChan) // shuts down sftpServerWorkers
	wg.Wait()      // wait for all workers to exit

	// close any still-open files
	for handle, file := range svr.openFiles {
		fmt.Fprintf(svr.debugStream, "sftp server file with handle %q left open: %v\n", handle, file.file.Name())
		file.file.Close()
	}
	return err // error from recvPacket
}

type ider interface {
	id() uint32
}

// The init packet has no ID, so we just return a zero-value ID
func (p sshFxInitPacket) id() uint32 { return 0 }

type sshFxpStatResponse struct {
	ID   uint32
	info os.FileInfo
}

func (p sshFxpStatResponse) MarshalBinary() ([]byte, error) {
	b := []byte{ssh_FXP_ATTRS}
	b = marshalUint32(b, p.ID)
	b = marshalFileInfo(b, p.info)
	return b, nil
}

var emptyFileStat = []interface{}{uint32(0)}

// translateErrno translates a syscall error number to a SFTP error code.
func translateErrno(errno syscall.Errno) uint32 {
	switch errno {
	case 0:
		return ssh_FX_OK
	case syscall.ENOENT:
		return ssh_FX_NO_SUCH_FILE
	case syscall.EPERM:
		return ssh_FX_PERMISSION_DENIED
	}

	return ssh_FX_FAILURE
}

func statusFromError(p ider, err error) sshFxpStatusPacket {
	ret := sshFxpStatusPacket{
		ID: p.id(),
		StatusError: StatusError{
			Code: ssh_FX_OK,
		},
	}
	if err == nil {
		return ret
	}

	debug("statusFromError: error is %T %#v", err, err)
	ret.StatusError.Code = ssh_FX_FAILURE
	ret.StatusError.msg = err.Error()

	switch e := err.(type) {
	case syscall.Errno:
		ret.StatusError.Code = translateErrno(e)
	case *os.PathError:
		debug("statusFromError,pathError: error is %T %#v", e.Err, e.Err)
		if errno, ok := e.Err.(syscall.Errno); ok {
			ret.StatusError.Code = translateErrno(errno)
		}
	default:
		switch e {
		case io.EOF:
			ret.StatusError.Code = ssh_FX_EOF
		case os.ErrNotExist:
			ret.StatusError.Code = ssh_FX_NO_SUCH_FILE
		}
	}

	return ret
}

func clamp(v, max uint32) uint32 {
	if v > max {
		return max
	}
	return v
}

func runLsStatt(dirent os.FileInfo, statt *syscall.Stat_t) string {
	// example from openssh sftp server:
	// crw-rw-rw-    1 root     wheel           0 Jul 31 20:52 ttyvd
	// format:
	// {directory / char device / etc}{rwxrwxrwx}  {number of links} owner group size month day [time (this year) | year (otherwise)] name

	typeword := runLsTypeWord(dirent)
	numLinks := statt.Nlink
	uid := statt.Uid
	gid := statt.Gid
	username := fmt.Sprintf("%d", uid)
	groupname := fmt.Sprintf("%d", gid)
	// TODO FIXME: uid -> username, gid -> groupname lookup for ls -l format output

	mtime := dirent.ModTime()
	monthStr := mtime.Month().String()[0:3]
	day := mtime.Day()
	year := mtime.Year()
	now := time.Now()
	isOld := mtime.Before(now.Add(-time.Hour * 24 * 365 / 2))

	yearOrTime := fmt.Sprintf("%02d:%02d", mtime.Hour(), mtime.Minute())
	if isOld {
		yearOrTime = fmt.Sprintf("%d", year)
	}

	return fmt.Sprintf("%s %4d %-8s %-8s %8d %s %2d %5s %s", typeword, numLinks, username, groupname, dirent.Size(), monthStr, day, yearOrTime, dirent.Name())
}

// ls -l style output for a file, which is in the 'long output' section of a readdir response packet
// this is a very simple (lazy) implementation, just enough to look almost like openssh in a few basic cases
func runLs(dirname string, dirent os.FileInfo) string {
	dsys := dirent.Sys()
	if dsys == nil {
	} else if statt, ok := dsys.(*syscall.Stat_t); !ok {
	} else {
		return runLsStatt(dirent, statt)
	}

	return path.Join(dirname, dirent.Name())
}

func runLsTypeWord(dirent os.FileInfo) string {
	// find first character, the type char
	// b     Block special file.
	// c     Character special file.
	// d     Directory.
	// l     Symbolic link.
	// s     Socket link.
	// p     FIFO.
	// -     Regular file.
	tc := '-'
	mode := dirent.Mode()
	if (mode & os.ModeDir) != 0 {
		tc = 'd'
	} else if (mode & os.ModeDevice) != 0 {
		tc = 'b'
		if (mode & os.ModeCharDevice) != 0 {
			tc = 'c'
		}
	} else if (mode & os.ModeSymlink) != 0 {
		tc = 'l'
	} else if (mode & os.ModeSocket) != 0 {
		tc = 's'
	} else if (mode & os.ModeNamedPipe) != 0 {
		tc = 'p'
	}

	// owner
	orc := '-'
	if (mode & 0400) != 0 {
		orc = 'r'
	}
	owc := '-'
	if (mode & 0200) != 0 {
		owc = 'w'
	}
	oxc := '-'
	ox := (mode & 0100) != 0
	setuid := (mode & os.ModeSetuid) != 0
	if ox && setuid {
		oxc = 's'
	} else if setuid {
		oxc = 'S'
	} else if ox {
		oxc = 'x'
	}

	// group
	grc := '-'
	if (mode & 040) != 0 {
		grc = 'r'
	}
	gwc := '-'
	if (mode & 020) != 0 {
		gwc = 'w'
	}
	gxc := '-'
	gx := (mode & 010) != 0
	setgid := (mode & os.ModeSetgid) != 0
	if gx && setgid {
		gxc = 's'
	} else if setgid {
		gxc = 'S'
	} else if gx {
		gxc = 'x'
	}

	// all / others
	arc := '-'
	if (mode & 04) != 0 {
		arc = 'r'
	}
	awc := '-'
	if (mode & 02) != 0 {
		awc = 'w'
	}
	axc := '-'
	ax := (mode & 01) != 0
	sticky := (mode & os.ModeSticky) != 0
	if ax && sticky {
		axc = 't'
	} else if sticky {
		axc = 'T'
	} else if ax {
		axc = 'x'
	}

	return fmt.Sprintf("%c%c%c%c%c%c%c%c%c%c", tc, orc, owc, oxc, grc, gwc, gxc, arc, awc, axc)
}

// Makes sure we have a clean POSIX (/) absolute path to work with
func cleanPath(p string) string {
	p = filepath.ToSlash(p)
	if !filepath.IsAbs(p) {
		p = "/" + p
	}
	return path.Clean(p)
}
