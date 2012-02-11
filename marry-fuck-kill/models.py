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

import logging
import random
import urllib2

from google.appengine.api import users
from google.appengine.ext import db

class Entity(db.Model):
  name = db.StringProperty(default=None)
  data = db.BlobProperty(default=None)
  query = db.StringProperty(default=None)
  creator = db.UserProperty()
  original_url = db.StringProperty(default=None)

  # Where images live on the server.
  BASE_URL = '/i/'

  @property
  def image_url(self):
    return Entity.BASE_URL + str(self.key())

  def __str__(self):
    return self.name

  def __repr__(self):
    return 'Entity(%r)' % self.name


class Triple(db.Model):
  # User who created the Triple, or None if no user was logged in.
  creator = db.UserProperty()
  # These fields are not required for DB backwards-compatibility.
  creatorip = db.StringProperty(default='')
  time = db.DateTimeProperty(auto_now_add=True)
  # A random value used to determine a randomized, repeatable display order.
  # Not necessarily unique. 'default' value is for backwards-compatibility.
  # If this is None, this Triple will not be returned by get_next_id() under
  # any cirumstances.
  rand = db.FloatProperty(default=0.0)

  one = db.ReferenceProperty(Entity,
      collection_name="triple_reference_one_set")
  two = db.ReferenceProperty(Entity,
      collection_name="triple_reference_two_set")
  three = db.ReferenceProperty(Entity,
      collection_name="triple_reference_three_set")

  reviewed = db.BooleanProperty(default=False)

  CONTEXT_COOKIE_NAME = 'mfkcontext'
  # week in seconds
  CONTEXT_COOKIE_MAX_AGE_SECS = 604800

  @property
  def id_string(self):
    return self.key().id()

  @property
  def creator_nickname(self):
    return self.creator and self.creator.nickname() or ''

  def disable(self):
    """Prevent this Triple from being picked in the random rotation.

    It will still be visible from a direct link.
    """
    self.rand = -1.0

  def enable(self):
    """Allow this Triple to be picked in the random rotation."""
    self.rand = random.random()

  @property
  def enabled(self):
    return self.rand != -1.0

  @staticmethod
  def get_next_id(request, response):
    """Gets the next Triple ID, apparently at random.

    We use cookies in the request, if available, to determine what the next
    Triple will be.

    Args:
      request: request object from the client
      response: the response object that will be sent to the client.

    Returns:
      (ID) The ID of the next Triple, and a cookie to set on the client
           to maintain state.
    """
    # We _could_ also store the random value on the client side if we want to
    # cut down on the number of DB requests we make.
    prev_id = None
    next_id = None
    try:
      # This fails if there is no cookie, or it isn't a valid long.
      prev_id = long(request.cookies.get(Triple.CONTEXT_COOKIE_NAME))
    except TypeError:
      pass
    logging.info('get_next_id: prev_id = %s', prev_id)
    if prev_id is not None:
      prev_triple = Triple.get_by_id(long(prev_id))
      if prev_triple is not None:
        prev_rand = prev_triple.rand
        next_id = Triple._get_greater_rand(prev_rand)

    # If we failed to intelligently determine the next triple, pick at random.
    if next_id is None:
      next_id = Triple._get_random_id()

    if next_id is not None:
      response.headers.add_header('Set-Cookie', '%s=%d; Max-Age=%d' % (
          Triple.CONTEXT_COOKIE_NAME, next_id,
          Triple.CONTEXT_COOKIE_MAX_AGE_SECS))

    return next_id

  @staticmethod
  def _get_greater_rand(rand, _or_equal=False):
    """Get the Triple ID with the next-greater rand value.

    This method contains all the real logic for picking random values.

    Args:
      rand: The value to compare against.
    """
    if _or_equal:
      logging.info('_get_greater_rand: rand >= %s', rand)
      query = db.Query(Triple, keys_only=True).filter(
          'rand >=', rand).order('rand').order('__key__')
    else:
      logging.info('_get_greater_rand: rand > %s', rand)
      query = db.Query(Triple, keys_only=True).filter(
          'rand >', rand).order('rand').order('__key__')
    result = query.get()
    if result is None:
      if not _or_equal:
        # If we were using a rand value > 0, we probably hit the end of the
        # list naturally. We want to wrap around to the lowest triple, so we
        # use rand >= 0.0.
        return Triple._get_greater_rand(0.0, _or_equal=True)
      else:
        # If we returned no results and yet were using a rand value < 0, the
        # database is empty or contains only invalid rand values. Give up. This
        # is the base case for the recursion above.
        logging.error('_get_greater_rand: Found no Triples! Giving up.')
        return None
    else:
      return result.id()

  @staticmethod
  def _get_random_id():
    """Returns a truly random Triple.

    This is a fallback used by get_next_id for when we don't have any context
    on previously-seen Triples from the client.
    """
    rand = random.random()
    logging.info('_get_random_id: rand = %f', rand)
    return Triple._get_greater_rand(rand)

  def __str__(self):
    return Triple.name_from_entities(self.one, self.two, self.three)

  def __repr__(self):
    return "Triple(one=%r, two=%r, three=%r)" % (self.one, self.two, self.three)

  @staticmethod
  def name_from_entities(one, two, three):
    """Given three Entities, generate the canonical key name for a Triple
    containing them.
    """
    keys = [one.name, two.name, three.name]
    keys.sort()
    return "%s.%s.%s" % (keys[0], keys[1], keys[2])


class Assignment(db.Model):
  user = db.UserProperty()
  userip = db.StringProperty(default='')
  time = db.DateTimeProperty(auto_now_add=True)

  triple = db.ReferenceProperty(Triple)

  marry = db.ReferenceProperty(Entity,
      collection_name="assignment_reference_marry_set")
  fuck = db.ReferenceProperty(Entity,
      collection_name="assignment_reference_fuck_set")
  kill = db.ReferenceProperty(Entity,
      collection_name="assignment_reference_kill_set")

  def __str__(self):
    return "Assignment(marry=%s, fuck=%s, kill=%s)" % (
        self.marry.key().name(),
        self.fuck.key().name(),
        self.kill.key().name())

  def __repr__(self):
    return str(self)
