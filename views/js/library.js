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
// Request to the server and refresh all files.
function loadAllFiles() {
  $.ajax({
    headers: {
      token: window.localStorage.getItem('token'),
    },
    url: '/file',
    type: 'get',
    data: null,
    dataType: 'json',
    success: function (data) {
      $('#file-list').empty();
      for (let i = 0; i < data.files.length; i++) {
        data.files[i].Time = changeDateFormat(data.files[i].Time);
        let addtion =
          '<li class="list-group-item"><div class="media"><div class="media-body"><h6><i class="far fa-calendar-alt"></i> ' +
          data.files[i].Time +
          '</h6><span>' +
          '<a href="' +
          data.files[i].Url +
          '" download="' +
          data.files[i].Name +
          '">' +
          data.files[i].Name +
          '</a>';
        data.files[i].Url + '</span></div></div></li>';
        $(addtion).prependTo('#file-list');
      }
    },
    error: function () {
      $('#upload').removeClass('btn-success').addClass('btn-danger');
    },
  });
}
// Function to control the refresh of the progress bar.
function updateProgress(per) {
  // FIXME: To solve the problem of dom.
  $('#progress').attr('style', 'width: ' + per + '%');
  $('#progress').attr('aria-valuenow', per);
}
// The real multi-part upload function.
function uploadFile(fileHash, fileName, sliceBuffer) {
  let flag = true;
  let chunkList = [];
  let state = false;
  $.ajax({
    headers: {
      token: window.localStorage.getItem('token'),
    },
    url: '/check',
    type: 'get',
    async: false,
    data: 'hash=' + fileHash + '&name=' + fileName,
    dataType: 'json',
    success: function (data) {
      chunkList = data.list;
      state = data.state;
    },
    error: function () {
      $('#upload').removeClass('btn-success').addClass('btn-danger');
      flag = false;
    },
  });
  // Exit early if there is a network error.
  if (flag == false) {
    return true;
  }
  // If the upload is complete, return directly.
  if (state == true) {
    updateProgress(100);
    return state;
  }
  chunkList = chunkList.map((e) => parseInt(e));
  sliceBuffer.forEach((buffer, i) => {
    if (state == true) {
      return state;
    }
    if (!chunkList.includes(String(i))) {
      const blob = new File([buffer], `${i}`);
      let formData = new FormData();
      formData.append('file', blob);
      formData.append('hash', fileHash);
      $.ajax({
        headers: {
          token: window.localStorage.getItem('token'),
        },
        url: '/upload',
        type: 'post',
        data: formData,
        async: false,
        cache: false,
        processData: false,
        contentType: false,
        success: function (data) {
          state = data.state;
          chunkList = data.list;
        },
        error: function () {
          $('#upload').removeClass('btn-success').addClass('btn-danger');
          state = false;
        },
      });
    }
    // Alert the server to merge the file.
    if (chunkList.length == sliceBuffer.length && state == false) {
      $.ajax({
        headers: {
          token: window.localStorage.getItem('token'),
        },
        url: '/merge',
        type: 'get',
        data: 'hash=' + fileHash + '&name=' + fileName,
        async: false,
        dataType: 'json',
        success: function (data) {
          state = data.state;
        },
        error: function () {
          $('#upload').removeClass('btn-success').addClass('btn-danger');
          state = false;
        },
      });
    }
  });
  // If the upload is incomplete, upload again.
  return state;
}
$(document).ready(function () {
  loadAllFiles();
  // Simulate clicking the select file button.
  $('#choose').on('click', function () {
    $('#file').click();
  });
  // Upload the file you had choosen.
  $('#upload').on('click', function () {
    updateProgress(0);
    const file = $('#file')[0].files[0];
    if (file.size == 0) {
      $('#upload').removeClass('btn-success').addClass('btn-danger');
      return;
    }
    const chunkSize = 2 * 1024 * 1024;
    const chunkTotal = Math.ceil(file.size / chunkSize);
    const sliceBuffer = [];
    for (let i = 0; i < chunkTotal; i++) {
      const blobPart = file.slice(
        sliceBuffer.length * chunkSize,
        Math.min((sliceBuffer.length + 1) * chunkSize, file.size)
      );
      sliceBuffer.push(blobPart);
    }
    let hash = window.localStorage.getItem(file.name + ' ' + file.size);
    // If there is already a calculated Hash, then transfer it directly.
    if (hash != null) {
      while (true) {
        let state = uploadFile(hash, file.name, sliceBuffer);
        if (state == true) {
          loadAllFiles();
          updateProgress(100);
          $('#upload').removeClass('btn-danger').addClass('btn-success');
          break;
        }
      }
    } else {
      const fileReader = new FileReader();
      const spark = new SparkMD5.ArrayBuffer();
      let index = 0;
      fileReader.onload = function () {
        spark.append(fileReader.result);
        index += 1;
        if (index < sliceBuffer.length) {
          loadNext();
        } else {
          hash = spark.end();
          window.localStorage.setItem(file.name + ' ' + file.size, hash);
          $('#upload').removeClass('btn-danger').addClass('btn-success');
          // Begin to check these chunk's state in server.
          while (true) {
            let state = uploadFile(hash, file.name, sliceBuffer);
            if (state == true) {
              loadAllFiles();
              updateProgress(100);
              $('#upload').removeClass('btn-danger').addClass('btn-success');
              break;
            }
          }
        }
      };
      function loadNext() {
        fileReader.readAsArrayBuffer(sliceBuffer[index]);
      }
      loadNext();
    }
  });
});
