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
      url: "http://127.0.0.1:9091/tradedobuy.json",
      success: function(data) {
        alert(data.msg);
        location.reload();
      },
      data: {
        buyprice: $('#buyprice').val(),
        buyamount: $('#buyamount').val(),
      }
    })

    return false;
  });

  $('input[name="dosell"]').click(function(e) {
    $.ajax({
      type: "POST",
      url: "http://127.0.0.1:9091/tradedosell.json",
      success: function(data) {
        alert(data.msg);
        location.reload();
      },
      data: {
        sellprice: $('#sellprice').val(),
        sellamount: $('#sellamount').val(),
      }
    })

    return false;
  });
})