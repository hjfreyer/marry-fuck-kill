'use strict';

/* Controllers */

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

Entity.prototype.isPlaceholder = function() {
    return this.imgIdx == -1;
};

Entity.prototype.search = function() {
    if (this.query == '') {
        return;
    }

    var imagesPromise = this.searcher(this.query);
    this.imagesPromise = imagesPromise;
    imagesPromise.then(function(result) {
        if (this.imagesPromise != imagesPromise) {
            return;
        }
        this.images = result.data.images;
    }.bind(this));
};

var Table = function(Searcher, StateStack) {
    this.entities = [];
    for (var i = 0; i < 3; i++) {
        this.entities[i] = new Entity(Searcher, StateStack);
    }

    this.errorCode = '';
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
    if (this.entities[0].imgIdx == -1 ||
	this.entities[1].imgIdx == -1 ||
	this.entities[2].imgIdx == -1) {
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
