// Function to control the color of submit button.
function sumbitColor(method) {
  if (method == 'success') {
    $('#submit').removeClass('btn-danger').addClass('btn-success');
  } else if (method == 'danger') {
    $('#submit').removeClass('btn-success').addClass('btn-danger');
  }
}

// Request to the server and refresh all notes.
function loadAllNotes() {
  $.ajax({
    headers: {
      token: window.localStorage.getItem('token'),
    },
    url: '/note',
    type: 'get',
    data: null,
    dataType: 'json',
    success: function (data) {
      sumbitColor('success');
      $('#note-list').empty();
      for (let i = 0; i < data.notes.length; i++) {
        data.notes[i].Time = changeDateFormat(data.notes[i].Time);
        let addtion =
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
      sumbitColor('danger');
    },
  });
}

// Parse the date format in mongodb into a simple date format.
function changeDateFormat(cellval) {
  let dateVal = cellval + '';
  if (cellval != null) {
    // Use regular expressions.
    let reg = new RegExp('.\\d{3}\\+\\d{4}$');
    let date = new Date(dateVal.replace(reg, '').replace('T', ' '));
    let month =
      date.getMonth() + 1 < 10
        ? '0' + (date.getMonth() + 1)
        : date.getMonth() + 1;
    let day = date.getDate() < 10 ? '0' + date.getDate() : date.getDate();
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
      let indent = '    ';
      let start = this.selectionStart;
      let end = this.selectionEnd;
      let selected = window.getSelection().toString();
      selected = indent + selected.replace(/\n/g, '\n' + indent);
      this.value =
        this.value.substring(0, start) + selected + this.value.substring(end);
      this.setSelectionRange(start + indent.length, start + selected.length);
    }
  });

  // Function to add log items to the log list.
  $('#submit').click(function () {
    let formData = $('#form-write').serialize();
    $('textarea').val('').height(24);
    $.ajax({
      headers: {
        token: window.localStorage.getItem('token'),
      },
      url: '/write',
      type: 'post',
      data: formData,
      dataType: 'json',
      success: function () {
        sumbitColor('success');
        // Get and refresh the log list after each submission.
        loadAllNotes();
      },
      error: function () {
        sumbitColor('danger');
      },
    });
  });
});
