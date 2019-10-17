function prepareChartOptions(cD){

    var options;

    switch(cD.type){
        case "line":
        case "radar":
        case "bar":
        case "area":

        options = {
            chart: {
                type: cD.type,
            },
            stroke: {
                curve: 'smooth',
                width: 2
            },
            series: cD.series,
            markers: {
                size: 4,
                strokeWidth: 0,
                hover: {
                size: 4
                }
            },
            grid: {
                show: true
            },
            labels: cD.periods,
            xaxis: {
                tooltip: {
                    enabled: false
                }
            },
            legend: {
                position: 'bottom',
            }
        }

        break;
        case "donut":
        case "radialBar":

        options = {
            chart: {
                type: cD.type,
            },
            stroke: {
                lineCap: 'round'
            },
            series: cD.series[1].data,
            labels: cD.periods,
            legend: {
                position: 'bottom',
            }
        }
        
        break;

        case "bubble":
        
        options = {
            chart: {
                type: cD.type,
            },
            series: cD.series,
            yaxis: {
                    max: 70
                }
            }
        break;

    }

    return options;

}