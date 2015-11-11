var liquidmanApp = angular.module('liquidmanApp', ['ngRoute']);

liquidmanApp.controller('MainCtrl', function($scope, $route, $routeParams, $location) {
  $scope.$route = $route;
  $scope.$location = $location;
  $scope.$routeParams = $routeParams;
  $scope.basehref = document.location.protocol + '//' + document.location.host;
});

liquidmanApp.config(function($routeProvider, $locationProvider) {
  $routeProvider
    .when('/playlists', {
      templateUrl: '/static/playlists.html',
      controller: 'PlaylistCtrl'
    })
    .otherwise({
      templateUrl: '/static/home.html',
      controller: 'HomeCtrl'
    });
  // $locationProvider.html5Mode(true);
});

liquidmanApp.controller('HomeCtrl', function($scope, $http, $routeParams) {
});

liquidmanApp.controller('PlaylistCtrl', function($scope, $http, $routeParams) {
});
