
goog.provide('mfk');

goog.require('util');
goog.require('goog.dom');
goog.require('goog.dom.query');
goog.require('goog.array');
goog.require('goog.events');
goog.require('goog.dom');
goog.require('goog.net.XhrIo');
goog.require('goog.structs.Map');
goog.require('goog.Uri.QueryData');

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
mfk.Triple = function(dom, triple_id, tallies) {
  this.dom_ = dom;
  this.id_ = triple_id;
  this.tallies_ = tallies;

  this.charts_ = goog.dom.query('.responses .others .chart', this.dom_);

  var voteEntities = goog.dom.query('.vote-area .vote-entity', this.dom_);
  for (var n = 0; n < 3; n++) {
    var e = goog.dom.getChildren(voteEntities[n]);

    util.click(e[0], this.selectButton_.bind(this, 'marry', n));
    util.click(e[1], this.selectButton_.bind(this, 'fuck', n));
    util.click(e[2], this.selectButton_.bind(this, 'kill', n));
  }

  if (this.isVoted()) {
    this.drawCharts_(this.getVote());
  }
};

mfk.Triple.prototype.isVoted = function() {
  return goog.dom.classes.has(this.dom_, 'voted');
};

mfk.Triple.prototype.getVote = function() {
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

  return vote;
};

mfk.Triple.prototype.setVoted = function() {
  if (this.isVoted()) {
    console.log('Error: already voted.');
    return;
  }

  var vote = this.getVote();
  if (vote.length != 3) {
    console.log('Error: Invalid vote: ' + vote);
    return;
  }

  this.drawCharts_(vote);

  goog.dom.classes.swap(this.dom_, 'unvoted', 'voted');
};

mfk.Triple.prototype.clearVotes = function() {
  goog.dom.classes.set(this.dom_, 'triple unvoted');
};

mfk.Triple.prototype.selectButton_ = function(action, num) {
  if (this.isVoted()) {
    return;
  }

  goog.dom.classes.add(this.dom_, 's' + action + num);
  for (var action2 in mfk.ACTIONS) {
    if (action2 != action) {
      goog.dom.classes.remove(this.dom_, 's' + action2 + num);
    }
  }
};

mfk.Triple.prototype.drawCharts_ = function(vote) {
  this.charts_.map(function(x) { goog.dom.removeChildren(x); });
  mfk.whenDepsReady(function() {
    var drawChart = function(tally, dom) {
      var data = google.visualization.arrayToDataTable([
        ['Vote', 'Votes'],
        ['Marry', tally['m']],
        ['Fuck', tally['f']],
        ['Kill', tally['k']]
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
            color: '#888'
          },
        },
        vAxis: {
          baselineColor : 'grey',
          textStyle : {
            fontSize: 12,
            color: '#BBB'
          },
        },
      };

      var chart = new google.visualization.ColumnChart(dom);
      chart.draw(data, options);
    };

    for (var n = 0; n < 3; n++) {
      drawChart(this.tallies_[n], this.charts_[n]);
    }
  }.bind(this));
};


mfk.sendVote = function(triple_id, vote, callback) {
  var postData = goog.Uri.QueryData.createFromMap(new goog.structs.Map(
    {
      triple_id : triple_id,
      vote : vote
    })).toString();

  goog.net.XhrIo.send('/api/v1/vote', callback, 'POST', postData);
};

mfk.SingleTriple = function(dom, triple_id, tallies) {
  this.triple_ = new mfk.Triple(dom, triple_id, tallies);
  this.voteButton_ = goog.dom.query('.vote', dom)[0];

  util.click(this.voteButton_, function() {
    this.triple_.setVoted();
  }.bind(this));
//  this.
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
goog.exportSymbol('mfk.SingleTriple', mfk.setupSingleTriple);
