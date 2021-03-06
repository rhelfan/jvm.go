package misc

import (
	"github.com/zxh0/jvm.go/jvmgo/jvm/rtda"
	rtc "github.com/zxh0/jvm.go/jvmgo/jvm/rtda/class"
	"github.com/zxh0/jvm.go/jvmgo/util"
)

func init() {
	_unsafe(allocateInstance, "allocateInstance", "(Ljava/lang/Class;)Ljava/lang/Object;")
	_unsafe(defineClass, "defineClass", "(Ljava/lang/String;[BIILjava/lang/ClassLoader;Ljava/security/ProtectionDomain;)Ljava/lang/Class;")
	_unsafe(ensureClassInitialized, "ensureClassInitialized", "(Ljava/lang/Class;)V")
	_unsafe(staticFieldOffset, "staticFieldOffset", "(Ljava/lang/reflect/Field;)J")
	_unsafe(staticFieldBase, "staticFieldBase", "(Ljava/lang/reflect/Field;)Ljava/lang/Object;")
}

// public native Object allocateInstance(Class<?> type) throws InstantiationException;
// (Ljava/lang/Class;)Ljava/lang/Object;
func allocateInstance(frame *rtda.Frame) {
	vars := frame.LocalVars()
	classObj := vars.GetRef(1)

	class := classObj.Extra().(*rtc.Class)
	obj := class.NewObj()

	stack := frame.OperandStack()
	stack.PushRef(obj)
}

// public native Class defineClass(String name, byte[] b, int off, int len,
//  		ClassLoader loader, ProtectionDomain protectionDomain)
// (Ljava/lang/String;[BIILjava/lang/ClassLoader;Ljava/security/ProtectionDomain;)Ljava/lang/Class;
func defineClass(frame *rtda.Frame) {
	vars := frame.LocalVars()
	nameObj := vars.GetRef(1)
	byteArr := vars.GetRef(2)
	off := vars.GetInt(3)
	_len := vars.GetInt(4)
	//loaderObj := vars.GetRef(5)
	//pd := vars.GetRef(6)

	name := rtda.GoString(nameObj)
	name = rtc.DotToSlash(name)
	int8s := byteArr.Fields().([]int8)
	data := util.CastInt8sToUint8s(int8s)
	data = data[off : off+_len]

	// todo
	class := frame.ClassLoader().DefineClass(name, data)
	stack := frame.OperandStack()
	stack.PushRef(class.JClass())
}

// public native void ensureClassInitialized(Class<?> c);
// (Ljava/lang/Class;)V
func ensureClassInitialized(frame *rtda.Frame) {
	vars := frame.LocalVars()
	// this := vars.GetRef(0)
	classObj := vars.GetRef(1)

	goClass := classObj.Extra().(*rtc.Class)
	if goClass.InitializationNotStarted() {
		// undo ensureClassInitialized()
		frame.RevertNextPC()
		// init
		frame.Thread().InitClass(goClass)
	}
}

// public native long staticFieldOffset(Field f);
// (Ljava/lang/reflect/Field;)J
func staticFieldOffset(frame *rtda.Frame) {
	vars := frame.LocalVars()
	// vars.GetRef(0) // this
	fieldObj := vars.GetRef(1)

	offset := fieldObj.GetFieldValue("slot", "I").(int32)
	stack := frame.OperandStack()
	stack.PushLong(int64(offset))
}

// public native Object staticFieldBase(Field f);
// (Ljava/lang/reflect/Field;)Ljava/lang/Object;
func staticFieldBase(frame *rtda.Frame) {
	vars := frame.LocalVars()
	// vars.GetRef(0) // this
	fieldObj := vars.GetRef(1)

	goField := _getGoField(fieldObj)
	goClass := goField.Class()
	obj := goClass.AsObj()

	stack := frame.OperandStack()
	stack.PushRef(obj)
}

func _getGoField(fieldObj *rtc.Obj) *rtc.Field {
	extra := fieldObj.Extra()
	if extra != nil {
		return extra.(*rtc.Field)
	}

	root := fieldObj.GetFieldValue("root", "Ljava/lang/reflect/Field;").(*rtc.Obj)
	return root.Extra().(*rtc.Field)
}
