//	菜单
$('li.dropdown').mouseover(function() {   
     $(this).addClass('open');    }).mouseout(function() {        $(this).removeClass('open');    }); 

function adjustFooter(){
return;
      if($(window).height()!=$(document).height()){
        $("#id_footer").removeClass("navbar-fixed-bottom");
      } else {
		if (!$("#id_footer").hasClass("navbar-fixed-bottom")) {
			$("#id_footer").addClass("navbar-fixed-bottom");
		}
	  }
}

$(document).ready(function(){
	adjustFooter();
})

$(window).resize(function(){
	adjustFooter();
})