
var apis = [];

$(document).ready(function () {
  $('#api-detail').on('show.bs.modal', function (e) {
    var tr = $(e.relatedTarget).closest('tr'),
        idx = parseInt(tr.find('.api-index').text(), 10),
        api = apis[idx - 1],
        url = '//' + location.hostname + (location.port ? ':'+location.port : '') + api.example;
    $('#api-detail-title').text(api.name);
    $('#api-detail-method').text(api.method);
    $('#api-detail-description').text(api.description);
    $('#api-detail-parameters').text(JSON.stringify(api.parameters, true, ' ').replace(/"/g, ''));
    $('#api-detail-example').attr('href', url).text(location.protocol + url);
  });
});

var TableRow = React.createClass({
  render: function() {
    return (
        <tr>
          <td className="api-index">{this.props.index+1}</td>
          <td className="api-name">
            <a data-toggle="modal" data-target="#api-detail">{this.props.content.name}</a>
          </td>
          <td className="api-name">{this.props.content.method}</td>
          <td className="api-name">{this.props.content.description}</td>
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
    app.func.ajax({type: 'GET', url: '/api-list', success: function (data) {
      apis = data;
      self.setState({data: apis});
    }});
  },
  render: function() {
    var rows = this.state.data.map(function(record, index) {
      return (
          <TableRow key={record.name} index={index} content={record} />
      );
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>#</th>
              <th>API</th>
              <th>Method</th>
              <th>Description</th>
              <th></th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

React.render(<Table />, document.getElementById('data'));
