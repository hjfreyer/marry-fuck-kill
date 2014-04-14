'use strict';

/* Services */

var StateStack = function() {
    this.stack_ = [];
};

StateStack.prototype.push = function(cb, finallyCb) {
    this.stack_.push([cb, finallyCb]);
};

StateStack.prototype.pop = function() {
    if (this.stack_.length > 0) {
	var top = this.stack_.pop();
	top[0]();
	if (top[1] != null) {
	    top[1]();
	}
    }
};

StateStack.prototype.release = function() {
    var top = this.stack_.pop();
	if (top[1] != null) {
	    top[1]();
	}
};

var FakeSearcher = ['$q', function($q) {
    return function(query) {
	var d = $q.defer();
	d.resolve({"data" :
	 {"images": [{"original": "http://mkalty.org/wp-content/uploads/2014/03/o-CATS-KILL-BILLIONS-facebook.jpg", "thumbnail": "https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcR2Q9KqDQLWmgytQCr4SxPt_kuVQeU6puzKDUkvY1Km_39heIyz3qeEcXn5Rw", "sig": "a833cfb0d02fdf63295bcc63662ea97d"}, {"original": "http://jasonlefkowitz.net/wp-content/uploads/2013/07/Cute-Cats-cats-33440930-1280-800.jpg", "thumbnail": "https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcTls144Go4PI4DMrYMzbgC56s6tSbtI21fd16PbDUDshhJlqIX29F536A0", "sig": "b1d69ab39a8b3ddd5d027ffdca68556a"}, {"original": "http://4.bp.blogspot.com/-MzZCzWI_6Xc/UIUQp1qPfzI/AAAAAAAAHpA/OTwHCJSWFAY/s1600/cats_animals_kittens_cat_kitten_cute_desktop_1680x1050_hd-wallpaper-753974.jpeg", "thumbnail": "https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcTcRAvhEUtY5Y_vJBggDDHLW-WSabMrGmEhI193uunmioyiQgSCrG8nX6Gl", "sig": "2ac08ea94524b97bce2bb57cff3f81f0"}, {"original": "http://jasonlefkowitz.net/wp-content/uploads/2013/07/cats-16140154-1920-1080.jpg", "thumbnail": "https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcTZwZYwlbuU9xJQ8IA4ZYqkouFHAlIHNuUu8WubwsP4Pjx7zW4rO5K5TUCr", "sig": "4af8745bc0c8bf79d66c96e5f547fc50"}, {"original": "http://lolomoda.com/wp-content/uploads/30192_1600x1200-4-cute-cats.jpg", "thumbnail": "https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcRuKanhuz1FsY5dJnoqwnF2cW8pfX7qVaSnUBJ0zoBxOztbrpletPZrz14o", "sig": "7e8fbb3bb7d5b5dc51b8fb8ee9b60c7e"}, {"original": "http://jasonlefkowitz.net/wp-content/uploads/2013/07/Cats.jpg", "thumbnail": "https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcQNh7S3A-QOUgpSSBNBtpFFFv1HJL1OevgDis2z9JQ06_I--3EJAQ75__Bv", "sig": "290f078814ec7509e95ae63ca0dec53c"}, {"original": "http://cdn.cstatic.net/images/gridfs/519cdcdff92ea153e404a54c/2013_critical_cat_!23.jpg", "thumbnail": "https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcTqP1BboFOnOjwcUSTNJ7PW0oFnY6S11qozUJvu0HgOYVtnwuzj3gAGwZs", "sig": "66f127b7f1ce9ce224ea5988218e605c"}, {"original": "http://www.andrew.cmu.edu/user/cfperron/cats/images/cat8.jpg", "thumbnail": "https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcRIQ8ioABJGiMr9_nOIjynSa6FjN61zqOY-oW2FRM81ALpv_ictmfNkSx2c", "sig": "91a08fc0453b82166463908c391e3abc"}, {"original": "http://affordablecatanddoghospital.com/cat/images/stories/slide/cats-blue.jpg", "thumbnail": "https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcSVOMKwOVLJXkCsOyjbecsDYLaNnrgZutKh0E5OlbL58jXMe7wIn_t1FKw", "sig": "95683f354af1d40fd29eb7756c3994e8"}, {"original": "http://sfoxwriting.files.wordpress.com/2013/09/cute-cat-pictures.jpeg", "thumbnail": "https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcT7iJ1b45ZnLuHTw2NNfzXDJX15YZjI_OKqN-Qjax0B1jRK11KCQbjQ04uO", "sig": "28bc647765393a8085be730dd475d2ee"}], "time": "2014-04-03 02:35:05.911086"}});
	return d.promise;
    }
}];

var Searcher = ['searchPrefix', '$http',  function(searchPrefix, $http) {
    return function(query) {
	return $http.get(searchPrefix + query);
    };
}];



// Demonstrate how to register services
// In this case it is a simple value service.
angular.module('mfkMaker.services', [])
    .value('version', '0.1')
    .value('searchPrefix', '/search?q=')
    .service('StateStack', StateStack)
    .factory('Searcher', Searcher);
