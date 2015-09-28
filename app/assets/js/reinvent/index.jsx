
var table, client = false, sessions = [], query = app.func.query('q'),
    filters = {type: -1, track: -1, level: -1, date: -1, text: ''};
if (query != '') {
  filters.text = query.replace(/\s/g,' ').replace(/　/g,' ');
  filters.text = filters.text.replace(/^\s+|\s+$/gm,'').toUpperCase();
}

$(document).ready(function () {
  if (query != '') $('#search-text').val(query);
  _resize();

  $('.dropdown-menu a').click(_setDropdownEvent);

  $('#session-detail').on('show.bs.modal', function (e) {
    var tr = $(e.relatedTarget).closest('tr');
    $.map(['abbreviation', 'title', 'abstract', 'length', 'type', 'track', 'level', 'start', 'room'], function (key) {
      $('#session-detail-'+key).text(tr.find('.sess_'+key).text());
    });
  });
});

$(window).keyup(function (e) {
  if ($('#search-text').is(':focus')) {
    var candidate = $('#search-text').val().replace(/\s/g,' ').replace(/　/g,' ');
    candidate = candidate.replace(/^\s+|\s+$/gm,'').toUpperCase();
    if (filters.text == candidate) return;
    filters.text = candidate;
    table.setProps();
  }
});
$(window).resize(_resize);

var windowWidth = 0;

function _setDropdownEvent(e) {
  var a = $(this), group = a.closest('.btn-group').removeClass('open');
  filters[group.attr('data-filter-key')] = parseInt(a.attr('href').substring(1), 10);
  group.find('.caption').text(a.text()).blur();
  table.setProps();
  return false;
}

function _titleSize() {
  if (windowWidth == 0) {
    windowWidth = $(window).width();
  }
  if (windowWidth <= 750) {
    return windowWidth-122;
  }
  return $('.container').width()-466;
}

function _resize() {
  var height = $(window).height();
  $('.table-inner').css({height: (height-220)+'px'});
  $('.sess_title > div').css({width: _titleSize()+'px'});
}

function _dateKey(value) {
  var date = new Date(value);
  var year = date.getYear();
  return ((year < 2000 ? (year+1900) : year)+''+
      _fill(date.getMonth()+1)+_fill(date.getDate()));
}
function _fill(value) {
  return (value.toString().length == 1 ? '0'+value : value);
}

var weekday = [];
weekday[0] = "Sunday";
weekday[1] = "Monday";
weekday[2] = "Tuesday";
weekday[3] = "Wednesday";
weekday[4] = "Thursday";
weekday[5] = "Friday";
weekday[6] = "Saturday";

var month = [];
month[0] = "Jan";
month[1] = "Feb";
month[2] = "Mar";
month[3] = "Apr";
month[4] = "May";
month[5] = "June";
month[6] = "July";
month[7] = "Aug";
month[8] = "Sept";
month[9] = "Oct";
month[10] = "Nov";
month[11] = "Dec";

function _dateTxt(d) {
  return weekday[d.getDay()]+', '+month[d.getMonth()]+' '+d.getDate();
}
function _datetime(value) {
  var date = new Date(value);
  return month[date.getMonth()]+' '+date.getDate()+' '+
         _fill(date.getHours())+':'+_fill(date.getMinutes());
}

function _setDayOptions() {
  var days = {};
  $.map(sessions, function (session) {
    var d = new Date(session.date*1000);
    days[''+_dateKey(session.date*1000)] = _dateTxt(d);
  });
  var keys = [];
  $.map(days, function (_, key) {
    keys.push(key);
  });
  keys.sort();
  var html = '<li><a href="#-1">All</a></li>';
  $.map(keys, function (key) {
    html += '<li><a href="#'+key+'">'+days[key]+'</a></li>';
  });
  $('#day-filter').html(html).find('a').click(_setDropdownEvent);
}

var TableRow = React.createClass({
  render: function() {
    var session = this.props.content,
        level = session.level.replace(/ \(.*\)/g, '');
    return (
        <tr data-rec-idx={this.props.index+1}>
          <td className="sess_id"><div>{session.id}</div></td>
          <td className="sess_abbreviation"><div>{session.abbreviation}</div></td>
          <td className="sess_title">
            <div style={{ width: _titleSize()+'px'}}>
              <a data-toggle="modal" data-target="#session-detail">{session.title}</a>
            </div>
          </td>
          <td className="sess_abstract">{session.abstract}</td>
          <td className="sess_length">{session.length}</td>
          <td className="sess_type"><div>{session.type}</div></td>
          <td className="sess_track"><div>{session.track}</div></td>
          <td className="sess_level_mod"><div>{level}</div></td>
          <td className="sess_level">{session.level}</td>
          <td className="sess_start">{_datetime(session.start)}</td>
          <td className="sess_room">at {session.room}</td>
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  componentDidMount: function() {
    var self = this;
    app.func.ajax({type: 'GET', url: 'sessions', success: function (data) {
      sessions = data.sessions;
      self.setState({data: self.filter()});
      _setDayOptions();
    }});
  },
  componentWillReceiveProps: function() {
    this.setState({data: this.filter()});
  },
  filter: function() {
    var data = [];
    $.map(sessions, function (session) {
      var date = (filters.date != -1) ? _dateKey(session.date*1000) : '',
          match = ((filters.date == -1) || (filters.date == date)) &&
                  ((filters.type == -1) || (filters.type == session.typeId)) &&
                  ((filters.track == -1) || (filters.track == session.trackId)) &&
                  ((filters.level == -1) || (filters.level == session.levelId));
      if (filters.text != '') {
        var title = session.title.toUpperCase(),
            abstract = session.abstract.toUpperCase(),
            room = session.room.toUpperCase(),
            type = session.type.toUpperCase(),
            track = session.track.toUpperCase(),
            level = session.level.toUpperCase();
        $.map(filters.text.split(' '), function (word) {
          match &= (session.id.indexOf(word) > -1) ||
            (session.abbreviation.indexOf(word) > -1) ||
            (title.indexOf(word) > -1) || (abstract.indexOf(word) > -1) ||
            (room.indexOf(word) > -1) || (type.indexOf(word) > -1) ||
            (track.indexOf(word) > -1) || (level.indexOf(word) > -1);
        });
      }
      if (match) data.push(session);
    });
    $('#count').text(data.length+' session'+(data.length > 1 ? 's' : ''));
    return data;
  },
  render: function() {
    var rows = this.state.data.map(function(record, index) {
      return (
          <TableRow key={record.id} index={index} content={record} />
      );
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th className="sess_id"><div>ID</div></th>
              <th className="sess_abbreviation"><div></div></th>
              <th className="sess_title"><div style={{ width: _titleSize()+'px'}}>Title</div></th>
              <th className="sess_type"><div>Type</div></th>
              <th className="sess_track"><div>Track</div></th>
              <th className="sess_level_mod"><div>Level</div></th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

table = React.render(<Table />, document.getElementById('data'));
