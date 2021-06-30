// Function to control the color of upload button. It running as same as the
// login button, the red is a symbol of error.
function sumbitColor(method) {
  if (method == 'success') {
    $('#upload').removeClass('btn-danger').addClass('btn-success');
  } else if (method == 'danger') {
    $('#upload').removeClass('btn-success').addClass('btn-danger');
  }
}

// Parse the date format which in go's time.Time into a simple date format as 2021-06-30.
function changeDateFormat(cellval) {
  let dateVal = cellval + '';
  if (cellval != null) {
    // use regular expressions
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

// Request to the server and refresh the file-list for all of files. There are
// a problem that if number of files is too large, it will load for a long time.
function loadAllFiles() {
  $.ajax({
    headers: {
      token: window.localStorage.getItem('token'),
    },
    url: '/library/list',
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
        // show the lastest file on the top
        $(addtion).prependTo('#file-list');
      }
    },
    error: function () {
      sumbitColor('danger');
    },
  });
}

// Function to control the refresh of the progress bar.
function updateProgress(now, all) {
  let per = ((now * 1.0) / all) * 100;
  per = Math.round(per);
  $('#progress').attr('style', 'width: ' + per + '%');
  $('#progress').attr('aria-valuenow', per);
}

// The real multi-part upload function, first request the file slices that the
// server has received. Then add the untransmitted file slices to the asynchronous
// transmission queue. The last step is check the stat of the file in server.
async function postChunks(fileHash, fileName, sliceBuffer) {
  let chunkList = [];
  let state = false;
  $.ajax({
    headers: {
      token: window.localStorage.getItem('token'),
    },
    url: '/library/check',
    type: 'get',
    async: false,
    data: 'hash=' + fileHash + '&name=' + fileName,
    dataType: 'json',
    success: function (data) {
      chunkList = data.list;
      state = data.state;
    },
    error: function () {
      sumbitColor('danger');
    },
  });
  // if the upload is complete return directly
  if (state == true) {
    return true;
  }
  postRequset = [];
  chunkList = chunkList.map((e) => parseInt(e));
  sliceBuffer.forEach((buffer, i) => {
    if (!chunkList.includes(i)) {
      const blob = new File([buffer], `${i}`);
      let formData = new FormData();
      formData.append('file', blob);
      formData.append('hash', fileHash);
      postRequset.push(
        new Promise((resolve, reject) => {
          $.ajax({
            headers: {
              token: window.localStorage.getItem('token'),
            },
            url: '/library/upload',
            type: 'post',
            data: formData,
            async: true,
            cache: false,
            processData: false,
            contentType: false,
            success: function (data) {
              resolve({
                nums: data.nums,
                state: data.state,
              });
            },
            error: function () {
              reject({
                data: null,
              });
            },
          }).then(function (data) {
            requestAnimationFrame(
              updateProgress(data.nums, sliceBuffer.length)
            );
          });
        })
      );
    }
  });
  // Use Promise to transfer file slices asynchronously and get the number of
  // file slices to determine whether to merge files.
  let received = 0;
  await Promise.all(postRequset)
    .then(function (result) {
      for (let i = 0; i < result.length; i++) {
        state = result[i].state;
        received = Math.max(received, result[i].nums);
      }
      sumbitColor('success');
    })
    .catch(function () {
      sumbitColor('danger');
    });
  // alert the server to merge the file
  if (received == sliceBuffer.length && state == false) {
    $.ajax({
      headers: {
        token: window.localStorage.getItem('token'),
      },
      url: '/library/merge',
      type: 'get',
      data: 'hash=' + fileHash + '&name=' + fileName,
      async: false,
      dataType: 'json',
      success: function (data) {
        state = data.state;
      },
      error: function () {
        sumbitColor('danger');
        state = false;
      },
    });
  }
  return state;
}

// Calculate the file hash value and upload the file, calculate MD5 code of file.
function uploadFile() {
  const file = $('#file')[0].files[0];
  if (file.size == 0) {
    sumbitColor('danger');
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
  // If there is already a calculated Hash, then transfer it directly
  if (hash != null) {
    (async () => {
      for (let now = 0; now < 3; now++) {
        let state = await postChunks(hash, file.name, sliceBuffer);
        if (state == true) {
          loadAllFiles();
          requestAnimationFrame(updateProgress(100, 100));
          sumbitColor('success');
          return;
        }
      }
    })();
  } else {
    const fileReader = new FileReader();
    const spark = new SparkMD5.ArrayBuffer();
    let index = 0;
    function loadSlice() {
      fileReader.readAsArrayBuffer(sliceBuffer[index]);
    }
    loadSlice();
    fileReader.onload = async () => {
      spark.append(fileReader.result);
      index += 1;
      if (index < sliceBuffer.length) {
        loadSlice();
      } else {
        hash = spark.end();
        window.localStorage.setItem(file.name + ' ' + file.size, hash);
        sumbitColor('success');
        // Begin to check these chunk's state in server, if the transmission
        // fails, try two more times.
        for (let now = 0; now < 3; now++) {
          let state = await postChunks(hash, file.name, sliceBuffer);
          if (state == true) {
            loadAllFiles();
            requestAnimationFrame(updateProgress(100, 100));
            sumbitColor('success');
            return;
          }
        }
      }
    };
  }
}

$(document).ready(function () {
  loadAllFiles();
  // Simulate clicking the select file button.
  $('#choose').on('click', function () {
    $('#file').click();
  });
  // Upload the file you had choosen and reset the progress bar.
  $('#upload').on('click', function () {
    updateProgress(0, 100);
    uploadFile();
  });
});
