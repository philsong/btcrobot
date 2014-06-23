function _onOKCoinAPIkeyUI(OKCoinAPIkey) {
  var divokcoinsingle = document.getElementById("divokcoinsingle");
  var divokcoincluster = document.getElementById("divokcoincluster");
  console.log(OKCoinAPIkey);


  if (OKCoinAPIkey == "single") {
    divokcoinsingle.style.display = "block";
    divokcoincluster.style.display = "none";
  } else {
    divokcoinsingle.style.display = "none";
    divokcoincluster.style.display = "block";
  }
}

function onOKCoinAPIkeyUI() {
  var OKCoinAPIkeyID = document.getElementById("OKCoinAPIkey");
  var OKCoinAPIkey = OKCoinAPIkeyID.value
  _onOKCoinAPIkeyUI(OKCoinAPIkey)
}


$(function() {
  $.getJSON('/secret', function(data) {

    console.log(data)

    $('#username').val(data.username);
    $('#password').val(data.password);

    $('#bitvc_email').val(data.bitvc_email);
    $('#bitvc_password').val(data.bitvc_password);

    $('#huobi_access_key').val(data.huobi_access_key);
    $('#huobi_secret_key').val(data.huobi_secret_key);

    $('#smtp_username').val(data.smtp_username);
    $('#smtp_password').val(data.smtp_password);
    $('#smtp_host').val(data.smtp_host);
    $('#smtp_addr').val(data.smtp_addr);

    var OKCoinAPIkey = data.OKCoinAPIkey

    $('#OKCoinAPIkey').val(OKCoinAPIkey);

    _onOKCoinAPIkeyUI(OKCoinAPIkey)

    $('#ok_partner').val(data.ok_partner);
    $('#ok_secret_key').val(data.ok_secret_key);

    $('#ok_partner1').val(data.ok_partner1);
    $('#ok_secret_key1').val(data.ok_secret_key1);
    $('#ok_partner2').val(data.ok_partner2);
    $('#ok_secret_key2').val(data.ok_secret_key2);
    $('#ok_partner3').val(data.ok_partner3);
    $('#ok_secret_key3').val(data.ok_secret_key3);
    $('#ok_partner4').val(data.ok_partner4);
    $('#ok_secret_key4').val(data.ok_secret_key4);
    $('#ok_partner5').val(data.ok_partner5);
    $('#ok_secret_key5').val(data.ok_secret_key5);
    $('#ok_partner6').val(data.ok_partner6);
    $('#ok_secret_key6').val(data.ok_secret_key6);
    $('#ok_partner7').val(data.ok_partner7);
    $('#ok_secret_key7').val(data.ok_secret_key7);
    $('#ok_partner8').val(data.ok_partner8);
    $('#ok_secret_key8').val(data.ok_secret_key8);
    $('#ok_partner9').val(data.ok_partner9);
    $('#ok_secret_key9').val(data.ok_secret_key9);
    $('#ok_partner10').val(data.ok_partner10);
    $('#ok_secret_key10').val(data.ok_secret_key10);

    $("select[name='OKCoinAPIkey']").selectpicker({
      style: 'btn-primary',
      menuStyle: 'dropdown-inverse'
    });
  });

  // 表单提交
  $('#update_conf').submit(function() {
    var self = $(this);
    $.post(self.attr('action'), self.serialize(), function(data) {
      console.log(data)
      alert(data);
      location.reload();
    });
    return false;
  });
})