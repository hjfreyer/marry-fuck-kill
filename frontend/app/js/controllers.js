'use strict';

/* Controllers */

var Page = ['$scope', 'StateStack', function($scope, StateStack) {
    $scope.page = {
	onKeyDown: function(event) {
	    if (event.keyCode == 27) {  // esc
		StateStack.pop();
	    }
	},
	popState: StateStack.pop.bind(StateStack)
    };
}];

var Entity = function(Searcher, StateStack) {
    this.searcher = Searcher;
    this.stateStack = StateStack;

    this.name = '';
    this.imgIdx = -1;
    this.query = '';
    this.images = [];
    this.imagesPromise = null;

    this.searching = false;
};

Entity.prototype.showSearch = function() {
    this.stateStack.push(
	function(imgIdx, query, images, imagesPromise) {
	    this.imgIdx = imgIdx;
	    this.query = query;
	    this.images = images;
	    this.imagesPromise = imagesPromise;

	    this.searching = false;
	}.bind(this, this.imgIdx, this.query, this.images, this.imagesPromise),
	function(activeElement) {
	    activeElement.focus();
	}.bind(this, document.activeElement));

    this.searching = true;
    this.query = this.name;
    this.search();
};

Entity.prototype.selectImage = function(idx) {
    this.imgIdx = idx;
    this.searching = false;
    this.stateStack.release();
};

Entity.prototype.getImage = function() {
    if (this.imgIdx == -1) {
	return '/s/mfk.png';
    }
    return this.images[this.imgIdx].thumbnail;
};

Entity.prototype.search = function() {
    console.log('foo');
    var imagesPromise = this.searcher(this.query);
    this.imagesPromise = imagesPromise;
    imagesPromise.then(function(result) {
	if (this.imagesPromise != imagesPromise) {
	    console.log('replaced');
	    return;
	}
	console.log(result);
	this.images = result.data.images;
    }.bind(this));
};

Entity.prototype.getReset = function() {
    return function(name, imgIdx, images, position) {
	this.name = name;
	this.imgIdx = imgIdx;
	this.images = images;
	this.position = position;
    }.bind(this, this.name, this.imgIdx, this.images, this.position);
};

Entity.prototype.imageUrls = function() {
    this.images.map(function() {

    }, this);
};

var MakerTable = ['$scope', 'Searcher', 'StateStack', function(
    $scope, Searcher, StateStack) {
    var entities = [];
    for (var i = 0; i < 3; i++) {
	entities[i] = new Entity(Searcher, StateStack);
    }

    $scope.table = {
	entities: entities,

	select: function(idx) {
	    console.log(idx);
	    console.log($scope.table.entities[idx]);
	    var reset = $scope.table.entities[idx].getReset();
	    $scope.table.entities[idx].searching = true;
	    StateStack.push(function() {
		$scope.table.entities[idx].searching = false;
		reset();
	    });
	},

	showSearch: function(entity) {
	    entity.position = 1;
	},

	search: function(entity) {
	},

	done: function() {
	    $scope.table.searchIdx = -1;
	    StateStack.release();
	}
    };
}];

angular.module('mfkMaker.controllers', [])
    .controller('Page', Page)
    .controller('MakerTable', MakerTable);
