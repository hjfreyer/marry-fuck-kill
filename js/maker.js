
goog.provide('mfk.maker');

goog.require('util');
goog.require('goog.dom');
goog.require('goog.dom.query');
goog.require('goog.array');
goog.require('goog.json');
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

  goog.net.XhrIo.send('/api/v1/imagesearch?query=' + query,
                      this.processResults.bind(this, query, callback));
};

mfk.ImageSearch.prototype.processResults = function(query, callback, event) {
  var results = event.target.getResponseJson();
  this.cache_[query] = results;
  callback(results);
};

mfk.clearLabelOnFocus = function(label, labelInput) {
  labelInput.setLabel(label);
  goog.events.listen(labelInput.getElement(),
                     goog.events.EventType.FOCUS,
                     function() {
                       labelInput.setLabel('');
                     });
  goog.events.listen(labelInput.getElement(),
                     goog.events.EventType.BLUR,
                     function() {
                       labelInput.setLabel(label);
                     });
};

/**
 * @constructor
 * @param {!mfk.ImageSearch} imageSearch
 */
mfk.Maker = function(dom, imageSearch) {
  this.dom_ = dom;

  this.createButton_ = Q('.create', this.dom_)[0];
  // this.next_.disabled = 'disabled';
  util.click(this.createButton_, this.onSubmit.bind(this));

  this.entityMakers_ = [];

  // this.names_ = [];
  // this.imageDivs_ = [];
  for (var i = 0; i < 3; i++) {
    var entityPreview = Q('.entity' + i, this.dom_)[0];
    var searchPanel = Q('.search-panel' + i, this.dom_)[0];

    this.entityMakers_.push(new mfk.EntityMaker(
      entityPreview,
      searchPanel,
      imageSearch,
      this.onChange.bind(this),
      this.onSearchSelect.bind(this, i)));
    // var entityDiv = Q('.entity' + i, this.dom_)[0];

    // var input = new goog.ui.LabelInput();
    // input.decorate(Q('textarea.name', entityDiv)[0]);
    // mfk.clearLabelOnFocus('Name me', input);

    // this.names_.push(input);

    // var imgDiv = Q('.image', entityDiv)[0];
    // util.click(imgDiv, this.entitySelected.bind(this, i));
    // this.imageDivs_.push(imgDiv);
    // this.entityMakers_.push(
    //   new mfk.EntityMaker(E('entity' + i),
    //                       imageSearch,
    //                       this.onChange.bind(this)));
    // this.entityMakers_[i].nameTextWrap_.setValue('' + i);
    // this.entityMakers_[i].onNameChange();
  }

  // this.imageForm_ = Q('.image-search .search-bar form', this.dom_)[0];
  // this.imageQuery_ = Q('input', this.imageForm_)[0];

  // this.review_ = E('review');
};

mfk.Maker.prototype.onSearchSelect = function(idx) {
  util.log(idx);
  this.displaySearch(idx);
};

mfk.Maker.prototype.displaySearch = function(idx) {
  goog.dom.classes.set(this.dom_, 'show-search-' + idx);
};

mfk.Maker.prototype.search = function() {
  var query = this.imageQuery_.value;

  this.imageSearch_.search(query, this.processResults.bind(this));
};

mfk.Maker.prototype.processResults = function(){

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
    this.createButton_.disabled = '';
  } else {
    this.createButton_.disabled = 'disabled';
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

mfk.Maker.prototype.onSubmit = function() {
  var data = {};

  data['a'] = {
    'name' : this.entityMakers_[0].getName(),
    'image' : this.entityMakers_[0].getImage()
  };
  data['b'] = {
    'name' : this.entityMakers_[1].getName(),
    'image' : this.entityMakers_[1].getImage()
  };
  data['c'] = {
    'name' : this.entityMakers_[2].getName(),
    'image' : this.entityMakers_[2].getImage()
  };

  goog.net.XhrIo.send('/api/v1/make',
                      this.processResults.bind(this),
                      'POST',
                      goog.json.serialize(data));
};

/**
 * @constructor
 */
mfk.EntityMaker = function(entityDom, searchDom, imageSearch,
                           onStateChange, onSearchSelect) {
  this.entityDom_ = entityDom;
  this.searchDom_ = searchDom;
  this.imageSearch_ = imageSearch;
  this.onStateChange_ = onStateChange;
  this.onSearchSelect_ = onSearchSelect;

  this.name_ = new goog.ui.LabelInput();
  this.name_.decorate(Q('textarea.name', this.entityDom_)[0]);
  this.imgPreview_ = Q('.image', this.entityDom_)[0];
  this.searchForm_ = Q('.search-bar form', this.searchDom_)[0];
  this.searchBox_ = Q('input', this.searchForm_)[0];
  this.throbber_ = Q('.throbber', this.searchDom_)[0];
  this.noResults_ = Q('.no-results', this.searchDom_)[0];
  this.resultSlots_ = Q('.result', this.searchDom_);

  util.cancelEnter(this.name_.getElement());

  goog.events.listen(this.name_.getElement(),
                     goog.events.EventType.KEYUP,
                     this.onNameChange.bind(this));

  mfk.clearLabelOnFocus('Give me a name', this.name_);
  util.click(this.imgPreview_, this.onShowSearch.bind(this));

  goog.events.listen(this.searchForm_,
                     goog.events.EventType.SUBMIT,
                     function(e) {
                       this.search();
                       e.preventDefault();
                     }.bind(this));

  this.firstTime_ = true;
  this.searchCount_ = 0;
  this.lastResults_ = null;
  this.selectedImage_ = null;

  for (var i = 0; i < this.resultSlots_.length; i++) {
    util.click(this.resultSlots_[i], this.selectImage.bind(this, i));
  }
};

mfk.EntityMaker.prototype.getName = function() {
  return this.name_.getValue();
};

mfk.EntityMaker.prototype.getImage = function() {
  return this.selectedImage_;
};

mfk.EntityMaker.prototype.onNameChange = function() {
  //this.searchBox_.value = this.name_.getValue();
};

mfk.EntityMaker.prototype.onShowSearch = function() {
  if (this.name_.getValue() == '') {
    util.log('Name not set');
    return;
  }

  if (this.firstTime_) {
    this.searchBox_.value = this.name_.getValue();
    this.search();
    this.firstTime_ = false;
  }
  this.onSearchSelect_();
};

// mfk.EntityMaker.prototype.onNameChange = function() {
//   var trimmed = goog.string.trim(this.nameTextWrap_.getValue());

//   if (trimmed == this.lastQuery_) {
//     return;
//   }
//   this.lastQuery_ = trimmed;

//   this.nameChangeCount_++;
//   this.hideAll();

//   if (trimmed != '') {
//     util.show(this.throbber_);
//     setTimeout(this.search.bind(this, trimmed, this.nameChangeCount_), 1000);
//   }
//   this.onStateChange_();
// };

mfk.EntityMaker.prototype.search = function() {
  var query = this.searchBox_.value;
  util.log('Searching: ' + query);

  this.searchCount_++;

  this.hideAll();
  if (query != '') {
    util.show(this.throbber_);
    this.imageSearch_.search(query, this.processResults.bind(
      this, query, this.searchCount_));
  }
};

mfk.EntityMaker.prototype.hideAll = function() {
  util.log(this);
  util.hideAll(this.resultSlots_);
  util.hide(this.throbber_);
  util.hide(this.noResults_);
};

mfk.EntityMaker.prototype.processResults = function(query, searchNum, results) {
  if (searchNum != this.searchCount_) {
    util.log("Results obsolete: " + query);
    return;
  }

  this.hideAll();
  util.log('Search results');
  util.log(results);

  this.results_ = results;
  if (results.image.length == 0) {
    util.show(this.noResults_);
    return;
  }

  for (var i = 0; i < this.resultSlots_.length; i++) {
    var slot = this.resultSlots_[i];
    if (i < results.image.length) {
      var result = results.image[i].metadata;
      util.show(slot);
      slot.style.backgroundImage = "url('"+ result.url +"')";
    }
  }
};

mfk.EntityMaker.prototype.selectImage = function(idx) {
  this.selectedImage_ = this.results_.image[idx];
  this.imgPreview_.style.backgroundImage =
    "url('"+ this.selectedImage_.metadata.url + "')";

  this.onStateChange_();
};

mfk.maker.main = function() {
  var imageSearch = new mfk.ImageSearch();
  new mfk.Maker(E('maker'), imageSearch);
};

goog.exportSymbol('mfk.maker.main', mfk.maker.main);
