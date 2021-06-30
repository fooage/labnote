// Function to control the color of submit button the, if there is an error
// happened it will turn red.
function sumbitColor(method) {
  if (method == 'success') {
    $('#submit').removeClass('btn-danger').addClass('btn-success');
  } else if (method == 'danger') {
    $('#submit').removeClass('btn-success').addClass('btn-danger');
  }
}

// Request to the server and refresh all notes. I want to display the log list
// in pages to prevent too much slow down the loading speed.
function loadAllNotes() {
  $.ajax({
    headers: {
      token: window.localStorage.getItem('token'),
    },
    url: '/journal/list',
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

// Parse the date format which in go's time.Time into a simple date format as 2021-06-30.
function changeDateFormat(cellval) {
  let dateVal = cellval + '';
  if (cellval != null) {
    // use regular expressions.
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
  loadAllNotes();
  // Determine the size of the input box according to the number of input lines.
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
  // Support the tab and enter key's use in textarea.
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
  // Function to add log items to the whole log list.
  $('#submit').click(function () {
    let formData = $('#form-write').serialize();
    $('textarea').val('').height(24);
    $.ajax({
      headers: {
        token: window.localStorage.getItem('token'),
      },
      url: '/journal/write',
      type: 'post',
      data: formData,
      dataType: 'json',
      success: function () {
        sumbitColor('success');
        // get the log list after each submission
        loadAllNotes();
      },
      error: function () {
        sumbitColor('danger');
      },
    });
  });
});
