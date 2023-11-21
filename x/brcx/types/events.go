package types

const (
	EventTypeBRCX           = ModuleName
	EventTypeManageContract = "manage_contract"
	EventTypeEntryPoint     = "entry_point"
	EventTypeCallEvm        = "call_evm"

	AttributeResult   = "result"
	AttributeProtocol = "protocol"

	AttributeManageContractOperation = "operation"
	AttributeManageContractAddress   = "contract_addrss"
	AttributeEvmOutput               = "evm_output"
	AttributeManageLog               = "log"
)
