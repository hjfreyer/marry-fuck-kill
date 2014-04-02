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

var Entity = function() {
    this.name = '';
    this.imgIdx = -1;
    this.query = '';
    this.images = [];

    this.animated = true;
    this.position = 0;
};

Entity.prototype.selectImage = function(idx) {
    this.imgIdx = idx;
};

Entity.prototype.getImage = function() {
    if (this.imgIdx == -1) {
	return '/s/mfk.png';
    }
    return this.images[this.imgIdx];
};

Entity.prototype.getReset = function() {
    return function(name, imgIdx, images, position) {
	this.name = name;
	this.imgIdx = imgIdx;
	this.images = images;
	this.position = position;
    }.bind(this, this.name, this.imgIdx, this.images, this.position);
};

var MakerTable = ['$scope', 'Searcher', 'StateStack', function(
    $scope, Searcher, StateStack) {
    $scope.table = {
	searchIdx: -1,
	entities: [new Entity(), new Entity(), new Entity()],

	select: function(idx) {
	    var reset = $scope.table.entities[idx].getReset();
	    $scope.table.searchIdx = idx;
	    StateStack.push(function() {
		$scope.table.searchIdx = -1;
		reset();
	    });
	},

	showSearch: function(entity) {
	    entity.position = 1;
	},

	search: function(entity) {
	    entity.images = Searcher('cats').then(function(x) {
		console.log(x);
		entity.images = x.data.images;
	    });
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
