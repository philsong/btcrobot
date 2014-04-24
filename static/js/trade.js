$(function() {
  $.getJSON('/trade', function(data) {
    $('#buyprice').val(data.buyprice);
    $('#sellprice').val(data.sellprice);
    $('#buytotalamount').val(data.buytotalamount);
    $('#selltotalamount').val(data.selltotalamount);
    $('#buyinterval').val(data.buyinterval);
    $('#sellinterval').val(data.sellinterval);
    $('#buytimes').val(data.buytimes);
    $('#selltimes').val(data.selltimes);
    $('#maxbuyamountratio').val(data.maxbuyamountratio);
    $('#maxsellamountratio').val(data.maxsellamountratio);
  });

  // 表单提交
  $('input[name="dobuy"]').click(function(e) {
    $.ajax({
      type: "POST",
      url: "/trade",
      success: function(data) {
        alert(data);
        location.reload();
      },
      data: {
        msgtype: "dobuy",
        buyprice: $('#buyprice').val(),
        buyamount: $('#buytotalamount').val(),
        buyinterval: $('#buyinterval').val(),
        buytimes: $('#buytimes').val(),
        maxbuyamountratio: $('#maxbuyamountratio').val(),
      }
    })

    return false;
  });

  $('input[name="dosell"]').click(function(e) {
    $.ajax({
      type: "POST",
      url: "/trade",
      success: function(data) {
        alert(data);
        location.reload();
      },
      data: {
        msgtype: "dosell",
        sellprice: $('#sellprice').val(),
        sellamount: $('#selltotalamount').val(),
        sellinterval: $('#sellinterval').val(),
        selltimes: $('#selltimes').val(),
        maxsellamountratio: $('#maxsellamountratio').val(),
      }
    })

    return false;
  });
})