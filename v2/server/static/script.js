
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

var mfk = {};

mfk.ACTIONS = { 'marry' : true, 'fuck' : true, 'kill' : true };

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

mfk.Maker = function(dom) {
  this.dom_ = $(dom);

  var imgSearch = new mfk.ImageSearch($('<div/>').appendTo(this.dom_));

  this.imgBox0_ = new mfk.ImageBox($('<div/>').appendTo(this.dom_),
                                   imgSearch);
  this.imgBox1_ = new mfk.ImageBox($('<div/>').appendTo(this.dom_),
                                   imgSearch);
  this.imgBox2_ = new mfk.ImageBox($('<div/>').appendTo(this.dom_),
                                   imgSearch);

  this.createButton_

  this.imgBox0_.name_.val('The Father');
};

mfk.ImageBox = function(dom, imgSearch) {
  this.dom_ = dom;
  this.imgSearch_ = imgSearch;

  this.name_ = $('<input/>').appendTo(this.dom_);

  this.link_ = $('<button>Pick an Image</button>').appendTo(this.dom_).click(this.pickImage.bind(this));
  this.img_ = $('<img/>').appendTo(this.dom_);

};

mfk.ImageBox.prototype.pickImage = function() {
  console.log(this.name_.val());
  this.imgSearch_.show(this.name_.val(), this.chooseImage.bind(this));
};

mfk.ImageBox.prototype.chooseImage = function(result) {
  console.log(result);
  this,result_ = result;
  this.img_.attr('src', result.url);
}

mfk.ImageSearch = function(dom) {
  // this.dom_ = $(dom);

  this.visible_ = false;

  this.dialog_ = new goog.ui.Dialog();
  this.dom_ = $(this.dialog_.getContentElement());
  this.dialog_.setTitle('Image Search');

  var buttonSet = goog.ui.Dialog.ButtonSet.createOkCancel()
  buttonSet.setDefault(null);
  this.dialog_.setButtonSet(buttonSet);


  this.form_ = $('<form>').appendTo(this.dom_).submit(function(){
      this.search();
      return false;
    }.bind(this));
  this.query_ = $('<input/>').appendTo(this.form_);
  $('<input type="submit"/>').appendTo(this.form_);


  $('<hr/>').appendTo(this.dom_);


  this.resultsDiv_ = $('<div/>').addClass('results').appendTo(this.dom_);
  this.results_ = [];
  this.selected_ = null;

  this.resultSlots_ = [];
  this.selectedSlot_ = null;

  for (var i = 0; i < 12; ++i) {
    var resultDiv = $('<a/>').attr('href', 'javascript:void(0)').addClass('result result' + i).appendTo(this.resultsDiv_).click(
                                                                                                                                this.select.bind(this, i));
    // resultDiv.click(function(div) {
    //     if (this.selected_ != null) {
    //       this.selected_.removeClass('selected');
    //     }
    //     this.selected_ = div;
    //     div.addClass('selected');
    //   }.bind(this, resultDiv));
    $('<div/>').addClass('centerbox').append($('<img/>')).appendTo(resultDiv);
    this.resultSlots_.push(resultDiv);
  }

  goog.events.listen(this.dialog_, goog.ui.Dialog.EventType.SELECT, function(e) {
      if (e.key == 'ok') {
        this.ok_cb_(this.selected_);
      }
      this.visible_ = false;
      console.log('You chose: ' + e.key);
}.bind(this));

};

mfk.ImageSearch.prototype.show = function(query, ok_cb) {
  if (this.visible_) {
    console.log("Already visible");
    return;
  }

  this.clear();
  this.visible_ = true;
  this.ok_cb_ = ok_cb;

  this.query_.val(query);
  this.search();
  this.dialog_.setVisible(true);
  this.query_.focus();
};

mfk.ImageSearch.prototype.clear = function() {
  this.dialog_.getButtonSet().setButtonEnabled('ok', false);

  this.deselect();
  this.selected_ = null;
  this.selectedSlot_ = null;

  $('.result').addClass('empty');
};

mfk.ImageSearch.prototype.deselect = function() {
  if (this.selectedSlot_ != null) {
    this.selectedSlot_.removeClass('selected');
  }
}

mfk.ImageSearch.prototype.select = function(i) {
  this.deselect();
  this.selected_ = this.results_.images[i];
  this.selectedSlot_ = this.resultSlots_[i];
  this.selectedSlot_.addClass('selected');

  this.dialog_.getButtonSet().setButtonEnabled('ok', true);
};

mfk.ImageSearch.prototype.search = function() {
  console.log(this.query_.val());
  if (this.query_.val() == '') {
    return;
  }

  $.getJSON('/api/v1/imagesearch',
{ 'query' : this.query_.val() },
            this.processResults.bind(this));
};

mfk.ImageSearch.prototype.processResults = function(results) {
  this.deselect();
  this.results_ = results;

  if (typeof(results.images) == 'undefined') {
    this.clear();
    //this.no_results_.toggleClass
    return;
  }

  console.log(results);

  for (var i = 0; i < this.resultSlots_.length; i++) {
    var slot = this.resultSlots_[i];
    if (i < results.images.length) {
      var result = results.images[i];
      slot.find('img').attr('src', result.url);
      slot.removeClass('empty');
    }
  }
};

function main() {
  goog.array.forEach(goog.dom.query('.triple'),
                     function(triple) {
                       var t = new mfk.Triple(triple);
                       t.decorate();
                     });

  goog.array.forEach(goog.dom.query('#maker'),
                     function(maker) {
                       var t = new mfk.Maker(maker);
                       //t.decorate();
                     });

}
