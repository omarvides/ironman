package bleve

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/ironman-project/ironman/template/model"
	"github.com/ironman-project/ironman/template/repository"
	uuid "github.com/satori/go.uuid"
)

func newBleveTestIndex(t *testing.T) (bleve.Index, func()) {
	dir, err := ioutil.TempDir("", "ironman-bleve-test")
	if err != nil {
		t.Fatal("Failed to create test bleve index directory", err)
	}

	indexPath := filepath.Join(dir, "index")
	indexMapping := bleve.NewIndexMapping()

	templateDocMapping := bleve.NewDocumentMapping()

	templateIDMapping := bleve.NewTextFieldMapping()

	templateDocMapping.AddFieldMappingsAt("ID", templateIDMapping)
	generatorsMapping := bleve.NewDocumentMapping()
	templateDocMapping.AddSubDocumentMapping("Generators", generatorsMapping)

	indexMapping.AddDocumentMapping("model.template", templateDocMapping)

	index, err := bleve.New(indexPath, indexMapping)

	if err != nil {
		t.Fatal("Failed to create test bleve index", err)
	}
	return index, func() {
		err := index.Close()

		if err != nil {
			t.Fatal("Failed to close bleve index", err)
		}
		err = os.RemoveAll(dir)
		if err != nil {
			t.Fatal("Failed to clean bleve index", err)
		}
	}
}

func newTestRepository(t *testing.T) (repository.Repository, bleve.Index, func()) {
	index, clean := newBleveTestIndex(t)
	r := New(SetIndex(index))
	return r, index, clean
}

func Test_bleeveRepository_Index(t *testing.T) {

	type args struct {
		template model.Template
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Index a template",
			args{
				model.Template{
					ID:          "test-template",
					Name:        "Test template",
					Description: "This is a test template",
					Generators: []*model.Generator{
						&model.Generator{
							ID:          "test-generator",
							Name:        "Test generator",
							Description: "This is a test generator",
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, index, clean := newTestRepository(t)
			defer clean()
			var id string
			var err error
			if id, err = r.Index(tt.args.template); (err != nil) != tt.wantErr {
				t.Errorf("bleeveRepository.Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			doc, err := index.Document(id)

			if err != nil {
				t.Errorf("bleeveRepository.Index() error = %v, wantErr %v", err, tt.wantErr)
			}

			if doc == nil {
				t.Errorf("bleeveRepository.Index() nil , want %v", tt.args.template)
			}
		})
	}
}

func Test_bleeveRepository_Update(t *testing.T) {

	type args struct {
		template model.Template
	}
	tests := []struct {
		name     string
		args     args
		template model.Template
		wantErr  bool
	}{
		{
			"Update template index",
			args{
				model.Template{
					ID:          "template-id",
					Name:        "Updated name",
					Description: "Updated description",
					Generators: []*model.Generator{
						&model.Generator{
							ID:          "test-generator",
							Name:        "Test generator",
							Description: "This is a test generator",
						},
						&model.Generator{
							ID:          "test-generator2",
							Name:        "Test generator2",
							Description: "This is a test generator 2",
						},
					},
				},
			},
			model.Template{
				ID: "template-id",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, index, clean := newTestRepository(t)
			defer clean()
			id := uuid.NewV4().String()
			err := index.Index(id, tt.template)

			if err != nil {
				t.Error("Failed to index template to update", err)
			}

			if err := r.Update(tt.args.template); (err != nil) != tt.wantErr {
				t.Errorf("bleeveRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			doc, err := index.Document(id)

			if doc == nil {
				t.Error("failed to retreive indexed document with ID", id)
			}

			if err != nil {
				t.Error("failed to retreive indexed document", tt.template, err)
			}

			for _, field := range doc.Fields {

				value := string(field.Value())

				switch field.Name() {

				case "ID":
					if string(value) == "" || (value != tt.args.template.ID) {
						t.Errorf("bleveRepository.Update() templateID = %v want %v", value, tt.args.template.ID)
					}
				case "Name":
					if string(value) == "" || (value != tt.args.template.Name) {
						t.Errorf("bleveRepository.Update() templateName = %v want %v", value, tt.args.template.Name)
					}
				case "Description":
					if string(value) == "" || (value != tt.args.template.Description) {
						t.Errorf("bleveRepository.Update() templateDescription = %v want %v", value, tt.args.template.Description)
					}
				case "Generators.ID":
					pos := field.ArrayPositions()[0]
					expectedID := tt.args.template.Generators[pos].ID
					if value != expectedID {
						t.Errorf("bleveRepository.Update() template.Generators[%d].ID = %v want %v", pos, value, expectedID)
					}
				case "Generators.Name":
					pos := field.ArrayPositions()[0]
					expectedName := tt.args.template.Generators[pos].Name
					if value != expectedName {
						t.Errorf("bleveRepository.Update() template.Generators[%d].Name = %v want %v", pos, value, expectedName)
					}
				case "Generators.Description":
					pos := field.ArrayPositions()[0]
					expectedDescription := tt.args.template.Generators[pos].Description
					if value != expectedDescription {
						t.Errorf("bleveRepository.Update() template.Generators[%d].Description = %v want %v", pos, value, expectedDescription)
					}
				default:
					t.Error("doc.Fields should assert field", field.Name(), string(field.Value()))
				}
			}

		})
	}
}