package nfsmanager

import (
	"os/exec"
	"reflect"
	"testing"
)

func Test_exportFSCommandLine(t *testing.T) {
	type args struct {
		path    string
		host    string
		options []nfsOption
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"No Options", args{"/foo/bar", "192.168.1.1", []nfsOption{}}, []string{"exportfs", "/foo/bar:192.168.1.1"}},
		{"One extra-less option", args{"/foo/bar", "192.168.1.1", []nfsOption{NoRootSquash}}, []string{"exportfs", "/foo/bar:192.168.1.1", "-o", "no_root_squash"}},
		{"Two extra-less options", args{"/foo/bar", "192.168.1.1", []nfsOption{NoRootSquash, InsecureLocks}}, []string{"exportfs", "/foo/bar:192.168.1.1", "-o", "no_root_squash,insecure_locks"}},
		{"Two option: one extra-less, one with extras", args{"/foo/bar", "192.168.1.1", []nfsOption{NoRootSquash, FsID("some-id")}}, []string{"exportfs", "/foo/bar:192.168.1.1", "-o", "no_root_squash,fsid=some-id"}},
		{"Option with multiple extras", args{"/foo/bar", "192.168.1.1", []nfsOption{Replicas("foo", "bar")}}, []string{"exportfs", "/foo/bar:192.168.1.1", "-o", "replicas=foo:bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exportFSCommandLine(tt.args.path, tt.args.host, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("exportFSCommandLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nfsOptions(t *testing.T) {
	tests := []struct {
		name   string
		option nfsOption
		want   string
	}{
		{"Secure", Secure, "secure"},
		{"RW", RW, "rw"},
		{"ASync", ASync, "async"},
		{"Sync", Sync, "sync"},
		{"NoWDelay", NoWDelay, "no_wdelay"},
		{"NoHide", NoHide, "nohide"},
		{"CrossMnt", CrossMnt, "crossmnt"},
		{"NoSubtreeCheck", NoSubtreeCheck, "no_subtree_check"},
		{"InsecureLocks", InsecureLocks, "insecure_locks"},
		{"NoAuthNLM", NoAuthNLM, "no_auth_nlm"},
		{"SecureLocks", SecureLocks, "secure_locks"},
		{"AuthNLM", AuthNLM, "auth_nlm"},
		{"MountPoint with empty string arg", MountPoint(""), "mountpoint"},
		{"MountPoint with arg", MountPoint("/foo/bar"), "mountpoint=/foo/bar"},
		{"MP with empty string arg", MP(""), "mp"},
		{"MP with arg", MP("/foo/bar"), "mp=/foo/bar"},
		{"FsID", FsID("the-fs-id"), "fsid=the-fs-id"},
		{"NoRDirPlus", NoRDirPlus, "nordirplus"},
		{"Refer with no args", Refer(), ""},
		{"Refer with args", Refer("foo", "bar"), "refer=foo:bar"},
		{"Refer with args, first is empty string", Refer("", "bar"), "refer=bar"},
		{"Refer with args, second is empty string", Refer("foo", ""), "refer=foo"},
		{"Refer with args, all are empty strings", Refer("", ""), ""},
		{"Replicas with no args", Replicas(), ""},
		{"Replicas with args", Replicas("foo", "bar"), "replicas=foo:bar"},
		{"Replicas with args, first is empty string", Replicas("", "bar"), "replicas=bar"},
		{"Replicas with args, second is empty string", Replicas("foo", ""), "replicas=foo"},
		{"Replicas with args, all are empty strings", Replicas("", ""), ""},
		{"PNFS", PNFS, "pnfs"},
		{"NoPNFS", NoPNFS, "no_pnfs"},
		{"RootSquash", RootSquash, "root_squash"},
		{"NoRootSquash", NoRootSquash, "no_root_squash"},
		{"AllSquash", AllSquash, "all_squash"},
		{"AnonUID", AnonUID(1234), "anonuid=1234"},
		{"AnonGID", AnonGID(2345), "anongid=2345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.option.string(); got != tt.want {
				t.Errorf("option.string() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nfsManager_wxportFs(t *testing.T) {
	succeed := func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	fail := func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	succeedOnlyWithSudo := func(name string, arg ...string) *exec.Cmd {
		if name == "sudo" {
			return succeed(name, arg...)
		}
		return fail(name, arg...)
	}
	type fields struct {
		Command func(name string, arg ...string) *exec.Cmd
	}
	type args struct {
		path    string
		host    string
		options []nfsOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Works without sudo", fields{succeed}, args{"path", "host", []nfsOption{NoRootSquash}}, false},
		{"Works with sudo", fields{succeedOnlyWithSudo}, args{"path", "host", []nfsOption{NoRootSquash}}, false},
		{"Fails either way", fields{fail}, args{"path", "host", []nfsOption{NoRootSquash}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := nfsManager{
				Command: tt.fields.Command,
			}
			if err := n.ExportFs(tt.args.path, tt.args.host, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("nfsManager.ExportFs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNFSManager(t *testing.T) {
	tests := []struct {
		name string
		want *nfsManager
	}{
		// TODO: Add test cases.
		{"NFSManager", &nfsManager{exec.Command}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := reflect.ValueOf(NFSManager().Command).Pointer()
			b := reflect.ValueOf(exec.Command).Pointer()
			if a != b {
				t.Errorf("NFSManager() = %v, want %v", a, b)
			}
		})
	}
}
