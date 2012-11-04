
goog.provide('mfk.maker');

goog.require('util');
goog.require('goog.dom');
goog.require('goog.dom.query');
goog.require('goog.array');
goog.require('goog.events');
goog.require('goog.debug.Logger');
goog.require('goog.debug.Console');

goog.require("goog.dom");
goog.require("goog.net.XhrIo");
goog.require("goog.structs.Map");
goog.require("goog.Uri.QueryData");

goog.require('goog.ui.Dialog');
goog.require('goog.ui.LabelInput');
goog.require('goog.ui.LinkButtonRenderer');

var E = goog.dom.getElement;
var Q = goog.dom.query;

/**
 * @constructor
 */
mfk.ImageSearch = function() {
  this.cache_ = {};
};

mfk.ImageSearch.prototype.search = function(query, callback) {
  if (query in this.cache_) {
    callback(this.cache_[query]);
    return;
  }

  $.getJSON('/api/v1/imagesearch',
            { 'query' : query },
            this.processResults.bind(this, query, callback));
};

mfk.ImageSearch.prototype.processResults = function(query, callback, results) {
  this.cache_[query] = results;
  callback(results);
};

/**
 * @constructor
 * @param {!mfk.ImageSearch} imageSearch
 */
mfk.Maker = function(imageSearch) {
  this.entityMakers_ = [];

  this.next_ = E('next');
  this.next_.disabled = 'disabled';
  util.click(this.next_, this.showReview.bind(this));

  for (var i = 0; i < 3; i++) {
    this.entityMakers_.push(
      new mfk.EntityMaker(E('entity' + i),
                          imageSearch,
                          this.onChange.bind(this)));
    this.entityMakers_[i].nameTextWrap_.setValue('' + i);
    this.entityMakers_[i].onNameChange();
  }

  this.review_ = E('review');
};

mfk.Maker.prototype.showReview = function() {
  for (var i = 0; i < 3; i++) {
    util.log(i);
    var reviewDom = E('entityReview' + i);
    Q('.name', reviewDom)[0].value = this.entityMakers_[i].getName();
  }
};

mfk.Maker.prototype.onChange = function() {
  if (this.readyToSubmit()) {
    this.next_.disabled = '';
  } else {
    this.next_.disabled = 'disabled';
  }
};

mfk.Maker.prototype.readyToSubmit = function() {
  for (var i = 0; i < 3; i++) {
    if (this.entityMakers_[i].getName() == '' ||
        this.entityMakers_[i].getImage() == null) {
      return false;
    }
  }
  return true;
};

/**
 * @constructor
 */
mfk.EntityMaker = function(dom, imageSearch, onStateChange) {
  this.dom_ = dom;
  this.imageSearch_ = imageSearch;
  this.onStateChange_ = onStateChange;

  this.nameText_ = Q('.preview .name', this.dom_)[0];
  this.nameTextWrap_ = new goog.ui.LabelInput;
  this.nameTextWrap_.decorate(this.nameText_);

  this.imgPreview_ = Q('.preview .image', this.dom_)[0];

  $(this.nameText_).focus(function () {
    console.log('foo');
    this.nameTextWrap_.setLabel('');
  }.bind(this));
  $(this.nameText_).blur(function () {
    this.nameTextWrap_.setLabel('Name me');
  }.bind(this));

  this.nameChangeCount_ = 0;
  this.lastSearch_ = null;
  $(this.nameText_).bind('keyup', this.onNameChange.bind(this));
  $(this.nameText_).change(this.onNameChange.bind(this));

  this.searchText_ = Q('.imagesearch .searchbar input', this.dom_);
  this.searchButton_ = Q('.imagesearch .searchbar button', this.dom_);
  $(dom).find('form').submit(this.search.bind(this));

//  $(this.searchText_).val('skittles');
  util.log("fooooo");
  util.log(this);
  var resultArea = Q('.imagesearch .result-area', this.dom_)[0];
  console.log(resultArea);

  this.resultUrls_ = null;
  this.resultImgs_ = goog.dom.getChildren(resultArea);
  for (var i = 0; i < this.resultImgs_.length; i++) {
    $(this.resultImgs_[i]).click(this.selectImage.bind(this, i));
  }

  this.throbber_ = Q('.throbber', this.dom_)[0];
  this.noResults_ = Q('.no-results', this.dom_)[0];

//  this.search();
};

mfk.EntityMaker.prototype.getName = function() {
  return goog.string.trim(this.nameTextWrap_.getValue());
};

mfk.EntityMaker.prototype.getImage = function() {
  return this.selectedImage_;
};

mfk.EntityMaker.prototype.onNameChange = function() {
  var trimmed = goog.string.trim(this.nameTextWrap_.getValue());

  if (trimmed == this.lastQuery_) {
    return;
  }
  this.lastQuery_ = trimmed;

  this.nameChangeCount_++;
  this.hideAll();

  if (trimmed != '') {
    util.show(this.throbber_);
    setTimeout(this.search.bind(this, trimmed, this.nameChangeCount_), 1000);
  }
  this.onStateChange_();
};

mfk.EntityMaker.prototype.search = function(query, changeNum) {
  if (changeNum != this.nameChangeCount_) {
    console.log("Search obsolete: " + query);
    return;
  }

  console.log('Searching: ' + query);
  this.imageSearch_.search(query,
                           this.processResults.bind(this, query, changeNum));
};

mfk.EntityMaker.prototype.hideAll = function() {
  console.log(this);
  util.hideAll(this.resultImgs_);
  util.hide(this.throbber_);
  util.hide(this.noResults_);
};

mfk.EntityMaker.prototype.processResults = function(query, changeNum, results) {
  if (changeNum != this.nameChangeCount_) {
    console.log("Results obsolete: " + query);
    return;
  }

  this.hideAll();
  console.log('Search results');
  console.log(results);

  if (results.images == null) {
    util.show(this.noResults_);
    return;
  }

//  this.deselect();
  this.results_ = results;

  if (typeof(results.images) == 'undefined') {
    this.clear();
    //this.no_results_.toggleClass
    return;
  }

  this.resultUrls_ = [];
  for (var i = 0; i < this.resultImgs_.length; i++) {
    var slot = this.resultImgs_[i];
    if (i < results.images.length) {
      var result = results.images[i];
      util.show(slot);
      this.resultUrls_.push(result.url);

      // goog.style.setTransparentBackgroundImage(slot, result.url);
      slot.style.backgroundImage= "url('"+ result.url +"')";
      // slot.find('img').attr('src', result.url);
      // slot.removeClass('empty');
    }
  }
};

mfk.EntityMaker.prototype.selectImage = function(idx) {
  this.selectedImage_ = this.resultUrls_[idx];
  this.imgPreview_.style.backgroundImage= "url('"+ this.resultUrls_[idx] +"')";

  this.onStateChange_();
};

mfk.maker.main = function() {
  var imageSearch = new mfk.ImageSearch();
  new mfk.Maker(imageSearch);
};

goog.exportSymbol('mfk.maker.main', mfk.maker.main);
