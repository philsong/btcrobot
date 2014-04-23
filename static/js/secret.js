$(function() {
  $.getJSON('/secret', function(data) {
    $('#username').val(data.username);
    $('#password').val(data.password);
    $('#huobi_access_key').val(data.huobi_access_key);
    $('#huobi_secret_key').val(data.huobi_secret_key);
    $('#ok_partner').val(data.ok_partner);
    $('#ok_secret_key').val(data.ok_secret_key);
    $('#smtp_username').val(data.smtp_username);
    $('#smtp_password').val(data.smtp_password);
    $('#smtp_host').val(data.smtp_host);
    $('#smtp_addr').val(data.smtp_addr);
  });

  // 表单提交
  $('#update_conf').submit(function() {
    var self = $(this);
    $.post(self.attr('action'), self.serialize(), function(data) {
      alert(data.msg);
      location.reload();
    });
    return false;
  });
})