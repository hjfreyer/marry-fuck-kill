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

import core
import models
from mapreduce import operation as op

from google.appengine.ext import db

def RecreateMapper(entity):
  """This sets all the default values on an entity's fields.

  It's useful when we add fields.
  """
  yield op.db.Put(entity)


def RandomizeMapper(triple):
  """This regenerates the 'rand' field in a triple.

  We run this periodically with a cronjob to mix up our sort order.
  """
  if triple.enabled:
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


def CalculateVoteCounts(triple):
  """Calculates vote counts and puts them directly in Triples.

  We used to always calculate them on the fly from Assignment counts, which is
  expensive.
  """
  votes = []
  for entity in [triple.one, triple.two, triple.three]:
    votes.append([entity.assignment_reference_marry_set.count(),
                  entity.assignment_reference_fuck_set.count(),
                  entity.assignment_reference_kill_set.count()])
  assert(len(votes) == 3)
  assert(len(votes[0]) == 3)
  assert(len(votes[1]) == 3)
  assert(len(votes[2]) == 3)
  logging.info('Calculated vote counts for Triple %s (%s) = %s', triple, triple.key(), votes)

  # This is the race -- someone else could have voted in the meantime.

  db.run_in_transaction(core._UpdateTripleVoteCounts, triple.key(), votes)

def ClearVoteCounts(triple):
  """Clears all calculated vote counts. For debugging."""
  logging.info('Clearing vote counts for Triple %s (%s)', triple, triple.key())

  def _ClearVotes(triple_key):
    triple = models.Triple.get(triple_key)

    (triple.votes_one_m, triple.votes_one_f, triple.votes_one_k,
        triple.votes_two_m, triple.votes_two_f, triple.votes_two_k,
        triple.votes_three_m, triple.votes_three_f, triple.votes_three_k) =  (
            None, None, None, None, None, None, None, None, None)
    triple.has_cached_votes = False

    triple.put()

  db.run_in_transaction(_ClearVotes, triple.key())
