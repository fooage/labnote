$(document).ready(function () {
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
          $('#login').removeClass('btn-danger').addClass('btn-success');
          window.localStorage.setItem('token', data.token);
          window.location.href = '/home';
        } else {
          // Change the color of the login button when the login fails.
          $('#login').removeClass('btn-success').addClass('btn-danger');
        }
      },
      error: function () {
        $('#login').removeClass('btn-success').addClass('btn-danger');
      },
    });
  });
});
