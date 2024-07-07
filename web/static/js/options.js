(function($){

"use strict";

 

$("#theme-color").bind(function(){

	 alert("changed");
});


 $(".setting-icons span").on("click", function(){

        $(".demo-wrapper").toggleClass("collapse-left");
 });




}(jQuery));