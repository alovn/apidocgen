package apidoc

import (
	"go/ast"
	"path/filepath"
	"sort"
	"strings"
)

type PackagesDefinitions struct {
	files             map[*ast.File]*AstFileInfo
	packages          map[string]*PackageDefinitions
	uniqueDefinitions map[string]*TypeSpecDef
}

// NewPackagesDefinitions create object PackagesDefinitions.
func NewPackagesDefinitions() *PackagesDefinitions {
	return &PackagesDefinitions{
		files:             make(map[*ast.File]*AstFileInfo),
		packages:          make(map[string]*PackageDefinitions),
		uniqueDefinitions: make(map[string]*TypeSpecDef),
	}
}

// CollectAstFile collect ast.file.
func (pkgDefs *PackagesDefinitions) CollectAstFile(packageDir, path string, astFile *ast.File) error {
	if pkgDefs.files == nil {
		pkgDefs.files = make(map[*ast.File]*AstFileInfo)
	}

	if pkgDefs.packages == nil {
		pkgDefs.packages = make(map[string]*PackageDefinitions)
	}

	// return without storing the file if we lack a packageDir
	if packageDir == "" {
		return nil
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dependency, ok := pkgDefs.packages[packageDir]
	if ok {
		// return without storing the file if it already exists
		_, exists := dependency.Files[path]
		if exists {
			return nil
		}

		dependency.Files[path] = astFile
	} else {
		pkgDefs.packages[packageDir] = &PackageDefinitions{
			Name:            astFile.Name.Name,
			Files:           map[string]*ast.File{path: astFile},
			TypeDefinitions: make(map[string]*TypeSpecDef),
		}
	}

	pkgDefs.files[astFile] = &AstFileInfo{
		File:        astFile,
		Path:        path,
		PackagePath: packageDir,
	}

	return nil
}

func (pkgDefs *PackagesDefinitions) findTypeSpec(pkgPath string, typeName string) *TypeSpecDef {
	if pkgDefs.packages == nil {
		return nil
	}

	pd, found := pkgDefs.packages[pkgPath]
	if found {
		typeSpec, ok := pd.TypeDefinitions[typeName]
		if ok {
			return typeSpec
		}
	}

	return nil
}

// RangeFiles for range the collection of ast.File in alphabetic order.
func rangeFiles(files map[*ast.File]*AstFileInfo, handle func(filename string, file *ast.File) error) error {
	sortedFiles := make([]*AstFileInfo, 0, len(files))
	for _, info := range files {
		sortedFiles = append(sortedFiles, info)
	}

	sort.Slice(sortedFiles, func(i, j int) bool {
		return strings.Compare(sortedFiles[i].Path, sortedFiles[j].Path) < 0
	})

	for _, info := range sortedFiles {
		err := handle(info.Path, info.File)
		if err != nil {
			return err
		}
	}

	return nil
}
