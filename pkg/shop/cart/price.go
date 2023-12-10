package cart

import (
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

var VisaPriceStandard = map[string]int{
	db.VisaType1MonthSingle:   55,
	db.VisaType1MonthMultiple: 136,
	db.VisaType3MonthSingle:   116,
	db.VisaType3MonthMultiple: 136,
}

var VisaPriceVOAStandard = map[string]int{
	db.VisaType1MonthSingle: 125,
}

var VisaPricePriority = map[string]int{
	db.VisaType1MonthSingle:   30,
	db.VisaType1MonthMultiple: 86,
	db.VisaType3MonthSingle:   91,
	db.VisaType3MonthMultiple: 86,
}

var ProcessingTimePrice1MonthSingle = map[string]int{
	db.ProcessingTime2Days:  45,
	db.ProcessingTime1Days:  60,
	db.ProcessingTimeUrgent: 90,
}

var ProcessingTimePrice1MonthMultiple = map[string]int{
	db.ProcessingTime2Days:  19,
	db.ProcessingTime1Days:  59,
	db.ProcessingTimeUrgent: 79,
}

var ProcessingTimePrice3MonthSingle = map[string]int{
	db.ProcessingTime2Days:  20,
	db.ProcessingTime1Days:  29,
	db.ProcessingTimeUrgent: 69,
}

var ProcessingTimePrice3MonthMultiple = map[string]int{
	db.ProcessingTime2Days:  29,
	db.ProcessingTime1Days:  59,
	db.ProcessingTimeUrgent: 89,
}

var ProcessingTimePrice = map[string]map[string]int{
	db.VisaType1MonthSingle:   ProcessingTimePrice1MonthSingle,
	db.VisaType1MonthMultiple: ProcessingTimePrice1MonthMultiple,
	db.VisaType3MonthSingle:   ProcessingTimePrice3MonthSingle,
	db.VisaType3MonthMultiple: ProcessingTimePrice3MonthMultiple,
}

var FastTrackPrice = map[string]int{
	db.FastTrackNormal: 65,
	db.FastTrackVIP:    95,
}

var CarPrice = map[string]int{
	db.CarYes: 35,
}
