var tableName= 'humstat';
var maxItems = 100;
//var yesterdayDateString =    '2018-02-16T17:00:00';

var d = new Date(8);
d.setDate(d.getDate() -1);
var yesterdayDateString =  d.getFullYear()+"-"+d.getMonth()+"-"+d.getDate()+"T00:00:00";
var eventToFind= 'UserSendMessages';
var hoursBack=0;

///////////////////////////////////////////////////////////////////
//Generating a string of the last X hours back
var ts = new Date().getTime();
var tsYesterday = (ts - (hoursBack * 3600) * 1000);
var d = new Date(tsYesterday);
var yesterdayDateString = d.getFullYear() + '-'
+ ('0' + (d.getMonth()+1)).slice(-2) + '-'
+ ('0' + d.getDate()).slice(-2);
// + 'T'
//+ ('0' + (d.getHours()+1)).slice(-2) + ':'
//+ ('0' + (d.getMinutes()+1)).slice(-2) + ':'
//+ ('0' + (d.getSeconds()+1)).slice(-2);

//Forming the DynamoDB Query
var params = {
	TableName: tableName,                
	Limit: maxItems,
	ConsistentRead: false,
	ScanIndexForward: true,
	ExpressionAttributeValues:{
		":start_date":yesterdayDateString,
		//":event_to_find":6
    },
    ExpressionAttributeNames: {
        "#TS": "id_timestamp"
     },
  //   FilterExpression :"#TS >= :start_date"
    KeyConditionExpression :
    "#TS = :start_date"
}
/*
/////////////////////////////////////////////////////

 
/////////////////////////////////////////////////////////////
*/
var docClient = new AWS.DynamoDB.DocumentClient();
docClient.query(params, function(err, data) {
if (err) console.log(err, err.stack); // an error occurred
else{
    //console.log(data)
    var recentEventsDateTime = [];
    var recentEventsCounter = [];
    var recentEventsCounter1 =[];
    var dateHour;
    data.Items.forEach(function(item) {
        dateHour = item.timestamp.toString();
       // recentEventsDateTime.push(dateHour.slice(0, -6));
        recentEventsDateTime.push(dateHour);
        if (typeof item.UserSendMessages === "undefined"){
            //console.log('the property is not available...'); // print into console
            item.UserSendMessages = 0;
        }
                if (typeof item['Post created'] === "undefined"){
            //console.log('the property is not available...'); // print into console
            item['Post created'] = 0;
        }
        recentEventsCounter.push(item.UserSendMessages.toString());
        recentEventsCounter1.push(item['Post created'].toString());
    });
//Chart.js code
var lineChartData = {
    labels : recentEventsDateTime,
    datasets : [
    {   
        title:  'User Send Messages',
        label: 'User Send Messages',
        text: 'User Send Messages',
        fillColor : "rgba(151,217,205,0.2)",
        strokeColor : "rgba(151,187,205,1)",
        pointColor : "rgba(151,197,255,1)",
        pointStrokeColor : "#fff",
        pointHighlightFill : "#fff",
        pointHighlightStroke : "rgba(151,217,205,1)",
        data : recentEventsCounter
     },{
        label: 'Post created',
        text: 'Post created',
        fillColor : "rgba(11,187,105,0.2)",
        strokeColor : "rgba(11,187,105,0.5)",
        pointColor : "rgba(11,197,255,0.5)",
        pointStrokeColor : "#fff",
        pointHighlightFill : "#fff",
        pointHighlightStroke : "rgba(11,187,105,1)",
        data : recentEventsCounter1
     }

    ]}
var ctx = document.getElementById("canvas").getContext("2d");
Chart.defaults.global.legend.display = true;
window.myLine = new Chart(ctx).Line(lineChartData, {
    responsive: false,
    pointDot: false
    });
}});