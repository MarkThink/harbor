/*
    Copyright (c) 2016 VMware, Inc. All Rights Reserved.
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
(function() {
  'use strict';

  angular
    .module('harbor.custom')
    .directive('listCustom', listCustom);

  ListCustomController.$inject = ['$scope', 'ListCustomService', '$filter', 'trFilter', '$location', 'getParameterByName'];

  function ListCustomController($scope, ListCustomService, $filter, trFilter, $location, getParameterByName) {

    $scope.subsTabPane = 30;

    var vm = this;

    vm.sectionHeight = {'min-height': '579px'};

    vm.filterInput = '';
    vm.toggleInProgress = [];

    vm.retrieve = retrieve;
    vm.deleteCustom = deleteCustom;
    vm.confirmToDelete = confirmToDelete;
    vm.showAddCustom = showAddCustom;
    vm.isOpen = false;

    var hashValue = $location.hash();
    if(hashValue) {
      var slashIndex = hashValue.indexOf('/');
      if(slashIndex >=0) {
        vm.filterInput = hashValue.substring(slashIndex + 1);
      }else{
        vm.filterInput = hashValue;
      }
    }

    vm.projectId = getParameterByName('project_id', $location.absUrl());
    vm.retrieve();

    $scope.$on('$locationChangeSuccess', function() {
      vm.projectId = getParameterByName('project_id', $location.absUrl());
      vm.filterInput = '';
      vm.retrieve();
    });

    //添加客户成功之后，刷新列表
    $scope.$on('addedSuccess', function(e, val) {
      vm.retrieve();
    });

    function retrieve(){
      //默认请求第0页
      ListCustomService(vm.projectId, vm.filterInput)
        .success(getCustomComplete)
        .error(getCustomFailed);
    }

    //根据根据配置初始化分页
    function getCustomComplete(data, status) {
      //获取客户列表
    }

    function getCustomFailed(response) {
      console.log('Failed to list repositories:' + response);
    }

    function showAddCustom() {
      if(vm.isOpen) {
        vm.isOpen = false;
      }else{
        vm.isOpen = true;
      }
    }

    function confirmToDelete(customId, custom) {
      vm.selectedCustomId = customId;

      $scope.$emit('modalTitle', $filter('tr')('confirm_delete_user_title'));
      $scope.$emit('modalMessage', $filter('tr')('confirm_delete_user', [custom]));

      var emitInfo = {
        'confirmOnly': false,
        'contentType': 'text/plain',
        'action': vm.deleteCustom
      };

      $scope.$emit('raiseInfo', emitInfo);
    }

    function deleteCustom() {
      DeleteCustomService(vm.selectedCustomId)
        .success(deleteCustomSuccess)
        .error(deleteCustomFailed);
    }

    function deleteCustomSuccess(data, status) {
      console.log('Successful delete user.');
      vm.retrieve();
    }

    function deleteCustomFailed(data, status) {
      $scope.$emit('modalTitle', $filter('tr')('error'));
      $scope.$emit('modalMessage', $filter('tr')('failed_to_delete_user'));
      $scope.$emit('raiseError', true);
      console.log('Failed to delete user.');
    }

  }

  function listCustom() {
    var directive = {
      'restrict': 'E',
      'templateUrl': '/static/resources/js/components/custom/list-custom.directive.html',
      'scope': {
        'sectionHeight': '='
      },
      'link': link,
      'controller': ListCustomController,
      'controllerAs': 'vm',
      'bindToController': true
    };

    return directive;

    function link(scope, element, attr, ctrl) {
      element.find('#txtSearchInput').on('keydown', function(e) {
        if($(this).is(':focus') && e.keyCode === 13) {
          ctrl.retrieve();
        }
      });
    }

  }

})();
