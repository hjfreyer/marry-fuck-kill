#!/usr/bin/env python
#
# Copyright 2010 Hunter Freyer and Michael Kelly
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import ezt


def WriteTemplate(template_name, template_dict, out_file):
  template = ezt.Template('templates/%s' % template_name, base_format='html')
  return template.generate(out_file, template_dict)


class EztItem(object):
  """A simple wrapper to convert dict keys to object attributes.

  This makes sending complex objects to EZT much easier.
  """
  def __init__(self, **kwargs):
    vars(self).update(kwargs)
