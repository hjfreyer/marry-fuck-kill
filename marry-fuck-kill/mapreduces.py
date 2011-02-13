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
  """This sets all the default values on an entity's fields.

  It's useful when we add fields.
  """
  yield op.db.Put(entity)


def RandomizeMapper(triple):
  """This regenerates the 'rand' field in a triple.

  We run this periodically with a cronjob to mix up our sort order.
  """
  if triple.is_enabled():
    triple.rand = random.random()
    logging.debug('generate_rand: id=%s rand=%.15f',
                  triple.key().id(), triple.rand)
    yield op.db.Put(triple)


def SetAssignmentTriple(assignment):
  """We didn't always have the 'triple' field in Assignments.

  It's straightforward, but a bit of a pain, to derive it -- this sets it on
  the Assignments that didn't already have it.
  """
  m, f, k = (assignment.marry, assignment.fuck,
             assignment.kill)
  logging.info('m = %s, f = %s, k = %s', m, f, k)
  
  for ref in [m.triple_reference_one_set.get(), 
              m.triple_reference_two_set.get(), 
              m.triple_reference_three_set.get()]:
    if ref is not None:
      triple = ref

  logging.info('triple = %s', triple)
  assignment.triple = triple

  yield op.db.Put(assignment)
