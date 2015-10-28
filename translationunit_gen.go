package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// A single translation unit, which resides in an index.
type TranslationUnit struct {
	c C.CXTranslationUnit
}

// Determine whether the given header is guarded against multiple inclusions, either with the conventional \#ifndef/\#define/\#endif macro guards or with \#pragma once.
func (tu TranslationUnit) IsFileMultipleIncludeGuarded(file File) bool {
	o := C.clang_isFileMultipleIncludeGuarded(tu.c, file.c)

	return o != C.uint(0)
}

// Retrieve a file handle within the given translation unit. \param tu the translation unit \param file_name the name of the file. \returns the file handle for the named file in the translation unit \p tu, or a NULL file handle if the file was not a part of this translation unit.
func (tu TranslationUnit) File(file_name string) File {
	c_file_name := C.CString(file_name)
	defer C.free(unsafe.Pointer(c_file_name))

	return File{C.clang_getFile(tu.c, c_file_name)}
}

// Retrieves the source location associated with a given file/line/column in a particular translation unit.
func (tu TranslationUnit) Location(file File, line uint16, column uint16) SourceLocation {
	return SourceLocation{C.clang_getLocation(tu.c, file.c, C.uint(line), C.uint(column))}
}

// Retrieves the source location associated with a given character offset in a particular translation unit.
func (tu TranslationUnit) LocationForOffset(file File, offset uint16) SourceLocation {
	return SourceLocation{C.clang_getLocationForOffset(tu.c, file.c, C.uint(offset))}
}

// Determine the number of diagnostics produced for the given translation unit.
func (tu TranslationUnit) NumDiagnostics() uint16 {
	return uint16(C.clang_getNumDiagnostics(tu.c))
}

// Retrieve a diagnostic associated with the given translation unit. \param Unit the translation unit to query. \param Index the zero-based diagnostic number to retrieve. \returns the requested diagnostic. This diagnostic must be freed via a call to \c clang_disposeDiagnostic().
func (tu TranslationUnit) Diagnostic(Index uint16) Diagnostic {
	return Diagnostic{C.clang_getDiagnostic(tu.c, C.uint(Index))}
}

// Retrieve the complete set of diagnostics associated with a translation unit. \param Unit the translation unit to query.
func (tu TranslationUnit) DiagnosticSetFromTU() DiagnosticSet {
	return DiagnosticSet{C.clang_getDiagnosticSetFromTU(tu.c)}
}

// Get the original translation unit source file name.
func (tu TranslationUnit) Spelling() string {
	o := cxstring{C.clang_getTranslationUnitSpelling(tu.c)}
	defer o.Dispose()

	return o.String()
}

// Returns the set of flags that is suitable for saving a translation unit. The set of flags returned provide options for \c clang_saveTranslationUnit() by default. The returned flag set contains an unspecified set of options that save translation units with the most commonly-requested data.
func (tu TranslationUnit) DefaultSaveOptions() uint16 {
	return uint16(C.clang_defaultSaveOptions(tu.c))
}

// Saves a translation unit into a serialized representation of that translation unit on disk. Any translation unit that was parsed without error can be saved into a file. The translation unit can then be deserialized into a new \c CXTranslationUnit with \c clang_createTranslationUnit() or, if it is an incomplete translation unit that corresponds to a header, used as a precompiled header when parsing other translation units. \param TU The translation unit to save. \param FileName The file to which the translation unit will be saved. \param options A bitmask of options that affects how the translation unit is saved. This should be a bitwise OR of the CXSaveTranslationUnit_XXX flags. \returns A value that will match one of the enumerators of the CXSaveError enumeration. Zero (CXSaveError_None) indicates that the translation unit was saved successfully, while a non-zero value indicates that a problem occurred.
func (tu TranslationUnit) SaveTranslationUnit(FileName string, options uint16) uint16 {
	c_FileName := C.CString(FileName)
	defer C.free(unsafe.Pointer(c_FileName))

	return uint16(C.clang_saveTranslationUnit(tu.c, c_FileName, C.uint(options)))
}

// Destroy the specified CXTranslationUnit object.
func (tu TranslationUnit) Dispose() {
	C.clang_disposeTranslationUnit(tu.c)
}

// Returns the set of flags that is suitable for reparsing a translation unit. The set of flags returned provide options for \c clang_reparseTranslationUnit() by default. The returned flag set contains an unspecified set of optimizations geared toward common uses of reparsing. The set of optimizations enabled may change from one version to the next.
func (tu TranslationUnit) DefaultReparseOptions() uint16 {
	return uint16(C.clang_defaultReparseOptions(tu.c))
}

// Return the memory usage of a translation unit. This object should be released with clang_disposeCXTUResourceUsage().
func (tu TranslationUnit) TUResourceUsage() TUResourceUsage {
	return TUResourceUsage{C.clang_getCXTUResourceUsage(tu.c)}
}

// Retrieve the cursor that represents the given translation unit. The translation unit cursor can be used to start traversing the various declarations within the given translation unit.
func (tu TranslationUnit) TranslationUnitCursor() Cursor {
	return Cursor{C.clang_getTranslationUnitCursor(tu.c)}
}

// Map a source location to the cursor that describes the entity at that location in the source code. clang_getCursor() maps an arbitrary source location within a translation unit down to the most specific cursor that describes the entity at that location. For example, given an expression \c x + y, invoking clang_getCursor() with a source location pointing to "x" will return the cursor for "x"; similarly for "y". If the cursor points anywhere between "x" or "y" (e.g., on the + or the whitespace around it), clang_getCursor() will return a cursor referring to the "+" expression. \returns a cursor representing the entity at the given source location, or a NULL cursor if no such entity can be found.
func (tu TranslationUnit) Cursor(sl SourceLocation) Cursor {
	return Cursor{C.clang_getCursor(tu.c, sl.c)}
}

// \param Module a module object. \returns the number of top level headers associated with this module.
func (tu TranslationUnit) Module_getNumTopLevelHeaders(Module Module) uint16 {
	return uint16(C.clang_Module_getNumTopLevelHeaders(tu.c, Module.c))
}

// \param Module a module object. \param Index top level header index (zero-based). \returns the specified top level header associated with the module.
func (tu TranslationUnit) Module_getTopLevelHeader(Module Module, Index uint16) File {
	return File{C.clang_Module_getTopLevelHeader(tu.c, Module.c, C.uint(Index))}
}

// Determine the spelling of the given token. The spelling of a token is the textual representation of that token, e.g., the text of an identifier or keyword.
func (tu TranslationUnit) TokenSpelling(t Token) string {
	o := cxstring{C.clang_getTokenSpelling(tu.c, t.c)}
	defer o.Dispose()

	return o.String()
}

// Retrieve the source location of the given token.
func (tu TranslationUnit) TokenLocation(t Token) SourceLocation {
	return SourceLocation{C.clang_getTokenLocation(tu.c, t.c)}
}

// Retrieve a source range that covers the given token.
func (tu TranslationUnit) TokenExtent(t Token) SourceRange {
	return SourceRange{C.clang_getTokenExtent(tu.c, t.c)}
}

// Find #import/#include directives in a specific file. \param TU translation unit containing the file to query. \param file to search for #import/#include directives. \param visitor callback that will receive pairs of CXCursor/CXSourceRange for each directive found. \returns one of the CXResult enumerators.
func (tu TranslationUnit) FindIncludesInFile(file File, visitor CursorAndRangeVisitor) Result {
	return Result(C.clang_findIncludesInFile(tu.c, file.c, visitor.c))
}
