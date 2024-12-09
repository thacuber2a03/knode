# Header

- `kronarknode` magic number [11]
- version number [1]

# Roots

- root positions [5]:
    - input root position x [10 bits] (offset by 500)
    - input root position y [10 bits] (offset by 500)
    - output root position x [10 bits] (offset by 500)
    - output root position y [10 bits] (offset by 500)
- output root connection count [1]
- output root connections:
    - connection node [1]
    - connection socket [1]

# Nodes

- table size [1]
- table items:
    - node name length [1]
    - node name [string]

# Types

- table size [1]
- table items:
    - type name length [1]
    - type name [string]

# Instances

- instance count [1]
- instances:
    - instance key [1]
    - instance type [1]
    - instance position and lengths [4]:
        - instance position x [10 bits] (offset by 500)
        - instance position y [10 bits] (offset by 500)
        - instance name length [6 bits]
        - instance socket count [6 bits]
    - instance name [string]
    - instance sockets:
        - socket flags [1]:
            - PADDING [2 bit]
            - type and direction [3 bits]:
                - 000 = outgoing named
                - 001 = incoming named
                - 010 = incoming number
                - 011 = incoming select
                - 100 = incoming switch
                - 101 = incoming text
            - repetitive [1 bit]
            - connected [1 bit]
            - switch value [1 bit]
        - socket type index [1]
        - socket port slot [1]
        - if incoming:
            - if connected:
                - connection node [1]
                - connection socket [1]
            - if not connected and not switch:
                - value length [4]
                - value [string]

# Builtin node types

```go
// BuiltinNodeType is an alias for the higher values of the type index field in an Instance
// that denote the nodes whose prototype is pre-defined within the compiler.
type BuiltinNodeType byte

const (
	PortBuiltin BuiltinNodeType = 0xff - iota
	SettingsBuiltin
	PathBuiltin
	BytesBuiltin
	JoinBuiltin
	OptionBuiltin
	ConditionBuiltin
	FormatBuiltin
	TypeBuiltin
	ApplyBuiltin
	SizeBuiltin
	FileBuiltin
	ReverseBuiltin
	ValueBuiltin
	MathBuiltin
	RepeatBuiltin
	TimeBuiltin
	SplitBuiltin
	CollectBuiltin
)
```

# Builtin value types

```go
// BuiltinValueType denotes the value types that built-in nodes output.
type BuiltinValueType byte

const (
	NoneType BuiltinValueType = 0xff - iota
	AnyType
	RepetitionType
	SettingsType
	OptionThenType
	OptionWhenType
	SelectionType
	BytesType
	TruthType
	NumberType
	TextType
	RepetitiveSelectionType
	RepetitiveBytesType
	RepetitiveTruthType
	RepetitiveNumberType
	RepetitiveTextType
	RepetitivePortDefaultType
	RepetitivePortValueType
	PortDefaultType
	PortChannelType
	PortValueType
	PathModuleType
	PathAbsoluteType
	RootOutputType
	RootInputType
)
```
