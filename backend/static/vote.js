// Copyright 2010 Hunter Freyer and Michael Kelly
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

function mfk_select(element) {
  var classList = element.attr('class').split(/\s+/);
  type = '';
  num = '';
  $.each( classList, function(index, item) {
      if (item === 'm' ||
          item === 'f' ||
          item === 'k') {
        type = item;
      }

      if (item === '1' ||
          item === '2' ||
          item === '3') {
        num = item;
      }
    });

  $('.' + type).removeClass('selected');
  $('.' + num).removeClass('selected');
  element.addClass('selected');

  $('#v' + num).val(type)

  if ($('.m.selected').length > 0 &&
      $('.f.selected').length > 0 &&
      $('.k.selected').length > 0) {
    $('#submit').attr('disabled', false);
  } else {
    $('#submit').attr('disabled', true);
  }
}

function submit(element) {
  if (!element.attr('disabled')) {
    $('#action').val(element.attr("id"));
    $('#voteform').submit();
  }
}

$(document).ready(function(){
    $('.selector').click(function(){ mfk_select($(this)); });
    $('#submit').click(function(){ submit($(this)); });
    $('#skip').click(function(){ submit($(this)); });
});
