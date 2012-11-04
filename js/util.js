
goog.provide('util');

goog.require('goog.array');
goog.require('goog.dom.classes');
goog.require('goog.events');

util.log = function(x) {
  console.log(x);
};

util.click = function(elem, cb) {
  goog.events.listen(elem, goog.events.EventType.CLICK, cb);
};

util.show = function(elem) {
  goog.dom.classes.remove(elem, 'hidden');
};

util.hide = function(elem) {
  goog.dom.classes.add(elem, 'hidden');
};

util.showAll = function(elems) {
  goog.array.forEach(elems, util.show);
};

util.hideAll = function(elems) {
  goog.array.forEach(elems, util.hide);
};
