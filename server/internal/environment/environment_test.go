//go:build !integration

package environment

import (
	"reflect"
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	t.Parallel()

	type args struct {
		env string
	}
	tests := []struct {
		name    string
		args    args
		want    Environment
		wantErr bool
	}{
		{
			name:    "local",
			args:    args{env: "local"},
			want:    Environment{EnvLocal},
			wantErr: false,
		},
		{
			name:    "dev",
			args:    args{env: "dev"},
			want:    Environment{EnvDev},
			wantErr: false,
		},
		{
			name:    "develop",
			args:    args{env: "develop"},
			want:    Environment{EnvDev},
			wantErr: false,
		},
		{
			name:    "development",
			args:    args{env: "development"},
			want:    Environment{EnvDev},
			wantErr: false,
		},
		{
			name:    "stg",
			args:    args{env: "stg"},
			want:    Environment{EnvStg},
			wantErr: false,
		},
		{
			name:    "stage",
			args:    args{env: "stage"},
			want:    Environment{EnvStg},
			wantErr: false,
		},
		{
			name:    "staging",
			args:    args{env: "staging"},
			want:    Environment{EnvStg},
			wantErr: false,
		},
		{
			name:    "prd",
			args:    args{env: "prd"},
			want:    Environment{EnvPrd},
			wantErr: false,
		},
		{
			name:    "prod",
			args:    args{env: "prod"},
			want:    Environment{EnvPrd},
			wantErr: false,
		},
		{
			name:    "production",
			args:    args{env: "production"},
			want:    Environment{EnvPrd},
			wantErr: false,
		},
		{
			name:    "uppercase LOCAL",
			args:    args{env: "LOCAL"},
			want:    Environment{EnvLocal},
			wantErr: false,
		},
		{
			name:    "mixed case Dev",
			args:    args{env: "Dev"},
			want:    Environment{EnvDev},
			wantErr: false,
		},
		{
			name:    "surrounding whitespace",
			args:    args{env: "  prd  "},
			want:    Environment{EnvPrd},
			wantErr: false,
		},
		{
			name:    "uppercase with whitespace",
			args:    args{env: "\tSTAGING\n"},
			want:    Environment{EnvStg},
			wantErr: false,
		},
		// 異常系
		{
			name:    "empty",
			args:    args{env: ""},
			want:    Environment{},
			wantErr: true,
		},
		{
			name:    "whitespace only",
			args:    args{env: "   "},
			want:    Environment{},
			wantErr: true,
		},
		{
			name:    "unknown",
			args:    args{env: "unknown"},
			want:    Environment{},
			wantErr: true,
		},
		{
			name:    "partial match prod-like",
			args:    args{env: "prdx"},
			want:    Environment{},
			wantErr: true,
		},
		{
			name:    "inner whitespace not trimmed",
			args:    args{env: "pr d"},
			want:    Environment{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewEnvironment(tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewEnvironment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Environment
		want string
	}{
		// 正常系
		{
			name: "local",
			e:    Environment{EnvLocal},
			want: "local",
		},
		{
			name: "dev",
			e:    Environment{EnvDev},
			want: "dev",
		},
		{
			name: "stg",
			e:    Environment{EnvStg},
			want: "stg",
		},
		{
			name: "prd",
			e:    Environment{EnvPrd},
			want: "prd",
		},
		// 境界値 / 異常系（ゼロ値）
		{
			name: "zero value",
			e:    Environment{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.e.String(); got != tt.want {
				t.Errorf("Environment.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_IsLocal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Environment
		want bool
	}{
		{
			name: "local is local",
			e:    Environment{EnvLocal},
			want: true,
		},
		{
			name: "dev is not local",
			e:    Environment{EnvDev},
			want: false,
		},
		{
			name: "stg is not local",
			e:    Environment{EnvStg},
			want: false,
		},
		{
			name: "prd is not local",
			e:    Environment{EnvPrd},
			want: false,
		},
		{
			name: "zero value is not local",
			e:    Environment{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.e.IsLocal(); got != tt.want {
				t.Errorf("Environment.IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_IsDev(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Environment
		want bool
	}{
		{
			name: "dev is dev",
			e:    Environment{EnvDev},
			want: true,
		},
		{
			name: "local is not dev",
			e:    Environment{EnvLocal},
			want: false,
		},
		{
			name: "stg is not dev",
			e:    Environment{EnvStg},
			want: false,
		},
		{
			name: "prd is not dev",
			e:    Environment{EnvPrd},
			want: false,
		},
		{
			name: "zero value is not dev",
			e:    Environment{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.e.IsDev(); got != tt.want {
				t.Errorf("Environment.IsDev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_IsStg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Environment
		want bool
	}{
		{
			name: "stg is stg",
			e:    Environment{EnvStg},
			want: true,
		},
		{
			name: "local is not stg",
			e:    Environment{EnvLocal},
			want: false,
		},
		{
			name: "dev is not stg",
			e:    Environment{EnvDev},
			want: false,
		},
		{
			name: "prd is not stg",
			e:    Environment{EnvPrd},
			want: false,
		},
		{
			name: "zero value is not stg",
			e:    Environment{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.e.IsStg(); got != tt.want {
				t.Errorf("Environment.IsStg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_IsPrd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Environment
		want bool
	}{
		{
			name: "prd is prd",
			e:    Environment{EnvPrd},
			want: true,
		},
		{
			name: "local is not prd",
			e:    Environment{EnvLocal},
			want: false,
		},
		{
			name: "dev is not prd",
			e:    Environment{EnvDev},
			want: false,
		},
		{
			name: "stg is not prd",
			e:    Environment{EnvStg},
			want: false,
		},
		{
			name: "zero value is not prd",
			e:    Environment{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.e.IsPrd(); got != tt.want {
				t.Errorf("Environment.IsPrd() = %v, want %v", got, tt.want)
			}
		})
	}
}
