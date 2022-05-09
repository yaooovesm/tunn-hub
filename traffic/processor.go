package traffic

import log "github.com/cihub/seelog"

//
// CipherFlowProcessor
// @Description:
//
type CipherFlowProcessor interface {
	FlowProcessor

	//
	// Measure
	// @Description:
	// @param mtu
	// @return int
	//
	Measure(mtu int) int
}

//
// FlowProcessor
// @Description: flow data processor interface
//
type FlowProcessor interface {
	//
	// Init
	// @Description: init processor
	//
	Init() bool
	//
	// Process
	// @Description: do process method
	// @param raw
	// @return []byte
	//
	Process(raw []byte) []byte
}

//
// ProcessorNode
// @Description: processor chain node
//
type ProcessorNode struct {
	Name      string
	Processor FlowProcessor
	Next      *ProcessorNode
	Last      *ProcessorNode
}

//
// FlowProcessors
// @Description: processors
//
type FlowProcessors struct {
	Name string
	head *ProcessorNode
	tail *ProcessorNode
}

//
// NewFlowProcessor
// @Description:
// @return *FlowProcessors
//
func NewFlowProcessor() *FlowProcessors {
	return &FlowProcessors{}
}

//
// ProcessReverse
// @Description: start process data by using processor reversely
// @receiver fp
//
func (fp FlowProcessors) ProcessReverse(bytes []byte) []byte {
	current := fp.tail
	for current != nil {
		//do process
		bytes = current.Processor.Process(bytes)
		//move to last
		current = current.Last
	}
	return bytes
}

//
// Process
// @Description: start process data by using processor positively
// @receiver fp
//
func (fp FlowProcessors) Process(bytes []byte) []byte {
	current := fp.head
	for current != nil {
		//do process
		bytes = current.Processor.Process(bytes)
		//move to next
		current = current.Next
	}
	return bytes
}

//
// Register
// @Description: register for data processor
// @receiver fp
// @param processor
//
func (fp *FlowProcessors) Register(processor FlowProcessor, name string) {
	if processor == nil {
		return
	}
	if node := fp.GetByName(name); node != nil {
		node.Processor = processor
		return
	}
	if !processor.Init() {
		_ = log.Warn("[", fp.Name, "][fp:", name, "] register failed")
	}
	node := &ProcessorNode{
		Name:      name,
		Processor: processor,
		Next:      nil,
		Last:      nil,
	}
	if fp.head == nil {
		fp.head = node
		fp.tail = node
	} else {
		fp.tail.Next = node
		node.Last = fp.tail
		fp.tail = node
	}
	log.Info("[", fp.Name, "][fp:", name, "] register success")
}

//
// GetByName
// @Description:
// @receiver fp
// @param name
// @return *ProcessorNode
//
func (fp *FlowProcessors) GetByName(name string) *ProcessorNode {
	current := fp.head
	for current != nil {
		if current.Name == name {
			return current
		}
		//move to next
		current = current.Next
	}
	return nil
}

//
// Delete
// @Description:
// @receiver fp
// @param name
//
func (fp *FlowProcessors) Delete(name string) {
	current := fp.head
	for current != nil {
		if current.Name == name {
			last := current.Last
			next := current.Next
			last.Next = next
			next.Last = last
		}
		//move to next
		current = current.Next
	}
}

//
// List
// @Description:
// @receiver fp
// @return []string
//
func (fp FlowProcessors) List() []string {
	var list []string
	current := fp.head
	for current != nil {
		list = append(list, current.Name)
		//move to next
		current = current.Next
	}
	return list
}
