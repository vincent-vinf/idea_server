package utils

var (
	RegisterVerify = Rules{
		"Username": {NotEmpty()},
		"Email":    {RegexpMatch("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$")},
		"Passwd":   {RegexpMatch("^[\\w_-]{6,16}$")},
		//"Code":     {NotEmpty()}
	}
	PageInfoVerify = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
)
