
goog.provide('mfk');

goog.require('util');
goog.require('goog.dom');
goog.require('goog.dom.query');
goog.require('goog.array');
goog.require('goog.events');
goog.require("goog.dom");
goog.require("goog.net.XhrIo");
goog.require("goog.structs.Map");
goog.require("goog.Uri.QueryData");

goog.require('goog.ui.Dialog');
goog.require('goog.ui.LabelInput');
goog.require('goog.ui.LinkButtonRenderer');

// var mfk = {};

mfk.ACTIONS = { 'marry' : true, 'fuck' : true, 'kill' : true };

mfk.depsReady = false;
mfk.depsWaiters = [];

mfk.depsInit = function() {
  google.load('visualization', '1.0', {'packages':['corechart']});

  google.setOnLoadCallback(function() {
    mfk.depsReady = true;
    for (var i = 0; i < mfk.depsWaiters.length; i++) {
      mfk.depsWaiters[i]();
    }
  });
};

mfk.whenDepsReady = function(func) {
  if (mfk.depsReady) {
    func();
  } else {
    mfk.depsWaiters.push(func);
  }
};

/**
 * @constructor
 */
mfk.Triple = function(dom, triple_id, tallyA, tallyB, tallyC) {
  this.dom_ = dom;
  this.id_ = triple_id;

  this.tallyA_ = tallyA;
  this.tallyB_ = tallyB;
  this.tallyC_ = tallyC;

  this.vote_ = '';

  if (!goog.dom.classes.has(this.dom_, 'unvoted')) {
    this.setVoteFromClasses();
  }

  this.voteButton_ = goog.dom.query('.vote', this.dom_)[0];
  this.edit_ = goog.dom.query('.edit', this.dom_)[0];

  this.chartA_ = goog.dom.query('.entity0 .chart', this.dom_)[0];
  this.chartB_ = goog.dom.query('.entity1 .chart', this.dom_)[0];
  this.chartC_ = goog.dom.query('.entity2 .chart', this.dom_)[0];

  this.decorate();

  console.log(this.vote_);
  if (this.vote_ != '') {
    this.drawCharts();
  }
};

mfk.Triple.prototype.decorate = function() {
  for (var action in mfk.ACTIONS) {
    for (var n = 0; n < 3; n++) {
      var button = goog.dom.query('.entity' + n + ' .vote-button.' + action,
                                  this.dom_)[0];

      goog.events.listen(button, goog.events.EventType.CLICK,
                         this.select.bind(this, action, n));
    }
  }

  goog.events.listen(this.voteButton_, goog.events.EventType.CLICK,
                     this.vote.bind(this));
  goog.events.listen(this.edit_, goog.events.EventType.CLICK,
                     this.edit.bind(this));
};

mfk.Triple.prototype.select = function(action, num) {
  if (goog.dom.classes.has(this.dom_, 'voted')) {
    return;
  }

  goog.dom.classes.add(this.dom_, 's' + action + num);
  for (var action2 in mfk.ACTIONS) {
    if (action2 != action) {
      goog.dom.classes.remove(this.dom_, 's' + action2 + num);
    }
  }
};

mfk.Triple.prototype.setVoteFromClasses = function() {
  var c = goog.dom.classes.has.bind(this, this.dom_);

  var vote = '';
  if (c('smarry0')) vote += 'm';
  if (c('sfuck0')) vote += 'f';
  if (c('skill0')) vote += 'k';

  if (c('smarry1')) vote += 'm';
  if (c('sfuck1')) vote += 'f';
  if (c('skill1')) vote += 'k';

  if (c('smarry2')) vote += 'm';
  if (c('sfuck2')) vote += 'f';
  if (c('skill2')) vote += 'k';

  this.vote_ = vote;
};

mfk.Triple.prototype.vote = function() {
  this.setVoteFromClasses();
  this.drawCharts();

  // create the xhr object
  var request = new goog.net.XhrIo();

  // create a QueryData object by initializing it with a simple Map
  var data = goog.Uri.QueryData.createFromMap(new goog.structs.Map(
    {
      triple_id : this.id_,
      vote : this.vote_
    }));

  request.send('/api/v1/vote', 'POST', data.toString());

  goog.dom.classes.swap(this.dom_, 'unvoted', 'voted');
};

mfk.Triple.prototype.drawCharts = function() {
  mfk.whenDepsReady(function() {
    console.log(this.vote_);
    var drawChart = function(baseTally, voteChar, dom) {
      var data = google.visualization.arrayToDataTable([
        ['Vote', 'Votes'],
        ['Marry', baseTally['m'] + (voteChar == 'm' ? 1 : 0)],
        ['Fuck', baseTally['f'] + (voteChar == 'f' ? 1 : 0)],
        ['Kill', baseTally['k'] + (voteChar == 'k' ? 1 : 0)]
      ]);

      WIDTH = 208;
      HEIGHT = 104;
      H_MARGIN = 30;

      var options = {
        width: WIDTH,
        height: HEIGHT,
        colors : ['#C76FDD'],  // or ['#9911BB'], or ['#63067A'],
        legend : {
          position : 'none'
        },
        chartArea : {
          top : 15,
          left : H_MARGIN,
          height : HEIGHT - 35,
          width : WIDTH - 2 * H_MARGIN
        },
        hAxis: {
          textPosition: 'out',
          textStyle : {
            fontSize: 12,
            color: "#888"
          },
        },
        vAxis: {
          baselineColor : 'grey',
          textStyle : {
            fontSize: 12,
            color: "#BBB"
          },
        },
      };

      var chart = new google.visualization.ColumnChart(dom);
      chart.draw(data, options);
    };

    drawChart(this.tallyA_, this.vote_[0], this.chartA_);
    drawChart(this.tallyB_, this.vote_[1], this.chartB_);
    drawChart(this.tallyC_, this.vote_[2], this.chartC_);
  }.bind(this));
};

mfk.Triple.prototype.edit = function() {
  goog.dom.classes.set(this.dom_, 'triple unvoted');
};

// mfk.makeTriples = function() {
//   goog.array.forEach(goog.dom.query('.triple'),
//                      function(triple) {
// 					   var t = new mfk.Triple(triple);
// 					   t.decorate();
//                      });
// };

mfk.main = function() {
 // goog.array.forEach(goog.dom.query('.triple'),
 //                     function(triple) {
 //                       var t = new mfk.Triple(triple);
 //                       t.decorate();
 //                     });

};

goog.exportSymbol('mfk.main', mfk.main);
goog.exportSymbol('mfk.depsInit', mfk.depsInit);
goog.exportSymbol('mfk.Triple', mfk.Triple);
