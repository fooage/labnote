// Function to control the color of login button.
function loginColor(method) {
  if (method == 'success') {
    $('#submit').removeClass('btn-danger').addClass('btn-success');
  } else if (method == 'danger') {
    $('#submit').removeClass('btn-success').addClass('btn-danger');
  }
}

$(document).ready(function () {
  // Post the login data to the server.
  $('#login').on('click', function () {
    let formData = $('#form-login').serialize();
    $('input').val('');
    $.ajax({
      url: '/login',
      type: 'post',
      data: formData,
      dataType: 'json',
      success: function (data) {
        if (data.pass == true) {
          // Login successfully and jump to the homepage.
          loginColor('success');
          window.localStorage.setItem('token', data.token);
          window.location.href = '/home';
        } else {
          // Change the color of the login button when the login fails.
          loginColor('danger');
        }
      },
      error: function () {
        loginColor('danger');
      },
    });
  });
});
