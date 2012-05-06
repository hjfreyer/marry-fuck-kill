
var mfk = {};

mfk.ACTIONS = { 'marry' : true, 'fuck' : true, 'kill' : true };

mfk.Triple = function(dom) {
  this.dom_ = dom;
  this.vote_ = this.dom_.find('.vote');
  this.edit_ = this.dom_.find('.edit');
};

mfk.Triple.prototype.decorate = function() {
  for (var action in mfk.ACTIONS) {
    for (var n = 0; n < 3; n++) {
      var button = this.dom_.find('.entity' + n + ' .vote-button.' + action);
      button.click(this.select.bind(this, action, n));
    }
  }

  this.vote_.click(this.vote.bind(this));
  this.edit_.click(this.edit.bind(this));
};

mfk.Triple.prototype.select = function(action, num) {
  if (this.dom_.hasClass('voted')) {
    return;
  }

  this.dom_.addClass('s' + action + num);
  for (var action2 in mfk.ACTIONS) {
    if (action2 != action) {
      this.dom_.removeClass('s' + action2 + num);
    }
  }
};

mfk.Triple.prototype.vote = function() {
  if (this.dom_.hasClass('voted')) {
    return;
  }

  var c = function(x) { return this.dom_.hasClass(x); }.bind(this);

  if ((c('smarry0') && c('sfuck1') && c('skill2')) ||
      (c('smarry0') && c('sfuck2') && c('skill1')) ||
      (c('smarry1') && c('sfuck0') && c('skill2')) ||
      (c('smarry1') && c('sfuck2') && c('skill0')) ||
      (c('smarry2') && c('sfuck0') && c('skill1')) ||
      (c('smarry2') && c('sfuck1') && c('skill0'))) {
    this.dom_.toggleClass('voted unvoted');
  }
};

mfk.Triple.prototype.edit = function() {
  if (!this.dom_.hasClass('voted')) {
    return;
  }
  this.dom_.attr('class', 'triple unvoted');
};

$(function() {
    $('.triple').each(function() {
        var t = new mfk.Triple($(this));
        t.decorate();
      });
  });
