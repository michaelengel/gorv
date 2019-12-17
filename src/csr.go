package main

// CSR definitions

var csr = map[uint32]uint32 {
	// user-level

	// user trap setup
	0x000: 0,	// URW ustatus User status register
	0x004: 0,	// URW uie     User interrupt-enable register
	0x005: 0,	// URW utvec   User trap handler base address

	// user trap handling
	0x040: 0,	// URW uscratchScratch register for user trap handlers.
	0x041: 0,	// URW uepcUser exception program counter.
	0x042: 0,	// URW ucauseUser trap cause.
	0x043: 0,	// URW ubadaddrUser bad address.
	0x044: 0,	// URW uipUser interrupt pending.

	// User Floating-Point CSRs
	0x001: 0,	// URW fflagsFloating-Point Accrued Exceptions.
	0x002: 0,	// URW frmFloating-Point Dynamic Rounding Mode.
	0x003: 0,	// URW fcsrFloating-Point Control and Status Register (frm+fflags).

	// User Counter/Timers
	0xC00: 0,	// URO cycleCycle counter for RDCYCLE instruction.
	0xC01: 0,	// URO timeTimer for RDTIME instruction.
	0xC02: 0,	// URO instretInstructions-retired counter for RDINSTRET instruction.
	0xC03: 0,	// URO hpmcounter3Performance-monitoring counter.
	0xC04: 0,	// URO hpmcounter4Performance-monitoring counter....
	0xC1F: 0,	// URO hpmcounter31Performance-monitoring counter.
	0xC80: 0,	// URO cyclehUpper 32 bits ofcycle, RV32I only.
	0xC81: 0,	// URO timehUpper 32 bits oftime, RV32I only.
	0xC82: 0,	// URO instrethUpper 32 bits ofinstret, RV32I only.
	0xC83: 0,	// URO hpmcounter3hUpper 32 bits ofhpmcounter3, RV32I only.
	0xC84: 0,	// URO hpmcounter4hUpper 32 bits ofhpmcounter4, RV32I only....
	0xC9F: 0,	// URO hpmcounter31hUpper 32 bits ofhpmcounter31, RV32I only.

	// supervisor-level

	// Supervisor Trap Setup
	0x100: 0,	// SRW sstatusSupervisor status register.
	0x102: 0,	// SRW sedelegSupervisor exception delegation register.
	0x103: 0,	// SRW sidelegSupervisor interrupt delegation register.
	0x104: 0,	// SRW sieSupervisor interrupt-enable register.
	0x105: 0,	// SRW stvecSupervisor trap handler base address.

	// Supervisor Trap Handling
	0x140: 0,	// SRW sscratchScratch register for supervisor trap handlers.
	0x141: 0,	// SRW sepcSupervisor exception program counter.
	0x142: 0,	// SRW scauseSupervisor trap cause.
	0x143: 0,	// SRW sbadaddrSupervisor bad address.
	0x144: 0,	// SRW sipSupervisor interrupt pending.

	// Supervisor Protection and Translation
	0x180: 0,	// SRW sptbr   Page-table base register

	// hypervisor-level

	// Hypervisor Trap Setup
	0x200: 0,	// HRW hstatusHypervisor status register.
	0x202: 0,	// HRW hedelegHypervisor exception delegation register.
	0x203: 0,	// HRW hidelegHypervisor interrupt delegation register.
	0x204: 0,	// HRW hieHypervisor interrupt-enable register.
	0x205: 0,	// HRW htvecHypervisor trap handler base address.Hypervisor Trap Handling
	0x240: 0,	// HRW hscratchScratch register for hypervisor trap handlers.
	0x241: 0,	// HRW hepcHypervisor exception program counter.
	0x242: 0,	// HRW hcauseHypervisor trap cause.
	0x243: 0,	// HRW hbadaddrHypervisor bad address.
	0x244: 0,	// HRW hipHypervisor interrupt pending.Hypervisor Protection and Translation

	// machine-level

	// Machine Information Registers
	0xF11: 0,	// MRO mvendorid Vendor ID.
	0xF12: 0,	// MRO marchid Architecture ID.
	0xF13: 0,	// MRO mimpid Implementation ID.
	0xF14: 0,	// MRO mhartid Hardware thread ID.Machine Trap Setup
	0x300: 0,	// MRW mstatus Machine status register.
	0x301: 0,	// MRW misa ISA and extensions
	0x302: 0,	// MRW medeleg Machine exception delegation register.
	0x303: 0,	// MRW mideleg Machine interrupt delegation register.
	0x304: 0,	// MRW mie Machine interrupt-enable register.
	0x305: 0,	// MRW mtvec Machine trap-handler base address.Machine Trap Handling
	0x340: 0,	// MRW mscratch Scratch register for machine trap handlers.
	0x341: 0,	// MRW mepc Machine exception program counter.
	0x342: 0,	// MRW mcause Machine trap cause.
	0x343: 0,	// MRW mbadaddr Machine bad address.
	0x344: 0,	// MRW mip Machine interrupt pending.Machine Protection and Translation
	0x380: 0,	// MRW mbase Base register.
	0x381: 0,	// MRW mbound Bound register.
	0x382: 0,	// MRW mibase Instruction base register.
	0x383: 0,	// MRW mibound Instruction bound register.
	0x384: 0,	// MRW mdbase Data base register.
	0x385: 0,	// MRW mdbound Data bound register.

	// Machine Counter/Timers
	0xB00: 0,	// MRW mcycleMachine cycle counter.
	0xB02: 0,	// MRW minstretMachine instructions-retired counter.
	0xB03: 0,	// MRW mhpmcounter3Machine performance-monitoring counter.
	0xB04: 0,	// MRW mhpmcounter4Machine performance-monitoring counter....
	0xB1F: 0,	// MRW mhpmcounter31Machine performance-monitoring counter.
	0xB80: 0,	// MRW mcyclehUpper 32 bits ofmcycle, RV32I only.
	0xB82: 0,	// MRW minstrethUpper 32 bits ofminstret, RV32I only.
	0xB83: 0,	// MRW mhpmcounter3hUpper 32 bits ofmhpmcounter3, RV32I only.
	0xB84: 0,	// MRW mhpmcounter4hUpper 32 bits ofmhpmcounter4, RV32I only....
	0xB9F: 0,	// MRW mhpmcounter31hUpper 32 bits ofmhpmcounter31, RV32I only.

	// Machine Counter Setup
	0x320: 0,	// MRW mucounterenUser-mode counter enable.
	0x321: 0,	// MRW mscounterenSupervisor-mode counter enable.
	0x322: 0,	// MRW mhcounterenHypervisor-mode counter enable.
	0x323: 0,	// MRW mhpmevent3Machine performance-monitoring event selector.
	0x324: 0,	// MRW mhpmevent4Machine performance-monitoring event selector....
	0x33F: 0,	// MRW mhpmevent31Machine performance-monitoring event selector.

	// Debug/Trace Registers (shared with Debug Mode)
	0x7A0: 0,	// MRW tselectDebug/Trace trigger register select.
	0x7A1: 0,	// MRW tdata1First Debug/Trace trigger data register.
	0x7A2: 0,	// MRW tdata2Second Debug/Trace trigger data register.
	0x7A3: 0,	// MRW tdata3Third Debug/Trace trigger data register.

	// Debug Mode Registers
	0x7B0: 0,	// DRW dcsrDebug control and status register.
	0x7B1: 0,	// DRW dpcDebug PC.
	0x7B2: 0,	// DRW dscratchDebug scratch register.
}

var csrnames = map[uint32]string {
	// user-level

	// user trap setup
	0x000: "ustatus",
	0x004: "uie",
	0x005: "utvec",

	// user trap handling
	0x040: "uscratch",
	0x041: "uepc",
	0x042: "ucause",
	0x043: "ubadaddr",
	0x044: "uip",

	// 
	0x001: "fflags",
	0x002: "frm",
	0x003: "fcsr",

	// 
	0xC00: "cycle",
	0xC01: "time",
	0xC02: "instret",
	0xC03: "hpmcounter3",
	0xC04: "hpmcounter4",
	0xC1F: "hpmcounter31",
	0xC80: "cycleh",
	0xC81: "timeh",
	0xC82: "instreth",
	0xC83: "hpmcounter3h",
	0xC84: "hpmcounter4h",
	0xC9F: "hpmcounter31h",

	// supervisor-level

	// 
	0x100: "sstatus",
	0x102: "sedeleg",
	0x103: "sideleg",
	0x104: "sie",
	0x105: "stvec",

	// 
	0x140: "sscratch",
	0x141: "sepc",
	0x142: "scause",
	0x143: "sbadaddr",
	0x144: "sip",

	// 
	0x180: "sptbr",

	// hypervisor-level

	// 
	0x200: "hstatus",
	0x202: "hedeleg",
	0x203: "hideleg",
	0x204: "hie",
	0x205: "htvec",
	0x240: "hscratch",
	0x241: "hepc",
	0x242: "hcause",
	0x243: "hbadaddr",
	0x244: "hip",

	// machine-level

	// 
	0xF11: "mvendorid",
	0xF12: "marchid",
	0xF13: "mimpid",
	0xF14: "mhartid",
	0x300: "mstatus",
	0x301: "misa ISA and extensions",
	0x302: "medeleg",
	0x303: "mideleg",
	0x304: "mie",
	0x305: "mtvec",
	0x340: "mscratch",
	0x341: "mepc",
	0x342: "mcause",
	0x343: "mbadaddr",
	0x344: "mip",
	0x380: "mbase",
	0x381: "mbound",
	0x382: "mibase",
	0x383: "mibound",
	0x384: "mdbase",
	0x385: "mdbound",

	// 
	0xB00: "mcycle",
	0xB02: "minstret",
	0xB03: "mhpmcounter3",
	0xB04: "mhpmcounter4",
	0xB1F: "mhpmcounter31",
	0xB80: "mcycleh",
	0xB82: "minstreth",
	0xB83: "mhpmcounter3h",
	0xB84: "mhpmcounter4h",
	0xB9F: "mhpmcounter31h",

	// 
	0x320: "mucounteren",
	0x321: "mscounteren",
	0x322: "mhcounteren",
	0x323: "mhpmevent3",
	0x324: "mhpmevent4",
	0x33F: "mhpmevent31",

	// 
	0x7A0: "tselect",
	0x7A1: "tdata1",
	0x7A2: "tdata2",
	0x7A3: "tdata3",

	// 
	0x7B0: "dcsr",
	0x7B1: "dpc",
	0x7B2: "dscratch",
}
