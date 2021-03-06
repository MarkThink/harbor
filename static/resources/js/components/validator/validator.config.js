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
    .module('harbor.validator')
    .constant('INVALID_CHARS', [",","~","#", "$", "%"])
    .constant('PASSWORD_REGEXP', /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?!.*\s).{7,20}$/)
    .constant('PROJECT_REGEXP', /^[a-z0-9](?:-*[a-z0-9])*(?:[._][a-z0-9](?:-*[a-z0-9])*)*$/)
    .constant('TAG_REGEXP', /^[a-z0-9](?:-*[a-z0-9])*(?:[._][a-z0-9](?:-*[a-z0-9])*)*$/) //标记的英文名
    .constant('CH_REGEXP', /[\u4e00-\u9fa5]/); //中文名
})();
