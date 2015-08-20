
var table, sessions = [], filters = {type: -1, track: -1, level: -1, text: app.func.query('q')};

$(document).ready(function () {
  if (filters.text != '') $('#search-text').val(filters.text);
  _resize();

  $('.dropdown-menu a').click(function() {
    var a = $(this), group = a.closest('.btn-group').removeClass('open');
    filters[group.attr('data-filter-key')] = parseInt(a.attr('href').substring(1), 10);
    group.find('.caption').text(a.text()).blur();
    table.setProps();
    return false;
  });

  $('#session-detail').on('show.bs.modal', function (e) {
    var tr = $(e.relatedTarget).closest('tr');
    $.map(['abbreviation', 'title', 'abstract', 'length', 'type', 'track', 'level'], function (key) {
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

function _titleSize() {
  return $('.container').width()-466;
}

function _resize() {
  var height = $(window).height();
  $('.table-inner').css({height: (height-220)+'px'});
  $('.sess_title > div').css({width: _titleSize()+'px'});
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
    app.func.ajax('GET', '/reinvent/sessions', '', function (data) {
      sessions = data.sessions;
      self.setState({data: self.filter()});
    });
  },
  componentWillReceiveProps: function() {
    this.setState({data: this.filter()});
  },
  filter: function() {
    var data = [];
    $.map(sessions, function (session) {
      var match = ((filters.type == -1) || (filters.type == session.typeId)) &&
                  ((filters.track == -1) || (filters.track == session.trackId)) &&
                  ((filters.level == -1) || (filters.level == session.levelId));
      if (filters.text != '') {
        var title = session.title.toUpperCase(),
            abstract = session.abstract.toUpperCase();
        $.map(filters.text.split(' '), function (word) {
          match &= (session.abbreviation.indexOf(word) > -1) ||
            (title.indexOf(word) > -1) || (abstract.indexOf(word) > -1);
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
