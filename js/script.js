
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

/**
 * @constructor
 */
mfk.Triple = function(dom) {
  this.dom_ = dom;
  this.id_ = dom.getAttribute('triple_id');
  this.vote_ = goog.dom.query('.vote', this.dom_)[0];
  this.edit_ = goog.dom.query('.edit', this.dom_)[0];
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

  goog.events.listen(this.vote_, goog.events.EventType.CLICK,
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

mfk.Triple.prototype.vote = function() {
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

  console.log(vote);

  // create the xhr object
  var request = new goog.net.XhrIo();

  // create a QueryData object by initializing it with a simple Map
  var data = goog.Uri.QueryData.createFromMap(new goog.structs.Map({
        'triple_id' : this.id_,
        'vote' : vote
      }));

  request.send('/api/v1/vote', 'POST', data.toString());

  goog.dom.classes.swap(this.dom_, 'unvoted', 'voted');
};

mfk.Triple.prototype.edit = function() {
  goog.dom.classes.set(this.dom_, 'triple unvoted');
};

mfk.main = function() {
 // goog.array.forEach(goog.dom.query('.triple'),
 //                     function(triple) {
 //                       var t = new mfk.Triple(triple);
 //                       t.decorate();
 //                     });

}

goog.exportSymbol('mfk.main', mfk.main);