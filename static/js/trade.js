$(function() {
  $.getJSON('/trade', function(data) {
    $('#buyprice').val(data.buyprice);
    $('#sellprice').val(data.sellprice);
    $('#buyamount').val(data.buyamount);
    $('#sellamount').val(data.sellamount);
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
        buyamount: $('#buyamount').val(),
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
        sellamount: $('#sellamount').val(),
      }
    })

    return false;
  });
})