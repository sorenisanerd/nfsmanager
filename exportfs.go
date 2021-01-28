package nfsmanager

import (
	"fmt"
	"os/exec"
	"strings"

	"log"
)

type nfsOption struct {
	optionString     string
	extra            []string
	omitIfExtraEmpty bool
}

// Secure requires that requests originate on an Internet port less than
// IPPORT_RESERVED (1024). This option is on by default. To turn it off,
// specify InSecure.
var Secure nfsOption = nfsOption{
	optionString: "secure",
}

// RW allows both read and write requests on this NFS volume. The
// default is to disallow any request which changes the filesystem. This
// can also be made  explicit by using the RO option.
var RW nfsOption = nfsOption{
	optionString: "rw",
}

// ASync allows the NFS server to violate the NFS protocol and reply to
// requests before any changes made by that request have been committed
// to stable storage (e.g. disc drive).
//
// Using this option usually improves performance, but at the cost that
// an unclean server restart (i.e. a crash) can cause data to be lost or
// corrupted.
var ASync nfsOption = nfsOption{
	optionString: "async",
}

// Sync causes the nfs server to reply to requests only after the
// changes have been committed to stable storage (see async above).
//
// In releases of nfs-utils up to and including 1.0.0, the ASync option
// was the default.  In all releases after 1.0.0, Sync is the default,
// and Async must be explicitly requested if needed.
// TODO: help make system administrators aware of this change, by
// insisting on either being set.
// exportfs will issue a. Update this doc per man page once done.
var Sync nfsOption = nfsOption{
	optionString: "sync",
}

// NoWDelay has no effect if ASync is also set.  The NFS server will
// normally delay committing a write request to disc slightly if it
// suspects that another related write  request may be in progress or
// may arrive soon.  This allows multiple write requests to be committed
// to disc with the one operation which can improve performance.  If an
// NFS server received mainly small unrelated requests, this behaviour
// could actually reduce performance, so no_wdelay is available to turn
// it off.  The default can be explicitly requested with the wdelay
// option.
var NoWDelay nfsOption = nfsOption{
	optionString: "no_wdelay",
}

// NoHide is based  on  the option of the same name provided in IRIX
// VNFS.  Normally, if a server exports two filesystems one of which
// is mounted on the other, then the client will have to mount both
// filesystems explicitly to get access to them.  If it just mounts the
// parent, it will see an empty  directory at the place where the other
// filesystem is mounted.  That filesystem is "hidden".
//
// Setting  the NoHide option on a filesystem causes it not to be hidden,
// and an appropriately authorised client will be able to move from the
// parent to that filesystem without noticing the change.
//
// However, some NFS clients do not cope well with this situation as,
// for instance, it is then possible for two files in the one apparent
// filesystem to have the same inode number.
//
// The nohide option is currently only effective on single host exports.
// It does not work reliably with netgroup, subnet, or wildcard exports.
//
// This option can be very useful in some situations, but it should be
// used with due care, and only after confirming that the client system
// copes with the situation effectively.
//
// The option can be explicitly disabled for NFSv2 and NFSv3 with hide.
//
// This option is not relevant when NFSv4 is use.  NFSv4 never hides
// subordinate filesystems.  Any filesystem that is exported will be
// visible where expected when using NFSv4.
var NoHide nfsOption = nfsOption{
	optionString: "nohide",
}

// CrossMnt is similar to NoHide but it makes it possible for clients to
// access all filesystems mounted on a filesystem marked with crossmnt.
// Thus when a child filesystem "B" is mounted on a parent "A", setting
// crossmnt on "A" has a similar effect to setting "nohide" on B.
//
// With NoHide the child filesystem needs to be explicitly exported.
// With CrossMnt it need not.  If a child of a CrossMnt file is not
// explicitly exported, then it will be implicitly exported with the
// same export options as the parent, except for FsId.  This makes it
// impossible to not export a child of a CrossMnt filesystem.  If some
// but not all subordinate filesystems of a parent are to be exported,
// then they must be explicitly exported and the parent should not have
// CrossMnt set.
//
// The NoCrossMnt option can explictly disable CrossMnt if it was
// previously set.  This is rarely useful.
var CrossMnt nfsOption = nfsOption{
	optionString: "crossmnt",
}

// NoSubtreeCheck option disables subtree checking, which has mild
// security implications, but can improve reliability in some
// circumstances.
//
// If a subdirectory of a filesystem is exported, but the whole
// filesystem isn't then whenever a NFS request arrives, the server must
// check not only that the accessed file is in the appropriate
// filesystem (which is easy) but also that it is in the exported tree
// (which is harder). This check is called the subtree_check.
//
// In order to perform this check, the server must include some
// information about the location of the file in the "filehandle" that
// is given to the client. This can cause problems with accessing files
// that are renamed while a client has them open (though in many simple
// cases it will still work).
//
// subtree checking is also used to make sure that files inside
// directories to which only root has access can only be accessed if the
// filesystem is exported with NoRootSquash (see below), even if the
// file itself allows more general access.
//
// As a general guide, a home directory filesystem, which is normally
// exported at the root and may see lots of file renames, should be
// exported with subtree checking disabled.  A filesystem which is
// mostly readonly, and at least doesn't see many file renames (e.g.
// /usr or /var) and for which subdirectories may be exported, should
// probably be exported with subtree checks enabled.
//
// The default of having subtree checks enabled, can be explicitly
// requested with SubtreeCheck.
//
// From release 1.1.0 of nfs-utils onwards, the default will be
// NoSubtreeCheck a SubtreeChecking tends to cause more problems than it
// is worth.  If you genuinely require subtree checking, you should
// explicitly put that option in the exports file.
var NoSubtreeCheck nfsOption = nfsOption{
	optionString: "no_subtree_check",
}

// InsecureLocks tells the NFS server not to require authentication of
// locking requests (i.e. requests which use the NLM protocol). Normally
// the NFS server will require a lock request to hold a credential for a
// user who has read access to the file.  With this flag no access
// checks will be performed.
//
// Early NFS client implementations did not send credentials with lock
// requests, and many current NFS clients still exist which are based on
// the old implementations.  Use this flag if you find that you can only
// lock files which are world readable.
//
// The default behaviour of requiring authentication for NLM requests
// can be explicitly requested with either of the (synonymous) AuthNlm,
// or SecureLocks.
var InsecureLocks nfsOption = nfsOption{
	optionString: "insecure_locks",
}

// NoAuthNLM is synonymous with InsecureLogs
var NoAuthNLM nfsOption = nfsOption{
	optionString: "no_auth_nlm",
}

// SecureLocks is the opposite of InsecureLocks
var SecureLocks nfsOption = nfsOption{
	optionString: "secure_locks",
}

// AuthNLM is the opposite of NoAuthNLM
var AuthNLM nfsOption = nfsOption{
	optionString: "auth_nlm",
}

// MountPoint makes it possible to only export a directory if it has
// successfully been mounted.  If no path is given (e.g. MountPoint or
// MP) then the export point must also be a mount point.  If it isn't
// then the export point is not exported.  This allows you to be sure
// that the directory underneath a mountpoint will never be exported by
// accident if, for example, the filesystem failed to mount due to a
// disc error.
//
// If a path is given (e.g. mountpoint=/path or mp=/path) then the
// nominated path must be a mountpoint for the exportpoint to be
// exported.
func MountPoint(path string) nfsOption {
	opt := nfsOption{
		optionString: "mountpoint",
	}
	if path != "" {
		opt.extra = []string{path}
	}
	return opt
}

// MP is synonymous with MountPoint
func MP(path string) nfsOption {
	opt := MountPoint(path)
	opt.optionString = "mp"
	return opt
}

// FsID defines an export's ID.
//
// NFS needs to be able to identify each filesystem that it
// exports. Normally it will use a UUID for the filesystem (if the
// filesystem has such a thing) or the device number of the device
// holding the filesystem (if the filesystem is stored on the device).
//
// As not all filesystems are stored on devices, and not all filesystems
// have UUIDs, it is sometimes necessary to explicitly tell NFS how to
// identify a filesystem.  This is done with the FsID option.
//
// For NFSv4, there is a distinguished filesystem which is the root of
// all exported filesystem.  This is specified with fsid=root or fsid=0
// both of which mean exactly the same thing.
//
// Other filesystems can be identified with a small integer, or a UUID
// which should contain 32 hex digits and arbitrary punctuation.
//
// Linux kernels version 2.6.20 and earlier do not understand the UUID
// setting so a small integer must be used if an fsid option needs to be
// set for such kernels.  Setting both a small number and a UUID is
// supported so the same configuration can be made to work on old and
// new kernels alike.
func FsID(id string) nfsOption {
	return nfsOption{
		optionString: "fsid",
		extra:        []string{id},
	}
}

// NoRDirPlus will disable READDIRPLUS request handling.  When set,
// READDIRPLUS requests from NFS clients return NFS3ERR_NOTSUPP, and
// clients fall back on READDIR.  This option affects only NFSv3
// clients.
var NoRDirPlus nfsOption = nfsOption{
	optionString: "nordirplus",
}

// Refer specifies relocations
//
// A client referencing the export point will be directed to choose from
// the given list an alternative location for the filesystem.  (Note
// that the server must have a mount‐point here, though a different
// filesystem is not required; so, for example, mount --bind /path /path
// is sufficient.)
func Refer(references ...string) nfsOption {
	return nfsOption{
		optionString:     "refer",
		extra:            references,
		omitIfExtraEmpty: true,
	}
}

// Replicas specifies replicas
//
// If the client asks for alternative locations for the export point, it
// will be given this list of alternatives. (Note that actual
// replication of the filesystem must be handled elsewhere.)
func Replicas(replicas ...string) nfsOption {
	return nfsOption{
		optionString:     "replicas",
		extra:            replicas,
		omitIfExtraEmpty: true,
	}
}

// PNFS allows enables the use of pNFS extension if protocol level is
// NFSv4.1 or higher, and the filesystem supports pNFS exports.  With
// pNFS clients can bypass the server and perform I/O directly to
// storage devices. The default can be explicitly requested with the
// NoPNFS option.
var PNFS nfsOption = nfsOption{
	optionString: "pnfs",
}

// NoPNFS is the opposite of PNFS
var NoPNFS nfsOption = nfsOption{
	optionString: "no_pnfs",
}

// RootSquash maps requests from uid/gid 0 to the anonymous uid/gid.
// Note that this does not apply to any other uids or gids that might be
// equally sensitive, such as user bin or group staff.
var RootSquash nfsOption = nfsOption{
	optionString: "root_squash",
}

// NoRootSquash turns off root squasing.  This option is mainly useful
// for diskless clients.
var NoRootSquash nfsOption = nfsOption{
	optionString: "no_root_squash",
}

// AllSquash maps all uids and gids to the anonymous user. Useful for
// NFS-exported public FTP directories, news spool directories, etc. The
// opposite option is no_all_squash, which is the default setting.
var AllSquash nfsOption = nfsOption{
	optionString: "all_squash",
}

// AnonUID explicitly set the uid of the anonymous account. This option
// is primarily useful for PC/NFS clients, where you might want all
// requests appear to be from one user.
func AnonUID(uid int) nfsOption {
	return nfsOption{
		optionString: "anonuid",
		extra:        []string{fmt.Sprintf("%d", uid)},
	}
}

// AnonGID explicitly set the gid of the anonymous account. This option
// is primarily useful for PC/NFS clients, where you might want all
// requests appear to be from one user.
func AnonGID(gid int) nfsOption {
	return nfsOption{
		optionString: "anongid",
		extra:        []string{fmt.Sprintf("%d", gid)},
	}
}

func (opt nfsOption) string() string {
	extrasString := opt.extrasString()
	if extrasString == "" && opt.omitIfExtraEmpty {
		return ""
	}
	return fmt.Sprintf("%s%s", opt.optionString, extrasString)
}

func optionsString(options []nfsOption) string {
	var optStrings []string
	for _, opt := range options {
		optStrings = append(optStrings, opt.string())
	}
	return strings.Join(optStrings, ",")
}

func (opt nfsOption) extrasString() string {
	extras := make([]string, 0)
	for _, extra := range opt.extra {
		if extra != "" {
			extras = append(extras, extra)
		}
	}
	if len(extras) > 0 {
		return "=" + strings.Join(extras, ":")
	}
	return ""
}

func exportFSCommandLine(path string, host string, options []nfsOption) []string {
	var exportString string = fmt.Sprintf("%s:%s", host, path)

	cmd := []string{"exportfs", exportString}
	if len(options) > 0 {
		cmd = append(cmd, []string{"-o", optionsString(options)}...)
	}
	return cmd
}

func unExportFSCommandLine(path string, host string) []string {
	var exportString string = fmt.Sprintf("%s:%s", host, path)

	return []string{"exportfs", "-u", exportString}
}

type execCommander func(name string, arg ...string) *exec.Cmd
type commandRetrierWithSudo func([]string, execCommander) error

type nfsManager struct {
	Command        execCommander
	commandRetrier commandRetrierWithSudo
}

func NFSManager() *nfsManager {
	return &nfsManager{
		Command:        exec.Command,
		commandRetrier: runAndRetryWithSudoOnFailure,
	}
}

// ExportFs will export path to host with the given options.
// Note: The export is not persisted to /etc/exports
func (n *nfsManager) ExportFs(path string, host string, options ...nfsOption) error {
	return n.commandRetrier(exportFSCommandLine(path, host, options), n.Command)
}

// UnExportFs will unexport path to host with the given options.
// Note: The export is not removed from /etc/exports if it's there
func (n *nfsManager) UnExportFs(path string, host string) error {
	return n.commandRetrier(unExportFSCommandLine(path, host), n.Command)
}

func runAndRetryWithSudoOnFailure(cmdLine []string, command execCommander) error {
	cmd := command(cmdLine[0], cmdLine[1:]...)
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Command %v failed: %s", cmd, err)
		log.Printf("Retrying with sudo")

		cmdLine = append(cmdLine, "", "")
		copy(cmdLine[2:], cmdLine)
		cmdLine[0] = "sudo"
		cmdLine[1] = "-n"

		cmd = command(cmdLine[0], cmdLine[1:]...)
		_, err := cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("Command %v failed with sudo as well: %s", cmd, err)
		}
	}
	return nil
}
