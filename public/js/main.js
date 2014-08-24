angular.module('CS', [ 'ngRoute' ]).

config(['$routeProvider', '$locationProvider',
  function($routeProvider, $locationProvider) {
  $routeProvider.
  when('/:CSID?', {
    templateUrl: '/views/main.html',
    controller: 'MainController as main',
    resolve: {
      cs: ['CommitStrip', function(CommitStrip) {
        return CommitStrip.get().then(function(res) {
          return res.data;
        });
      }]
    }
  }).
  otherwise({
    redirectTo: '/'
  });

  $locationProvider.html5Mode(false);
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

service('CommitStrip', ['$http', function($http) {
  this.get = function() {
    return $http.get('/commitstrip.json');
  };
  this.update = function(index, content) {
    var data = { content: content };
    return $http.post('/update/' + index, data);
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
