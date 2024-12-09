// Package knode (pronounced just 'node') is a Kronark Node (.knode) file parser.
package knode

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// The magic number of a .knode file.
	MAGIC_NUMBER = "kronarknode"
	// The latest version of the Kronark node format.
	LATEST_VERSION = 1
)

type (
	// Position represents a node's position within a project or another node.
	// Even though it is encoded as two int16 values, the values actually hold 10 bit numbers.
	Position struct{ X, Y int16 }

	// ValueType denotes the types of all the possible values
	// that any of the sockets of the nodes can output.
	ValueType string
)

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

	builtinNodeStart
)

// BuiltinNodeTypeNames holds the names of all the built-in node types.
// This array should be indexed by members of the BuiltinNodeType enum.
var BuiltinNodeTypeNames = []string{
	PortBuiltin:      "builtin:port",
	SettingsBuiltin:  "builtin:settings",
	PathBuiltin:      "builtin:path",
	BytesBuiltin:     "builtin:bytes",
	JoinBuiltin:      "builtin:join",
	OptionBuiltin:    "builtin:option",
	ConditionBuiltin: "builtin:condition",
	FormatBuiltin:    "builtin:format",
	TypeBuiltin:      "builtin:type",
	ApplyBuiltin:     "builtin:apply",
	SizeBuiltin:      "builtin:size",
	FileBuiltin:      "builtin:file",
	ReverseBuiltin:   "builtin:reverse",
	ValueBuiltin:     "builtin:value",
	MathBuiltin:      "builtin:math",
	RepeatBuiltin:    "builtin:repeat",
	TimeBuiltin:      "builtin:time",
	SplitBuiltin:     "builtin:split",
	CollectBuiltin:   "builtin:collect",
}

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

	builtinTypeStart
)

// BuiltinValueTypeNames holds the names of all the built-in value types.
// This array should be indexed by members of the BuiltinValueType enum.
var BuiltinValueTypeNames = []ValueType{
	NoneType:                  "builtin-type:none",
	AnyType:                   "builtin-type:any",
	RepetitionType:            "builtin-type:repetition",
	SettingsType:              "builtin-type:settings",
	OptionThenType:            "builtin-type:option then",
	OptionWhenType:            "builtin-type:option when",
	SelectionType:             "builtin-type:selection",
	BytesType:                 "builtin-type:bytes",
	TruthType:                 "builtin-type:truth",
	NumberType:                "builtin-type:number",
	TextType:                  "builtin-type:text",
	RepetitiveSelectionType:   "builtin-type:repetitive selection",
	RepetitiveBytesType:       "builtin-type:repetitive bytes",
	RepetitiveTruthType:       "builtin-type:repetitive truth",
	RepetitiveNumberType:      "builtin-type:repetitive number",
	RepetitiveTextType:        "builtin-type:repetitive text",
	RepetitivePortDefaultType: "builtin-type:repetitive port default",
	RepetitivePortValueType:   "builtin-type:repetitive port value",
	PortDefaultType:           "builtin-type:port default",
	PortChannelType:           "builtin-type:port channel",
	PortValueType:             "builtin-type:port value",
	PathModuleType:            "builtin-type:path module",
	PathAbsoluteType:          "builtin-type:path absolute",
	RootOutputType:            "builtin-type:root output",
	RootInputType:             "builtin-type:root input",
}

// SocketType denotes the type of a node socket.
type SocketType byte

const (
	OutgoingNamed SocketType = iota
	IncomingNamed
	IncomingNumber
	IncomingSelect
	IncomingSwitch
	IncomingText
)

// Socket represents a node's input/output socket.
type Socket struct {
	// The type of this socket.
	Type SocketType

	// The index into the type list of the type this socket outputs.
	ValueType byte

	// Where will this socket be placed outside of the node.
	PortSlot byte

	// If set, allows the socket to show multiple times in the node.
	// Useful for array inputs. Invalid (always false) for switch sockets.
	Repetitive bool

	// Whether this socket is connected. Invalid for outgoing sockets.
	Connected bool

	// What other socket is this socket connected to.
	ConnectedSocket byte

	// In what node instance is the socket this one is connected to placed.
	ConnectedNode byte

	// The value this socket contains.
	Value string
}

// Instance represents an instance of a node prototype inside or outside the project.
type Instance struct {
	// The key of the instance.
	Key byte

	// An index into the list of paths to the node protoype this instance is derived from.
	Type byte

	// The name of this instance.
	Name string

	// The position of this instance.
	Position Position

	// The list of sockets going into or out of this instance.
	Sockets []Socket
}

type (

	// Node represents the structure of a node file.
	Node struct {
		// The format version of this node.
		Version byte

		// The position of the node's input root. Invalid for root nodes.
		InputRootPosition Position

		// The node's output root. Invalid for root nodes.
		OutputRoot struct {
			Position    Position
			Connections [][2]byte
		}

		// The paths to all nodes this node refers to.
		Nodes []string

		// The list of non-reserved types this node uses.
		Types []ValueType

		// The list of node instances this node contains.
		Instances []Instance
	}
)

// Space, NodeSpace, RootNode and RootSpace are aliases for [knode.Node].
// Helpful for showcasing intent.
type (
	Space     = Node
	NodeSpace = Node
	RootNode  = Node
	RootSpace = Node
)

// ParseError gets returned when any single one of the file parsing functions
// is unable to successfully parse a [knode.Node] out of a file or buffer.
type ParseError struct {
	Context string
	Message string
	Index   int
}

// Error implements error for ParseError.
func (pe *ParseError) Error() string {
	return fmt.Sprintf(
		"error while %s: %s (at index %d [0x%x])",
		pe.Context, pe.Message, pe.Index, pe.Index,
	)
}

type parser struct {
	buf []byte
	i   int
	ctx string
	n   *Node
}

func (p *parser) error(msg string, a ...any) *ParseError {
	return &ParseError{
		Context: p.ctx,
		Message: fmt.Sprintf(msg, a...),
		Index:   p.i,
	}
}

func (p *parser) read(v any) *ParseError {
	amt, e := binary.Decode(p.buf[p.i:], binary.BigEndian, v)
	if e != nil {
		return p.error("internal read error (%v)", e)
	}
	p.i += amt
	return nil
}

func (p *parser) parseHeader() (e *ParseError) {
	p.ctx = "parsing magic"
	magic := make([]byte, 11)
	if e = p.read(&magic); e != nil {
		return e
	}

	if string(magic) != MAGIC_NUMBER {
		return p.error("invalid magic %v", magic)
	}

	p.ctx = "reading version number"
	if e = p.read(&p.n.Version); e != nil {
		return
	}

	if p.n.Version > LATEST_VERSION {
		return p.error(
			"invalid version number %d (higher than latest [%d])",
			p.n.Version, LATEST_VERSION,
		)
	}

	return
}

func (p *parser) parseRoots() (e *ParseError) {
	p.ctx = "unpacking root positions"

	pos := make([]byte, 5)
	if e = p.read(&pos); e != nil {
		return
	}

	// there's probably a better way to do this
	// in the meantime, please read the file format for this one
	p.n.InputRootPosition.X = (int16(pos[0])<<2 | int16((pos[1]>>(8-2))&0b11)) - 500
	p.n.InputRootPosition.Y = (int16(pos[1]&0b111111)<<4 | int16((pos[2]>>4)&0b1111)) - 500
	p.n.OutputRoot.Position.X = (int16(pos[2]&0b1111)<<6 | int16((pos[3]>>(8-6))&0b111111)) - 500
	p.n.OutputRoot.Position.Y = (int16(pos[3]&0b11)<<8 | int16(pos[4])) - 500

	p.ctx = "reading amount of connections"
	var connCount byte
	if e = p.read(&connCount); e != nil {
		return
	}
	for i := range connCount {
		p.ctx = fmt.Sprintf("reading outgoing connection [%d]", i)
		var nodeAndSock [2]byte
		if e = p.read(&nodeAndSock); e != nil {
			return
		}

		p.n.OutputRoot.Connections =
			append(p.n.OutputRoot.Connections, nodeAndSock)
	}

	return
}

func (p *parser) parseNodesAndTypes() (e *ParseError) {
	p.ctx = "reading amount of nodes"
	var nodeCount byte
	if e = p.read(&nodeCount); e != nil {
		return
	}
	for range nodeCount {
		var l byte
		if e = p.read(&l); e != nil {
			return
		}
		nodeType := make([]byte, l)
		if e = p.read(&nodeType); e != nil {
			return
		}
		p.n.Nodes = append(p.n.Nodes, string(nodeType))
	}

	p.ctx = "reading amount of types"
	var typeCount byte
	if e = p.read(&typeCount); e != nil {
		return
	}
	for range typeCount {
		var l byte
		if e = p.read(&l); e != nil {
			return
		}
		valueType := make([]byte, l)
		if e = p.read(&valueType); e != nil {
			return
		}
		p.n.Types = append(p.n.Types, ValueType(valueType))
	}

	return
}

func (p *parser) parseSocket(i *Instance) (e *ParseError) {
	var s Socket

	var socketFlags byte
	if e = p.read(&socketFlags); e != nil {
		return
	}

	s.Type = SocketType(socketFlags >> 3 & 0b111)
	s.Connected = socketFlags&0b00000010 != 0
	s.Repetitive = socketFlags&0b00000100 != 0
	switchValue := socketFlags&0b00000001 != 0

	if e = p.read(&s.ValueType); e != nil {
		return
	}

	if e = p.read(&s.PortSlot); e != nil {
		return
	}

	var valueLen uint32

	if s.Type != OutgoingNamed {
		if s.Connected {
			if e = p.read(&s.ConnectedNode); e != nil {
				return
			}
			if e = p.read(&s.ConnectedSocket); e != nil {
				return
			}
		} else if s.Type != IncomingSwitch {
			if e = p.read(&valueLen); e != nil {
				return
			}
			value := make([]byte, valueLen)
			if e = p.read(&value); e != nil {
				return
			}
			s.Value = string(value)
		} else {
			s.Value = fmt.Sprintf("%v", switchValue)
		}
	}

	i.Sockets = append(i.Sockets, s)

	return
}

func (p *parser) parseInstances() (e *ParseError) {
	p.ctx = "reading amount of instances"

	var instanceCount byte
	if e = p.read(&instanceCount); e != nil {
		return
	}

	for ii := range instanceCount {
		var i Instance
		p.ctx = fmt.Sprintf("parsing instance [%d]", ii)

		if e = p.read(&i.Key); e != nil {
			return
		}

		if e = p.read(&i.Type); e != nil {
			return
		}

		sizeInfo := make([]byte, 4)
		if e = p.read(&sizeInfo); e != nil {
			return
		}

		i.Position.X = (int16(sizeInfo[0])<<2 | int16((sizeInfo[1]>>(8-2))&0b11)) - 500
		i.Position.Y = (int16(sizeInfo[1]&0b111111)<<4 | int16((sizeInfo[2]>>(8-4))&0b1111)) - 500

		instNameLen := (sizeInfo[2]&0b1111)<<2 | (sizeInfo[3]>>(8-2))&0b11
		name := make([]byte, instNameLen)
		if e = p.read(&name); e != nil {
			return
		}
		i.Name = string(name)

		instSockCount := (sizeInfo[3] & 0b111111)
		for si := range instSockCount {
			p.ctx = fmt.Sprintf("parsing socket [%d] in instance [%d]", si, ii)
			if e := p.parseSocket(&i); e != nil {
				return e
			}
		}

		p.n.Instances = append(p.n.Instances, i)
	}
	return
}

func (p *parser) parse() (_ *Node, err *ParseError) {
	p.n = &Node{}
	if err = p.parseHeader(); err != nil {
		return
	} else if err = p.parseRoots(); err != nil {
		return
	} else if err = p.parseNodesAndTypes(); err != nil {
		return
	} else if err = p.parseInstances(); err != nil {
		return
	}
	return p.n, nil
}

// ParseFromSlice decodes a node from buf.
func ParseFromSlice(buf []byte) (n *Node, err *ParseError) {
	p := parser{buf: buf}
	return p.parse()
}

// ParseFromReader decodes a node from reader r.
// The error is not guaranteed to be a [knode.ParseError]
func ParseFromReader(r io.Reader) (n *Node, err error) {
	var buf []byte
	if buf, err = io.ReadAll(r); err != nil {
		return
	}
	return ParseFromSlice(buf)
}
