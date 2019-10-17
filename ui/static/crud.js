function list(data, type){

    var listTable = $("#" + type + "List");
    var chartSources = $("#chartSources");

    listTable.html("");
                            
    if (type == "source"){
        $("#chartSources").empty();
    }
    
    if (
        data === "Not found any documents" || 
        data === "null" ||
        data === null
    ){
        switch(type){
            case "chart":
            listTable.append( '<tr><td colspan="6">Not found any documents</td></tr>' );
            break;
            case "source":
            listTable.append( '<tr><td colspan="7">Not found any documents</td></tr>' );
            break;
        }
        
    }else {

        for(var i = 0; i < data.length; i++) {
            var obj = data[i];

            var row = "";

            switch(type){
                case "source":

                    chartSources.append($('<option>', {
                        value: obj.unique,
                        text: obj.name
                    }));

                    row = `<tr> \
                    <td>`+ obj.name +`</td> \
                    <td>`+ obj.address +`</td> \
                    <td>`+ obj.user +`</td> \
                    <td>`+ obj.password +`</td> \
                    <td>`+ obj.database +`</td> \
                    <td>`+ obj.collection +`</td> \
                    <td> \
                        <button \ 
                            data-unique="`+ obj.unique +`" \
                            data-name="`+ obj.name +`" \
                            data-address="`+ obj.address +`" \
                            data-user="`+ obj.user +`" \
                            data-password="`+ obj.password +`" \
                            data-database="`+ obj.database +`" \
                            data-collection="`+ obj.collection +`" \
                            data-d="source" \
                            type="button" \
                            class="updateButton btn btn-info"> \
                            Update \
                        </button> \
                        <button \ 
                            data-unique="`+ obj.unique +`" \
                            data-d="source" \
                            type="button" \
                            class="deleteButton btn btn-danger"> \
                            Delete \
                        </button> \
                    </td> \
                    </tr>`;
                
                break;
                case "chart":

                    row = `<tr> \
                    <td>`+ obj.source +`</td> \
                    <td>`+ obj.title +`</td> \
                    <td>`+ obj.type +`</td> \
                    <td>`+ obj.interval +`</td> \
                    <td>\
                        <button \ 
                            data-unique="`+ obj.unique +`"\
                            target="_blank"\
                            type="button" \
                            class="btn btn-primary viewChart"> \
                            View \
                        </button> \
                    </td> \
                    <td> \
                        <button \ 
                            data-unique="`+ obj.unique +`" \
                            data-source="`+ obj.source +`" \
                            data-title="`+ obj.title +`" \
                            data-type="`+ obj.type +`" \
                            data-query='`+ obj.query +`' \
                            data-interval="`+ obj.interval +`" \
                            data-d="chart" \
                            type="button" \
                            class="updateButton btn btn-info"> \
                            Update \
                        </button> \
                        <button \ 
                            data-unique="`+ obj.unique +`" \
                            data-d="chart" \
                            type="button" \
                            class="deleteButton btn btn-danger"> \
                            Delete \
                        </button> \
                    </td> \
                    </tr>`;
                
                break;
            }

            listTable.append( row );

        }
        $(".viewChart").on("click",chartNav);
        $(".updateButton").on("click",updateModal);
        $(".deleteButton").on("click",deleteObject);
        
    }

};


function mapValues($inputs, d, obj){

    $inputs.each(function() {

        if (this.name == "query"){
            $(this).val(JSON.stringify(obj.data(this.name)));
        }else{
            $(this).val(obj.data(this.name));
        }

    });

}

function mapInputs($inputs, d ){

    var data = {};

    $inputs.each(function() {
        
        switch(d){
            case "source":

                switch(this.name){
                    case "unique":
                        data.unique = $(this).val();
                    break;
                    case "name":
                        data.name = $(this).val();
                    break;
                    case "address":
                        data.address = $(this).val();
                    break;
                    case "user":
                        data.user = $(this).val();
                    break;
                    case "password":
                        data.password = $(this).val();
                    break;
                    case "database":
                        data.database = $(this).val();
                    break;
                    case "collection":
                        data.collection = $(this).val();
                    break;
                }

            break;
            case "chart":

                switch(this.name){
                    case "unique":
                        data.unique = $(this).val();
                    break;
                    case "source":
                        data.source = $(this).val();
                    break;
                    case "title":
                        data.title = $(this).val();
                    break;
                    case "type":
                        data.type = $(this).val();
                    break;
                    case "query":
                        data.query = $(this).val();
                    break;
                    case "interval":
                        data.interval = $(this).val();
                    break;
                }

            break;

        }

        $(this).val("");
    });

    return data

};

// Add object
var addModal = function (){

    var d = $(this).data("d");
    
    $("#" + d + "Action").val("create");

    $("#" + d + "Modal").modal('show');
};

var updateModal = function (){

    var d = $(this).data("d");
    
    $("#" + d + "Action").val("update");

    var $inputs = $("." + d + "Crud");

    var ilen = $inputs.length;

    if (ilen != 0){

        mapValues($inputs, d, $(this));               

    }

    $("#" + d + "Modal").modal('show');

};

var crudObject = function (){

    var d = $(this).data("d");
    var action = $("#" + d + "Action").val();

    var $inputs = $("." + d + "Crud");

    var ilen = $inputs.length;

    if (ilen != 0){
        
        var data = mapInputs($inputs, d);

        var message = {
            sender:"client",
            action: action,
            type: d
        };  

        switch (d){
            case "source":
            message.source = data;
            break;
            case "chart":
            message.chart = data;
            break;
        }

        var json = JSON.stringify(message);
        wsConn.send(json);

    }
    
    $("#" + d + "Action").val("");
    $("#" + d + "Modal").modal('hide');

};

// Delete object
var deleteObject = function (){

    var d = $(this).data("d");
    var unique = $(this).data("unique");
    
    var message = {
        sender:"client",
        action:"delete",
        type:d,
    };  

    switch (d){
        case "source":
        message.source = {unique:unique};
        break;
        case "chart":
        message.chart =  {unique:unique};
        break;
    }

    var json = JSON.stringify(message);
    wsConn.send(json);

};