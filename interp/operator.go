package interp

import (
	"go/token"
	"log"
	"reflect"

	"golang.org/x/tools/go/exact"
	"golang.org/x/tools/go/types"
)

// Map from assignment operation token (op=) to operator (op)
var assignOps = map[token.Token]token.Token{
	// Arithmetic
	token.ADD_ASSIGN: token.ADD,
	token.SUB_ASSIGN: token.SUB,
	token.MUL_ASSIGN: token.MUL,
	token.QUO_ASSIGN: token.QUO,
	token.REM_ASSIGN: token.REM,

	// Bitwise
	token.AND_ASSIGN:     token.AND,
	token.AND_NOT_ASSIGN: token.AND_NOT,
	token.OR_ASSIGN:      token.OR,
	token.SHL_ASSIGN:     token.SHL,
	token.SHR_ASSIGN:     token.SHR,
	token.XOR_ASSIGN:     token.XOR,
}

func doBinaryOp(env *environ, left, right Object, op token.Token) Object {
	var obj Object
	switch op {
	case token.ADD:
		obj = operatorAdd(env, left, right)
	case token.SUB:
		obj = operatorSubtract(env, left, right)
	case token.MUL:
		obj = operatorMultiply(env, left, right)
	case token.QUO:
		obj = operatorQuotient(env, left, right)
	case token.REM:
		obj = operatorRemainder(env, left, right)
	case token.AND:
		obj = operatorAnd(env, left, right)
	case token.OR:
		obj = operatorOr(env, left, right)
	case token.XOR:
		obj = operatorXor(env, left, right)
	case token.AND_NOT:
		obj = operatorAndNot(env, left, right)
	case token.SHR:
		obj = operatorShiftRight(env, left, right)
	case token.SHL:
		obj = operatorShiftLeft(env, left, right)
	default:
		// TODO: Implement other binary operators
		log.Fatalf("Binary operator %v not implemented yet", op)
	}
	return obj
}

func doBinaryComparisonOp(env *environ, left, right Object, op token.Token, typ types.Type) Object {
	var obj Object
	switch op {
	case token.LSS:
		obj = operatorLess(env, left, right, typ)
	case token.GTR:
		obj = operatorGreater(env, left, right, typ)
	case token.LEQ:
		obj = operatorLessEqual(env, left, right, typ)
	case token.GEQ:
		obj = operatorGreaterEqual(env, left, right, typ)
	case token.EQL:
		obj = operatorEqual(env, left, right, typ)
	default:
		// TODO: Implement other binary operators
		log.Fatalf("Binary comparison operator %v not implemented yet", op)
	}
	return obj
}

func getTypedObject(obj Object) Object {
	if isTyped(obj.Typ) {
		return obj
	}
	t := obj.Typ.Underlying().(*types.Basic)
	ev := obj.Value.(exact.Value)
	switch t.Kind() {
	case types.UntypedBool:
		b := exact.BoolVal(ev)
		return Object{
			Value: reflect.ValueOf(b),
			Typ:   types.Typ[types.Bool],
		}
	case types.UntypedInt:
		i64, _ := exact.Int64Val(ev)
		return Object{
			Value: reflect.ValueOf(int(i64)),
			Typ:   types.Typ[types.Int],
		}
	case types.UntypedFloat:
		f64, _ := exact.Float64Val(ev)
		return Object{
			Value: reflect.ValueOf(f64),
			Typ:   types.Typ[types.Float64],
		}
	case types.UntypedComplex:
		real64, _ := exact.Float64Val(exact.Real(ev))
		imag64, _ := exact.Float64Val(exact.Imag(ev))
		c128 := complex(real64, imag64)
		return Object{
			Value: reflect.ValueOf(c128),
			Typ:   types.Typ[types.Complex128],
		}
	case types.UntypedRune:
		i64, _ := exact.Int64Val(ev)
		r := rune(i64)
		return Object{
			Value: reflect.ValueOf(r),
			Typ:   types.Typ[types.Rune],
		}
	case types.UntypedString:
		s := exact.StringVal(ev)
		return Object{
			Value: reflect.ValueOf(s),
			Typ:   types.Typ[types.String],
		}
	case types.UntypedNil:
		log.Fatal("getTypedObject: Got untyped nil")
	}
	return obj
}

// operatorAdd implements the binary operation '+'.
// If this is being called at all, then the left and right objects
// have a value that's a "reflect.Value". Also, both have the same type,
// since the expression passed type checking.
func operatorAdd(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorAdd: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		sum := int(lv.Int()) + int(rv.Int())
		newVal.SetInt(int64(sum))
	case reflect.Int8:
		sum := int8(lv.Int()) + int8(rv.Int())
		newVal.SetInt(int64(sum))
	case reflect.Int16:
		sum := int16(lv.Int()) + int16(rv.Int())
		newVal.SetInt(int64(sum))
	case reflect.Int32:
		sum := int32(lv.Int()) + int32(rv.Int())
		newVal.SetInt(int64(sum))
	case reflect.Int64:
		sum := int64(lv.Int()) + int64(rv.Int())
		newVal.SetInt(int64(sum))
	case reflect.Uint:
		sum := uint(lv.Uint()) + uint(rv.Uint())
		newVal.SetUint(uint64(sum))
	case reflect.Uint8:
		sum := uint8(lv.Uint()) + uint8(rv.Uint())
		newVal.SetUint(uint64(sum))
	case reflect.Uint16:
		sum := uint16(lv.Uint()) + uint16(rv.Uint())
		newVal.SetUint(uint64(sum))
	case reflect.Uint32:
		sum := uint32(lv.Uint()) + uint32(rv.Uint())
		newVal.SetUint(uint64(sum))
	case reflect.Uint64:
		sum := uint64(lv.Uint()) + uint64(rv.Uint())
		newVal.SetUint(uint64(sum))
	case reflect.Uintptr:
		sum := uintptr(lv.Uint()) + uintptr(rv.Uint())
		newVal.SetUint(uint64(sum))
	case reflect.Float32:
		sum := float32(lv.Float()) + float32(rv.Float())
		newVal.SetFloat(float64(sum))
	case reflect.Float64:
		sum := float64(lv.Float()) + float64(rv.Float())
		newVal.SetFloat(float64(sum))
	case reflect.Complex64:
		sum := complex64(lv.Complex()) + complex64(rv.Complex())
		newVal.SetComplex(complex128(sum))
	case reflect.Complex128:
		sum := complex128(lv.Complex()) + complex128(rv.Complex())
		newVal.SetComplex(complex128(sum))
	case reflect.String:
		sum := lv.String() + rv.String()
		newVal.SetString(sum)
	default:
		panic("Type error: Invalid operands to addition: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}

	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorSubtract implements the binary operation '-'.
func operatorSubtract(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorSubtract: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		diff := int(lv.Int()) - int(rv.Int())
		newVal.SetInt(int64(diff))
	case reflect.Int8:
		diff := int8(lv.Int()) - int8(rv.Int())
		newVal.SetInt(int64(diff))
	case reflect.Int16:
		diff := int16(lv.Int()) - int16(rv.Int())
		newVal.SetInt(int64(diff))
	case reflect.Int32:
		diff := int32(lv.Int()) - int32(rv.Int())
		newVal.SetInt(int64(diff))
	case reflect.Int64:
		diff := int64(lv.Int()) - int64(rv.Int())
		newVal.SetInt(int64(diff))
	case reflect.Uint:
		diff := uint(lv.Uint()) - uint(rv.Uint())
		newVal.SetUint(uint64(diff))
	case reflect.Uint8:
		diff := uint8(lv.Uint()) - uint8(rv.Uint())
		newVal.SetUint(uint64(diff))
	case reflect.Uint16:
		diff := uint16(lv.Uint()) - uint16(rv.Uint())
		newVal.SetUint(uint64(diff))
	case reflect.Uint32:
		diff := uint32(lv.Uint()) - uint32(rv.Uint())
		newVal.SetUint(uint64(diff))
	case reflect.Uint64:
		diff := uint64(lv.Uint()) - uint64(rv.Uint())
		newVal.SetUint(uint64(diff))
	case reflect.Uintptr:
		diff := uintptr(lv.Uint()) - uintptr(rv.Uint())
		newVal.SetUint(uint64(diff))
	case reflect.Float32:
		diff := float32(lv.Float()) - float32(rv.Float())
		newVal.SetFloat(float64(diff))
	case reflect.Float64:
		diff := float64(lv.Float()) - float64(rv.Float())
		newVal.SetFloat(float64(diff))
	case reflect.Complex64:
		diff := complex64(lv.Complex()) - complex64(rv.Complex())
		newVal.SetComplex(complex128(diff))
	case reflect.Complex128:
		diff := complex128(lv.Complex()) - complex128(rv.Complex())
		newVal.SetComplex(complex128(diff))
	default:
		panic("Type error: Invalid operands to subtraction: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorMultiply implements the binary operation '*'.
func operatorMultiply(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorMultiply: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		prod := int(lv.Int()) * int(rv.Int())
		newVal.SetInt(int64(prod))
	case reflect.Int8:
		prod := int8(lv.Int()) * int8(rv.Int())
		newVal.SetInt(int64(prod))
	case reflect.Int16:
		prod := int16(lv.Int()) * int16(rv.Int())
		newVal.SetInt(int64(prod))
	case reflect.Int32:
		prod := int32(lv.Int()) * int32(rv.Int())
		newVal.SetInt(int64(prod))
	case reflect.Int64:
		prod := int64(lv.Int()) * int64(rv.Int())
		newVal.SetInt(int64(prod))
	case reflect.Uint:
		prod := uint(lv.Uint()) * uint(rv.Uint())
		newVal.SetUint(uint64(prod))
	case reflect.Uint8:
		prod := uint8(lv.Uint()) * uint8(rv.Uint())
		newVal.SetUint(uint64(prod))
	case reflect.Uint16:
		prod := uint16(lv.Uint()) * uint16(rv.Uint())
		newVal.SetUint(uint64(prod))
	case reflect.Uint32:
		prod := uint32(lv.Uint()) * uint32(rv.Uint())
		newVal.SetUint(uint64(prod))
	case reflect.Uint64:
		prod := uint64(lv.Uint()) * uint64(rv.Uint())
		newVal.SetUint(uint64(prod))
	case reflect.Uintptr:
		prod := uintptr(lv.Uint()) * uintptr(rv.Uint())
		newVal.SetUint(uint64(prod))
	case reflect.Float32:
		prod := float32(lv.Float()) * float32(rv.Float())
		newVal.SetFloat(float64(prod))
	case reflect.Float64:
		prod := float64(lv.Float()) * float64(rv.Float())
		newVal.SetFloat(float64(prod))
	case reflect.Complex64:
		prod := complex64(lv.Complex()) * complex64(rv.Complex())
		newVal.SetComplex(complex128(prod))
	case reflect.Complex128:
		prod := complex128(lv.Complex()) * complex128(rv.Complex())
		newVal.SetComplex(complex128(prod))
	default:
		panic("Type error: Invalid operands to multiplication: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorQuotient implements the binary operation '/'.
func operatorQuotient(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorQuotient: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		quot := int(lv.Int()) / int(rv.Int())
		newVal.SetInt(int64(quot))
	case reflect.Int8:
		quot := int8(lv.Int()) / int8(rv.Int())
		newVal.SetInt(int64(quot))
	case reflect.Int16:
		quot := int16(lv.Int()) / int16(rv.Int())
		newVal.SetInt(int64(quot))
	case reflect.Int32:
		quot := int32(lv.Int()) / int32(rv.Int())
		newVal.SetInt(int64(quot))
	case reflect.Int64:
		quot := int64(lv.Int()) / int64(rv.Int())
		newVal.SetInt(int64(quot))
	case reflect.Uint:
		quot := uint(lv.Uint()) / uint(rv.Uint())
		newVal.SetUint(uint64(quot))
	case reflect.Uint8:
		quot := uint8(lv.Uint()) / uint8(rv.Uint())
		newVal.SetUint(uint64(quot))
	case reflect.Uint16:
		quot := uint16(lv.Uint()) / uint16(rv.Uint())
		newVal.SetUint(uint64(quot))
	case reflect.Uint32:
		quot := uint32(lv.Uint()) / uint32(rv.Uint())
		newVal.SetUint(uint64(quot))
	case reflect.Uint64:
		quot := uint64(lv.Uint()) / uint64(rv.Uint())
		newVal.SetUint(uint64(quot))
	case reflect.Uintptr:
		quot := uintptr(lv.Uint()) / uintptr(rv.Uint())
		newVal.SetUint(uint64(quot))
	case reflect.Float32:
		quot := float32(lv.Float()) / float32(rv.Float())
		newVal.SetFloat(float64(quot))
	case reflect.Float64:
		quot := float64(lv.Float()) / float64(rv.Float())
		newVal.SetFloat(float64(quot))
	case reflect.Complex64:
		quot := complex64(lv.Complex()) / complex64(rv.Complex())
		newVal.SetComplex(complex128(quot))
	case reflect.Complex128:
		quot := complex128(lv.Complex()) / complex128(rv.Complex())
		newVal.SetComplex(complex128(quot))
	default:
		panic("Type error: Invalid operands to division: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorRemainder implements the binary operation '%'.
func operatorRemainder(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorRemainder: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		rem := int(lv.Int()) % int(rv.Int())
		newVal.SetInt(int64(rem))
	case reflect.Int8:
		rem := int8(lv.Int()) % int8(rv.Int())
		newVal.SetInt(int64(rem))
	case reflect.Int16:
		rem := int16(lv.Int()) % int16(rv.Int())
		newVal.SetInt(int64(rem))
	case reflect.Int32:
		rem := int32(lv.Int()) % int32(rv.Int())
		newVal.SetInt(int64(rem))
	case reflect.Int64:
		rem := int64(lv.Int()) % int64(rv.Int())
		newVal.SetInt(int64(rem))
	case reflect.Uint:
		rem := uint(lv.Uint()) % uint(rv.Uint())
		newVal.SetUint(uint64(rem))
	case reflect.Uint8:
		rem := uint8(lv.Uint()) % uint8(rv.Uint())
		newVal.SetUint(uint64(rem))
	case reflect.Uint16:
		rem := uint16(lv.Uint()) % uint16(rv.Uint())
		newVal.SetUint(uint64(rem))
	case reflect.Uint32:
		rem := uint32(lv.Uint()) % uint32(rv.Uint())
		newVal.SetUint(uint64(rem))
	case reflect.Uint64:
		rem := uint64(lv.Uint()) % uint64(rv.Uint())
		newVal.SetUint(uint64(rem))
	case reflect.Uintptr:
		rem := uintptr(lv.Uint()) % uintptr(rv.Uint())
		newVal.SetUint(uint64(rem))
	default:
		panic("Type error: Invalid operands to '%' operator: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorAnd implements the binary operation '&'.
func operatorAnd(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorAnd: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		and := int(lv.Int()) & int(rv.Int())
		newVal.SetInt(int64(and))
	case reflect.Int8:
		and := int8(lv.Int()) & int8(rv.Int())
		newVal.SetInt(int64(and))
	case reflect.Int16:
		and := int16(lv.Int()) & int16(rv.Int())
		newVal.SetInt(int64(and))
	case reflect.Int32:
		and := int32(lv.Int()) & int32(rv.Int())
		newVal.SetInt(int64(and))
	case reflect.Int64:
		and := int64(lv.Int()) & int64(rv.Int())
		newVal.SetInt(int64(and))
	case reflect.Uint:
		and := uint(lv.Uint()) & uint(rv.Uint())
		newVal.SetUint(uint64(and))
	case reflect.Uint8:
		and := uint8(lv.Uint()) & uint8(rv.Uint())
		newVal.SetUint(uint64(and))
	case reflect.Uint16:
		and := uint16(lv.Uint()) & uint16(rv.Uint())
		newVal.SetUint(uint64(and))
	case reflect.Uint32:
		and := uint32(lv.Uint()) & uint32(rv.Uint())
		newVal.SetUint(uint64(and))
	case reflect.Uint64:
		and := uint64(lv.Uint()) & uint64(rv.Uint())
		newVal.SetUint(uint64(and))
	case reflect.Uintptr:
		and := uintptr(lv.Uint()) & uintptr(rv.Uint())
		newVal.SetUint(uint64(and))
	default:
		panic("Type error: Invalid operands to '&' operator: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorOr implements the binary operation '|'.
func operatorOr(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorOr: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		or := int(lv.Int()) | int(rv.Int())
		newVal.SetInt(int64(or))
	case reflect.Int8:
		or := int8(lv.Int()) | int8(rv.Int())
		newVal.SetInt(int64(or))
	case reflect.Int16:
		or := int16(lv.Int()) | int16(rv.Int())
		newVal.SetInt(int64(or))
	case reflect.Int32:
		or := int32(lv.Int()) | int32(rv.Int())
		newVal.SetInt(int64(or))
	case reflect.Int64:
		or := int64(lv.Int()) | int64(rv.Int())
		newVal.SetInt(int64(or))
	case reflect.Uint:
		or := uint(lv.Uint()) | uint(rv.Uint())
		newVal.SetUint(uint64(or))
	case reflect.Uint8:
		or := uint8(lv.Uint()) | uint8(rv.Uint())
		newVal.SetUint(uint64(or))
	case reflect.Uint16:
		or := uint16(lv.Uint()) | uint16(rv.Uint())
		newVal.SetUint(uint64(or))
	case reflect.Uint32:
		or := uint32(lv.Uint()) | uint32(rv.Uint())
		newVal.SetUint(uint64(or))
	case reflect.Uint64:
		or := uint64(lv.Uint()) | uint64(rv.Uint())
		newVal.SetUint(uint64(or))
	case reflect.Uintptr:
		or := uintptr(lv.Uint()) | uintptr(rv.Uint())
		newVal.SetUint(uint64(or))
	default:
		panic("Type error: Invalid operands to '|' operator: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorXor implements the binary operation '^'.
func operatorXor(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorXor: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		xor := int(lv.Int()) ^ int(rv.Int())
		newVal.SetInt(int64(xor))
	case reflect.Int8:
		xor := int8(lv.Int()) ^ int8(rv.Int())
		newVal.SetInt(int64(xor))
	case reflect.Int16:
		xor := int16(lv.Int()) ^ int16(rv.Int())
		newVal.SetInt(int64(xor))
	case reflect.Int32:
		xor := int32(lv.Int()) ^ int32(rv.Int())
		newVal.SetInt(int64(xor))
	case reflect.Int64:
		xor := int64(lv.Int()) ^ int64(rv.Int())
		newVal.SetInt(int64(xor))
	case reflect.Uint:
		xor := uint(lv.Uint()) ^ uint(rv.Uint())
		newVal.SetUint(uint64(xor))
	case reflect.Uint8:
		xor := uint8(lv.Uint()) ^ uint8(rv.Uint())
		newVal.SetUint(uint64(xor))
	case reflect.Uint16:
		xor := uint16(lv.Uint()) ^ uint16(rv.Uint())
		newVal.SetUint(uint64(xor))
	case reflect.Uint32:
		xor := uint32(lv.Uint()) ^ uint32(rv.Uint())
		newVal.SetUint(uint64(xor))
	case reflect.Uint64:
		xor := uint64(lv.Uint()) ^ uint64(rv.Uint())
		newVal.SetUint(uint64(xor))
	case reflect.Uintptr:
		xor := uintptr(lv.Uint()) ^ uintptr(rv.Uint())
		newVal.SetUint(uint64(xor))
	default:
		panic("Type error: Invalid operands to '^' operator: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorAndNot implements the binary operation '&^'.
func operatorAndNot(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorAndNot: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		andNot := int(lv.Int()) &^ int(rv.Int())
		newVal.SetInt(int64(andNot))
	case reflect.Int8:
		andNot := int8(lv.Int()) &^ int8(rv.Int())
		newVal.SetInt(int64(andNot))
	case reflect.Int16:
		andNot := int16(lv.Int()) &^ int16(rv.Int())
		newVal.SetInt(int64(andNot))
	case reflect.Int32:
		andNot := int32(lv.Int()) &^ int32(rv.Int())
		newVal.SetInt(int64(andNot))
	case reflect.Int64:
		andNot := int64(lv.Int()) &^ int64(rv.Int())
		newVal.SetInt(int64(andNot))
	case reflect.Uint:
		andNot := uint(lv.Uint()) &^ uint(rv.Uint())
		newVal.SetUint(uint64(andNot))
	case reflect.Uint8:
		andNot := uint8(lv.Uint()) &^ uint8(rv.Uint())
		newVal.SetUint(uint64(andNot))
	case reflect.Uint16:
		andNot := uint16(lv.Uint()) &^ uint16(rv.Uint())
		newVal.SetUint(uint64(andNot))
	case reflect.Uint32:
		andNot := uint32(lv.Uint()) &^ uint32(rv.Uint())
		newVal.SetUint(uint64(andNot))
	case reflect.Uint64:
		andNot := uint64(lv.Uint()) &^ uint64(rv.Uint())
		newVal.SetUint(uint64(andNot))
	case reflect.Uintptr:
		andNot := uintptr(lv.Uint()) &^ uintptr(rv.Uint())
		newVal.SetUint(uint64(andNot))
	default:
		panic("Type error: Invalid operands to '&^' operator: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorShiftRight implements the binary operation '>>'.
func operatorShiftRight(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	lv := left.Value.(reflect.Value)

	rv, ok := right.Value.(reflect.Value)
	if !ok {
		// We have an exact.Value. Turn it into a uint64 and set rv from that
		ev := right.Value.(exact.Value)
		v64, _ := exact.Uint64Val(ev)
		rv = reflect.ValueOf(v64)
	}
	amt := rv.Uint()

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorShiftRight: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		shifted := int(lv.Int()) >> amt
		newVal.SetInt(int64(shifted))
	case reflect.Int8:
		shifted := int8(lv.Int()) >> amt
		newVal.SetInt(int64(shifted))
	case reflect.Int16:
		shifted := int16(lv.Int()) >> amt
		newVal.SetInt(int64(shifted))
	case reflect.Int32:
		shifted := int32(lv.Int()) >> amt
		newVal.SetInt(int64(shifted))
	case reflect.Int64:
		shifted := int64(lv.Int()) >> amt
		newVal.SetInt(int64(shifted))
	case reflect.Uint:
		shifted := uint(lv.Uint()) >> amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uint8:
		shifted := uint8(lv.Uint()) >> amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uint16:
		shifted := uint16(lv.Uint()) >> amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uint32:
		shifted := uint32(lv.Uint()) >> amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uint64:
		shifted := uint64(lv.Uint()) >> amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uintptr:
		shifted := uintptr(lv.Uint()) >> amt
		newVal.SetUint(uint64(shifted))
	default:
		panic("Type error: Invalid operands to shift: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorShiftLeft implements the binary operation '<<'.
func operatorShiftLeft(env *environ, left, right Object) Object {
	left = getTypedObject(left)
	lv := left.Value.(reflect.Value)

	rv, ok := right.Value.(reflect.Value)
	if !ok {
		// We have an exact.Value. Turn it into a uint64 and set rv from that
		ev := right.Value.(exact.Value)
		v64, _ := exact.Uint64Val(ev)
		rv = reflect.ValueOf(v64)
	}
	amt := rv.Uint()

	newTyp := left.Typ
	newRtyp, _ := getReflectType(env.interp.typeMap, newTyp)
	if newRtyp == nil {
		log.Fatal("operatorShiftLeft: Couldn't get reflect.Type from types.Type")
	}
	newVal := getSettableZeroVal(newRtyp)

	switch lv.Kind() {
	case reflect.Int:
		shifted := int(lv.Int()) << amt
		newVal.SetInt(int64(shifted))
	case reflect.Int8:
		shifted := int8(lv.Int()) << amt
		newVal.SetInt(int64(shifted))
	case reflect.Int16:
		shifted := int16(lv.Int()) << amt
		newVal.SetInt(int64(shifted))
	case reflect.Int32:
		shifted := int32(lv.Int()) << amt
		newVal.SetInt(int64(shifted))
	case reflect.Int64:
		shifted := int64(lv.Int()) << amt
		newVal.SetInt(int64(shifted))
	case reflect.Uint:
		shifted := uint(lv.Uint()) << amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uint8:
		shifted := uint8(lv.Uint()) << amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uint16:
		shifted := uint16(lv.Uint()) << amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uint32:
		shifted := uint32(lv.Uint()) << amt
		newVal.SetUint(uint64(shifted))
	case reflect.Uintptr:
		shifted := uintptr(lv.Uint()) << amt
		newVal.SetUint(uint64(shifted))
	default:
		panic("Type error: Invalid operands to shift: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}
	return Object{
		Value: newVal,
		Typ:   newTyp,
	}
}

// operatorLess implements the binary operation '<'.
// If this is being called at all, then the left and right objects
// have a value that's a "reflect.Value". Also, they can be compared,
// since the expression passed type checking.
func operatorLess(env *environ, left, right Object, typ types.Type) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	less := false
	switch lv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		less = lv.Int() < rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		less = lv.Uint() < rv.Uint()
	case reflect.Float32, reflect.Float64:
		less = lv.Float() < rv.Float()
	case reflect.String:
		less = lv.String() < rv.String()
	default:
		panic("Type error: Invalid operands to ordered comparison: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}

	var newVal reflect.Value
	if _, isNamed := typ.(*types.Named); isNamed {
		// Type is not "bool" but some other named boolean type.
		newRtyp, _ := getReflectType(env.interp.typeMap, typ)
		if newRtyp == nil {
			log.Fatal("operatorLess: Couldn't get reflect.Type from types.Type")
		}
		newVal = getSettableZeroVal(newRtyp)
		newVal.SetBool(less)
	} else {
		// Type is "bool" or "untyped bool". Use "bool".
		newVal = reflect.ValueOf(less)
	}

	return Object{
		Value: newVal,
		Typ:   typ,
	}
}

// operatorGreater implements the binary operation '>'.
// If this is being called at all, then the left and right objects
// have a value that's a "reflect.Value". Also, they can be compared,
// since the expression passed type checking.
func operatorGreater(env *environ, left, right Object, typ types.Type) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	greater := false
	switch lv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		greater = lv.Int() > rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		greater = lv.Uint() > rv.Uint()
	case reflect.Float32, reflect.Float64:
		greater = lv.Float() > rv.Float()
	case reflect.String:
		greater = lv.String() > rv.String()
	default:
		panic("Type error: Invalid operands to ordered comparison: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}

	var newVal reflect.Value
	if _, isNamed := typ.(*types.Named); isNamed {
		// Type is not "bool" but some other named boolean type.
		newRtyp, _ := getReflectType(env.interp.typeMap, typ)
		if newRtyp == nil {
			log.Fatal("operatorGreater: Couldn't get reflect.Type from types.Type")
		}
		newVal = reflect.New(newRtyp).Elem()
		newVal.SetBool(greater)
	} else {
		// Type is "bool" or "untyped bool". Use "bool".
		newVal = reflect.ValueOf(greater)
	}

	return Object{
		Value: newVal,
		Typ:   typ,
	}
}

// operatorLessEqual implements the binary operation '<='.
// If this is being called at all, then the left and right objects
// have a value that's a "reflect.Value". Also, they can be compared,
// since the expression passed type checking.
func operatorLessEqual(env *environ, left, right Object, typ types.Type) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	lessEqual := false
	switch lv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		lessEqual = lv.Int() <= rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		lessEqual = lv.Uint() <= rv.Uint()
	case reflect.Float32, reflect.Float64:
		lessEqual = lv.Float() <= rv.Float()
	case reflect.String:
		lessEqual = lv.String() <= rv.String()
	default:
		panic("Type error: Invalid operands to ordered comparison: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}

	var newVal reflect.Value
	if _, isNamed := typ.(*types.Named); isNamed {
		// Type is not "bool" but some other named boolean type.
		newRtyp, _ := getReflectType(env.interp.typeMap, typ)
		if newRtyp == nil {
			log.Fatal("operatorLessEqual: Couldn't get reflect.Type from types.Type")
		}
		newVal = reflect.New(newRtyp).Elem()
		newVal.SetBool(lessEqual)
	} else {
		// Type is "bool" or "untyped bool". Use "bool".
		newVal = reflect.ValueOf(lessEqual)
	}

	return Object{
		Value: newVal,
		Typ:   typ,
	}
}

// operatorGreaterEqual implements the binary operation '>='.
// If this is being called at all, then the left and right objects
// have a value that's a "reflect.Value". Also, they can be compared,
// since the expression passed type checking.
func operatorGreaterEqual(env *environ, left, right Object, typ types.Type) Object {
	left = getTypedObject(left)
	right = getTypedObject(right)
	lv := left.Value.(reflect.Value)
	rv := right.Value.(reflect.Value)

	greaterEqual := false
	switch lv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		greaterEqual = lv.Int() >= rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		greaterEqual = lv.Uint() >= rv.Uint()
	case reflect.Float32, reflect.Float64:
		greaterEqual = lv.Float() >= rv.Float()
	case reflect.String:
		greaterEqual = lv.String() >= rv.String()
	default:
		panic("Type error: Invalid operands to ordered comparison: " + TypeString(left.Typ) + ", " + TypeString(right.Typ))
	}

	var newVal reflect.Value
	if _, isNamed := typ.(*types.Named); isNamed {
		// Type is not "bool" but some other named boolean type.
		newRtyp, _ := getReflectType(env.interp.typeMap, typ)
		if newRtyp == nil {
			log.Fatal("operatorGreaterEqual: Couldn't get reflect.Type from types.Type")
		}
		newVal = reflect.New(newRtyp).Elem()
		newVal.SetBool(greaterEqual)
	} else {
		// Type is "bool" or "untyped bool". Use "bool".
		newVal = reflect.ValueOf(greaterEqual)
	}

	return Object{
		Value: newVal,
		Typ:   typ,
	}
}

// operatorGreaterEqual implements the binary operation '=='.
// If this is being called at all, then the left and right objects
// have a value that's a "reflect.Value". Also, they can be compared,
// since the expression passed type checking.
func operatorEqual(env *environ, left, right Object, typ types.Type) Object {
	isUntypedNil := func(t types.Type) bool {
		if isTyped(t) {
			return false
		}
		k := t.Underlying().(*types.Basic).Kind()
		if k != types.UntypedNil {
			return false
		}
		return true
	}

	var lv, rv reflect.Value

	leftIsUntypedNil, rightIsUntypedNil := true, true
	if !isUntypedNil(left.Typ) {
		leftIsUntypedNil = false
		left = getTypedObject(left)
		lv = left.Value.(reflect.Value)
	}
	if !isUntypedNil(right.Typ) {
		rightIsUntypedNil = false
		right = getTypedObject(right)
		rv = right.Value.(reflect.Value)
	}

	var equal bool
	switch {
	case leftIsUntypedNil:
		equal = rightIsUntypedNil || rv.IsNil()
	case rightIsUntypedNil:
		equal = lv.IsNil()
	case types.Identical(left.Typ, right.Typ):
		equal = lv.Interface() == rv.Interface()
	case types.AssignableTo(left.Typ, right.Typ):
		clv := lv.Convert(rv.Type())
		equal = clv.Interface() == rv.Interface()
	default:
		crv := rv.Convert(lv.Type())
		equal = lv.Interface() == crv.Interface()
	}

	var newVal reflect.Value
	if _, isNamed := typ.(*types.Named); isNamed {
		// Type is not "bool" but some other named boolean type.
		newRtyp, _ := getReflectType(env.interp.typeMap, typ)
		if newRtyp == nil {
			log.Fatal("operatorEqual: Couldn't get reflect.Type from types.Type")
		}
		newVal = reflect.New(newRtyp).Elem()
		newVal.SetBool(equal)
	} else {
		// Type is "bool" or "untyped bool". Use "bool".
		newVal = reflect.ValueOf(equal)
	}

	return Object{
		Value: newVal,
		Typ:   typ,
	}
}
