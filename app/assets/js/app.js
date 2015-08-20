var app = {};

(function($) {
	app.func = {
		ajax : _ajax,
		link : _link,
		stop : _stop,
		trim : _trim,
		query : _query
	};
	app.storage = {
		set : _set,
		get : _get
	};

	var baseurl = '';

	$(document).ready(function() {
		baseurl = $('#baseurl').val();

		$('.js_placeholder').each(function () {
			var that = $(this), font = that.css('fontSize'),
				placeholder = that.attr('data-placeholder').replace(/\|/g,'\n');
			that.focus(function () {
				that.css({color: 'black', fontSize: font});
				if (that.val() == placeholder) {
					that.val('');
				}
			}).blur(function () {
				if (that.val() == ''){
					that.css({color: 'silver', fontSize:'12px'}).val(placeholder);
				}
			}).trigger('blur');
		});
	});
	$(window).keyup(function (e) {
		if (e.keyCode == 13) {
			var active = $(document.activeElement);
			if (active.hasClass('imgless-btn')) {
				active.trigger('click');
			}
		}
		return true;
	});
	function _ajax(type, url, data, onSuccess, dataType) {
		var dt = dataType ? dataType : 'json';
		$.ajax({
			url: url, type: type,
			data: data, dataType: dt,
			success: function (data) {
				onSuccess && onSuccess(data);
			},
			error: function(xhr, status, err) {
				console.error(url, status, err.toString());
			}
		});
	}
	function _link(href, e) {
		e = e || window.event;
		if (!e)
			return false;
		if (href.indexOf(baseurl) == -1)
			href = baseurl + href;
		if (e && (e.ctrlKey || e.metaKey)) {
			window.open(href, '_blank');
		} else {
			location.href = href;
		}
	}
	function _stop(e) {
		e = e || window.event;
		if (!e)
			return false;
		e.cancelBubble = true;
		if (e.stopPropagation)
			e.stopPropagation();
		e.returnValue = false;
		if (e.preventDefault)
			e.preventDefault();
		return e;
	}
	var ls = false;
	try {ls = window.localStorage;} catch (e) {}
	if (! ls) {
		ls = window.addBehavior ? (function() {
			var storage = {}, prefix = 'data-userdata', attrs = document.body, mark = function(
					key, isRemove, temp, reg) {
				attrs.load(prefix);
				var temp = attrs.getAttribute(prefix) || '', reg = RegExp('\\b'
						+ key + '\\b,?', 'i'), hasKey = reg.test(temp) ? 1 : 0;
				temp = isRemove ? temp.replace(reg, '') : hasKey ? temp
						: temp === '' ? key : temp.split(',').concat(key).join(',');
				attrs.setAttribute(prefix, temp);
				attrs.save(prefix);
			};
			// add IE behavior support
			attrs.addBehavior('#default#userData');

			storage.getItem = function(key) {
				attrs.load(key);
				return attrs.getAttribute(key);
			};
			storage.setItem = function(key, value) {
				attrs.setAttribute(key, value);
				attrs.save(key);
				mark(key);
			};
			storage.removeItem = function(key) {
				attrs.removeAttribute(key);
				attrs.save(key);
				mark(key, 1);
			};
			return storage;
		})()
		: (function() {
			var storage = {}, cache = {};
			storage.getItem = function(key) {
				return cache[key];
			};
			storage.setItem = function(key, value) {
				cache[key] = value;
			};
			storage.removeItem = function(key) {
				cache[key] = null;
			};
			return storage;
		})();
	}
	function _set(key, value) {
		try {ls.setItem(key, JSON.stringify(value));} catch (e) {}
	}
	function _get(key, def) {
		var value = ls.getItem(key), candidate = (value != null) ? JSON.parse(value) : undefined;
		return (candidate == undefined) ? def : candidate;
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
