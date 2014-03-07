/**
 * 更新票数
 * 
 * @param val 1为正反，0为反方
 */
function vote2(tid, val){
    dn = "#dn-" + tid
    up = "#up-" + tid
  
  //提交一个Json请求，获取数据库中的内容
  $.getJSON("/topics/vote." + tid +"_" + val+".json",
  function(data){
  	if (data.errno == 0) {
        $(up).text(data.like);
        $(dn).text("-" + data.hate);
  	};
  }); 
}