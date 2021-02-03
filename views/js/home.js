// Request to the server and refresh all notes.
function loadAllNotes() {
  $.ajax({
    url: '/data',
    type: 'get',
    data: null,
    dataType: 'json',
    success: function (data) {
      $('#submit').removeClass('btn-danger').addClass('btn-success');
      $('#note-list').empty();
      for (var i = 0; i < data.notes.length; i++) {
        data.notes[i].Time = changeDateFormat(data.notes[i].Time);
        var addtion =
          '<li class="list-group-item"><div class="media"><div class="media-body"><h6><i class="far fa-calendar-alt"></i> ' +
          data.notes[i].Time +
          '</h6><span>' +
          data.notes[i].Content +
          '</span></div></div></li>';
        // FIXME: Optimize the refresh effect of the list.
        $(addtion).prependTo('#note-list');
      }
    },
    error: function () {
      $('#submit').removeClass('btn-success').addClass('btn-danger');
    },
  });
}
// Parse the date format in mongodb into a simple date format.
function changeDateFormat(cellval) {
  var dateVal = cellval + '';
  if (cellval != null) {
    // Use regular expressions.
    var reg = new RegExp('.\\d{3}\\+\\d{4}$');
    var date = new Date(dateVal.replace(reg, '').replace('T', ' '));
    var month =
      date.getMonth() + 1 < 10
        ? '0' + (date.getMonth() + 1)
        : date.getMonth() + 1;
    var day = date.getDate() < 10 ? '0' + date.getDate() : date.getDate();
    return date.getFullYear() + '-' + month + '-' + day;
  }
}
$(document).ready(function () {
  // Load all notes when the interface is loaded.
  loadAllNotes();
  // Textarea adaptive height is supported.
  $('textarea')
    .each(function () {
      this.setAttribute(
        'style',
        'height:' + this.scrollHeight + 'px;overflow-y:hidden;'
      );
    })
    .on('input', function () {
      this.style.height = 'auto';
      this.style.height = this.scrollHeight + 'px';
    });
  // Support the tab key's use in textarea.
  $('textarea').on('keydown', function (e) {
    if (e.keyCode == 9) {
      e.preventDefault();
      var indent = '    ';
      var start = this.selectionStart;
      var end = this.selectionEnd;
      var selected = window.getSelection().toString();
      selected = indent + selected.replace(/\n/g, '\n' + indent);
      this.value =
        this.value.substring(0, start) + selected + this.value.substring(end);
      this.setSelectionRange(start + indent.length, start + selected.length);
    }
  });
  // Function to add log items to the log list.
  $('#submit').click(function () {
    var formParam = $('#form-write').serialize();
    $('textarea').val('').height(24);
    $.ajax({
      url: '/data',
      type: 'post',
      data: formParam,
      dataType: 'json',
      success: function () {
        $('#submit').removeClass('btn-danger').addClass('btn-success');
        // Get and refresh the log list after each submission.
        loadAllNotes();
      },
      error: function () {
        $('#submit').removeClass('btn-success').addClass('btn-danger');
      },
    });
  });
});
