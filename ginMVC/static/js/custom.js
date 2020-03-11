function ajaxJsonRequest(method, url, jsonData, handleData) {
	$.ajax({
		type : method,
		url : url,
		data : jsonData,
		dataType : "json",
		contentType : "application/json; charset=utf-8",
		success : handleData,
		error : function(XMLHttpRequest, textStatus, errorThrown) {
			alert(errorThrown);
		}
	});
	return false;
}

function ajaxGetRequest(url, id, handleData) {
	$.ajax({
		type : "GET",
		url : url + "?id=" + id,
		dataType : "json",
		contentType : "application/json; charset=utf-8",
		success : handleData,
		error : function(XMLHttpRequest, textStatus, errorThrown) {
			alert(errorThrown);
		}
	});
	return false;
}

function ajaxFormRequest(method, url, formData, handleData) {
	$.ajax({
		type : method,
		url : url,
		data : formData,
		contentType : false,
		processData : false,
		cache : false,
		success : handleData,
		error : function(XMLHttpRequest, textStatus, errorThrown) {
			alert(errorThrown);
		}
	});
	return false;
}