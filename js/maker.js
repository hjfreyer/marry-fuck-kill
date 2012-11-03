
goog.provide('mfk.maker');

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

/**
 * @constructor
 */
mfk.Maker = function(dom, imageSearch) {
  this.dom_ = dom;

  this.entityMakers_ = goog.array.map(
    goog.dom.getChildren(this.dom_),
    function(x) {
      return new mfk.EntityMaker(x, imageSearch);
    });
};

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
 */
mfk.EntityMaker = function(dom, imageSearch) {
  this.dom_ = dom;
  this.imageSearch_ = imageSearch;

  this.nameText_ = goog.dom.query('.preview .name', this.dom_)[0];
  this.nameTextWrap_ = new goog.ui.LabelInput;
  this.nameTextWrap_.decorate(this.nameText_);

  this.imgPreview_ = goog.dom.query('.preview .image', this.dom_)[0];

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

//  this.nameText_.rows = 3;
//  var foo= new mfk.AutoHeight(this.nameText_);

  this.searchText_ = goog.dom.query('.imagesearch .searchbar input', this.dom_);
  this.searchButton_ = goog.dom.query('.imagesearch .searchbar button',
                                      this.dom_);
  $(dom).find('form').submit(this.search.bind(this));

//  $(this.searchText_).val('skittles');
  console.log(this);
  var resultArea = goog.dom.query('.imagesearch .result-area', this.dom_)[0];
  console.log(resultArea);

  this.resultUrls_ = null;
  this.resultImgs_ = goog.dom.getChildren(resultArea);
  for (var i = 0; i < this.resultImgs_.length; i++) {
    $(this.resultImgs_[i]).click(this.selectImage.bind(this, i));
  }

  this.throbber_ = goog.dom.query('.throbber', this.dom_)[0];
  this.noResults_ = goog.dom.query('.no-results', this.dom_)[0];

//  this.search();
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
  console.log(idx);
  this.selectedImage_ = this.resultUrls_[idx];
  this.imgPreview_.style.backgroundImage= "url('"+ this.resultUrls_[idx] +"')";
};

mfk.MakeMaker = function() {
  var imageSearch = new mfk.ImageSearch();
  new mfk.Maker(goog.dom.getElement('maker'), imageSearch);
};

goog.exportSymbol('mfk.maker.main', mfk.MakeMaker);