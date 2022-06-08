package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/saltosystems/winrt-go/winmd"
	"github.com/tdakkota/win32metadata/md"
	"github.com/tdakkota/win32metadata/types"
)

type classNotFoundError struct {
	class string
}

func (e *classNotFoundError) Error() string {
	return fmt.Sprintf("class %s was not found", e.class)
}

type generator struct {
	class       string
	skipFactory bool
	skipStatics bool

	logger log.Logger

	genData  *genData
	winmdCtx *types.Context
}

// Generate generates the code for the given config.
func Generate(cfg *Config, logger log.Logger) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	g := &generator{
		class:       cfg.Class,
		skipFactory: cfg.SkipFactory,
		skipStatics: cfg.SkipStatics,
		logger:      logger,
	}
	return g.run()
}

func (g *generator) run() error {
	_ = level.Debug(g.logger).Log("msg", "starting code generation", "class", g.class)

	winmdFiles, err := winmd.AllFiles()
	if err != nil {
		return err
	}

	// we don't know which winmd file contains the class, so we have to iterate over all of them
	for _, f := range winmdFiles {
		winmdCtx, err := parseWinMDFile(f.Name())
		if err != nil {
			return err
		}
		g.winmdCtx = winmdCtx

		classIndex, err := g.findClass(g.class)
		if err != nil {
			// class not found errors are ok
			if _, ok := err.(*classNotFoundError); ok {
				continue
			}

			return err
		}

		return g.generateClass(g.class, classIndex)
	}

	return fmt.Errorf("class %s was not found", g.class)

}

func (g *generator) findClass(class string) (uint32, error) {
	typeDefTable := g.winmdCtx.Table(md.TypeDef)
	for i := uint32(0); i < typeDefTable.RowCount(); i++ {
		var t types.TypeDef
		if err := t.FromRow(typeDefTable.Row(i)); err != nil {
			return 0, err
		}

		if t.TypeNamespace+"."+t.TypeName == class {
			return i, nil
		}
	}
	return 0, &classNotFoundError{class: class}
}

func (g *generator) generateClass(class string, i uint32) error {
	typeDef, err := g.typeDefByIndex(i)
	if err != nil {
		return err
	}

	// we only support runtime classes: check the tdWindowsRuntime flag (0x4000)
	// https://docs.microsoft.com/en-us/uwp/winrt-cref/winmd-files#runtime-classes
	if typeDef.Flags&0x4000 == 0 {
		return fmt.Errorf("%s.%s is not a runtime class", typeDef.TypeNamespace, typeDef.TypeName)
	}

	_ = level.Info(g.logger).Log("msg", "generating class", "class", class)

	// get templates
	tmpl, err := getTemplates()
	if err != nil {
		return err
	}

	// get data & execute templates

	if err := g.initGenData(typeDef); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "file.tmpl", g.genData); err != nil {
		return err
	}

	// create file & write contents
	folder := typeToFolder(typeDef.TypeNamespace, typeDef.TypeName)
	filename := folder + "/" + strings.ToLower(typeDef.TypeName) + ".go"
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// format the output source code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// write unformatted source code to file as a debugging mechanism
		_, _ = file.Write(buf.Bytes())
		return err
	}

	// and write it to file
	_, err = file.Write(formatted)

	return err
}

func typeToFolder(ns, name string) string {
	fullName := ns + "." + name
	return strings.ToLower(strings.Replace(fullName, ".", "/", -1))
}

func typePackage(ns, name string) string {
	return strings.ToLower(name)
}

func (g *generator) typeDefByIndex(i uint32) (types.TypeDef, error) {
	var typeDef types.TypeDef
	typeDefTable := g.winmdCtx.Table(md.TypeDef)
	if err := typeDef.FromRow(typeDefTable.Row(i)); err != nil {
		return types.TypeDef{}, err
	}
	return typeDef, nil
}

func parseWinMDFile(path string) (*types.Context, error) {
	f, err := winmd.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return types.FromPE(f)
}

func (g *generator) initGenData(runtimeClass types.TypeDef) error {
	// runtime classes have their methods split between three interfaces:
	// Buffer for example
	// - IBuffer (methods)
	// - IBufferFactory (constructors)
	// - IBufferStatics (static methods)
	// and we need to generate types to hold all the methods in separate structures

	g.genData = &genData{
		Package: strings.ToLower(runtimeClass.TypeName),
	}

	// the interface should always exist
	classIndex, err := g.findClass(runtimeClass.TypeNamespace + ".I" + runtimeClass.TypeName)
	if err != nil {
		return err
	}
	typeDefIntf, err := g.typeDefByIndex(classIndex)
	if err != nil {
		return err
	}

	if err := g.addGenType(typeDefIntf, runtimeClass, false); err != nil {
		return err
	}

	if !g.skipFactory {
		// the factory (may not exist)
		classIndex, err = g.findClass(runtimeClass.TypeNamespace + "." + factoryTypeName(typeDefIntf))
		if err == nil {
			typeDefFactory, err := g.typeDefByIndex(classIndex)
			if err != nil {
				return err
			}

			if err := g.addGenType(typeDefFactory, runtimeClass, true); err != nil {
				return err
			}
		}
	}

	if !g.skipStatics {
		// statics (may not exist)
		classIndex, err = g.findClass(runtimeClass.TypeNamespace + "." + staticsTypeName(typeDefIntf))
		if err == nil {
			typeDefStatics, err := g.typeDefByIndex(classIndex)
			if err != nil {
				return err
			}
			if err := g.addGenType(typeDefStatics, runtimeClass, true); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) addType(gt genType) {
	g.genData.Types = append(g.genData.Types, gt)
}

func (g *generator) addImportFor(ns, name string) {
	folder := typeToFolder(ns, name)
	i := "github.com/saltosystems/winrt-go/" + folder
	g.genData.Imports = append(g.genData.Imports, i)
}

func (g *generator) addGenType(typeDef, runtimeClass types.TypeDef, activatable bool) error {
	// class interface
	funcs, err := g.getGenFuncs(typeDef, runtimeClass, activatable)
	if err != nil {
		return err
	}

	g.addType(genType{
		Name:  typeDef.TypeName,
		Funcs: funcs,
	})
	return nil
}

func (g *generator) getGenFuncs(typeDef, runtimeClass types.TypeDef, activatable bool) ([]genFunc, error) {
	var genFuncs []genFunc

	if activatable {
		var parentGUID string
		parentGUID, err := g.typeGUID(typeDef)
		if err != nil {
			// the type may not contain this information, just ignore it
			_ = level.Warn(g.logger).Log("msg", "failed to get type GUID", "type", typeDef.TypeNamespace+"."+typeDef.TypeName, "err", err)
		}

		// add the constructor
		genFuncs = append(genFuncs, genFunc{
			Name:          "Activate" + typeDef.TypeName,
			IsConstructor: true,
			InParams:      []genParam{},
			ReturnParam: &genParam{
				Type:         "*" + typeDef.TypeName,
				DefaultValue: "nil",
			},
			Signature:      nil,
			ParentType:     typeDef,
			ParentTypeGUID: parentGUID,
			RuntimeClass:   runtimeClass,
			FuncOwner:      "",
		})
	}

	methods, err := typeDef.ResolveMethodList(g.winmdCtx)
	if err != nil {
		return nil, err
	}

	for _, m := range methods {
		generatedFunc, err := g.genFuncFromMethod(typeDef, runtimeClass, m)
		if err != nil {
			return nil, err
		}
		genFuncs = append(genFuncs, *generatedFunc)
	}

	return genFuncs, nil
}

func (g *generator) genFuncFromMethod(typeDef, runtimeClass types.TypeDef, m types.MethodDef) (*genFunc, error) {
	params, err := g.getInParameters(typeDef, runtimeClass, m)
	if err != nil {
		return nil, err
	}

	retParam, err := g.getReturnParameters(typeDef, runtimeClass, m)
	if err != nil {
		return nil, err
	}

	var parentGUID string
	parentGUID, err = g.typeGUID(typeDef)
	if err != nil {
		// the type may not contain this information, just ignore it
		_ = level.Warn(g.logger).Log("msg", "failed to get type GUID", "type", typeDef.TypeNamespace+"."+typeDef.TypeName, "err", err)
	}

	return &genFunc{
		Name:           m.Name,
		IsConstructor:  false,
		InParams:       params,
		ReturnParam:    retParam,
		Signature:      m.Signature,
		ParentType:     typeDef,
		ParentTypeGUID: parentGUID,
		RuntimeClass:   runtimeClass,
		FuncOwner:      typeDef.TypeName,
	}, nil
}

func (g *generator) typeGUID(t types.TypeDef) (string, error) {
	// GUIDs are stored in custom attributes.
	// To find the GUID of the given type, we need to iterate
	// through all the custom attributes and find the one that
	// matches the type
	tableCustomAttributes := g.winmdCtx.Table(md.CustomAttribute)
	for i := uint32(0); i < tableCustomAttributes.RowCount(); i++ {
		var cAttr types.CustomAttribute
		if err := cAttr.FromRow(tableCustomAttributes.Row(i)); err != nil {
			continue
		}

		if cAttrParentTable, ok := cAttr.Parent.Table(); !ok || cAttrParentTable != md.TypeDef {
			continue
		}
		row, ok := cAttr.Parent.Row(g.winmdCtx)
		if !ok {
			continue // something failed
		}
		var parentTypeDef types.TypeDef
		if err := parentTypeDef.FromRow(row); err != nil {
			continue
		}

		if parentTypeDef.TypeNamespace == t.TypeNamespace && parentTypeDef.TypeName == t.TypeName {
			// a type may have more than one blob, so do not immediately fail
			guid, err := guidBlobToString(cAttr.Value)
			if err != nil {
				continue
			}
			return guid, nil
		}
	}
	return "", fmt.Errorf("GUID not found for type %s.%s", t.TypeNamespace, t.TypeName)
}

// guidBlobToString converts an array into the textual representation of a GUID
func guidBlobToString(b types.Blob) (string, error) {
	// the blob is surrounded by a header (0x01, 0x00) and a footer (0x00, 0x00). Remove them
	guid := b[2 : len(b)-2]
	if len(guid) != 16 {
		return "", fmt.Errorf("invalid GUID blob length: %d", len(guid))
	}
	// the string version has 5 parts separated by '-'
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%04x%08x",
		// The first 3 are encoded as little endian
		uint32(guid[0])|uint32(guid[1])<<8|uint32(guid[2])<<16|uint32(guid[3])<<24,
		uint16(guid[4])|uint16(guid[5])<<8,
		uint16(guid[6])|uint16(guid[7])<<8,
		//the rest is not
		uint16(guid[8])<<8|uint16(guid[9]),
		uint16(guid[10])<<8|uint16(guid[11]),
		uint32(guid[12])<<24|uint32(guid[13])<<16|uint32(guid[14])<<8|uint32(guid[15])), nil
}

func factoryTypeName(t types.TypeDef) string {
	return t.TypeName + "Factory"
}

func staticsTypeName(t types.TypeDef) string {
	return t.TypeName + "Statics"
}

func (g *generator) getInParameters(t, rt types.TypeDef, m types.MethodDef) ([]genParam, error) {

	params, err := m.ResolveParamList(g.winmdCtx)
	if err != nil {
		return nil, err
	}

	// the signature contains the parameter
	// types and return type of the method
	r := m.Signature.Reader()
	mr, err := r.Method(g.winmdCtx)
	if err != nil {
		return nil, err
	}

	genParams := []genParam{}
	for i, e := range mr.Params {
		genParams = append(genParams, genParam{
			Name: getParamName(params, uint16(i+1)),
			Type: g.elementType(e, typePackage(rt.TypeNamespace, rt.TypeName)),
		})
	}

	return genParams, nil
}

func (g *generator) getReturnParameters(t, rt types.TypeDef, m types.MethodDef) (*genParam, error) {
	// the signature contains the parameter
	// types and return type of the method
	r := m.Signature.Reader()
	methodSignature, err := r.Method(g.winmdCtx)
	if err != nil {
		return nil, err
	}

	// ignore void types
	if methodSignature.Return.Type.Kind == types.ELEMENT_TYPE_VOID {
		return nil, nil
	}

	return &genParam{
		Name:         "",
		Type:         g.elementType(methodSignature.Return, typePackage(rt.TypeNamespace, rt.TypeName)),
		DefaultValue: g.elementDefaultValue(methodSignature.Return),
	}, nil
}

func getParamName(params []types.Param, i uint16) string {
	for _, p := range params {
		if p.Flags.In() && p.Sequence == i {
			return p.Name
		}
	}
	return "__ERROR__"
}

func (g *generator) elementType(e types.Element, currentPkg string) string {
	switch e.Type.Kind {
	case types.ELEMENT_TYPE_BOOLEAN:
		return "bool"
	case types.ELEMENT_TYPE_CHAR:
		return "byte"
	case types.ELEMENT_TYPE_I1:
		return "int8"
	case types.ELEMENT_TYPE_U1:
		return "uint8"
	case types.ELEMENT_TYPE_I2:
		return "int16"
	case types.ELEMENT_TYPE_U2:
		return "uint16"
	case types.ELEMENT_TYPE_I4:
		return "int32"
	case types.ELEMENT_TYPE_U4:
		return "uint32"
	case types.ELEMENT_TYPE_I8:
		return "int64"
	case types.ELEMENT_TYPE_U8:
		return "uint64"
	case types.ELEMENT_TYPE_R4:
		return "float32"
	case types.ELEMENT_TYPE_R8:
		return "float64"
	case types.ELEMENT_TYPE_STRING:
		return "string"
	case types.ELEMENT_TYPE_CLASS:
		// return class name
		namespace, name, err := g.winmdCtx.ResolveTypeDefOrRefName(e.Type.TypeDef.Index)
		if err != nil {
			return "__ERROR_ELEMENT_TYPE_CLASS__"
		}

		// this may be the runtime class itself, but we only know the interfaces
		// TODO: make this smarter
		if !strings.HasPrefix(name, "I") {
			name = "I" + name
		}

		// name is always an interface, so we need to remove the initial I
		typePkg := typePackage(namespace, name[1:])
		if currentPkg != typePkg {
			g.addImportFor(namespace, name[1:])
			name = typePkg + "." + name
		}

		return "*" + name
	default:
		return "__ERROR_" + e.Type.Kind.String() + "__"
	}
}

func (g *generator) elementDefaultValue(e types.Element) string {
	switch e.Type.Kind {
	case types.ELEMENT_TYPE_BOOLEAN:
		return "false"
	case types.ELEMENT_TYPE_CHAR:
		fallthrough
	case types.ELEMENT_TYPE_I1:
		fallthrough
	case types.ELEMENT_TYPE_U1:
		fallthrough
	case types.ELEMENT_TYPE_I2:
		fallthrough
	case types.ELEMENT_TYPE_U2:
		fallthrough
	case types.ELEMENT_TYPE_I4:
		fallthrough
	case types.ELEMENT_TYPE_U4:
		fallthrough
	case types.ELEMENT_TYPE_I8:
		fallthrough
	case types.ELEMENT_TYPE_U8:
		return "0"
	case types.ELEMENT_TYPE_R4:
		fallthrough
	case types.ELEMENT_TYPE_R8:
		return "0.0"
	case types.ELEMENT_TYPE_STRING:
		return "\"\""
	case types.ELEMENT_TYPE_CLASS:
		return "nil"
	default:
		return "__ERROR_" + e.Type.Kind.String() + "__"
	}
}