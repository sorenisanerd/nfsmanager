package nfsmanager

import (
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
		// TODO: Add test cases.
		{"No Options", args{"/foo/bar", "192.168.1.1", []nfsOption{}}, []string{"exportfs", "/foo/bar:192.168.1.1"}},
		{"One extra-less option", args{"/foo/bar", "192.168.1.1", []nfsOption{NoRootSquash}}, []string{"exportfs", "/foo/bar:192.168.1.1", "-o", "no_root_squash"}},
		{"Two extra-less options", args{"/foo/bar", "192.168.1.1", []nfsOption{NoRootSquash, InsecureLocks}}, []string{"exportfs", "/foo/bar:192.168.1.1", "-o", "no_root_squash,insecure_locks"}},
		{"Two option: one extra-less, one with extras", args{"/foo/bar", "192.168.1.1", []nfsOption{NoRootSquash, FsID("some-id")}}, []string{"exportfs", "/foo/bar:192.168.1.1", "-o", "no_root_squash,fsid=some-id"}},
		{"Option with multiple extras", args{"/foo/bar", "192.168.1.1", []nfsOption{Replicas([]string{"foo", "bar"})}}, []string{"exportfs", "/foo/bar:192.168.1.1", "-o", "replicas=foo:bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exportFSCommandLine(tt.args.path, tt.args.host, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("exportFSCommandLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
