// Function to control the color of login button. If there are something wrong
// with the login, it will turn red to alert the user.
function loginColor(method) {
  if (method == 'success') {
    $('#submit').removeClass('btn-danger').addClass('btn-success');
  } else if (method == 'danger') {
    $('#submit').removeClass('btn-success').addClass('btn-danger');
  }
}

$(document).ready(function () {
  // Post the login data to the server and control the redirect to the main pages.
  $('#login').on('click', function () {
    let formData = $('#form-login').serialize();
    $('input').val('');
    $.ajax({
      url: '/login/submit',
      type: 'post',
      data: formData,
      dataType: 'json',
      success: function (data) {
        if (data.pass == true) {
          // login successfully and jump to the journal page
          loginColor('success');
          window.localStorage.setItem('token', data.token);
          window.location.href = '/journal';
        } else {
          // change the color of the login button when the login fails
          loginColor('danger');
        }
      },
      error: function () {
        loginColor('danger');
      },
    });
  });
});
