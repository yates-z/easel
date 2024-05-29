package backend

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	backupTimeFormat = "2006-01-02T15-04-05.000"
	compressSuffix   = ".gz"
	defaultMaxSize   = 100
)

var (
	// megabyte is the conversion factor between MaxSize and bytes.  It is a
	// variable so tests can mock it out and not need to write megabytes of data
	// to disk.
	megabyte = 1024 * 1024
)

var _ Backend = (*defaultFileBackend)(nil)

type defaultFileBackend struct {
	mu   sync.Mutex
	opts FileBackendOptions
	file *os.File
	size int64

	millCh    chan bool
	startMill sync.Once
}

type defaultFileBackendBuilder struct {
	backend *defaultFileBackend
}

func (b *defaultFileBackendBuilder) Filename(filename string) *defaultFileBackendBuilder {
	b.backend.opts.filename = filename
	return b
}

func (b *defaultFileBackendBuilder) MaxSize(maxSize int) *defaultFileBackendBuilder {
	b.backend.opts.maxSize = maxSize
	return b
}

func (b *defaultFileBackendBuilder) MaxAge(maxAge int) *defaultFileBackendBuilder {
	b.backend.opts.maxAge = maxAge
	return b
}

func (b *defaultFileBackendBuilder) MaxBackups(maxBackups int) *defaultFileBackendBuilder {
	b.backend.opts.maxBackups = maxBackups
	return b
}

func (b *defaultFileBackendBuilder) Compress(compress bool) *defaultFileBackendBuilder {
	b.backend.opts.compress = compress
	return b
}

func (b *defaultFileBackendBuilder) LocalTime(localTime bool) *defaultFileBackendBuilder {
	b.backend.opts.localTime = localTime
	return b
}

func (b *defaultFileBackendBuilder) Build() *defaultFileBackend {
	return b.backend
}

func DefaultFileBackend() *defaultFileBackendBuilder {

	return &defaultFileBackendBuilder{
		backend: &defaultFileBackend{},
	}
}

type FileBackendOptions struct {
	// filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	filename string

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	maxSize int

	// maxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	maxAge int

	// maxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	maxBackups int

	// compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	compress bool

	// localTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	localTime bool
}

func (d *defaultFileBackend) Sync() error {
	return nil
}

func (d *defaultFileBackend) AllowANSI() bool {
	return false
}

func (d *defaultFileBackend) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	writeLen := int64(len(p))
	if writeLen > d.max() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, d.max(),
		)
	}

	if d.file == nil {
		if err = d.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	if d.size+writeLen > d.max() {
		if err := d.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = d.file.Write(p)
	d.size += int64(n)

	return n, err
}

// filename generates the name of the logfile from the current time.
func (d *defaultFileBackend) filename() string {
	if d.opts.filename != "" {
		return d.opts.filename
	}
	name := filepath.Base(os.Args[0]) + "-easel.log"
	return filepath.Join(os.TempDir(), name)
}

// openExistingOrNew opens the logfile if it exists and if the current write
// would not put it over MaxSize.  If there is no such file or the write would
// put it over the MaxSize, a new file is created.
func (d *defaultFileBackend) openExistingOrNew(writeLen int) error {
	d.mill()

	filename := d.filename()
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return d.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if info.Size()+int64(writeLen) >= d.max() {
		return d.rotate()
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// if we fail to open the old log file for some reason, just ignore
		// it and open a new log file.
		return d.openNew()
	}
	d.file = file
	d.size = info.Size()
	return nil
}

// millRunOnce performs compression and removal of stale log files.
// Log files are compressed if enabled via configuration and old log
// files are removed, keeping at most d.MaxBackups files, as long as
// none of them are older than MaxAge.
func (d *defaultFileBackend) millRunOnce() error {
	if d.opts.maxBackups == 0 && d.opts.maxAge == 0 && !d.opts.compress {
		return nil
	}

	files, err := d.oldLogFiles()
	if err != nil {
		return err
	}

	var compress, remove []logInfo

	if d.opts.maxBackups > 0 && d.opts.maxBackups < len(files) {
		preserved := make(map[string]bool)
		var remaining []logInfo
		for _, f := range files {
			// Only count the uncompressed log file or the
			// compressed log file, not both.
			fn := f.Name()
			if strings.HasSuffix(fn, compressSuffix) {
				fn = fn[:len(fn)-len(compressSuffix)]
			}
			preserved[fn] = true

			if len(preserved) > d.opts.maxBackups {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}
	if d.opts.maxAge > 0 {
		diff := time.Duration(int64(24*time.Hour) * int64(d.opts.maxAge))
		cutoff := time.Now().Add(-1 * diff)

		var remaining []logInfo
		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}

	if d.opts.compress {
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), compressSuffix) {
				compress = append(compress, f)
			}
		}
	}

	for _, f := range remove {
		errRemove := os.Remove(filepath.Join(d.dir(), f.Name()))
		if err == nil && errRemove != nil {
			err = errRemove
		}
	}
	for _, f := range compress {
		fn := filepath.Join(d.dir(), f.Name())
		errCompress := compressLogFile(fn, fn+compressSuffix)
		if err == nil && errCompress != nil {
			err = errCompress
		}
	}

	return err
}

// millRun runs in a goroutine to manage post-rotation compression and removal
// of old log files.
func (d *defaultFileBackend) millRun() {
	for range d.millCh {
		// what am I going to do, log this?
		_ = d.millRunOnce()
	}
}

// mill performs post-rotation compression and removal of stale log files,
// starting the mill goroutine if necessary.
func (d *defaultFileBackend) mill() {
	d.startMill.Do(func() {
		d.millCh = make(chan bool, 1)
		go d.millRun()
	})
	select {
	case d.millCh <- true:
	default:
	}
}

// oldLogFiles returns the list of backup log files stored in the same
// directory as the current log file, sorted by ModTime
func (d *defaultFileBackend) oldLogFiles() ([]logInfo, error) {
	files, err := ioutil.ReadDir(d.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	var logFiles []logInfo

	prefix, ext := d.prefixAndExt()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if t, err := d.timeFromName(f.Name(), prefix, ext); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
		if t, err := d.timeFromName(f.Name(), prefix, ext+compressSuffix); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
		// error parsing means that the suffix at the end was not generated
		// by lumberjack, and therefore it's not a backup file.
	}

	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

// openNew opens a new log file for writing, moving any old log file out of the
// way.  This methods assumes the file has already been closed.
func (d *defaultFileBackend) openNew() error {
	err := os.MkdirAll(d.dir(), 0755)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := d.filename()
	mode := os.FileMode(0600)
	info, err := os.Stat(name)
	if err == nil {
		// Copy the mode off the old logfile.
		mode = info.Mode()
		// move the existing file
		newname := backupName(name, d.opts.localTime)
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}

		// this is a no-op anywhere but linux
		if err := chown(name, info); err != nil {
			return err
		}
	}

	// we use truncate here because this should only get called when we've moved
	// the file ourselves. if someone else creates the file in the meantime,
	// just wipe out the contents.
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	d.file = f
	d.size = 0
	return nil
}

// Close implements io.Closer, and closes the current logfile.
func (d *defaultFileBackend) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.close()
}

// close closes the file if it is open.
func (d *defaultFileBackend) close() error {
	if d.file == nil {
		return nil
	}
	err := d.file.Close()
	d.file = nil
	return err
}

// Rotate causes Logger to close the existing log file and immediately create a
// new one.  This is a helper function for applications that want to initiate
// rotations outside of the normal rotation rules, such as in response to
// SIGHUP.  After rotating, this initiates compression and removal of old log
// files according to the configuration.
func (d *defaultFileBackend) Rotate() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.rotate()
}

// rotate closes the current file, moves it aside with a timestamp in the name,
// (if it exists), opens a new file with the original filename, and then runs
// post-rotation processing and removal.
func (d *defaultFileBackend) rotate() error {
	if err := d.close(); err != nil {
		return err
	}
	if err := d.openNew(); err != nil {
		return err
	}
	d.mill()
	return nil
}

// dir returns the directory for the current filename.
func (d *defaultFileBackend) dir() string {
	return filepath.Dir(d.filename())
}

// max returns the maximum size in bytes of log files before rolling.
func (d *defaultFileBackend) max() int64 {
	if d.opts.maxSize == 0 {
		return int64(defaultMaxSize * megabyte)
	}
	return int64(d.opts.maxSize) * int64(megabyte)
}

// timeFromName extracts the formatted time from the filename by stripping off
// the filename's prefix and extension. This prevents someone's filename from
// confusing time.parse.
func (d *defaultFileBackend) timeFromName(filename, prefix, ext string) (time.Time, error) {
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, errors.New("mismatched extension")
	}
	ts := filename[len(prefix) : len(filename)-len(ext)]
	return time.Parse(backupTimeFormat, ts)
}

// prefixAndExt returns the filename part and extension part from the Logger's
// filename.
func (d *defaultFileBackend) prefixAndExt() (prefix, ext string) {
	filename := filepath.Base(d.filename())
	ext = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(ext)] + "-"
	return prefix, ext
}

// backupName creates a new filename from the given name, inserting a timestamp
// between the filename and the extension, using the local time if requested
// (otherwise UTC).
func backupName(name string, local bool) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	t := time.Now()
	if !local {
		t = t.UTC()
	}

	timestamp := t.Format(backupTimeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

func chown(_ string, _ os.FileInfo) error {
	return nil
}

// compressLogFile compresses the given log file, removing the
// uncompressed log file if successful.
func compressLogFile(src, dst string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer f.Close()

	fi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %v", err)
	}

	if err := chown(dst, fi); err != nil {
		return fmt.Errorf("failed to chown compressed log file: %v", err)
	}

	// If this file already exists, we presume it was created by
	// a previous attempt to compress the log file.
	gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed to open compressed log file: %v", err)
	}
	defer gzf.Close()

	gz := gzip.NewWriter(gzf)

	defer func() {
		if err != nil {
			os.Remove(dst)
			err = fmt.Errorf("failed to compress log file: %v", err)
		}
	}()

	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := gzf.Close(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}

// logInfo is a convenience struct to return the filename and its embedded
// timestamp.
type logInfo struct {
	timestamp time.Time
	os.FileInfo
}

// byFormatTime sorts by newest time formatted in the name.
type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
