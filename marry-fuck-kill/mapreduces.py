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
#

import random
import logging

from mapreduce import operation as op

def RecreateMapper(entity):
  yield op.db.Put(entity)


def RandomizeMapper(triple):
  if triple.is_enabled():
    triple.rand = random.random()
    logging.debug('generate_rand: id=%s rand=%.15f',
                  triple.key().id(), triple.rand)
    yield op.db.Put(triple)

