var radiomanApp = angular.module('radiomanApp', ['ngRoute', 'angularMoment', 'angular-loading-bar', 'ngAnimate']);

radiomanApp.controller('MainCtrl', function($scope, $route, $routeParams, $location) {
  $scope.$route = $route;
  $scope.$location = $location;
  $scope.$routeParams = $routeParams;
  $scope.basehref = document.location.protocol + '//' + document.location.host;
});

radiomanApp.controller('ShortcutsCtrl', function($scope, $http) {
  $scope.skipSong = function() {
    $http.post('/api/radios/default/skip-song', {}).success(function (data) {
      console.log('song skipped');
    });
  };
});

radiomanApp.config(function($routeProvider, $locationProvider) {
  $routeProvider
    .when('/playlists', {
      templateUrl: '/static/playlists.html',
      controller: 'PlaylistListCtrl'
    })
    .when('/playlists/:name', {
      templateUrl: '/static/playlist.html',
      controller: 'PlaylistViewCtrl'
    })
    .when('/tracks/:name', {
      templateUrl: '/static/track.html',
      controller: 'TrackViewCtrl'
    })
    .otherwise({
      templateUrl: '/static/home.html',
      controller: 'HomeCtrl'
    });
  // $locationProvider.html5Mode(true);
});

radiomanApp.controller('HomeCtrl', function($scope, $http, $routeParams) {
  $http.get('/api/radios/default').success(function (data) {
    $scope.radio = data.radio;
  });
});

radiomanApp.controller('PlaylistListCtrl', function($scope, $http, $routeParams) {
  $scope.orderByField = 'name';
  $scope.reverseSort = true;
  $http.get('/api/playlists').success(function (data) {
    $scope.playlists = data.playlists;
  });
  $scope.makeDefault = function(playlist) {
    var input = {
      default: true
    };
    $http.patch('/api/playlists/' + playlist.name, input).success(function (data) {
      console.log(data);
    });
  };
});

radiomanApp.controller('PlaylistViewCtrl', function($scope, $http, $routeParams) {
  $scope.orderByField = 'path';
  $scope.reverseSort = true;
  $http.get('/api/playlists/' + $routeParams.name).success(function (data) {
    $scope.playlist = data.playlist;
  });
  $http.get('/api/playlists/' + $routeParams.name + '/tracks').success(function (data) {
    $scope.tracks = data.tracks;
  });
});

radiomanApp.controller('TrackViewCtrl', function($scope, $http, $routeParams) {
  $http.get('/api/tracks/' + $routeParams.name).success(function (data) {
    $scope.track = data.track;
  });
});

radiomanApp.filter('dictToArray', function() {
  return function (obj) {
    if (!(obj instanceof Object)) return obj;
    return _.map(obj, function(val, key) {
      return Object.defineProperty(val, '$key', {__proto__: null, value: key});
    });
  };
});
