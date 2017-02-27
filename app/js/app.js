'use strict';

var reload = setInterval(load, 5000)

function load() {
    fetch('/api/watchman')
        .then(res => res.json())
        .then(createCharts)
        .catch(console.error)
}

load()

function createCharts(data) {
    createRespTimeChart(data)
    createEventCountChart(data)
}

// https://jsfiddle.net/jonathansampson/m7G64/
function throttle(callback, limit) {
    var wait = false;                  // Initially, we're not waiting
    return function () {               // We return a throttled function
        if (!wait) {                   // If we're not waiting
            callback.call();           // Execute users function
            wait = true;               // Prevent future invocations
            setTimeout(function () {   // After a period of time
                wait = false;          // And allow future invocations
            }, limit);
        }
    }
}

function cancelReload() {
    clearInterval(reload)
}

function createRespTimeChart(data) {
    d3.selectAll('#chart1 > *').remove();
    var chart = d3.timeseries()
        .addSerie(
        data,
        { x: 'created', y: 'resp_time' },
        { interpolate: 'linear', width: 3 }
        )
        .addHandlers([{ name: 'brush', cb: throttle(cancelReload, 1500) }])
        .yscale.domain([0]) // show 0 on y axis
        .margin.left(40)
        .width(1200)

    chart('#chart1')
}


function createEventCountChart(data) {
    d3.selectAll('#chart2 > *').remove();
    var chart = d3.timeseries()
        .addSerie(
        data,
        { x: 'created', y: 'event_count' },
        { interpolate: 'linear', width: 3 }
        )
        .addHandlers([{ name: 'brush', cb: throttle(cancelReload, 1500) }])
        .yscale.domain([0]) // show 0 on y axis
        .margin.left(40)
        .width(1200)

    chart('#chart2')
}