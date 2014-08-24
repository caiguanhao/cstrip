angular.module('CS', [ 'ngRoute' ]).

config(['$routeProvider', '$locationProvider',
  function($routeProvider, $locationProvider) {
  $routeProvider.
  when('/', {
    redirectTo: '/0'
  }).
  when('/:CSID?', {
    templateUrl: '/views/main.html',
    controller: 'MainController as main',
    resolve: {
      cs: ['CommitStrip', function(CommitStrip) {
        return CommitStrip.get();
      }]
    }
  }).
  otherwise({
    redirectTo: '/'
  });

  $locationProvider.html5Mode(false);
}]).

directive('textareaWithHotkeys', ['$location', function($location) {
  return {
    link: function(scope, elem, attrs) {
      elem[0].focus();
      elem.on('keydown', function(e) {
        var isPrev = (e.which === 219 && (event.metaKey || event.ctrlKey));
        var isNext = (e.which === 221 && (event.metaKey || event.ctrlKey));
        var isSave = (e.which === 13  && (event.metaKey || event.ctrlKey));
        if (isPrev) {
          e.preventDefault();
          $location.path(attrs.twhPrev)
          scope.$apply();
        } else if (isNext) {
          e.preventDefault();
          $location.path(attrs.twhNext)
          scope.$apply();
        } else if (isSave) {
          e.preventDefault();
          scope.$eval(attrs.twhSave)
        }
      });
    }
  };
}]).

directive('showWhenLoadCompletes', [function() {
  return {
    scope: {
      src: '@ngSrc'
    },
    link: function(scope, elem, attrs) {
      attrs.$observe('src', function(src) {
        if (!src) return;
        elem.addClass('ng-hide');
      });
      elem.on('load', function() {
        elem.removeClass('ng-hide');
      });
    }
  };
}]).

service('CommitStrip', ['$http', '$q', function($http, $q) {
  var that = this;
  this.data = undefined;
  this.get = function() {
    if (this.data) return $q.when(this.data);
    return $http.get('/commitstrip.json').then(function(res) {
      that.data = res.data
      return that.data;
    });
  };
  this.update = function(index, content) {
    var data = { content: content };
    return $http.post('/update/' + index, data).then(function() {
      delete that.data[index]._changed;
    });
  };
}]).

controller('MainController', ['$routeParams', 'cs', 'CommitStrip',
  function($routeParams, cs, CommitStrip) {
  var that = this;
  this.cs = cs;
  this.index = +$routeParams.CSID || 0;
  this.save = function() {
    CommitStrip.update(that.index, that.cs[that.index].content);
  };
}]);
