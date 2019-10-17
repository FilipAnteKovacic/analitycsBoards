var wsConn;
var chartOpt;
var chartCanvas;

// Admin nav
var adminNav = function (){

    $(".view").hide();
    $("#" + $(this).data("nav")).show();
    $( "#navbarToggler" ).removeClass( "show" );
    
};

// Admin nav
var chartNav = function (){

    $(".view").hide();
    $("#chart").show();

    unique = $(this).data("unique");

    var connected = {
        sender:"client",
        action:"read",
        type:"chart",
        chart:{unique:unique}
    }

    var json = JSON.stringify(connected);
    wsConn.send(json);

};

$( document ).ready(function() {
    
    $(".view").hide();
    $("#sources").show();

    $(".adminLink").on("click",adminNav);
    $(".viewChart").on("click",chartNav);

    ws = $("#ws").val();

    chartRender = $("#chartCanvas");

    wsConn = new WebSocket(ws);

    wsConn.onclose = function(evt) {
        console.log('Connection closed')
    }
    wsConn.onmessage = function(evt) {

        var msg = JSON.parse(evt.data);

        switch (msg.type){
            case "list":

            list(msg.sources,"source");
            list(msg.charts,"chart");

            break;
            case "chart":

            if (msg.action == "updateData"){

                if (chartRender.data("unique") == msg.chart.unique){
                    chartOpt.updateSeries(msg.chart.series)
                }

            }else{

                chartRender.empty();
                chartRender.data("unique",msg.chart.unique);

                $(".chartHeader").html(msg.chart.title);
    
                chartOpt = new ApexCharts(
                    document.querySelector('#chartCanvas'), 
                    prepareChartOptions(msg.chart)
                );
    
                chartOpt.render();

            }

            break;
        }

    }
    wsConn.onopen = function(evt) {

        var connected = {
            sender:"client",
            action:"list",
        }

        var json = JSON.stringify(connected);
        wsConn.send(json);
    }   

    $(".addButton").on("click",addModal);
    $(".saveButton").on("click",crudObject);

});