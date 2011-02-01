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

from google.appengine.ext import db

class EntityValidationError(Exception):
  """There was an error validating an entity."""
  pass

class Entity(db.Model):
  name = db.StringProperty(required=True)
  data = db.BlobProperty(required=True)
  query = db.StringProperty(required=False)
  creator = db.UserProperty()

  # Where images live on the server.
  BASE_URL = '/i/'

  # TODO(mjkelly): do something with this or get rid of it
  type = db.StringProperty(choices=set(["abstract", "person", "object"]))

  def __str__(self):
    return self.name

  def __repr__(self):
    return "Entity(%r)" % self.name

  def json(self):
    return {'name': self.name, 'url': self.get_full_url()}

  def get_full_url(self):
    logging.info("get_full_url: returning %s",
        Entity.BASE_URL + str(self.key()))
    return Entity.BASE_URL + str(self.key())

  def get_stats_url(self, w=160, h=85):
    """Returns the URL to a chart of MFK stats for this Entity.

    Args:
      w: (int) image width
      h: (int) image height
    """
    m = self.assignment_reference_marry_set.count()
    f = self.assignment_reference_fuck_set.count()
    k = self.assignment_reference_kill_set.count()

    url = ('http://chart.apis.google.com/chart'
        '?chxr=0,0,%(max)d'
        '&chxt=y'
        '&chbh=a'
        '&chs=%(w)dx%(h)d'
        '&cht=bvg'
        '&chco=9911BB,C76FDD,63067A'
        '&chds=0,%(max)d,0,%(max)d,0,%(max)d'
        '&chd=t:%(m)d|%(f)d|%(k)d'
        '&chdl=Marry|Fuck|Kill'
        '&chdlp=r' % (dict(m=m, f=f, k=k, max=max([m,f,k]), w=w, h=h)))
    return url

def PutEntity(name, url, query):
  if not url.startswith('http://images.google.com/images?q'):
    raise EntityValidationError("URL must come from Google image search.")

  # This could throw URLError. We'll let it bubble up.
  fh = urllib2.urlopen(url)
  logging.info("Downloading %s" % url)
  data = fh.read()
  logging.info("Downloaded %s bytes" % len(data))

  entity = Entity(name=name, data=data, query=query)

  entity.put()
  return entity

class Triple(db.Model):
  creator = db.UserProperty()

  quality = db.FloatProperty(default=1.0)

  one = db.ReferenceProperty(Entity,
                 collection_name="triple_reference_one_set")
  two = db.ReferenceProperty(Entity,
                 collection_name="triple_reference_two_set")
  three = db.ReferenceProperty(Entity,
                 collection_name="triple_reference_three_set")

  @staticmethod
  def get_random_id():
    # TODO(hjfreyer): Change to a constant-time randomization scheme.
    keys = [k for k in db.Query(Triple, keys_only=True).filter(
        'quality >', 0.0)]
    if not keys:
      return None
    else:
      return random.choice(keys).id()

  @staticmethod
  def get_random():
    id_ = Triple.get_random_id()
    if not id_:
      return None
    else:
      return Triple.get_by_id(id_)

  def __str__(self):
    return Triple.name_from_entities(self.one, self.two, self.three)

  def __repr__(self):
    return "Triple(one=%r, two=%r, three=%r)" % (self.one, self.two, self.three)

  def json(self):
    return {'one': self.one.json(),
        'two': self.two.json(),
        'three': self.three.json(),
        'key': str(self.key())}

  @staticmethod
  def name_from_entities(one, two, three):
    """Given three Entities, generate the canonical key name for a Triple
    containing them.
    """
    keys = [one.name, two.name, three.name]
    keys.sort()
    return "%s.%s.%s" % (keys[0], keys[1], keys[2])

  @staticmethod
  def validate(one, two, three):
    """Verify that the 3 given entities can make a valid triple.

    Raises EntityValidationError if there was a problem.
    """

    if len(set([one.name, two.name, three.name])) < 3:
      raise EntityValidationError("All item names must be distinct.")
    if len(set([one.data, two.data, three.data])) < 3:
      raise EntityValidationError("All item URLs must be distinct.")
    # TODO(mjkelly): implement me
    #raise EntityValidationError("test error")
    pass

def PutTriple(one, two, three):
  """Put a triple in the DB with canonical key.

  See Triple.key_name_from_entities().
  """
  triple = Triple(one=one,
          two=two,
          three=three)
  triple.put()
  return triple

class Assignment(db.Model):
  user = db.UserProperty()

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

  @staticmethod
  def make_assignment(request):
    triple_key = request.get('key')
    values = [request.get('v1'), request.get('v2'), request.get('v3')]

    logging.debug("triple_key = %s", triple_key)
    logging.debug("values = %s", values)

    if set(values) != set(['m', 'f', 'k']):
      return None

    # We get an entity->action map from the client, but we need to reverse
    # it to action->entity to update the DB.
    triple = Triple.get_by_id(int(triple_key))
    logging.debug("triple = %s", triple)
    if triple is None:
      logging.error("No triple with key %s", triple_key)
      return None

    triple_entities = [triple.one,
                       triple.two,
                       triple.three]
    entities = {}
    for i in range(len(values)):
      # Items in values are guaranteed to be 'm', 'f', 'k' (check above)
      entities[values[i]] = triple_entities[i]

    if (entities['m'] is None or entities['f'] is None
        or entities['k'] is None):
      logging.error("Not all non-None: marry = %s, fuck = %s, kill = %s",
              entities['m'], entities['f'], entities['k'])
      return None

    assign = Assignment(marry=entities['m'],
                        fuck=entities['f'],
                        kill=entities['k'])
    assign.put()
    logging.info("Assigned m=%s, f=%s, k=%s to %s", entities['m'],
        entities['f'], entities['k'], triple)
    return assign
