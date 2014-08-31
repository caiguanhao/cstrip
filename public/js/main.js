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
        if (event.metaKey || event.ctrlKey) {
          var path;
          switch (e.which) {
          case 72:
            path = attrs.twhFirst;
            break;
          case 76:
            path = attrs.twhLast;
            break;
          case 219:
            path = attrs.twhPrev;
            break;
          case 221:
            path = attrs.twhNext;
            break;
          case 69:
            e.preventDefault();
            scope.$eval(attrs.twhGotoempty);
            return;
          case 13:
            e.preventDefault();
            scope.$eval(attrs.twhSave);
            return;
          }
          if (path) {
            e.preventDefault();
            $location.path(path);
            scope.$apply();
          }
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
      var load = function() {
        elem.removeClass('ng-hide');
        if (angular.isString(attrs.assignHeightTo)) {
          scope.$parent.$apply(function() {
            var exp = attrs.assignHeightTo + '=' + elem[0].height;
            scope.$parent.$eval(exp);
          });
        }
        elem.off('load', load);
      };
      elem.on('load', load);
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

factory('last', [function() {
  return {
    textareaHeight: 0
  };
}]).

controller('MainController', ['$routeParams', 'cs', 'CommitStrip', 'last', '$location', '$scope',
  function($routeParams, cs, CommitStrip, last, $location, $scope) {
  var that = this;
  this.cs = cs;
  this.last = last;
  this.index = +$routeParams.CSID || 0;
  this.save = function() {
    CommitStrip.update(that.index, that.cs[that.index].content);
  };
  this.gotoempty = function() {
    var index = -1;
    for (var i = 0; i < that.cs.length; i++) {
      if (!that.cs[i].content) {
        index = i;
        break;
      }
    }
    if (index > -1) {
      $location.path(index + '');
      $scope.$apply();
    }
  };
}]);
