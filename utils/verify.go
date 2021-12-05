package utils

var (
	IdVerify       = Rules{"ID": {NotEmpty()}}
	PageInfoVerify = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	RegisterVerify = Rules{
		"Username": {NotEmpty()},
		"Email":    {RegexpMatch("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$")},
		"Passwd":   {RegexpMatch("^[\\w_-]{6,16}$")},
		//"Code":     {NotEmpty()}
	}
	CreateCommentVerify = Rules{
		"IdeaId":  {NotEmpty()},
		"UserId":  {NotEmpty()},
		"Content": {NotEmpty()},
	}
)
