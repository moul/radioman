var radiomanApp = angular.module('radiomanApp', ['ngRoute']);

radiomanApp.controller('MainCtrl', function($scope, $route, $routeParams, $location) {
  $scope.$route = $route;
  $scope.$location = $location;
  $scope.$routeParams = $routeParams;
  $scope.basehref = document.location.protocol + '//' + document.location.host;
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
    .otherwise({
      templateUrl: '/static/home.html',
      controller: 'HomeCtrl'
    });
  // $locationProvider.html5Mode(true);
});

radiomanApp.controller('HomeCtrl', function($scope, $http, $routeParams) {
});

radiomanApp.controller('PlaylistListCtrl', function($scope, $http, $routeParams) {
  $http.get('/api/playlists').success(function (data) {
    $scope.playlists = data.playlists;
  });
});

radiomanApp.controller('PlaylistViewCtrl', function($scope, $http, $routeParams) {
  $http.get('/api/playlists/' + $routeParams.name).success(function (data) {
    $scope.playlist = data.playlist;
  });
});
