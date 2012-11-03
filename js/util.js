
goog.provide('util');

goog.require('goog.array');
goog.require('goog.dom.classes');

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
