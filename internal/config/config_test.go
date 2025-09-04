package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/Blustak/bootdev-gator/internal/config"
)

func isType (T any,x any) bool {
    return reflect.TypeOf(x) == T
}

func TestRead(t *testing.T) {
    confDir,err := os.UserConfigDir()
    if err != nil {
        t.Errorf("Could not get user config directory")
    }
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		configFilePath string
		want           config.Config
		wantErr        bool
	}{
        {
            name: "Default case",
            configFilePath: confDir + "/gatorconfig.json",
            want: config.Config{},
            wantErr: false,

        },
        {
            name: "Malformed path",
            configFilePath: "",
            want: config.Config{},
            wantErr: true,
        },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := config.Read(tt.configFilePath)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Read() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Read() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if  ! isType(reflect.TypeOf(config.Config{}),got) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

