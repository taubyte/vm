package service

import (
	"context"

	wazero "github.com/taubyte/go-interfaces/vm/wazero"
	"github.com/tetratelabs/wazero/api"
)

/******************************** wazero.Module <> api.Module ****************************************/

var _ wazero.Module = &callContext{}

type callContext struct {
	wazero api.Module
}

func (c *callContext) Name() string {
	return c.wazero.Name()
}

func (c *callContext) Close(ctx context.Context) error {
	return c.wazero.Close(ctx)
}

func (c *callContext) CloseWithExitCode(exitCode uint32) error {
	return c.wazero.CloseWithExitCode(context.TODO(), exitCode)
}

func (c *callContext) Memory() wazero.Memory {
	return &memory{wazero: c.wazero.Memory()}
}

func (c *callContext) ExportedFunction(name string) wazero.Function {
	return &importedFn{wazero: c.wazero.ExportedFunction(name)}
}

func (c *callContext) ExportedMemory(name string) wazero.Memory {
	return &memory{wazero: c.wazero.ExportedMemory(name)}
}

func (c *callContext) ExportedGlobal(name string) wazero.Global {
	return &global{wazero: c.wazero.ExportedGlobal(name)}
}

func (c *callContext) String() string {
	return c.wazero.String()
}

/******************************** wazero.Memory <> api.Memory ****************************************/
type memory struct {
	wazero api.Memory
}

func (m *memory) Size() uint32 {
	return m.wazero.Size()
}

func (m *memory) Grow(deltaPages uint32) (previousPages uint32, ok bool) {
	return m.wazero.Grow(deltaPages)
}

func (m *memory) ReadByte(offset uint32) (byte, bool) {
	return m.wazero.ReadByte(offset)
}

func (m *memory) ReadUint16Le(offset uint32) (uint16, bool) {
	return m.wazero.ReadUint16Le(offset)
}

func (m *memory) ReadUint32Le(offset uint32) (uint32, bool) {
	return m.wazero.ReadUint32Le(offset)
}

func (m *memory) ReadFloat32Le(offset uint32) (float32, bool) {
	return m.wazero.ReadFloat32Le(offset)
}

func (m *memory) ReadUint64Le(offset uint32) (uint64, bool) {
	return m.wazero.ReadUint64Le(offset)
}

func (m *memory) ReadFloat64Le(offset uint32) (float64, bool) {
	return m.wazero.ReadFloat64Le(offset)
}

func (m *memory) Read(offset, byteCount uint32) ([]byte, bool) {
	return m.wazero.Read(offset, byteCount)
}

func (m *memory) ReadString(offset, byteCount uint32) (string, bool) {
	v, ok := m.wazero.Read(offset, byteCount)
	if !ok {
		return "", false
	}
	return string(v), true
}

func (m *memory) WriteByte(offset uint32, v byte) bool {
	return m.wazero.WriteByte(offset, v)
}

func (m *memory) WriteUint16Le(offset uint32, v uint16) bool {
	return m.wazero.WriteUint16Le(offset, v)
}

func (m *memory) WriteUint32Le(offset, v uint32) bool {
	return m.wazero.WriteUint32Le(offset, v)
}

func (m *memory) WriteFloat32Le(offset uint32, v float32) bool {
	return m.wazero.WriteFloat32Le(offset, v)
}

func (m *memory) WriteUint64Le(offset uint32, v uint64) bool {
	return m.wazero.WriteUint64Le(offset, v)
}

func (m *memory) WriteFloat64Le(offset uint32, v float64) bool {
	return m.wazero.WriteFloat64Le(offset, v)
}

func (m *memory) Write(offset uint32, v []byte) bool {
	return m.wazero.Write(offset, v)
}

/******************************** wazero.Function <> api.Function ****************************************/
type importedFn struct {
	wazero api.Function
}

func (f *importedFn) Definition() wazero.FunctionDefinition {
	return &importedFnDef{f.wazero.Definition()}
}

func (f *importedFn) Call(ctx context.Context, params ...uint64) ([]uint64, error) {
	return f.wazero.Call(ctx, params...)
}

type importedFnDef struct {
	wazero api.FunctionDefinition
}

func (f *importedFnDef) ModuleName() string {
	return f.wazero.ModuleName()
}

func (f *importedFnDef) Index() uint32 {
	return f.wazero.Index()
}

func (f *importedFnDef) Name() string {
	return f.wazero.Name()
}

func (f *importedFnDef) Import() (string, string, bool) {
	return f.wazero.Import()
}

func (f *importedFnDef) ExportNames() []string {
	return f.wazero.ExportNames()
}

func (f *importedFnDef) DebugName() string {
	return f.wazero.DebugName()
}

func (f *importedFnDef) ParamTypes() []wazero.ValueType {
	_vt := f.wazero.ParamTypes()
	count := len(_vt)
	vt := make([]wazero.ValueType, count)
	for i := 0; i < count; i++ {
		vt[i] = wazero.ValueType(_vt[i])
	}
	return vt
}

func (f *importedFnDef) ResultTypes() []wazero.ValueType {
	_vt := f.wazero.ResultTypes()
	count := len(_vt)
	vt := make([]wazero.ValueType, count)
	for i := 0; i < count; i++ {
		vt[i] = wazero.ValueType(_vt[i])
	}
	return vt
}

func (f *importedFnDef) ParamNames() []string {
	return f.wazero.ParamNames()
}

/******************************** wazero.Global <> api.Global ****************************************/
type global struct {
	wazero api.Global
}

func (g *global) Type() wazero.ValueType {
	return wazero.ValueType(g.wazero.Type())
}

func (g *global) Get() uint64 {
	return g.wazero.Get()
}

func (g *global) String() string {
	return g.wazero.String()
}
