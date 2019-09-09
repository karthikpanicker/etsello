 $( document ).ready(function() {
    $('#boards button').on("click",function(){ 
        var boardId = $(this).attr("id");
        loadBoardLists(boardId,$(this));
     });
     if ($('#boards button').length > 0) {
         var selectedElem = $('#boards button')[0];
        if ($('#boards button.active').length > 0){
            selectedElem = $('#boards button.active')[0];
        }
        var boardId = $(selectedElem).attr('id');
        loadBoardLists(boardId,$(selectedElem));
     }

     $("#saveList").on("click",function(){
        var $this = $(this);
        var loadingText = '<i class="fa fa-circle-o-notch fa-spin"></i> Saving...';
        if ($(this).html() !== loadingText) {
            $this.data('original-text', $(this).html());
            $this.html(loadingText);
        }
        var bId = $('#boards button').filter('.active').attr('id');
        var lId = $('#boardLists button').filter('.active').attr('id');
        var fields = []
        if ($('#listing_desc').prop("checked")){
            fields.push("listing_desc");
        }
        if ($('#listing_image').prop("checked")){
            fields.push("listing_image");
        }
        if ($('#listing_buy_profile').prop("checked")){
            fields.push("listing_buy_profile");
        }
        if ($('#listing_link').prop("checked")){
            fields.push("listing_link");
        }
        var transactionFilter = parseInt($('#trello-radio-set input:radio:checked').val(),10);
        $.ajax({
            type: "POST",
            url: "api/user-info",
            data: JSON.stringify({"boardId": bId, "listId": lId,
                "fieldsToUse": fields, "transactionFilter": transactionFilter}),
            success: function(data){
                setTimeout(function () {
                    $this.html($this.data('original-text'));
                }, 1000);
            }
          });
    });
    var left  = ($(window).width()/2)-(500/2);
    var top   = ($(window).height()/2)-(600/2);
    $("#etsy-authorize").on("click",function(){
        window.open('/authorize-etsy','EtsyAuthorize', "width=500, height=600, top=" + top + ", left=" + left);
    });
    $("#trello-authorize").on("click",function(){
        window.open('/authorize-trello','TrelloAuthorize', "width=500, height=600, top=" + top + ", left=" + left);
    });

    $('.github-link').on("click",function(e){
        e.preventDefault(); 
        window.open('https://github.com/karthikpanicker/etsello','_blank');
    });

    $("#gtask-authorize").on("click",function(){
        window.open('/authorize-gtask','GoogleAuthorize', "width=500, height=600, top=" + top + ", left=" + left);
    });
});

function loadBoardLists(boardId,selectedBoard) {
    $('#spinner-board-list').show();
    $('#boards button').removeClass('active');
    selectedBoard.addClass('active');
    $.get( "api/trello-boards/" + boardId + "/lists", function( boardLists ) {
        $('#boardLists button').remove();
        $.each(boardLists, function( index, list ) {
            var activeClass = (list.isSelected) ? ' active' : '' 
            $('#boardLists').append('<button id="' + list.id 
            + '" type="button" class="list-group-item' + activeClass +'">' + list.name + '</button>');
        });
        if ($('#boardLists button').length > 0) {
            // Bind click action on newly added items
            $('#boardLists button').on("click",function(){ 
                $('#boardLists button').removeClass('active');
                $(this).addClass('active');
                $('#spinner-board-list').hide();
            });
        }
      }); 
}