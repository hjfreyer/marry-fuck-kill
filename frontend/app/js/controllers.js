'use strict';

/* Controllers */

var Entity = function(Searcher, StateStack) {
    this.searcher = Searcher;
    this.stateStack = StateStack;

    this.name = '';
    this.query = '';
    this.images = [];
    this.imageIdx = -1;
    this.imagesPromise = null;

    this.pristine = true;
    this.searchShown = false;
    this.searching = false;
    this.error = false;
};

Entity.FIELDS = ['name', 'query', 'images', 'imageIdx',
		 'imagesPromise', 'pristine', 'searchShown',
		 'searching', 'error'];

Entity.prototype.dump = function() {
    var result = {}
    Entity.FIELDS.forEach(function(attr) {
	result[attr] = this[attr];
    }, this);
    return result;
};

Entity.prototype.load = function(source) {
    Entity.FIELDS.forEach(function(attr) {
	this[attr] = source[attr];
    }, this);
};

Entity.prototype.showSearch = function() {
    var saved = this.dump();
    this.stateStack.push(
	function() {
	    this.load(saved);
	}.bind(this),
	function(activeElement) {
	    activeElement.focus();
	}.bind(this, document.activeElement));

    this.searchShown = true;
    this.query = this.name;
    this.search();
};

Entity.prototype.selectImage = function(idx) {
    this.imageIdx = idx;
    this.searchShown = false;
    this.stateStack.release();
};

Entity.prototype.getImage = function() {
    if (this.imageIdx == -1) {
        return '/s/mfk.png';
    }
    return this.images[this.imageIdx].thumbnail;
};

Entity.prototype.isPlaceholder = function() {
    return this.imageIdx == -1;
};

Entity.prototype.hasNoResults = function() {
    return !this.searching
	&& !this.pristine
	&& !this.error
	&& this.images.length == 0;
};

Entity.prototype.search = function() {
    if (this.query == '') {
        return;
    }

    this.pristine = false;
    this.searching = true;

    var imagesPromise = this.searcher(this.query);
    this.imagesPromise = imagesPromise;
    imagesPromise
	.success(function(data) {
            if (this.imagesPromise != imagesPromise) {
		return;
            }
            this.images = data.images;
	    this.searching = false;
	}.bind(this))
	.error(function() {
	    this.searching = false;
	    this.error = true;
	}.bind(this));
};

var Table = function(Searcher, StateStack) {
    this.entities = [];
    for (var i = 0; i < 3; i++) {
        this.entities[i] = new Entity(Searcher, StateStack);
    }

    this.errorCode = '';
};

Table.prototype.submittable = function() {
    for (var ii = 0; ii < 3; ii++) {
	var entity = this.entities[ii];
	if (entity.name == '') {
	    return false;
	}
	if (entity.imageIdx == -1) {
	    return false;
	}
    }
    return true;
};

Table.prototype.submit = function(event) {
    if (this.entities[0].name == '' ||
	this.entities[1].name == '' ||
	this.entities[2].name == '') {
	this.errorCode = 'NO_NAME';
	console.log(event);
	event.stopPropagation();
	return;
    }
    if (this.entities[0].imageIdx == -1 ||
	this.entities[1].imageIdx == -1 ||
	this.entities[2].imageIdx == -1) {
	this.errorCode = 'NO_IMAGE';
	return;
    }
};

var Page = ['$scope', 'Searcher', 'StateStack', function($scope, Searcher, StateStack) {
    $scope.page = {
        onKeyDown: function(event) {
            if (event.keyCode == 27) {  // esc
                StateStack.pop();
            }
        },
        popState: StateStack.pop.bind(StateStack),
	reset: function() {
	    $scope.table = new Table(Searcher, StateStack);
	},
    };
    $scope.table = new Table(Searcher, StateStack);
}];


// var MakerTable = ['$scope', '$timeout', 'Searcher', 'StateStack', function(
//     $scope, $timeout, Searcher, StateStack) {
//     var entities = [];

//     $scope.table = {
//         entities: entities,
// 	errorCode: '',

// 	reset: function() {
// 	    $scope.table.entities = [];
// 	    for (var i = 0; i < 3; i++) {
// 		$scope.table.entities[i] = new Entity(Searcher, StateStack);
// 	    }
// 	},

// 	submit: function() {
// 	    if (entities[0].name == '' || entities[1].name == '' || entities[2].name == '') {
// 		$scope
// 	    }

// 	    $scope.table.nameNeeded = true;
// 	    $timeout(function() { $scope.table.nameNeeded = false; }, 5000);
// 	},
//     };
// }];

angular.module('mfkMaker.controllers', [])
    .controller('Page', Page);
//    .controller('MakerTable', MakerTable);
