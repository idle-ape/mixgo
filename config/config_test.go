package config

import (
	"sync"
	"testing"
	"time"
)

var conf = &yamlStruct{}

type yamlStruct struct {
	Name string `yaml:"name"`
}

func TestMixgoConfigLoader_Load(t *testing.T) {
	type fields struct {
		cache map[string]Config
		rwl   sync.RWMutex
	}
	type args struct {
		path string
		opts []LoadOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				cache: make(map[string]Config),
				rwl:   sync.RWMutex{},
			},
			args: args{
				path: "../testdata/test.yaml",
				opts: []LoadOption{
					WithReciver(conf),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &MixgoConfigLoader{
				cache: tt.fields.cache,
				rwl:   tt.fields.rwl,
			}
			got, err := cl.Load(tt.args.path, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MixgoConfigLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got: %#v", got)
			t.Logf("conf: %#v", conf)
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		path string
		opts []LoadOption
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				path: "../testdata/test.yaml",
				opts: []LoadOption{
					WithReciver(conf),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(tt.args.path, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("conf: %+v", conf)
			time.Sleep(10 * time.Second)
			t.Logf("conf: %+v", conf)
		})
	}
}

func TestMixgoConfig_Load(t *testing.T) {
	type fields struct {
		p            Provider
		r            interface{}
		path         string
		disableWatch bool
		decoder      Codec
		rawData      []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MixgoConfig{
				p:            tt.fields.p,
				r:            tt.fields.r,
				path:         tt.fields.path,
				disableWatch: tt.fields.disableWatch,
				decoder:      tt.fields.decoder,
				rawData:      tt.fields.rawData,
			}
			if err := m.Load(); (err != nil) != tt.wantErr {
				t.Errorf("MixgoConfig.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
