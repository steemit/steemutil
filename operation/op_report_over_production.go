package operation

// FC_REFLECT( steemit::chain::report_over_production_operation,
//             (reporter)
//             (first_block)
//             (second_block) )

type ReportOverProductionOperation struct {
	Reporter string `json:"reporter"`
}

func (op *ReportOverProductionOperation) Type() OpType {
	return TypeReportOverProduction
}

func (op *ReportOverProductionOperation) Data() interface{} {
	return op
}
