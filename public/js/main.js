angular.module('CS', [ 'ngRoute' ]).

config(['$routeProvider', '$locationProvider',
  function($routeProvider, $locationProvider) {
  $routeProvider.
  when('/:CSID?', {
    templateUrl: '/views/main.html',
    controller: 'MainController as main',
    resolve: {
      cs: ['CommitStrip', function(CommitStrip) {
        return CommitStrip.then(function(res) {
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

factory('CommitStrip', ['$http', function($http) {
  return $http.get('/commitstrip.json');
}]).

controller('MainController', ['$routeParams', 'cs',
  function($routeParams, cs) {
  this.cs = cs;
  this.index = +$routeParams.CSID || 0;
}]);
