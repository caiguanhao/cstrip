angular.module('CS', [ 'ngRoute' ]).

directive('body', [function() {
  return {
    templateUrl: '/views/app.html'
  };
}]).

config(['$routeProvider', '$locationProvider', '$injector',
  function($routeProvider, $locationProvider, $injector) {
  $routeProvider.
  when('/', {
    templateUrl: '/views/home.html',
    controller: 'HomeController as home',
    resolve: {
      cs: ['CommitStrip', function(CommitStrip) {
        return CommitStrip.get();
      }]
    }
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

  var production = false;
  try {
    production = !!$injector.get('PRODUCTION');
  } catch(e) {}

  if (production) {
    $locationProvider.html5Mode(true);
  } else {
    $locationProvider.html5Mode(false);
  }
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

service('CommitStrip', ['$http', '$q', 'current', function($http, $q, current) {
  var that = this;
  this.data = undefined;
  this.get = function() {
    if (this.data) return $q.when(this.data);
    return $http.get('/commitstrip.json').then(function(res) {
      that.data = res.data
      var l = that.data.length;
      for (var i = 0; i < l; i++) {
        that.data[i].image = that.data[i].image.split('\n');
        that.data[i].id = l - i;
      }
      return that.data;
    });
  };
  this.update = function(index, content) {
    var data = { content: content };
    return $http.post('/update/' + index, data, {
      headers: {
        USERNAME: current.username,
        PASSWORD: current.password,
      }
    }).then(function() {
      delete that.data[index]._changed;
    });
  };
}]).

factory('last', [function() {
  return {
    textareaHeight: 0
  };
}]).

factory('current', [function() {
  return {
    username: localStorage['username'],
    password: localStorage['password']
  };
}]).

controller('GlobalController', ['current', function(current) {
  this.current = current;
}]).

controller('HomeController', ['cs', 'current', function(cs, current) {
  this.CommitStrips = cs;
  current.position = {};
}]).

controller('MainController', ['$routeParams', 'cs', 'CommitStrip', 'last',
  '$location', '$scope', 'current', '$timeout',
  function($routeParams, cs, CommitStrip, last, $location, $scope, current, $timeout) {
  var that = this;
  this.cs = cs;
  this.last = last;
  this.id = +$routeParams.CSID || 0;
  this.index = this.cs.length - this.id;

  var dI = localStorage['draft-index'];
  if (dI && this.index === +dI) {
    this.cs[this.index].content = localStorage['draft-content'];
    this.cs[this.index]._changed = true;
  }

  $scope.$watch(function() {
    return that.index;
  }, function() {
    angular.extend(current, that.cs[that.index]);
    current.position = {
      index: that.index,
      id: that.id,
      total: that.cs.length
    };
  });

  this.save = function() {
    CommitStrip.update(that.index, that.cs[that.index].content).then(function() {
      delete localStorage['draft-index'];
      delete localStorage['draft-content'];
    }, function(err) {
      alert(err.statusText);
    });
  };
  this.change = function() {
    that.cs[that.index]._changed = true;
    localStorage['draft-index'] = that.index;
    localStorage['draft-content'] = that.cs[that.index].content;
  };
  this.lS = function(name) {
    localStorage[name] = current[name];
  };
  this.gotoempty = function() {
    var id = 0;
    for (var i = 0; i < that.cs.length; i++) {
      if (!that.cs[i].content) {
        id = that.cs[i].id;
        break;
      }
    }
    if (id > 0) {
      $timeout(function() {
        $location.path(id + '');
      });
    }
  };
}]);
