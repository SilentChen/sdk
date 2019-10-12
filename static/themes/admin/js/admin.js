//requirement |
//			  |- jquery
//			  |- layer


const G_SUCCESS 			=	"success"
const G_FAIL				=	"fail"
const G_ERROR				=	"error"
const G_MSG_CONFIRM			=	"Are You Sure To Do This?"
const G_MSG_CONFIRM_YES     =	"Yes"
const G_MSG_CONFIRM_NO      =	"No"
const G_MSG_CONFIRM_DEL     =	"Delete?"
const G_MSG_ALERT_NOOPTSEL	=	"Plz Select An Option"
const G_MSG_ALERT_NOROWSEL 	=   "Plz Select An Row"


const G_APIRET_OK	=	0
const G_APIRET_DESC = {
	"-5"	:		"error, no data row affected.",
	"-4"	:		"error, query fail.",
	"-3"	:		"error, involid operation.",
	"-2"	:		"error, check params.",
	"-1"	:		"error, program go wrong.",
	"0"		:		"success.",
	"1"		:		"fail, try again.",
}

function ReqAndShowAuto (url, params, datatype) {
	$.post(url, params, function (ret) {
		if (G_APIRET_OK == ret.state) {
			OnNetretReloadAuto()
		}else{
			layer.alert(G_APIRET_DESC[ret.state])
		}
	}, datatype)
}

function OnNetretReloadAuto (second) {
	if (!second) {
		second = 3000
	}
	layer.msg(G_SUCCESS)
	setTimeout(function(){
		location.reload(location.href);
	}, second)
}

function RefreshWithParams() {
	location.reload(location.href);
}

function DoubleCheck(msg) {
	var question = G_MSG_CONFIRM
	if(msg) {
		question += "\r\n" + msg
	}

	return confirm(question)
}

function Checkall(name, obj) {
	$(":checkbox[name='"+name+"']").each(function(o) {
		$(this).prop('checked', obj.checked);
	});
}

function GetChkVal(chkname) {
	var ret = []
	$(":checkbox[name='"+chkname+"']").each(function(k, v) {
		if ($(this).is(":checked")) {
			ret.push($(this).val())
		}
	})

	return ret
}

function del_confirm() {
	return confirm('一旦删除将不可恢复，确定吗？');
}