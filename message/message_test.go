package message

import (
	"testing"
)

func TestGetCommendArg(t *testing.T) {
	type args struct {
		command string
		str     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Test 1", args: args{command: "!a", str: "!a hello world"}, want: "hello world", wantErr: false},
		{name: "Test 2", args: args{command: "!a", str: "!b hello world"}, want: "", wantErr: true},
		{name: "Test 3", args: args{command: "!b", str: "!b hello "}, want: "hello", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCommendArg(tt.args.command, tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommendArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCommendArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
		want2 string
	}{
		{name: "Test 1", args: args{text: "/name hello"}, want: true, want1: "name", want2: "hello"},
		{name: "Test 2", args: args{text: "!name jack"}, want: false, want1: "!name", want2: "jack"},
		{name: "Test 3", args: args{text: "/name  hello"}, want: true, want1: "name", want2: "hello"},
		{name: "Test 4", args: args{text: "/name hello "}, want: true, want1: "name", want2: "hello"},
		{name: "Test 5", args: args{text: " /name hello"}, want: true, want1: "name", want2: "hello"},
		{name: "Test 6", args: args{text: " !name jack"}, want: false, want1: "!name", want2: "jack"},
		{name: "Test 7", args: args{text: "/name hello world"}, want: true, want1: "name", want2: "hello world"},
		//{name: "Test 8", args: args{text: "!com"}, want: false, want1: "!com", want2: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := TrimText(tt.args.text)
			if got != tt.want {
				t.Errorf("TrimText() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TrimText() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("TrimText() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
