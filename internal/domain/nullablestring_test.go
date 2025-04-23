package domain

import (
	"database/sql/driver"
	"testing"
)

func TestNullableString_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    NullableString
		wantErr bool
	}{
		{
			name:    "nil input",
			input:   nil,
			want:    "",
			wantErr: false,
		},
		{
			name:    "valid string",
			input:   "test string",
			want:    "test string",
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   123,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ns NullableString
			err := ns.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullableString.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ns != tt.want {
				t.Errorf("NullableString.Scan() = %v, want %v", ns, tt.want)
			}
		})
	}
}

func TestNullableString_Value(t *testing.T) {
	tests := []struct {
		name    string
		ns      NullableString
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "empty string",
			ns:      "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "non-empty string",
			ns:      "test string",
			want:    "test string",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ns.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullableString.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("NullableString.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
