var app = {};

(function($) {
	app.func = {
		ajax : _ajax,
		trim : _trim,
		query : _query
	};

	function _ajax(arg) {
		var dt = arg.dataType ? arg.dataType : 'json';
		$.ajax({
			url: arg.url, type: arg.type,
			data: arg.data, dataType: dt,
			success: function (data) {
				arg.success && arg.success(data);
			},
			error: function(xhr, status, err) {
				if (arg.error) {
					arg.error(xhr, status, err);
					return;
				}
				console.error(arg.url, status, err.toString());
			}
		});
	}
	function _trim(value) {
		return value.replace(/\s/g,'').replace(/　/g,'')
	}
	function _query(key, def) {
		key = key.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
		def = def ? def : "";
		var regex = new RegExp("[\\?&]" + key + "=([^&#]*)"),
				results = regex.exec(location.search);
		return results === null ? def : decodeURIComponent(results[1].replace(/\+/g, " "));
	}
})(jQuery);
