package git

import (
	"os"
	"testing"

	"github.com/ironman-project/ironman/template/repository"
)

func newTestGitRepository() repository.Repository {
	baseRepository := repository.NewBaseRepository("testing")
	return New(baseRepository)
}

func TestRepository_Install(t *testing.T) {

	type args struct {
		location string
	}
	tests := []struct {
		name               string
		args               args
		expectedTemplateID string
		wantErr            bool
	}{
		{
			"Install template",
			args{"https://github.com/ottogiron/wizard-hello-world.git"},
			"wizard-hello-world",
			false,
		},
		{
			"Install template",
			args{"https://github.com/ottogiron/unexisting-template"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestGitRepository()
			defer func() {
				_ = r.Uninstall(tt.expectedTemplateID)
			}()
			if err := r.Install(tt.args.location); (err != nil) != tt.wantErr {
				t.Errorf("Repository.Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			expectedTemplatePath := r.TemplatePath(tt.expectedTemplateID)
			if _, err := os.Stat(expectedTemplatePath); os.IsNotExist(err) {
				t.Errorf("Repository.Install() template was not installed want path %v", expectedTemplatePath)
			}
		})
	}
}

func TestRepository_Upgrade(t *testing.T) {

	type args struct {
		id       string
		location string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Upgrade template",
			args{"wizard-hello-world", "https://github.com/ottogiron/wizard-hello-world.git"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestGitRepository()

			err := r.Install(tt.args.location)

			defer func() {
				r.Uninstall(tt.args.id)
			}()

			if err != nil {
				t.Errorf("Repository.Upgrade() error = %v", err)
			}

			if err := r.Upgrade(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Repository.Upgrade() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
